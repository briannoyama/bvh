# Online BVH

### Intro

This code (hereby called onlineBVH) is a golang implementation of a binary self-balancing bounding Volume Hierarchy inspired by the tree rotations in [Fast, Effective BVH Updates for Animated Scenes](https://www.cs.utah.edu/~aek/research/tree.pdf). The hierarchies created via this algorithm have the following properties for _n_ volumes (defined by orthotopes) of any integer dimension (greater than 0):

- Average _log(n)_ addition/insertion time.
- Average _log(n)_ removal time.
- Average _mlog(n)_ query time where m is the number of volumes found.

Example Use Cases:

- Collisions between objects in a game or for ray tracing.
- Dynamically updating n-dimentional vectors (e.g. word-vectors).

### How it works

The algorithm uses integers (personal preference) to define the points of volumes. Queries are thread-safe; however, additions and removals are not. The animations below show the algorithm in action (they are pixelated, save them and look at them on your computer to get rid of the blur): 

<table>
  <tr>
    <td>
      Adding volumes
    </td>
    <td>
      Removing volumes #1
    </td>
    <td>
      Removing volumes #2
    </td>
  </tr>
  <tr>
    <td>
      <img style="image-rendering: pixelated;" alt="Animated steps of showing addition of volumes to the BVH" width="200" src="http://briannoyama.github.io/assets/images/bvh-steps/add.gif">
    </td>
    <td>
      <img style="image-rendering: pixelated;" alt="Animated steps of showing removal of volumes to the BVH" width="200" src="http://briannoyama.github.io/assets/images/bvh-steps/remove0.gif">
    </td>
    <td>
      <img style="image-rendering: pixelated;" alt="Animated steps of showing an alernative removal of volumes to the BVH" width="200" src="http://briannoyama.github.io/assets/images/bvh-steps/remove1.gif">
    </td>
  </tr>
</table>

Here's an example of how to use the code:

```golang
import "github.com/briannoyama/bvh/rect"

// Change the DIMENSIONS constant in orthotope.go for your use case.
orth := &rect.Orthotope{Point: [2]int{10, -20}, Delta: [2]int{30, 30}}
bvol := &rect.BVol{}
bvol.Add(orth)
bvol.Remove(orth)

// Use an iterator to reduce the amount of Garbage Collection
iter := bvol.Iterator()

iter.Add(orth)
for r := iter.Query(q); r != nil; r = iter.Query(q) {
    // Do something with each orthtope r
}
iter.Remove(orth)
```

To ensure _log(n)_ access along with close to ideal performance, the algorithm swaps child nodes within the BVH tree both to balance the tree and to reduce the Surface Area of the generated bounding volumes. Below one can see the output of onlineBVH vs an offline algorithm (hereby offlineBVH) that attempts to create "ideal" binary BVHs. The offline algorithm tries to create an ideal tree by sorting all of the volumes in each of their dimensions and comparing the surface areas of half the volumes at a time. Rinse and repeat recursively. This takes _O(dnlog<sup>2</sup>(n))_ for the offline method compared to the _O(nlog(n))_ time for the online method. (I'm not presenting a formal proof of big O. There may be a tighter big O bound, but that should be close enough.) In short, the offline method takes way more time to construct.

<table>
  <tr>
    <td>
      Online BVH
    </td>
    <td>
      Offline BVH
    </td>
  </tr>
  <tr>
    <td>
      <img style="image-rendering: pixelated;" alt="Output of online algorithm for generating BVH" width="200" src="http://briannoyama.github.io/assets/images/bvh-steps/online.png">
    </td>
    <td>
      <img style="image-rendering: pixelated;" alt="Output of offline algorithm for generating BVH" width="200" src="http://briannoyama.github.io/assets/images/bvh-steps/offline.png">
    </td>
  </tr>
</table>

### Performance Test

For those who plan to use onlineBVH for an application with strict runtime requirements, I conducted a small experiment on my Intel Core i5-7440HQ CPU @ 2.80GHz × 4. The test generated random cubes in a 3d space to add (100,000) remove (50,000) and query (100,000) such that the final BVH would contain 50,000 items. I ran this test 20 times and combined the data to get the below graphs:

The different colored lines represent the different percentiles. Interestingly, the first addition took way longer than any of the following additions. The rest of the additions hovered around 0.01 ms. The per size graph had a lot of noise, most likely because I did not run enough tests =P. Instead of running more tests (perhaps like I should have) I did a moving window average of 100 points before and average each point plotted in the graph above (and it still had a lot of noise). The large number of branches (if-statements) in the code may explain some of the observed variance. 

![Speed of adding an object per depth](http://briannoyama.github.io/assets/images/bvh-steps/AddRuntimePerDepth.svg)
![Speed of adding an object per number of volumes](http://briannoyama.github.io/assets/images/bvh-steps/AddRuntimePerSize.svg)

Subtractions worked closer to what I expected. The performance seems to increase linearly. It takes approximately 100 times as long to remove an item after 50,000 volumes have been added versus removing an item when there's only one in the hierarchy. (Note, since the BVH is also a binary tree, there are log<sub>2</sub>(50,000) parent volumes.) Other than height, the runtime performance also depends on the surface area of the volumes. The surface area is a good metric for the odds that a random query (or subtraction) volume will have to search multiple paths in the tree. It may also explain why we do not have _log_ performance.

![Speed of removing an object per number of volumes](http://briannoyama.github.io/assets/images/bvh-steps/SubRuntimePerSize.svg)
![Speed of removing an object per depth](http://briannoyama.github.io/assets/images/bvh-steps/SubRuntimePerDepth.svg)

Due to the random nature of the test, there does not exist data for smaller BVH's with smaller depths for all of the possible return values. (Query 2 means two volumes were returned or intersected by the query volume.) The query speed seems to mirror the subtraction speed. After a depth of ~10 the speed of both subtraction and query is slower than that for add (~0.01 ms). Surprisingly, the number of volumes affected seemed to have very little affect on the performance. This is likely because the query method does not have to recurse back to the root of the tree to find more things to return. 

![Speed of querying an object per depth](http://briannoyama.github.io/assets/images/bvh-steps/QueryPerDepth.svg)
![Speed of querying an object per number of volumes](http://briannoyama.github.io/assets/images/bvh-steps/QueryPerSize.svg)

As mentioned earlier there are two things that should (in theory) determine the performance of a BVH. One is the depth of the tree, and the other is the surface area of the tree. The different sizes of the parent volumes affects the total surface area. This test only relied on additions for the online method. The offline represents an approximate best possible tree.

![Depth of the online vs offline algorithms](http://briannoyama.github.io/assets/images/bvh-steps/Depth.svg)
![Surface area of the online vs offline algorithms](http://briannoyama.github.io/assets/images/bvh-steps/SurfaceArea.svg)

As one can see, the offline tree does not create as good of a tree (which we expect), but! It is close, and it grows at the same rate. For this study we ended at around 32000 added volumes due to the time it took to create an offline tree. For applications where some of the data can be preprocessed online, one can use the offline method to construct an initial hierarchy.

A few thoughts about the performance: There are a large number of relatively small method calls that are not likely inlined (which ones? I leave this as an activity for the reader). Currently for moving an existing volume, one needs to do a removal followed by an addition. The results from the query study suggests that for volumes that only need to be moved a small amount, it may be possible to make a better movement method that would take approximately half the time.

I did not do studies for the memory usage, though one can probably get a good estimate from looking at the code (fairly minimal). If one has questions, feel free to email me.

*This is not an officially supported Google product.
