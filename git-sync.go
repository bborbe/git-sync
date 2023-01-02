/*
Copyright 2014 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// git-sync is a command that pull a git repository to a local directory.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	defaultWait = 300
	errorLimit  = 5
)

var gitRepo = flag.String("repo", envString("GIT_SYNC_REPO", ""), "git repo url")
var gitBranch = flag.String("branch", envString("GIT_SYNC_BRANCH", "master"), "git branch")
var gitRev = flag.String("rev", envString("GIT_SYNC_REV", "HEAD"), "git rev")
var gitDepthSync = flag.Int("depth", envInt("GIT_SYNC_DEPTH", 0), "shallow clone with a history truncated to the specified number of commits")
var gitUser = flag.String("username", envString("GIT_SYNC_USERNAME", ""), "username")
var gitPassword = flag.String("password", envString("GIT_SYNC_PASSWORD", ""), "password")
var targetDirectory = flag.String("dest", envString("GIT_SYNC_DEST", ""), "destination path")
var wait = flag.Int("wait", envInt("GIT_SYNC_WAIT", defaultWait), "number of seconds to wait before next sync")
var oneTime = flag.Bool("one-time", envBool("GIT_SYNC_ONE_TIME", false), "exit after the initial checkout")
var permissions = flag.Int("change-permissions", envInt("GIT_SYNC_PERMISSIONS", 0), `If set it will change the permissions of the directory 
		that contains the git repository. Example: 744`)
var callbackUrl = flag.String("callback-url", envString("CALLBACK_URL", ""), "url to call after each git pull")

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string, def bool) bool {
	if env := os.Getenv(key); env != "" {
		res, err := strconv.ParseBool(env)
		if err != nil {
			return def
		}

		return res
	}
	return def
}

func envInt(key string, def int) int {
	if env := os.Getenv(key); env != "" {
		val, err := strconv.Atoi(env)
		if err != nil {
			glog.V(2).Infof("invalid value for %q: using default: %q", key, def)
			return def
		}
		return val
	}
	return def
}

const usage = "usage: GIT_SYNC_REPO= GIT_SYNC_DEST= [GIT_SYNC_BRANCH= GIT_SYNC_WAIT= GIT_SYNC_DEPTH= GIT_SYNC_USERNAME= GIT_SYNC_PASSWORD= GIT_SYNC_ONE_TIME=] git-sync -repo GIT_REPO_URL -dest PATH [-branch -wait -username -password -depth -one-time]"

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *gitRepo == "" || *targetDirectory == "" {
		flag.Usage()
		glog.Error(usage)
		glog.Flush()
		os.Exit(1)
	}
	glog.V(0).Infof("sync repo %v to %v", *gitRepo, *targetDirectory)
	if _, err := exec.LookPath("git"); err != nil {
		glog.Exitf("required git executable not found: %v", err)
	}

	if *gitUser != "" && *gitPassword != "" {
		if err := setupGitAuth(*gitUser, *gitPassword, *gitRepo); err != nil {
			glog.Exitf("error creating .netrc file: %v", err)
		}
	}

	var errorCounter int
	for {
		if err := syncRepo(*gitRepo, *targetDirectory, *gitBranch, *gitRev, *gitDepthSync); err != nil {
			glog.Errorf("error syncing repo: %v", err)
			errorCounter++
		} else {
			errorCounter = 0
		}
		if errorCounter > errorLimit {
			glog.Exitf("error limit of %d exceeded", errorLimit)
		}

		if *oneTime {
			glog.Flush()
			os.Exit(0)
		}

		glog.V(2).Infof("wait %d seconds", *wait)
		time.Sleep(time.Duration(*wait) * time.Second)
		glog.V(2).Infof("done")
	}
}

// syncRepo syncs the branch of a given repository to the destination at the given rev.
func syncRepo(repo, dest, branch, rev string, depth int) error {
	gitRepoPath := path.Join(dest, ".git")
	_, err := os.Stat(gitRepoPath)
	switch {
	case os.IsNotExist(err):
		// clone repo
		args := []string{"clone", "--no-checkout", "-b", branch}
		if depth != 0 {
			args = append(args, "-depth")
			args = append(args, strconv.Itoa(depth))
		}
		args = append(args, repo)
		args = append(args, dest)
		output, err := runCommand("git", "", args)
		if err != nil {
			return err
		}

		glog.V(2).Infof("clone %q: %s", repo, string(output))
	case err != nil:
		return fmt.Errorf("error checking if repo exist %q: %v", gitRepoPath, err)
	}

	// set remote url
	output, err := runCommand("git", dest, []string{"remote", "set-url", "origin", repo})
	if err != nil {
		return err
	}

	glog.V(2).Infof("set remote-url to %s: %s", repo, string(output))

	// fetch branch
	output, err = runCommand("git", dest, []string{"pull", "origin", branch})
	if err != nil {
		return err
	}

	glog.V(2).Infof("fetch %q: %s", branch, string(output))

	// reset working copy
	output, err = runCommand("git", dest, []string{"reset", "--hard", rev})
	if err != nil {
		return err
	}

	glog.V(2).Infof("reset %q: %v", rev, string(output))

	if *permissions != 0 {
		// set file permissions
		_, err = runCommand("chmod", "", []string{"-R", strconv.Itoa(*permissions), dest})
		if err != nil {
			return err
		}
	}

	if *callbackUrl != "" {
		glog.V(4).Infof("get url %s", *callbackUrl)
		resp, err := http.Get(*callbackUrl)
		if err != nil {
			return errors.Wrapf(err, "get url %s failed", *callbackUrl)
		}
		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("request to %s failed with statusCode %d", *callbackUrl, resp.StatusCode)
		}
		glog.V(1).Infof("url %s called successful", *callbackUrl)
	}

	return nil
}

func runCommand(command, cwd string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("error running command %q : %v: %s", strings.Join(cmd.Args, " "), err, string(output))
	}

	return output, nil
}

func setupGitAuth(username, password, gitURL string) error {
	glog.V(2).Infof("setting up the git credential cache")
	cmd := exec.Command("git", "config", "--global", "credential.helper", "cache")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error setting up git credentials %v: %s", err, string(output))
	}

	glog.V(2).Infof("git credential approve")
	cmd = exec.Command("git", "credential", "approve")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	glog.V(4).Infof("url=%s", gitURL)
	fmt.Fprintf(stdin, "url=%s\n", gitURL)
	glog.V(4).Infof("username=%s", username)
	fmt.Fprintf(stdin, "username=%s\n", username)
	glog.V(4).Infof("password=%s", password)
	fmt.Fprintf(stdin, "password=%s\n", password)
	glog.V(4).Infof("write creds finished")
	stdin.Close()
	glog.V(4).Infof("stdin closed")

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error setting up git credentials %v", err)
	}
	glog.V(2).Infof("setting up the git credential cache completed")

	return nil
}
