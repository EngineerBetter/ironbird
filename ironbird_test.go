package ironbird_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	. "github.com/EngineerBetter/ironbird"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("", func() {
	defer GinkgoRecover()

	for _, specX := range specs {
		// Re-assign as otherwise the single pointer is shared by the anonymous functions, and by the time that they
		// execute, range will have made the single pointer point to something else. Possibly.
		spec := specX

		Describe(spec.Config, func() {
			for _, specCaseX := range spec.Cases {
				specCase := specCaseX

				Describe("when "+specCase.When, func() {
					inputDirs := make(map[string]string)

					if spec.MainInput.Name != "" {
						sourceDir := spec.SpecDir
						if spec.MainInput.RelativeRoot != "" {
							sourceDir = spec.MainInput.RelativeRoot
						}

						tmpInputPath := createTmpInputDir(spec.MainInput.Name, sourceDir, spec.SpecDir)
						inputDirs[spec.MainInput.Name] = tmpInputPath
					}

					for _, input := range specCase.It.HasInputs {
						tmpInputPath := createTmpInputDir(input.Name, input.From, spec.SpecDir)
						inputDirs[input.Name] = tmpInputPath

						if input.Setup != "" {
							MustBash(input.Setup, tmpInputPath)
						}
					}

					outputDirs := make(map[string]string)
					for _, outputExpectation := range specCase.It.HasOutputs {
						outputPath, err := ioutil.TempDir("", outputExpectation.Name)
						expectErrToNotHaveOccurred(err)
						outputDirs[outputExpectation.Name] = outputPath
					}

					var session *Session
					It(fmt.Sprintf("exits %d", specCase.It.Exits), func() {
						within := specCase.Within
						if within == "" {
							within = defaultTimeout
						}
						timeout, err := time.ParseDuration(within)
						Expect(err).ToNot(HaveOccurred())
						timeout = timeout * time.Duration(timeoutFactorArg)
						session = FlyExecute(targetArg, spec.SpecDir, spec.Config, specCase.Params, inputDirs, outputDirs, timeout)

						interceptMessage := ""
						if specCase.It.Exits != 0 && session.Out != nil {
							pattern := regexp.MustCompile(`executing build (\d*) at http`)
							buildNumber := pattern.Find(session.Out.Contents())
							interceptMessage = fmt.Sprintf("\nTask failed unexpectedly, debug with:\nfly -t %s intercept -b %s", targetArg, string(buildNumber))
						}

						Expect(session).To(Exit(specCase.It.Exits), OutErrMessage(session)+interceptMessage)
						Expect(session).To(Say("executing build"))
						Expect(session).To(Say("initializing"))
					})

					for _, sayExpectationX := range specCase.It.Says {
						sayExpectation := sayExpectationX

						It("says "+sayExpectation, func() {
							Expect(session).To(Say(sayExpectation))
						})
					}

					for _, outputExpectationX := range specCase.It.HasOutputs {
						outputExpectation := outputExpectationX

						Describe(fmt.Sprintf(", it has an output '%s'", outputExpectation.Name), func() {
							for _, forWhichX := range outputExpectation.ForWhich {
								forWhich := forWhichX

								Describe(fmt.Sprintf("for which '%s'", forWhich.Bash), func() {
									var assertionSession *Session
									It(fmt.Sprintf("exits %d", forWhich.Exits), func() {
										// THE REDIRECT IS ABSOLUTE CHEDDAR
										assertionSession = Bash(forWhich.Bash+" 2>&1", outputDirs[outputExpectation.Name])
										Expect(assertionSession).To(Exit(forWhich.Exits), OutErrMessage(assertionSession))
									})

									for _, sayExpectationX := range forWhich.Says {
										sayExpectation := sayExpectationX
										It("says "+sayExpectation, func() {
											Expect(assertionSession).To(Say(sayExpectation))
										})
									}
								})
							}
						})
					}
				})
			}
		})
	}
})

func createTmpInputDir(name, copyFrom, specDir string) string {
	tmpInputPath, err := ioutil.TempDir("", name)
	expectErrToNotHaveOccurred(err)

	// Allow all users access, otherwise Concourse can't read dir when executing
	MustBash("chmod 0777 "+tmpInputPath, "")

	if copyFrom != "" {
		MustBash("cp -r "+copyFrom+"/. "+tmpInputPath, specDir)
	}

	return tmpInputPath
}

func expectErrToNotHaveOccurred(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
