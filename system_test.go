// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var pathToServerBinary string

var serverSession *gexec.Session

var _ = BeforeSuite(func() {
	var err error
	pathToServerBinary, err = gexec.Build("github.com/bborbe/git-sync")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

type args map[string]string

func (a args) list() []string {
	var result []string
	for k, v := range a {
		if len(v) == 0 {
			result = append(result, fmt.Sprintf("-%s", k))
		} else {
			result = append(result, fmt.Sprintf("-%s=%s", k, v))
		}
	}
	return result
}

var _ = Describe("git-sync", func() {
	var err error
	It("returns with exitcode != 0 if no parameters have been given", func() {
		serverSession, err = gexec.Start(exec.Command(pathToServerBinary), GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		serverSession.Wait(time.Second)
		Expect(serverSession.ExitCode()).NotTo(Equal(0))
	})
	Context("when validating parameters", func() {
		var validargs args
		var targetDirectory string
		BeforeEach(func() {
			targetDirectory, err = ioutil.TempDir("", "git-sync")
			Expect(err).To(BeNil())
			validargs = map[string]string{
				"logtostderr": "",
				"v":           "0",
				"repo":        "http://github.com/bborbe/dotfiles.git",
				"one-time":    "",
				"dest":        targetDirectory,
			}
		})
		AfterEach(func() {
			_ = os.RemoveAll(targetDirectory)
		})
		It("returns with exitcode == 0", func() {
			serverSession, err = gexec.Start(exec.Command(pathToServerBinary, validargs.list()...), GinkgoWriter, GinkgoWriter)
			Expect(err).To(BeNil())
			serverSession.Wait(5 * time.Second)
			Expect(serverSession.ExitCode()).To(Equal(0))
			_, err = os.Stat(targetDirectory)
			Expect(os.IsNotExist(err)).To(BeFalse())
		})
		Context("and url parameter", func() {
			var server *ghttp.Server
			BeforeEach(func() {
				server = ghttp.NewServer()
				server.RouteToHandler(http.MethodGet, "/", ghttp.RespondWith(http.StatusOK, "OK"))
			})
			AfterEach(func() {
				serverSession.Interrupt()
				Eventually(serverSession).Should(gexec.Exit())
				server.Close()
			})
			It("calls the url", func() {
				validargs["callback-url"] = server.URL()
				serverSession, err = gexec.Start(exec.Command(pathToServerBinary, validargs.list()...), GinkgoWriter, GinkgoWriter)
				Expect(err).To(BeNil())
				serverSession.Wait(5 * time.Second)
				Expect(serverSession.ExitCode()).To(Equal(0))
				Expect(len(server.ReceivedRequests())).To(Equal(1))
			})
			It("calls the url with env", func() {
				delete(validargs, "callback-url")
				command := exec.Command(pathToServerBinary, validargs.list()...)
				command.Env = []string{fmt.Sprintf("CALLBACK_URL=%s", server.URL()), fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
				serverSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).To(BeNil())
				serverSession.Wait(5 * time.Second)
				Expect(serverSession.ExitCode()).To(Equal(0))
				Expect(len(server.ReceivedRequests())).To(Equal(1))
			})
		})

	})
})
