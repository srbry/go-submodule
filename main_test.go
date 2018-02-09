package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var pathToCLI string

var _ = BeforeSuite(func() {
	var err error
	pathToCLI, err = gexec.Build("github.com/srbry/go-submodule")
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("go_submodule", func() {
	var (
		err     error
		session *gexec.Session
	)

	BeforeEach(func() {
		command := exec.Command(pathToCLI, "--source=fixtures/example.toml")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
	})

	AfterEach(func() {
		session.Kill()
	})

	It("shows some toml", func() {
		Ω(err).ShouldNot(HaveOccurred())
		Eventually(session.Out).Should(gbytes.Say(`mkdir -p github.com/ryanuber
pushd github.com/ryanuber
  git submodule add https://github.com/ryanuber/go-glob.git
  cd go-glob
  git checkout master
popd`))
		Eventually(string(session.Out.Contents())).Should(ContainSubstring(`mkdir -p github.com/docker
pushd github.com/docker
  git submodule add https://github.com/containous/leadership.git
  cd leadership
  git checkout master
popd`))
		Eventually(string(session.Out.Contents())).Should(ContainSubstring(`mkdir -p github.com/cenk
pushd github.com/cenk
  git submodule add https://github.com/cenk/backoff.git
  cd backoff
  git checkout master
popd`))
		Eventually(string(session.Out.Contents())).Should(ContainSubstring(`mkdir -p github.com/abbot
pushd github.com/abbot
  git submodule add https://github.com/containous/go-http-auth.git
  cd go-http-auth
  git checkout containous-fork
popd`))
		Eventually(string(session.Out.Contents())).Should(ContainSubstring(`mkdir -p github.com/mailgun
pushd github.com/mailgun
  git submodule add https://github.com/mailgun/timetools.git
  cd timetools
  git checkout 7e6055773c5137efbeb3bd2410d705fe10ab6bfd
popd`))
		Eventually(string(session.Out.Contents())).Should(ContainSubstring(`mkdir -p github.com/eapache
pushd github.com/eapache
  git submodule add https://github.com/eapache/channels.git
  cd channels
  git checkout v1.1.0
popd`))
		Eventually(string(session.Out.Contents())).Should(ContainSubstring(`mkdir -p github.com/mesosphere
pushd github.com/mesosphere
  git submodule add https://github.com/containous/mesos-dns.git.git
popd`))
		Eventually(session.Err).Should(gbytes.Say(""))
	})
})
