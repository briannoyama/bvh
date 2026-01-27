# Online Bounding Volume Hierarchy

_Improvements compared to v0.*.*: Uses generic shapes + Add, Sub, and Query latency reduced by ~15%, 70% and 50%_

### Intro

This code is a golang implementation of a binary self-balancing Bounding Volume Hierarchy (BVH) inspired by the tree rotations in [Fast, Effective BVH Updates for Animated Scenes](https://www.cs.utah.edu/~aek/research/tree.pdf). The BVH can be used with orthotopes (ie. Axis Aligned Bounding Boxes or AABB), spheres or other data types that implement the `volume.Volume` interface. The hierarchies created via this algorithm have the following properties for _n_ volumes of any natural dimension:

- Average _log(n)_ addition/insertion time.
- Average _log(n)_ removal time.
- Average _mlog(n)_ query time where m is the number of volumes found.

Note that this is the _Average_ as the big _O_ depends on the input (similar to hashmaps).

#### Example Use Cases:

- Collisions between objects in a game or for ray tracing.
- Dynamically updating a search index for n-dimentional vectors (e.g. word-vectors).

#### Available Tree implementations:

|         | Orthotope | Sphere  |
| ------- | --------- | ------- |
| int32   | &#9989;   | &#9989; |
| int     | &#9989;   | &#9989; |
| float32 | &#9989;   | &#9989; |

- Note: this implementation is generic. BYO implementation of `volume.Volume`. See `volume.Orthotope` or `volume.Sphere` for an example.

### How it Works

Queries are thread-safe; however, additions and removals are not. The animations below show the algorithm in action:

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
      <img style="image-rendering: pixelated;" alt="Animated steps of showing addition of volumes to the BVH" width="200" src="https://github.com/briannoyama/bvhstats/blob/main/assets/add.gif">
    </td>
    <td>
      <img style="image-rendering: pixelated;" alt="Animated steps of showing removal of volumes to the BVH" width="200" src="https://github.com/briannoyama/bvhstats/blob/main/assets/remove0.gif">
    </td>
    <td>
      <img style="image-rendering: pixelated;" alt="Animated steps of showing an alernative removal of volumes to the BVH" width="200" src="https://github.com/briannoyama/bvhstats/blob/main/assets/remove1.gif">
    </td>
  </tr>
</table>

See `bvh/tree_test.go` for an example of how to use a `bvh.Tree`. Below is a short code snippet showing most of the interesting functionality:

```golang
import "github.com/briannoyama/bvh/bvh"
import "github.com/briannoyama/bvh/volume"

...

    // Pre allocate for 10 key (volume.Orthotope[int32]) value (string) pairs.
    tree := bvh.NewInt32OrthTree[string](10)

    // Add a cube of length 2 to the coordinate (2, 2, 4)
    // Change DIM as needed for your use-case
    ref := tree.Add(volume.Orthotope[int32]{
      P0: [volume.DIM]int32{2, 2, 4},
      P1: [volume.DIM]int32{4, 4, 6}},
    }, "Hello BVH!")

    q := volume.Orthotope[int32]{
      P0: [volume.DIM]int32{1, 1, 3},
      P1: [volume.DIM]int32{3, 3, 5},
    }

    for k, v := range tree.Query(q) {
      // Do something cool when things overlap
    }

    q.P1[0] = -1
    var td float32
    for k, v := range tree.Intersects(q, [volume.DIM]int32{4, 0, 0}, &td) {
      // Do something cool when a volume with a velocity intersects at time td.
    }

    tree.Remove(ref)
```

### Performance

To improve performance, the algorithm swaps child nodes within the BVH tree both to balance the tree and to reduce the Surface Area of the generated bounding volumes. This leads to trees that with generally lower (more efficient) Surface Area Heuristics when compared to simple sort and split offline algorithms.

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
      <img style="image-rendering: pixelated;" alt="Output of online algorithm for generating BVH" width="200" src="http://github.com/briannoyama/bvhstats/blob/main/assets/online.png">
    </td>
    <td>
      <img style="image-rendering: pixelated;" alt="Output of offline algorithm for generating BVH" width="200" src="http://github.com/briannoyama/bvhstats/blob/main/assets/offline.png">
    </td>
  </tr>
</table>

For those who plan to use onlineBVH for an application with strict runtime requirements, I conducted a small experiment on an AMD Ryzen 7 7735HS. The test generated random int32 Orthotopes in a 3D space to add (100,000) remove (50,000) and query (100,000) such that the final BVH would contain 50,000 items. Running this test 5 times and combining the data gave the following performance graphs:

![Latency of adding an object per number of volumes](http://github.com/briannoyama/bvhstats/blob/main/assets/AddRuntimePerSize.svg)
![Latency of adding an object per depth](http://github.com/briannoyama/bvhstats/blob/main/assets/AddRuntimePerDepth.svg)

I'm not 100 on why the algorithm improves in efficiency ~12K volumes (golang magic? Please email me if you have any ideas). Runtime per depth shows (with the exception of needing to warm up the cache for initial adds), that the latency of adds/subtracts increases linearly with the depth of the tree, or logarithmically with the number of volumes added.

![Latency of removing an object per number of volumes](http://github.com/briannoyama/bvhstats/blob/main/assets/SubRuntimePerSize.svg)
![Latency of removing an object per depth](http://github.com/briannoyama/bvhstats/blob/main/assets/SubRuntimePerDepth.svg)

Querying follows this pattern, and exhibits similar latency largely independent of how many objects are returned.

![Latency of querying an object per depth](http://github.com/briannoyama/bvhstats/blob/main/assets/QueryPerDepth.svg)
![Latency of querying an object per number of volumes](http://github.com/briannoyama/bvhstats/blob/main/assets/QueryPerSize.svg)

For any questions, feel free to email me.
