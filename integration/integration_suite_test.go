package integration_test

import (
	"flag"
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var targetArg string

func init() {
	flag.StringVar(&targetArg, "target", "", "fly target")
}

func TestIntegration(t *testing.T) {
	if targetArg == "" {
		log.Fatal("--target argument must be provided")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
