# ginkgo-perf

Trying to get to the bottom of `ginkgo` vs `go test` performance.

## Usage

`go run .` will run a few different scenarios and save the test run timings in `test_runs`.  If a scenario is already saved in `test_runs` it is skipped.
`go run . force` will forcibly rerun all the scenarios regardless of what is stored in `test_runs`

Either command spits out a summary of the various scenarios.

`python visualize.py` generates a plot

## Finding (To Date)

There are four scenarios, each runs the tests stored in `sample`:

- `go test`: simply runs `go test sample`
- `go test compile`: mimics Ginkgo's behavior of compiling the test first (`go test -c -output="test"`) then running the test (`./test`).
- `ginkgo_cli_base`: runs the default Ginkgo cli, which first compiles the test and stores it in a temp directory then runs the test.
- `ginkgo_cli_no_tmp`: runs a modified Ginkgo cli, which first compiles the test and stores it in the package directory then runs the test.

### Run-time Statistics

For each scenario (where possible) we capture the Compile Time, Run Time, and Total Time.  The following are stats from a 100-iteration run on my macbook pro:

| Scenario | Metric | Min | Max | Mean | Median | Std Dev |
| --- | --- | --- | --- | --- | --- | --- |
| `go test` | Run Time | 0.9271 | 1.0830 | 0.9998s | 0.9999s | 0.0186s |
| `go test` | Total Time | 0.9271 | 1.0830 | 0.9998s | 0.9999s | 0.0186s |
| `go test compile` | Compile Time | 0.8978s | 1.1267 | 0.9642s | 0.9564s | 0.0421s |
| `go test compile` | Run Time | 0.1264s | 0.2085 | 0.1353s | 0.1333s | 0.0105s |
| `go test compile` | Total Time | 1.0339s | 1.2637 | 1.1007s | 1.0903s | 0.0438s |
| `ginkgo_cli_base` | Compile Time | 0.8865s | 1.0358 | 0.9345s | 0.9306s | 0.0283s |
| `ginkgo_cli_base` | Run Time | 0.1752s | 0.9746 | 0.4585s | 0.1866s | 0.3380s |
| `ginkgo_cli_base` | Total Time | 1.0782s | 1.9955 | 1.4081s | 1.1674s | 0.3380s |
| `ginkgo_cli_no_tmp` | Compile Time | 0.9008s | 1.1703 | 0.9794s | 0.9646s | 0.0568s |
| `ginkgo_cli_no_tmp` | Run Time | 0.1262s | 0.1754 | 0.1370s | 0.1339s | 0.0104s |
| `ginkgo_cli_no_tmp` | Total Time | 1.0468s | 1.3352 | 1.1308s | 1.1161s | 0.0636s |


And here's the plot:

![timings plot](https://github.com/onsi/ginkgo_perf/blob/master/timings.png)

### What it means:

- `go test` is 100ms faster than `go test compile`
- `ginkgo_cli_no_tmp` is equivalent to `go test compile`
- `ginkgo_cli_base` - i.e. the default ginkgo behavior that compiles the test into a tmp directory - sees notable latency in the Run Time - of about 500ms on average with a wide variance.

So.  Apparently.  Don't try to run binaries from a temp directory (at least on MacOS) as that introduces latency with a long-tail distribution (see the plot).
