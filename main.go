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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TestRuns    TestRuns `json:"test_runs"`
}

type ScenarioRunner interface {
	Name() string
	Description() string
	RunTest(testPath string) TestRun
	Prep()
	Cleanup()
}

const NUM_RUNS = 30

var ScenarioRunners = []ScenarioRunner{
	GoTestScenarioRunner{},
	GoTestCompileScenarioRunner{},
	GinkgoTestScenarioRunner{"ginkgo_cli_base", "Original Ginkgo CLI"},
	GinkgoTestScenarioRunner{"ginkgo_cli_run", "Ginkgo CLI that calls cmd.Run() instead of cmd.Start() then cmd.Wait()"},
}

func main() {
	force := false
	fmt.Println("Running Scenarios")
	if len(os.Args) > 1 && os.Args[1] == "force" {
		fmt.Println("Forcing Rerun")
		force = true
	}
	runScenarios(force)
	analyzeScenarios()
}

func analyzeScenarios() {
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
	fmt.Println(scenario.Description)
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
		description := scenarioRunner.Description()
		jsonPath := filepath.Join(testRunPath, name+".json")

		if !force {
			_, err := os.Stat(jsonPath)
			if err == nil {
				fmt.Println("Skipping", name)
				continue
			}
		}

		scenario := Scenario{
			Name:        name,
			Description: description,
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
