from mpl_toolkits.mplot3d import Axes3D  

from matplotlib import cbook
from matplotlib import cm
from matplotlib.colors import LightSource
from statistics import mean

import matplotlib.pyplot as plt
import matplotlib
import numpy as np

PERCENTILES = [0.5, 0.9, 0.99]
matplotlib.use('SVG') 

def percentile2d(data, title, x_label, flatten = 0):
    x, p_y = _percentiles(data)
    fig, ax = plt.subplots()
    
    # Graph each percentile
    _graph(x, p_y, ax, title, x_label, flatten)
    
    
    fig.savefig(''.join(x for x in title.title().split()) + ".svg")

def _percentiles(data):
    x = []
    p_y = {str(p):[] for p in PERCENTILES}
    for i, distribution in enumerate(data):
        distribution.sort()
        if len(distribution) > 0:
            x.append(i)

            for p in p_y:
                p_y[p].append(distribution[int(len(distribution)*float(p))])
    return x, p_y

def _graph(x, p_y, ax, title, x_label, flatten):
    x_array = np.array(x)
    for p, percents in p_y.items():
        ax.semilogy(x_array, np.array([mean(percents[max(0, p-flatten):min(len(percents), p+1)]) 
            for p in range(len(percents))]), label=p)
        ax.set_xlabel(x_label)
        ax.set_ylabel("speed (ns)")
        ax.set_title(title)
        ax.legend()

def percentile3d(data, title, x_label, max_returns = 5, flatten=0):
    max_y = 0
    max_x = 0
    x = []
    p_y = []

    #print(len(data))
    #exit()
    
    graphs = min(len(data), max_returns)

    # Get percent maps.
    for dist in data[: graphs]:
        #print("Sorting")
        max_y = max(max(dist[-1]), max_y)
        max_x = max(len(dist), max_x)
        x_dist, p_y_dist = _percentiles(dist)
        x.append(x_dist)
        p_y.append(p_y_dist)

    fig, ax = plt.subplots(ncols=graphs, sharey=True, squeeze=True)
    
    # Graph each percentile
    for index, ax_ in enumerate(ax):
        print("Plotting")
        _graph(x[index], p_y[index], ax_, title, x_label, flatten)
    
    fig.savefig("query_" + ''.join(x for x in title.title().split()) + p + ".svg")


def comparison(data1, data2, y_label):
    x_array = np.array(range(len(data1)))
    y_offline = np.array(data1)
    y_online = np.array(data2)
    fig, ax = plt.subplots()

    for p, percents in p_y.items():
        ax.plot(x_array, y_offline, label="offline")
        ax.plot(x_array, y_online, label="online")
        ax.set_xlabel("size")
        ax.set_ylabel(y_label)
        ax.set_title(title)
        ax.legend()

    fig.savefig(''.join(x for x in y_label.title().split()) + ".svg")

