package integration_test

import (
	. "github.com/EngineerBetter/ironbird"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var _ = Describe("Integration", func() {
	var tmpDir, executablePath string

	BeforeSuite(func(){
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

		MustBashIn("../", "cp -r integration/fixtures "+tmpDir)
		MustBash("chmod 0777 "+tmpDir)
	})

	AfterSuite(func(){
		err := os.Remove(executablePath)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("running all the specs", func(){
		It("works", func(){
			cmd := exec.Command(executablePath, "--specs", "fixtures/echo_spec.yml", "--target", "eb")
			cmd.Dir = tmpDir
			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 1*time.Minute).Should(Exit(0))
		})
	})
})
