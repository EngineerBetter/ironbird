package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/EngineerBetter/ironbird"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("running ironbird", func() {
	var tmpDir, executablePath string

	invoke := func(args ...string) *Session {
		args = append(args, "--target", targetArg)
		cmd := exec.Command(executablePath, args...)
		cmd.Dir = tmpDir
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		return session
	}

	BeforeSuite(func() {
		tmpDir, err := ioutil.TempDir("", "ironbird-build")
		MustBash("cp -r . "+tmpDir, "../")

		cmd := exec.Command("ginkgo", "build", ".")
		cmd.Dir = tmpDir
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 20*time.Second).Should(Exit(0))
		executablePath = filepath.Join(tmpDir, "ironbird.test")
		MustBash("mv *.test ironbird.test", tmpDir)
		Expect(filepath.Join(tmpDir, "ironbird.test")).To(BeAnExistingFile())

		MustBash("cp -r integration/fixtures "+tmpDir, "../")
		MustBash("chmod 0777 "+tmpDir, "")
	})

	AfterSuite(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("running all the specs", func() {
		It("works", func() {
			session := invoke("--specs", "fixtures/file_write_spec.yml,fixtures/input_spec.yml", "--timeout-factor", "100")
			Eventually(session, 2*time.Minute).Should(Exit(0))
		})
	})

	Describe("testing task exit codes", func() {
		When("the spec expects failure", func() {
			Context("and the task exits 1", func() {
				It("passes", func() {
					session := invoke("--specs", "fixtures/passing_exit1_spec.yml")
					Eventually(session, 2*time.Minute).Should(Exit(0))
				})
			})

			Context("and the task exits 0", func() {
				It("fails", func() {
					session := invoke("--specs", "fixtures/failing_exit1_spec.yml")
					Eventually(session, 2*time.Minute).Should(Exit(1))
				})

				It("logs the fly intercept command to debug with", func() {
					session := invoke("--specs", "fixtures/failing_exit1_spec.yml")
					Eventually(session, 2*time.Minute).Should(Exit(1))
					Expect(session).To(Say(`fly -t ` + targetArg + ` intercept -b \d*`))
				})
			})
		})
	})

	Describe("timing out", func() {
		When("a task times out", func() {
			It("fails", func() {
				session := invoke("--specs", "fixtures/failing_sleep_spec.yml")
				Eventually(session, 1*time.Minute).Should(Exit(1))
			})

			Context("but timeout-factor allows more time", func() {
				It("passes", func() {
					session := invoke("--specs", "fixtures/failing_sleep_spec.yml", "--timeout-factor", "2")
					Eventually(session, 1*time.Minute).Should(Exit(0))
				})
			})
		})
	})

	Describe("testing task output", func() {
		It("works", func() {
			session := invoke("--specs", "fixtures/echo_spec.yml")
			Eventually(session, 1*time.Minute).Should(Exit(0))
		})

		It("matches against STDOUT and STDERR", func() {
			session := invoke("--specs", "fixtures/echo_redirect_spec.yml")
			Eventually(session, 1*time.Minute).Should(Exit(0))
		})
	})

	Describe("testing a task in a repo subdirectory", func() {
		It("works", func() {
			session := invoke("--specs", "fixtures/pretend-repo/ci/tasks/nested-task/nested-task_spec.yml")
			Eventually(session, 1*time.Minute).Should(Exit(0))
		})
	})
})
