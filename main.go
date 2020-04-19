package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type TestRun struct {
	CompileTime float64 `json:"compile_time"`
	RunTime     float64 `json:"run_time"`
	TotalTime   float64 `json:"total_time"`
}

type TestRuns []TestRun

type Scenario struct {
	Name     string   `json:"name"`
	TestRuns TestRuns `json:"test_runs"`
}

type ScenarioRunner interface {
	Name() string
	RunTest(testPath string) TestRun
	Prep()
	Cleanup()
}

const NUM_RUNS = 30

var ScenarioRunners = []ScenarioRunner{
	GoTestScenarioRunner{},
	GoTestCompileScenarioRunner{},
	GinkgoTestScenarioRunner{CLI: "ginkgo_cli_base"},
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "run" {
		force := false
		fmt.Println("Running Scenarios")
		if len(os.Args) > 2 && os.Args[2] == "force" {
			fmt.Println("Forcing Rerun")
			force = true
		}

		runScenarios(force)
	}

	scenarios := []Scenario{}

	f, _ := os.Open("scenarios.json")
	json.NewDecoder(f).Decode(&scenarios)

	analyzeScenarios(scenarios)
}

func analyzeScenarios(scenarios []Scenario) {
	wd, _ := os.Getwd()
	testRunPath := filepath.Join(wd, "test_runs")

	for _, scenarioRunner := range ScenarioRunners {
		name := scenarioRunner.Name()
		jsonPath := filepath.Join(testRunPath, name+".json")

		f, err := os.Open(jsonPath)
		if err != nil {
			fmt.Printf("Skipping %s - no json found", name)
			continue
		}
		scenario := Scenario{}
		json.NewDecoder(f).Decode(&scenario)
		f.Close()
		analyzeScenario(scenario)
	}
}

func analyzeScenario(scenario Scenario) {
	fmt.Println(scenario.Name)
	for i, testRun := range scenario.TestRuns {
		fmt.Printf("  # %3d | C: %7.4fs R: %7.4fs T: %7.4fs\n", i+1, testRun.CompileTime, testRun.RunTime, testRun.TotalTime)
	}
}

func runScenarios(force bool) {
	wd, _ := os.Getwd()
	testPath := filepath.Join(wd, "sample")
	testRunPath := filepath.Join(wd, "test_runs")

	for _, scenarioRunner := range ScenarioRunners {
		name := scenarioRunner.Name()
		jsonPath := filepath.Join(testRunPath, name+".json")

		if !force {
			_, err := os.Stat(jsonPath)
			if err == nil {
				fmt.Println("Skipping", name)
				continue
			}
		}

		scenario := Scenario{
			Name: name,
		}
		fmt.Println("Running", name)
		scenarioRunner.Prep()
		for i := 0; i < NUM_RUNS; i += 1 {
			testRun := scenarioRunner.RunTest(testPath)
			fmt.Printf("  # %3d | C: %7.4fs R: %7.4fs T: %7.4fs\n", i+1, testRun.CompileTime, testRun.RunTime, testRun.TotalTime)
			scenario.TestRuns = append(scenario.TestRuns, testRun)
		}
		scenarioRunner.Cleanup()

		f, _ := os.Create(jsonPath)
		json.NewEncoder(f).Encode(scenario)
		f.Close()
	}
}
