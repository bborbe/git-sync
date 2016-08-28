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
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"runtime"

	"github.com/golang/glog"
)

const (
	DEFAULT_WAIT = 300
	ERROR_LIMIT  = 5
)

var flRepo = flag.String("repo", envString("GIT_SYNC_REPO", ""), "git repo url")
var flBranch = flag.String("branch", envString("GIT_SYNC_BRANCH", "master"), "git branch")
var flRev = flag.String("rev", envString("GIT_SYNC_REV", "HEAD"), "git rev")
var flDest = flag.String("dest", envString("GIT_SYNC_DEST", ""), "destination path")
var flWait = flag.Int("wait", envInt("GIT_SYNC_WAIT", DEFAULT_WAIT), "number of seconds to wait before next sync")
var flOneTime = flag.Bool("one-time", envBool("GIT_SYNC_ONE_TIME", false), "exit after the initial checkout")
var flDepth = flag.Int("depth", envInt("GIT_SYNC_DEPTH", 0), "shallow clone with a history truncated to the specified number of commits")

var flUsername = flag.String("username", envString("GIT_SYNC_USERNAME", ""), "username")
var flPassword = flag.String("password", envString("GIT_SYNC_PASSWORD", ""), "password")

var flChmod = flag.Int("change-permissions", envInt("GIT_SYNC_PERMISSIONS", 0), `If set it will change the permissions of the directory 
		that contains the git repository. Example: 744`)

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

	if *flRepo == "" || *flDest == "" {
		flag.Usage()
		glog.Error(usage)
		glog.Flush()
		os.Exit(0)
	}
	if _, err := exec.LookPath("git"); err != nil {
		glog.Exitf("required git executable not found: %v", err)
	}

	if *flUsername != "" && *flPassword != "" {
		if err := setupGitAuth(*flUsername, *flPassword, *flRepo); err != nil {
			glog.Exitf("error creating .netrc file: %v", err)
		}
	}

	var errorCounter int
	for {
		if err := syncRepo(*flRepo, *flDest, *flBranch, *flRev, *flDepth); err != nil {
			glog.Errorf("error syncing repo: %v", err)
			errorCounter++
		} else {
			errorCounter = 0
		}
		if errorCounter > ERROR_LIMIT {
			glog.Exitf("error limit of %d exceeded", ERROR_LIMIT)
		}

		if *flOneTime {
			glog.Flush()
			os.Exit(0)
		}

		glog.V(2).Infof("wait %d seconds", *flWait)
		time.Sleep(time.Duration(*flWait) * time.Second)
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
			args = append(args, string(depth))
		}
		args = append(args, repo)
		args = append(args, dest)
		output, err := runCommand("git", "", args)
		if err != nil {
			return err
		}

		glog.Infof("clone %q: %s", repo, string(output))
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

	glog.Infof("fetch %q: %s", branch, string(output))

	// reset working copy
	output, err = runCommand("git", dest, []string{"reset", "--hard", rev})
	if err != nil {
		return err
	}

	glog.V(2).Infof("reset %q: %v", rev, string(output))

	if *flChmod != 0 {
		// set file permissions
		_, err = runCommand("chmod", "", []string{"-R", string(*flChmod), dest})
		if err != nil {
			return err
		}
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
