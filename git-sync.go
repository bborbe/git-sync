// Copyright 2018 The Git-Sync Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	flag "github.com/bborbe/flagenv"
	"runtime"
	"github.com/golang/glog"
	"github.com/bborbe/git-sync/git"
	"os"
	"net/http"
	"context"
	"time"
)

const (
	DEFAULT_WAIT = 300
	ERROR_LIMIT  = 5
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())

	gitSync := &git.Sync{
		Transport: http.DefaultTransport,
	}
	flag.StringVar(&gitSync.GitRepo, "repo", "", "git repo url")
	flag.StringVar(&gitSync.GitBranch, "branch", "master", "git branch")
	flag.StringVar(&gitSync.GitRev, "rev", "HEAD", "git rev")
	flag.IntVar(&gitSync.GitDepthSync, "depth", 0, "shallow clone with a history truncated to the specified number of commits")
	flag.StringVar(&gitSync.GitUser, "username", "", "username")
	flag.StringVar(&gitSync.GitPassword, "password", "", "password")
	flag.StringVar(&gitSync.TargetDirectory, "dest", "", "destination path")
	flag.IntVar(&gitSync.Permissions, "permissions", 0, "If set it will change the permissions of the directory that contains the git repository. Example: 744")
	flag.StringVar(&gitSync.CallbackUrl, "callback-url", "", "url to call after each git pull")
	wait := flag.Int("wait", DEFAULT_WAIT, "number of seconds to wait before next sync")
	oneTime := flag.Bool("one-time", false, "exit after the initial checkout")
	flag.Parse()

	if err := gitSync.Validate(); err != nil {
		os.Exit(1)
	}

	var errorCounter int
	for {
		if err := gitSync.Run(context.Background()); err != nil {
			glog.Errorf("error syncing repo: %v", err)
			errorCounter++
		} else {
			errorCounter = 0
		}
		if errorCounter > ERROR_LIMIT {
			glog.Exitf("error limit of %d exceeded", ERROR_LIMIT)
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
