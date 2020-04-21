package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/montanaflynn/stats"
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
	GinkgoTestScenarioRunner{"ginkgo_cli_no_tmp", "Ginkgo CLI that compiles tests into the current package, not a tempdir"},
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
	compileTime := stats.Float64Data{}
	runTime := stats.Float64Data{}
	totalTime := stats.Float64Data{}
	for _, testRun := range scenario.TestRuns {
		compileTime = append(compileTime, testRun.CompileTime)
		runTime = append(runTime, testRun.RunTime)
		totalTime = append(totalTime, testRun.TotalTime)
		// fmt.Printf("  # %3d | C: %7.4fs R: %7.4fs T: %7.4fs\n", i+1, testRun.CompileTime, testRun.RunTime, testRun.TotalTime)
	}

	compileTimeMin, _ := compileTime.Min()
	compileTimeMean, _ := compileTime.Mean()
	compileTimeMedian, _ := compileTime.Median()
	compileTimeMax, _ := compileTime.Max()
	compileTimeStandardDeviation, _ := compileTime.StandardDeviation()

	runTimeMin, _ := runTime.Min()
	runTimeMean, _ := runTime.Mean()
	runTimeMedian, _ := runTime.Median()
	runTimeMax, _ := runTime.Max()
	runTimeStandardDeviation, _ := runTime.StandardDeviation()

	totalTimeMin, _ := totalTime.Min()
	totalTimeMean, _ := totalTime.Mean()
	totalTimeMedian, _ := totalTime.Median()
	totalTimeMax, _ := totalTime.Max()
	totalTimeStandardDeviation, _ := totalTime.StandardDeviation()

	fmt.Printf("    Compile Time: %6.4fs < <%6.4fs> [%6.4fs] ± %6.4fs < %6.4f\n", compileTimeMin, compileTimeMean, compileTimeMedian, compileTimeStandardDeviation, compileTimeMax)
	fmt.Printf("    Run Time: %6.4fs < <%6.4fs> [%6.4fs] ± %6.4fs < %6.4f\n", runTimeMin, runTimeMean, runTimeMedian, runTimeStandardDeviation, runTimeMax)
	fmt.Printf("    Total Time: %6.4fs < <%6.4fs> [%6.4fs] ± %6.4fs < %6.4f\n", totalTimeMin, totalTimeMean, totalTimeMedian, totalTimeStandardDeviation, totalTimeMax)
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
