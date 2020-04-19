package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type GoTestScenarioRunner struct{}

func (g GoTestScenarioRunner) Prep()    {}
func (g GoTestScenarioRunner) Cleanup() {}

func (g GoTestScenarioRunner) Name() string {
	return "go test"
}

func (g GoTestScenarioRunner) Description() string {
	return "Vanilla go test"
}

func (g GoTestScenarioRunner) RunTest(testPath string) TestRun {
	cmd := exec.Command("go", "test")
	cmd.Dir = testPath
	t := time.Now()
	cmd.Run()
	dt := time.Since(t).Seconds()

	return TestRun{
		CompileTime: -1,
		RunTime:     dt,
		TotalTime:   dt,
	}
}

type GoTestCompileScenarioRunner struct{}

func (g GoTestCompileScenarioRunner) Prep()    {}
func (g GoTestCompileScenarioRunner) Cleanup() {}

func (g GoTestCompileScenarioRunner) Name() string {
	return "go test compile"
}

func (g GoTestCompileScenarioRunner) Description() string {
	return "go test -c, followed by invoking the resulting binary"
}

func (g GoTestCompileScenarioRunner) RunTest(testPath string) TestRun {
	tTotal := time.Now()

	cmd := exec.Command("go", "test", "-c", "-o", "./test")
	cmd.Dir = testPath
	tCompile := time.Now()
	cmd.Run()
	dtCompile := time.Since(tCompile).Seconds()

	cmd = exec.Command("test")
	cmd.Dir = testPath
	tRun := time.Now()
	cmd.Run()
	dtRun := time.Since(tRun).Seconds()

	os.Remove(filepath.Join(testPath, "test"))

	dtTotal := time.Since(tTotal).Seconds()

	return TestRun{
		CompileTime: dtCompile,
		RunTime:     dtRun,
		TotalTime:   dtTotal,
	}
}
