import numpy as np
import numpy.ma as ma
import matplotlib.pyplot as plt
import subprocess
import json

SCENARIOS = ["go test", "go test compile", "ginkgo ginkgo_cli_base", "ginkgo ginkgo_cli_run", "ginkgo ginkgo_cli_no_tmp"]
COLORS = {"go test": "#aaaaff", "go test compile": "#0000ff", "ginkgo ginkgo_cli_base": "#ff0000", "ginkgo ginkgo_cli_run": "#ffaaaa", "ginkgo ginkgo_cli_no_tmp": "#00aa00"}


plt.ioff()

def load_scenarios():
    scenarios = {}
    for scenario in SCENARIOS:
        with open('test_runs/'+scenario+'.json') as f:
            data = json.load(f)
            data["test_runs"]
            scenarios[scenario] = {
                "compile_time": np.array(list(map(lambda x: x["compile_time"], data["test_runs"]))),
                "run_time": np.array(list(map(lambda x: x["run_time"], data["test_runs"]))),
                "total_time": np.array(list(map(lambda x: x["total_time"], data["test_runs"])))
            }

    return scenarios

def plot(d):
    plt.style.use('seaborn')

    columns = 3
    rows = 1

    fig = plt.figure(constrained_layout=True, figsize = [columns*8,rows*4], dpi=200)
    gs = fig.add_gridspec(rows, columns)

    bins = np.linspace(0, 2, 50)

    ax = fig.add_subplot(gs[0, 0])
    for scenario in SCENARIOS:
        ax.hist(d[scenario]["compile_time"],bins=bins,linewidth=2,color=COLORS[scenario],histtype='step', label=scenario)
    ax.set_xlim(0, 2)
    ax.set_xlabel("Compile Time (s)")
    ax.legend()

    ax = fig.add_subplot(gs[0, 1])
    for scenario in SCENARIOS:
        ax.hist(d[scenario]["run_time"],bins=bins,linewidth=2,color=COLORS[scenario],histtype='step')
    ax.set_xlim(0, 2)
    ax.set_xlabel("Run Time (s)")

    ax = fig.add_subplot(gs[0, 2])
    for scenario in SCENARIOS:
        ax.hist(d[scenario]["total_time"],bins=bins,linewidth=2,color=COLORS[scenario],histtype='step')
    ax.set_xlim(0, 2)
    ax.set_xlabel("Total Time (s)")


    fig.savefig('timings.png')


def main():
    d = load_scenarios()
    plot(d)

if __name__=="__main__":
    main()
