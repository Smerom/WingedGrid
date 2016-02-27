package wingedGrid

import (
    
)
const goldenRatio = 1.61803398875
type WingedFace struct {
    // three for a triangular tiling, but support others
    // index of edge in wingedMap
    edges []int32
}

type WingedEdge struct {
    // index of vertices in winged map
    vertex1, vertex2 int32
    // index of faces in winged map
    faceA, faceB int32
    // index of edges in winged map for face A and face B
    prevA, nextA, prevB, nextB int32
}

type WingedVertex struct {
    coords [3]float64
}

type WingedMap struct {
    faces []WingedFace
    edges []WingedEdge
    vertices []WingedVertex
}


func BaseIcosahedron() (WingedMap, error){
    var icosahedron WingedMap
    
    icosahedron.faces = []WingedFace {
        // cap 1
        {edges: []int32{ 0,  4,  9}}, // 0
        {edges: []int32{ 0,  1,  5}}, // 1
        {edges: []int32{ 1,  2,  6}}, // 2
        {edges: []int32{ 2,  3,  7}}, // 3
        {edges: []int32{ 3,  4,  8}}, // 4
        // cap 2
        {edges: []int32{10, 14, 19}}, // 5
        {edges: []int32{10, 11, 15}}, // 6
        {edges: []int32{11, 12, 16}}, // 7
        {edges: []int32{12, 13, 17}}, // 8
        {edges: []int32{13, 14, 18}}, // 9
        // ring of 10 between caps, clockwise around cap 1
        //   starting with the first face adjacent to face 0
        {edges: []int32{ 9, 28, 29}}, // 10
        {edges: []int32{17, 20, 29}}, // 11
        {edges: []int32{ 5, 20, 21}}, // 12
        {edges: []int32{18, 21, 22}}, // 13
        {edges: []int32{ 6, 22, 23}}, // 14
        {edges: []int32{19, 23, 24}}, // 15
        {edges: []int32{ 7, 24, 25}}, // 16
        {edges: []int32{15, 25, 26}}, // 17
        {edges: []int32{ 8, 26, 27}}, // 18
        {edges: []int32{16, 27, 28}}, // 19
    }
    icosahedron.vertices = []WingedVertex{
        // y-z plane rectangle
        {coords: [3]float64{ 0, 1, goldenRatio}}, // 0
        {coords: [3]float64{ 0, 1,-goldenRatio}}, // 1
        {coords: [3]float64{ 0,-1, goldenRatio}}, // 2
        {coords: [3]float64{ 0,-1,-goldenRatio}}, // 3
        // x-z plane rectangle
        {coords: [3]float64{ goldenRatio, 0, 1}}, // 4
        {coords: [3]float64{-goldenRatio, 0, 1}}, // 5
        {coords: [3]float64{ goldenRatio, 0,-1}}, // 6
        {coords: [3]float64{-goldenRatio, 0,-1}}, // 7
        // x-y plane rectangle
        {coords: [3]float64{ 1, goldenRatio, 0}}, // 8
        {coords: [3]float64{ 1,-goldenRatio, 0}}, // 9
        {coords: [3]float64{-1, goldenRatio, 0}}, // 10
        {coords: [3]float64{-1,-goldenRatio, 0}}, // 11
    }
    icosahedron.edges = []WingedEdge{
        // cap 1 around vertex 0
        // 5 spokes starting between face 0 and 1,
        //   clockwise around starting with short edge
        //   of the rectangle (y-z)
        {
            vertex1: 0, vertex2: 2,
            faceA: 0, faceB: 1,
            prevA: 9, nextA: 4,
            prevB: 1, nextB: 5,
        },  // 0
        {
            vertex1: 0, vertex2: 4,
            faceA: 1, faceB: 2,
            prevA: 5, nextA: 0,
            prevB: 2, nextB: 6,
        }, // 1
        {
            vertex1: 0, vertex2: 8,
            faceA: 2, faceB: 3,
            prevA: 6, nextA: 1,
            prevB: 3, nextB: 7,
        }, // 2
        {
            vertex1: 0, vertex2: 10,
            faceA: 3, faceB: 4,
            prevA: 7, nextA: 2,
            prevB: 4, nextB: 8,
        }, // 3
        {
            vertex1: 0, vertex2: 5,
            faceA: 4, faceB: 0,
            prevA: 8, nextA: 3,
            prevB: 0, nextB: 9,
        }, // 4

        
        // ring of 5
        {
            vertex1: 2, vertex2: 4,
            faceA: 1, faceB: 12,
            prevA: 0, nextA: 1,
            prevB: 21, nextB: 20,
        }, // 5
        {
            vertex1: 4, vertex2: 8,
            faceA: 2, faceB: 14,
            prevA: 1, nextA: 2,
            prevB: 23, nextB: 22,
        }, // 6
        {
            vertex1: 8, vertex2: 10,
            faceA: 3, faceB: 16,
            prevA: 2, nextA: 3,
            prevB: 25, nextB: 24,
        }, // 7
        {
            vertex1: 10, vertex2: 5,
            faceA: 4, faceB: 18,
            prevA: 3, nextA: 4,
            prevB: 27, nextB: 26,
        }, // 8
        {
            vertex1: 5, vertex2: 2,
            faceA: 0, faceB: 10,
            prevA: 4, nextA: 0,
            prevB: 29, nextB: 28,
        }, // 9
        
        // cap 2 around vertex 3
        // 5 spokes starting between face 5 and 6,
        //  counter-clockwise from short edge (y-z) rectangle
        {
            vertex1: 3, vertex2: 1,
            faceA: 5, faceB: 6,
            prevA: 14, nextA: 19,
            prevB: 15, nextB: 11,
        }, // 10
        {
            vertex1: 3, vertex2: 7,
            faceA: 6, faceB: 7,
            prevA: 10, nextA: 15,
            prevB: 16, nextB: 12,
        }, // 11
        {
            vertex1: 3, vertex2: 11,
            faceA: 7, faceB: 8,
            prevA: 11, nextA: 16,
            prevB: 17, nextB: 13,
        }, // 12
        {
            vertex1: 3, vertex2: 9,
            faceA: 8, faceB: 9,
            prevA: 12, nextA: 17,
            prevB: 18, nextB: 14,
        }, // 13
        {
            vertex1: 3, vertex2: 6,
            faceA: 9, faceB: 5,
            prevA: 13, nextA: 18,
            prevB: 19, nextB: 10,
        }, // 14
        
        // ring of 5
        {
            vertex1: 1, vertex2: 7,
            faceA: 6, faceB: 17,
            prevA: 11, nextA: 10,
            prevB: 25, nextB: 26,
        }, // 15
        {
            vertex1: 7, vertex2: 11,
            faceA: 7, faceB: 19,
            prevA: 12, nextA: 11,
            prevB: 27, nextB: 28,
        }, // 16
        {
            vertex1: 11, vertex2: 9,
            faceA: 8, faceB: 11,
            prevA: 13, nextA: 12,
            prevB: 29, nextB: 20,
        }, // 17
        {
            vertex1: 9, vertex2: 6,
            faceA: 9, faceB: 13,
            prevA: 14, nextA: 13,
            prevB: 21, nextB: 22,
        }, // 18
        {
            vertex1: 6, vertex2: 1,
            faceA: 5, faceB: 15,
            prevA: 10, nextA: 14,
            prevB: 23, nextB: 24,
        }, // 19
        
        
        // zig-zag down the middle
        // 10 triangles, 10 new edges
        //   starting clockwise from end of edge 0
        //   
        {
            vertex1: 2, vertex2: 9,
            faceA: 11, faceB: 12,
            prevA: 17, nextA: 29,
            prevB: 5, nextB: 21,
        }, // 20
        {
            vertex1: 9, vertex2: 4,
            faceA: 12, faceB: 13,
            prevA: 20, nextA: 5,
            prevB: 22, nextB: 18,
        }, // 21
        {
            vertex1: 4, vertex2: 6,
            faceA: 13, faceB: 14,
            prevA: 18, nextA: 21,
            prevB: 6, nextB: 23,
        }, // 22
        {
            vertex1: 6, vertex2: 8,
            faceA: 14, faceB: 15,
            prevA: 22, nextA: 6,
            prevB: 24, nextB: 19,
        }, // 23
        {
            vertex1: 8, vertex2: 1,
            faceA: 15, faceB: 16,
            prevA: 19, nextA: 23,
            prevB: 7, nextB: 25,
        }, // 24
        {
            vertex1: 1, vertex2: 10,
            faceA: 16, faceB: 17,
            prevA: 24, nextA: 7,
            prevB: 26, nextB: 15,
        }, // 25
        {
            vertex1: 10, vertex2: 7,
            faceA: 17, faceB: 18,
            prevA: 15, nextA: 25,
            prevB: 8, nextB: 27,
        }, // 26
        {
            vertex1: 7, vertex2: 5,
            faceA: 18, faceB: 19,
            prevA: 26, nextA: 8,
            prevB: 28, nextB: 16,
        }, // 27
        {
            vertex1: 5, vertex2: 11,
            faceA: 19, faceB: 10,
            prevA: 16, nextA: 27,
            prevB: 9, nextB: 29,
        }, // 28
        {
            vertex1: 11, vertex2: 2,
            faceA: 10, faceB: 11,
            prevA: 28, nextA: 9,
            prevB: 20, nextB: 17,
        }, // 29
        
    }
    return icosahedron, nil
}