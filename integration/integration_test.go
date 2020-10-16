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
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("running ironbird", func() {
	var tmpDir, executablePath string

	BeforeSuite(func() {
		ginkgoArgs := []string{"build", "../"}
		cmd := exec.Command("ginkgo", ginkgoArgs...)
		session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 20*time.Second).Should(Exit(0))
		Expect("../ironbird.test").To(BeAnExistingFile())

		tmpDir, err = ioutil.TempDir("", "ironbird")
		Expect(err).ToNot(HaveOccurred())
		executablePath = filepath.Join(tmpDir, "ironbird")
		err = os.Rename("../ironbird.test", executablePath)
		Expect(err).ToNot(HaveOccurred())

		MustBash("cp -r integration/fixtures "+tmpDir, "../")
		MustBash("chmod 0777 "+tmpDir, "")
	})

	AfterSuite(func() {
		err := os.RemoveAll(tmpDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("running all the specs", func() {
		It("works", func() {
			cmd := exec.Command(executablePath, "--specs", "fixtures/file_write_spec.yml,fixtures/input_spec.yml", "--target", "eb", "--timeout-factor", "100")
			cmd.Dir = tmpDir
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 2*time.Minute).Should(Exit(0))
		})
	})

	Describe("testing task exit codes", func() {
		When("the spec expects failure", func() {
			Context("and the task exits 1", func() {
				It("passes", func() {
					cmd := exec.Command(executablePath, "--specs", "fixtures/passing_exit1_spec.yml", "--target", "eb")
					cmd.Dir = tmpDir
					session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())
					Eventually(session, 2*time.Minute).Should(Exit(0))
				})
			})

			Context("and the task exits 0", func() {
				It("fails", func() {
					cmd := exec.Command(executablePath, "--specs", "fixtures/failing_exit1_spec.yml", "--target", "eb")
					cmd.Dir = tmpDir
					session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())
					Eventually(session, 2*time.Minute).Should(Exit(1))
				})
			})
		})
	})

	Describe("timing out", func() {
		When("a task times out", func() {
			It("fails", func() {
				cmd := exec.Command(executablePath, "--specs", "fixtures/failing_sleep_spec.yml", "--target", "eb")
				cmd.Dir = tmpDir
				session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(Exit(1))
			})

			Context("but timeout-factor allows more time", func() {
				It("passes", func() {
					cmd := exec.Command(executablePath, "--specs", "fixtures/failing_sleep_spec.yml", "--target", "eb", "--timeout-factor", "2")
					cmd.Dir = tmpDir
					session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())
					Eventually(session, 1*time.Minute).Should(Exit(0))
				})
			})
		})
	})

	Describe("testing task output", func() {
		It("works", func() {
			cmd := exec.Command(executablePath, "--specs", "fixtures/echo_spec.yml", "--target", "eb")
			cmd.Dir = tmpDir
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 1*time.Minute).Should(Exit(0))
		})
	})

	Describe("testing a task in a repo subdirectory", func() {
		It("works", func() {
			cmd := exec.Command(executablePath, "--specs", "fixtures/pretend-repo/ci/tasks/nested-task/nested-task_spec.yml", "--target", "eb")
			cmd.Dir = tmpDir
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 1*time.Minute).Should(Exit(0))
		})
	})
})
