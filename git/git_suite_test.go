// Copyright 2018 The Git-Sync Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package git_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/bborbe/git-sync/git"
	"context"
	"github.com/bborbe/git-sync/mocks"
	"net/http"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

var _ = Describe("sync", func() {
	var err error
	var gitSync *git.Sync
	var transport *mocks.RoundTripper
	var targetDirectory string
	BeforeEach(func() {
		transport = &mocks.RoundTripper{}
		targetDirectory, err = ioutil.TempDir("", "git-sync")
		Expect(err).To(BeNil())
		gitSync = &git.Sync{
			Transport:       transport,
			GitBranch:       "master",
			GitRepo:         "http://github.com/bborbe/dotfiles.git",
			TargetDirectory: targetDirectory,
			GitRev:          "HEAD",
		}
	})
	AfterEach(func() {
		os.RemoveAll(targetDirectory)
	})
	Context("with callback url", func() {
		BeforeEach(func() {
			gitSync.CallbackUrl = "http://localhost:1234"
		})
		It("call the callback if defined", func() {
			transport.RoundTripReturns(&http.Response{StatusCode: 200}, nil)
			err = gitSync.Run(context.Background())
			Expect(err).To(BeNil())
			Expect(transport.RoundTripCallCount()).To(Equal(1))
		})
		It("return error if http call failed", func() {
			transport.RoundTripReturns(&http.Response{}, errors.New("banana"))
			err = gitSync.Run(context.Background())
			Expect(err).NotTo(BeNil())
		})
	})
	It("no http call is made if callback is not defined", func() {
		gitSync.CallbackUrl = ""
		err = gitSync.Run(context.Background())
		Expect(err).To(BeNil())
		Expect(transport.RoundTripCallCount()).To(Equal(0))
	})
	It("return no error on validate", func() {
		err = gitSync.Validate()
		Expect(err).To(BeNil())
	})
	It("return error if branch is missing", func() {
		gitSync.GitBranch = ""
		err = gitSync.Validate()
		Expect(err).NotTo(BeNil())
	})
	It("return error if rev is missing", func() {
		gitSync.GitRev = ""
		err = gitSync.Validate()
		Expect(err).NotTo(BeNil())
	})
	It("return error if repo is missing", func() {
		gitSync.GitRepo = ""
		err = gitSync.Validate()
		Expect(err).NotTo(BeNil())
	})
})

func TestGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Suite")
}
