package main

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

var COMPILED_RE = regexp.MustCompile(`Ω COMPILED IN ([\d\.]+)`)
var RAN_RE = regexp.MustCompile(`Ω RAN IN ([\d\.]+)`)

type GinkgoTestScenarioRunner struct {
	CLI  string
	Desc string
}

func (g GinkgoTestScenarioRunner) Prep() {
	cmd := exec.Command("go", "build", "-o", "./ginkgo-cli", "./ginkgo/"+g.CLI)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func (g GinkgoTestScenarioRunner) Cleanup() {
	os.Remove("./ginkgo-cli")
}

func (g GinkgoTestScenarioRunner) Name() string {
	return "ginkgo " + g.CLI
}

func (g GinkgoTestScenarioRunner) Description() string {
	return g.Desc
}

func (g GinkgoTestScenarioRunner) RunTest(testPath string) TestRun {
	cmd := exec.Command("./ginkgo-cli", testPath)
	t := time.Now()
	output, _ := cmd.Output()
	dt := time.Since(t).Seconds()

	match := COMPILED_RE.FindSubmatch(output)
	if match == nil || len(match) != 2 {
		panic("couldn't get compile time")
	}
	dtCompile, err := strconv.ParseFloat(string(match[1]), 64)
	if err != nil {
		panic(err)
	}

	match = RAN_RE.FindSubmatch(output)
	if match == nil || len(match) != 2 {
		panic("couldn't get compile time")
	}
	dtRan, err := strconv.ParseFloat(string(match[1]), 64)
	if err != nil {
		panic(err)
	}

	return TestRun{
		CompileTime: dtCompile,
		RunTime:     dtRan,
		TotalTime:   dt,
	}
}
