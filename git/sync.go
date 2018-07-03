// Copyright 2018 The Git-Sync Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git

import (
	"github.com/pkg/errors"
	"net/http"
	"context"
	"github.com/golang/glog"
	"fmt"
	"time"
	"os/exec"
	"strings"
	"path"
	"os"
)

//go:generate counterfeiter -o ../mocks/roundtripper.go --fake-name RoundTripper . roundTripper
type roundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}

type Sync struct {
	GitRepo         string
	GitBranch       string
	GitRev          string
	GitDepthSync    int
	GitUser         string
	GitPassword     string
	TargetDirectory string
	Permissions     int
	CallbackUrl     string
	Transport       http.RoundTripper
}

func (s *Sync) Validate() error {
	if s.GitRepo == "" {
		return errors.New("git repo missing")
	}
	if s.GitBranch == "" {
		return errors.New("git branch missing")
	}
	if s.GitRev == "" {
		return errors.New("git rev missing")
	}
	if _, err := exec.LookPath("git"); err != nil {
		return errors.Errorf("required git executable not found: %v", err)
	}
	return nil
}

func (s *Sync) Run(ctx context.Context) error {
	glog.V(0).Infof("sync repo %v to %v", s.GitRepo, s.TargetDirectory)

	if s.GitUser != "" && s.GitPassword != "" {
		if err := setupGitAuth(s.GitUser, s.GitPassword, s.GitRepo); err != nil {
			return errors.Errorf("error creating .netrc file: %v", err)
		}
	}

	gitRepoPath := path.Join(s.TargetDirectory, ".git")
	_, err := os.Stat(gitRepoPath)
	switch {
	case os.IsNotExist(err):
		// clone repo
		args := []string{"clone", "--no-checkout", "-b", s.GitBranch}
		if s.GitDepthSync != 0 {
			args = append(args, "-depth")
			args = append(args, string(s.GitDepthSync))
		}
		args = append(args, s.GitRepo)
		args = append(args, s.TargetDirectory)
		output, err := runCommand("git", "", args)
		if err != nil {
			return err
		}

		glog.V(2).Infof("clone %q: %s", s.GitRepo, string(output))
	case err != nil:
		return fmt.Errorf("error checking if repo exist %q: %v", gitRepoPath, err)
	}

	// set remote url
	output, err := runCommand("git", s.TargetDirectory, []string{"remote", "set-url", "origin", s.GitRepo})
	if err != nil {
		return err
	}

	glog.V(2).Infof("set remote-url to %s: %s", s.GitRepo, string(output))

	// fetch branch
	output, err = runCommand("git", s.TargetDirectory, []string{"pull", "origin", s.GitBranch})
	if err != nil {
		return err
	}

	glog.V(2).Infof("fetch %q: %s", s.GitBranch, string(output))

	// reset working copy
	output, err = runCommand("git", s.TargetDirectory, []string{"reset", "--hard", s.GitRev})
	if err != nil {
		return err
	}

	glog.V(2).Infof("reset %q: %v", s.GitRev, string(output))

	if s.Permissions != 0 {
		// set file permissions
		_, err = runCommand("chmod", "", []string{"-R", string(s.Permissions), s.TargetDirectory})
		if err != nil {
			return err
		}
	}

	if s.CallbackUrl != "" {
		req, err := http.NewRequest("GET", s.CallbackUrl, nil)
		if err != nil {
			return errors.Errorf("create http request failed: %v", err)
		}
		resp, err := s.Transport.RoundTrip(req)
		if err != nil {
			return errors.Errorf("perform http request failed: %v", err)
		}
		if resp.StatusCode/100 != 2 {
			return fmt.Errorf("request to %s failed with statusCode %d", s.CallbackUrl, resp.StatusCode)
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
