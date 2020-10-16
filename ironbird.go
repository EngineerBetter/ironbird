package ironbird

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type TaskTestSuite struct {
	SpecDir   string
	Config    string `yaml:"config"`
	MainInput struct {
		Name         string `yaml:"name"`
		RelativeRoot string `yaml:"path_relative_to_spec"`
	} `yaml:"enclosed_in_input"`
	Cases []struct {
		When   string `yaml:"when"`
		Within string `yaml:"within"`
		It     struct {
			Exits      int      `yaml:"exits"`
			Says       []string `yaml:"says"`
			HasOutputs []struct {
				Name     string `yaml:"name"`
				ForWhich []struct {
					Bash  string   `yaml:"bash"`
					Exits int      `yaml:"exits"`
					Says  []string `yaml:"says"`
				} `yaml:"for_which"`
			} `yaml:"has_outputs,omitempty"`
			HasInputs []struct {
				Name  string `yaml:"name"`
				From  string `yaml:"from"`
				Setup string `yaml:"setup"`
			} `yaml:"has_inputs,omitempty"`
		} `yaml:"it,omitempty"`
		Params map[string]string `yaml:"params,omitempty"`
	} `yaml:"cases"`
}

func FlyExecute(target, specDir, configPath string, params map[string]string, inputDirs, outputDirs map[string]string, timeout time.Duration) *gexec.Session {
	gomega.Expect(specDir).To(gomega.BeADirectory())

	flyArgs := []string{"-t", target, "execute", "-c", configPath, "--include-ignored"}

	for name, dir := range inputDirs {
		flyArgs = append(flyArgs, "--input="+name+"="+dir)
	}

	for name, dir := range outputDirs {
		flyArgs = append(flyArgs, "--output="+name+"="+dir)
	}

	cmd := exec.Command("fly", flyArgs...)
	cmd.Dir = specDir
	cmd.Env = os.Environ()
	for key, value := range params {
		setEnv(key, value, cmd)
	}

	session, err := gexec.Start(cmd, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(session, timeout).Should(gexec.Exit())
	return session
}

func OutErrMessage(session *gexec.Session) string {
	return fmt.Sprintf("---\nSTDOUT:\n%v\nSTDERR:\n%v\n---", string(session.Out.Contents()), string(session.Err.Contents()))
}

func setEnv(key, value string, cmd *exec.Cmd) {
	cmd.Env = append(cmd.Env, key+"="+value)
}

func Bash(command, dir string) *gexec.Session {
	cmd := exec.Command("bash", "-x", "-e", "-u", "-c", command)
	cmd.Dir = dir
	session, err := gexec.Start(cmd, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("command: %s\ndir: %s\n", command, dir))
	gomega.Eventually(session, 20*time.Second).Should(gexec.Exit())
	return session
}

func MustBash(command, dir string) *gexec.Session {
	session := Bash(command, dir)
	absWorkingDir, err := filepath.Abs(session.Command.Dir)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(session.ExitCode()).To(gomega.BeZero(), "Bash command:\n%v\nWorking dir:\n%s\nSTDOUT:\n%v\nSTDERR:\n%v", command, absWorkingDir, string(session.Out.Contents()), string(session.Err.Contents()))
	return session
}
