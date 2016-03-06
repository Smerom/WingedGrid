package wingedGrid

import (
    
)
const goldenRatio = 1.61803398875

func BaseIcosahedron() (WingedGrid, error){
    var icosahedron WingedGrid
    
    icosahedron.Faces = []WingedFace {
        // cap 1
        {Edges: []int32{ 0,  4,  9}}, // 0
        {Edges: []int32{ 0,  1,  5}}, // 1
        {Edges: []int32{ 1,  2,  6}}, // 2
        {Edges: []int32{ 2,  3,  7}}, // 3
        {Edges: []int32{ 3,  4,  8}}, // 4
        // cap 2
        {Edges: []int32{10, 14, 19}}, // 5
        {Edges: []int32{10, 11, 15}}, // 6
        {Edges: []int32{11, 12, 16}}, // 7
        {Edges: []int32{12, 13, 17}}, // 8
        {Edges: []int32{13, 14, 18}}, // 9
        // ring of 10 between caps, clockwise around cap 1
        //   starting with the first face adjacent to face 0
        {Edges: []int32{ 9, 28, 29}}, // 10
        {Edges: []int32{17, 20, 29}}, // 11
        {Edges: []int32{ 5, 20, 21}}, // 12
        {Edges: []int32{18, 21, 22}}, // 13
        {Edges: []int32{ 6, 22, 23}}, // 14
        {Edges: []int32{19, 23, 24}}, // 15
        {Edges: []int32{ 7, 24, 25}}, // 16
        {Edges: []int32{15, 25, 26}}, // 17
        {Edges: []int32{ 8, 26, 27}}, // 18
        {Edges: []int32{16, 27, 28}}, // 19
    }
    icosahedron.Vertices = []WingedVertex{
        // y-z plane rectangle
        {Coords: [3]float64{ 0, 1, goldenRatio},
         Edges:  []int32{}}, // 0
        {Coords: [3]float64{ 0, 1,-goldenRatio},
         Edges:  []int32{}}, // 1
        {Coords: [3]float64{ 0,-1, goldenRatio},
         Edges:  []int32{}}, // 2
        {Coords: [3]float64{ 0,-1,-goldenRatio},
         Edges:  []int32{}}, // 3
        // x-z plane rectangle
        {Coords: [3]float64{ goldenRatio, 0, 1},
         Edges:  []int32{}}, // 4
        {Coords: [3]float64{-goldenRatio, 0, 1},
         Edges:  []int32{}}, // 5
        {Coords: [3]float64{ goldenRatio, 0,-1},
         Edges:  []int32{}}, // 6
        {Coords: [3]float64{-goldenRatio, 0,-1},
         Edges:  []int32{}}, // 7
        // x-y plane rectangle
        {Coords: [3]float64{ 1, goldenRatio, 0},
         Edges:  []int32{}}, // 8
        {Coords: [3]float64{ 1,-goldenRatio, 0},
         Edges:  []int32{}}, // 9
        {Coords: [3]float64{-1, goldenRatio, 0},
         Edges:  []int32{}}, // 10
        {Coords: [3]float64{-1,-goldenRatio, 0},
         Edges:  []int32{}}, // 11
    }
    icosahedron.Edges = []WingedEdge{
        // cap 1 around vertex 0
        // 5 spokes starting between face 0 and 1,
        //   clockwise around starting with short edge
        //   of the rectangle (y-z)
        {
            Vertex1: 0, Vertex2: 2,
            FaceA: 0, FaceB: 1,
            PrevA: 9, NextA: 4,
            PrevB: 1, NextB: 5,
        },  // 0
        {
            Vertex1: 0, Vertex2: 4,
            FaceA: 1, FaceB: 2,
            PrevA: 5, NextA: 0,
            PrevB: 2, NextB: 6,
        }, // 1
        {
            Vertex1: 0, Vertex2: 8,
            FaceA: 2, FaceB: 3,
            PrevA: 6, NextA: 1,
            PrevB: 3, NextB: 7,
        }, // 2
        {
            Vertex1: 0, Vertex2: 10,
            FaceA: 3, FaceB: 4,
            PrevA: 7, NextA: 2,
            PrevB: 4, NextB: 8,
        }, // 3
        {
            Vertex1: 0, Vertex2: 5,
            FaceA: 4, FaceB: 0,
            PrevA: 8, NextA: 3,
            PrevB: 0, NextB: 9,
        }, // 4
        
        // ring of 5
        {
            Vertex1: 2, Vertex2: 4,
            FaceA: 1, FaceB: 12,
            PrevA: 0, NextA: 1,
            PrevB: 21, NextB: 20,
        }, // 5
        {
            Vertex1: 4, Vertex2: 8,
            FaceA: 2, FaceB: 14,
            PrevA: 1, NextA: 2,
            PrevB: 23, NextB: 22,
        }, // 6
        {
            Vertex1: 8, Vertex2: 10,
            FaceA: 3, FaceB: 16,
            PrevA: 2, NextA: 3,
            PrevB: 25, NextB: 24,
        }, // 7
        {
            Vertex1: 10, Vertex2: 5,
            FaceA: 4, FaceB: 18,
            PrevA: 3, NextA: 4,
            PrevB: 27, NextB: 26,
        }, // 8
        {
            Vertex1: 5, Vertex2: 2,
            FaceA: 0, FaceB: 10,
            PrevA: 4, NextA: 0,
            PrevB: 29, NextB: 28,
        }, // 9
        
        // cap 2 around vertex 3
        // 5 spokes starting between face 5 and 6,
        //  counter-clockwise from short edge (y-z) rectangle
        {
            Vertex1: 3, Vertex2: 1,
            FaceA: 5, FaceB: 6,
            PrevA: 14, NextA: 19,
            PrevB: 15, NextB: 11,
        }, // 10
        {
            Vertex1: 3, Vertex2: 7,
            FaceA: 6, FaceB: 7,
            PrevA: 10, NextA: 15,
            PrevB: 16, NextB: 12,
        }, // 11
        {
            Vertex1: 3, Vertex2: 11,
            FaceA: 7, FaceB: 8,
            PrevA: 11, NextA: 16,
            PrevB: 17, NextB: 13,
        }, // 12
        {
            Vertex1: 3, Vertex2: 9,
            FaceA: 8, FaceB: 9,
            PrevA: 12, NextA: 17,
            PrevB: 18, NextB: 14,
        }, // 13
        {
            Vertex1: 3, Vertex2: 6,
            FaceA: 9, FaceB: 5,
            PrevA: 13, NextA: 18,
            PrevB: 19, NextB: 10,
        }, // 14
        
        // ring of 5
        {
            Vertex1: 1, Vertex2: 7,
            FaceA: 6, FaceB: 17,
            PrevA: 11, NextA: 10,
            PrevB: 25, NextB: 26,
        }, // 15
        {
            Vertex1: 7, Vertex2: 11,
            FaceA: 7, FaceB: 19,
            PrevA: 12, NextA: 11,
            PrevB: 27, NextB: 28,
        }, // 16
        {
            Vertex1: 11, Vertex2: 9,
            FaceA: 8, FaceB: 11,
            PrevA: 13, NextA: 12,
            PrevB: 29, NextB: 20,
        }, // 17
        {
            Vertex1: 9, Vertex2: 6,
            FaceA: 9, FaceB: 13,
            PrevA: 14, NextA: 13,
            PrevB: 21, NextB: 22,
        }, // 18
        {
            Vertex1: 6, Vertex2: 1,
            FaceA: 5, FaceB: 15,
            PrevA: 10, NextA: 14,
            PrevB: 23, NextB: 24,
        }, // 19
        
        
        // zig-zag down the middle
        // 10 triangles, 10 new edges
        //   starting clockwise from end of edge 0
        //   
        {
            Vertex1: 2, Vertex2: 9,
            FaceA: 11, FaceB: 12,
            PrevA: 17, NextA: 29,
            PrevB: 5, NextB: 21,
        }, // 20
        {
            Vertex1: 9, Vertex2: 4,
            FaceA: 12, FaceB: 13,
            PrevA: 20, NextA: 5,
            PrevB: 22, NextB: 18,
        }, // 21
        {
            Vertex1: 4, Vertex2: 6,
            FaceA: 13, FaceB: 14,
            PrevA: 18, NextA: 21,
            PrevB: 6, NextB: 23,
        }, // 22
        {
            Vertex1: 6, Vertex2: 8,
            FaceA: 14, FaceB: 15,
            PrevA: 22, NextA: 6,
            PrevB: 24, NextB: 19,
        }, // 23
        {
            Vertex1: 8, Vertex2: 1,
            FaceA: 15, FaceB: 16,
            PrevA: 19, NextA: 23,
            PrevB: 7, NextB: 25,
        }, // 24
        {
            Vertex1: 1, Vertex2: 10,
            FaceA: 16, FaceB: 17,
            PrevA: 24, NextA: 7,
            PrevB: 26, NextB: 15,
        }, // 25
        {
            Vertex1: 10, Vertex2: 7,
            FaceA: 17, FaceB: 18,
            PrevA: 15, NextA: 25,
            PrevB: 8, NextB: 27,
        }, // 26
        {
            Vertex1: 7, Vertex2: 5,
            FaceA: 18, FaceB: 19,
            PrevA: 26, NextA: 8,
            PrevB: 28, NextB: 16,
        }, // 27
        {
            Vertex1: 5, Vertex2: 11,
            FaceA: 19, FaceB: 10,
            PrevA: 16, NextA: 27,
            PrevB: 9, NextB: 29,
        }, // 28
        {
            Vertex1: 11, Vertex2: 2,
            FaceA: 10, FaceB: 11,
            PrevA: 28, NextA: 9,
            PrevB: 20, NextB: 17,
        }, // 29
        
    }
    return icosahedron, nil
}