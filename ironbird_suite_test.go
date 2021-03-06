package ironbird_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/EngineerBetter/ironbird"
	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

const defaultTimeout = time.Second * time.Duration(20)

var specs []*ironbird.TaskTestSuite

var specsArg, targetArg string
var timeoutFactorArg int

func init() {
	flag.StringVar(&specsArg, "specs", "", "Comma-separated list of spec files to execute")
	flag.StringVar(&targetArg, "target", "", "fly target")
	flag.IntVar(&timeoutFactorArg, "timeout-factor", 1, "multiplier for timeouts")
}

func TestIronbird(t *testing.T) {
	if specsArg == "" {
		log.Fatal("--specs must be provided")
	}

	if targetArg == "" {
		log.Fatal("--target must be provided")
	}

	if timeoutFactorArg < 1 {
		log.Fatal("--timeout-factor must be >= 1")
	}

	modifier := .9
	slowSpecThreshold := time.Duration(defaultTimeout.Seconds()*float64(timeoutFactorArg)*modifier) * time.Second
	config.DefaultReporterConfig.SlowSpecThreshold = slowSpecThreshold.Seconds()

	specFiles := strings.Split(specsArg, ",")
	for _, specFile := range specFiles {
		loadSpec(specFile)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Ironbird: "+specsArg)
}

func loadSpec(filename string) {
	if filename == "" {
		log.Fatalf("Spec file list (%s) contained empty element", specsArg)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("Spec file '%s' does not exist", filename)
	}

	yamlFile, setupErr := ioutil.ReadFile(filename)
	expectErrToNotHaveOccurred(setupErr)

	var spec *ironbird.TaskTestSuite
	setupErr = yaml.Unmarshal(yamlFile, &spec)
	expectErrToNotHaveOccurred(setupErr)

	absSpecDir, err := filepath.Abs(filepath.Dir(filename))
	expectErrToNotHaveOccurred(err)
	spec.SpecDir = absSpecDir
	specs = append(specs, spec)
}

var _ = BeforeSuite(func() {
	// Do validation that the test spec is valid here, where we can use assertions
	Expect(specs).ToNot(BeEmpty())
	for _, spec := range specs {
		Expect(spec.Cases).ToNot(BeEmpty(), fmt.Sprintf("%s had no cases", spec.Config))
		for _, specCase := range spec.Cases {
			if specCase.Within != "" {
				_, err := time.ParseDuration(specCase.Within)
				Expect(err).ToNot(HaveOccurred())
			}
			for _, input := range specCase.It.HasInputs {
				if input.From != "" {
					inputPath := filepath.Join(spec.SpecDir, input.From)
					Expect(inputPath).To(BeADirectory())
				}
			}
		}
	}
})
