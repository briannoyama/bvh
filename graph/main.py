import argparse, os
from os.path import join, isfile, basename
from collections import defaultdict
from plot import percentile2d, percentile3d, comparison

def main():
    parser = argparse.ArgumentParser(description="")
    parser.add_argument("path", metavar='P', type=str, help="The path to the folder holding the statistics.")
    
    args = parser.parse_args()
    
    # Look for files in path
    files = [join(args.path, f) for f in os.listdir(args.path)]
    files = [f for f in files if isfile(f)]

    # Gather stats into a dictionary
    stats = defaultdict(lambda: [])
    for f in files:
        if basename(f).startswith("test"):
            for text in open(f):
                line = [s.strip() for s in text.split(",")]
                if line[0] == "que":
                    timePerSizeReturns(stats, line, 1)
                    timePerSizeReturns(stats, line, 2)
                else:
                    timePerSize(stats, line, 1)
                    timePerSize(stats, line, 2)

        elif basename(f).startswith("compare"):
            for text in open(f):
                line = [s.strip() for s in text.split(",")]
                saPerSize(stats, line, 1)
                saPerSize(stats, line, 2)
    
    # 2D graph of percentiles
    percentile2d(stats["add_size"], "Add Runtime per Size", "size", 100)
    percentile2d(stats["sub_size"], "Sub Runtime per Size", "size", 100)
    percentile2d(stats["add_depth"], "Add Runtime per Depth", "depth")
    percentile2d(stats["sub_depth"], "Sub Runtime per Depth", "depth")

    # 3D graph of percentiles data, output, title, x_label,
    percentile3d(stats["que_size"], "Query Runtime per Size", "size")
    percentile3d(stats["que_depth"], "Query Runtime per Depth", "depth")

    # Comparison graphs
    comparison(stats["offline_sa"], stats["online_sa"], "surface area")
    comparison(stats["offline_depth"], stats["online_depth"], "depth")
    

def timePerSize(stats, line, size_index=1):
    data = stats[line[0] + ("_size" if size_index == 1 else "_depth")]
    size = int(line[size_index])
    _setLength(data, size)
    data[size].append(int(line[3]))


def timePerSizeReturns(stats, line, size_index=1):    
    data = stats[line[0] + ("_size" if size_index == 1 else "_depth")]
    size = int(line[size_index])
    returns = int(line[4])
    _setLength(data, returns)
    _setLength(data[returns], size)
    data[returns][size].append(int(line[4]))


def saPerSize(stats, line, sa_index=1):
    data = stats["offline" + ("_sa" if sa_index == 2 else "_depth")]
    size = int(line[0])
    _setLength(data, size)    
    data[size] = int(line[sa_index])

    data = stats["online" + ("_sa" if sa_index == 2 else "_depth")]
    _setLength(data, size)
    data[size] = int(line[sa_index+2])


def _setLength(list_, size):
    if len(list_) <= size :
        list_.extend([[]]* (size + 1 - len(list_)))

if __name__ == "__main__":
    main()
