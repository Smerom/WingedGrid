package wingedGrid

import (
    
)
const goldenRatio = 1.61803398875
// Sets up our base icosahedron from three golden retangles, resulting in edges
// or length 2
func BaseIcosahedron() (WingedGrid, error){
    var icosahedron WingedGrid
    
    icosahedron.Faces = []WingedFace {
        // Duplicate info could be done pragmatically after contructing edges
        // must be in clockwise order
        // cap 1
        {Edges: []int32{ 0,  4,  9}}, // 0
        {Edges: []int32{ 0,  5,  1}}, // 1
        {Edges: []int32{ 1,  6,  2}}, // 2
        {Edges: []int32{ 2,  7,  3}}, // 3
        {Edges: []int32{ 3,  8,  4}}, // 4
        // cap 2
        {Edges: []int32{10, 19, 14}}, // 5
        {Edges: []int32{10, 11, 15}}, // 6
        {Edges: []int32{11, 12, 16}}, // 7
        {Edges: []int32{12, 13, 17}}, // 8
        {Edges: []int32{13, 14, 18}}, // 9
        // Ring of 10 between caps, clockwise around cap 1
        //   starting with the first face adjacent to face 0
        {Edges: []int32{ 9, 28, 29}}, // 10
        {Edges: []int32{17, 20, 29}}, // 11
        {Edges: []int32{ 5, 20, 21}}, // 12
        {Edges: []int32{18, 22, 21}}, // 13
        {Edges: []int32{ 6, 22, 23}}, // 14
        {Edges: []int32{19, 24, 23}}, // 15
        {Edges: []int32{ 7, 24, 25}}, // 16
        {Edges: []int32{15, 26, 25}}, // 17
        {Edges: []int32{ 8, 26, 27}}, // 18
        {Edges: []int32{16, 28, 27}}, // 19
    }
    icosahedron.Vertices = []WingedVertex{
        // Edges are duplicate info, could be done pragmatically after constucting
        // edges
        // Must be in clockwise order
        // y-z plane rectangle
        {Coords: [3]float64{ 0, 1, goldenRatio},
         Edges:  []int32{ 0,  1,  2,  3,  4}}, // 0
        {Coords: [3]float64{ 0, 1,-goldenRatio},
         Edges:  []int32{10, 15, 25, 24, 19}}, // 1
        {Coords: [3]float64{ 0,-1, goldenRatio},
         Edges:  []int32{ 0,  9, 29, 20,  5}}, // 2
        {Coords: [3]float64{ 0,-1,-goldenRatio},
         Edges:  []int32{10, 14, 13, 12, 11}}, // 3
        // x-z plane rectangle
        {Coords: [3]float64{-goldenRatio, 0, 1},
         Edges:  []int32{ 1,  5, 21, 22,  6}}, // 4
        {Coords: [3]float64{ goldenRatio, 0, 1},
         Edges:  []int32{ 4,  8, 27, 28,  9}}, // 5
        {Coords: [3]float64{-goldenRatio, 0,-1},
         Edges:  []int32{14, 19, 23, 22, 18}}, // 6
        {Coords: [3]float64{ goldenRatio, 0,-1},
         Edges:  []int32{11, 16, 27, 26, 15}}, // 7
        // x-y plane rectangle
        {Coords: [3]float64{-1, goldenRatio, 0},
         Edges:  []int32{ 2,  6, 23, 24,  7}}, // 8
        {Coords: [3]float64{-1,-goldenRatio, 0},
         Edges:  []int32{13, 18, 21, 20, 17}}, // 9
        {Coords: [3]float64{ 1, goldenRatio, 0},
         Edges:  []int32{ 3,  7, 25, 26,  8}}, // 10
        {Coords: [3]float64{ 1,-goldenRatio, 0},
         Edges:  []int32{12, 17, 29, 28, 16}}, // 11
    }
    icosahedron.Edges = []WingedEdge{
        // cap 1 around vertex 0
        // 5 spokes starting between face 0 and 1,
        // clockwise around starting with short edge
        // of the rectangle (y-z)
        {
            FirstVertexA: 2, FirstVertexB: 0,
            FaceA: 0, FaceB: 1,
            PrevA: 9, NextA: 4,
            PrevB: 1, NextB: 5,
        },  // 0
        {
            FirstVertexA: 4, FirstVertexB: 0,
            FaceA: 1, FaceB: 2,
            PrevA: 5, NextA: 0,
            PrevB: 2, NextB: 6,
        }, // 1
        {
            FirstVertexA: 8, FirstVertexB: 0,
            FaceA: 2, FaceB: 3,
            PrevA: 6, NextA: 1,
            PrevB: 3, NextB: 7,
        }, // 2
        {
            FirstVertexA: 10, FirstVertexB: 0,
            FaceA: 3, FaceB: 4,
            PrevA: 7, NextA: 2,
            PrevB: 4, NextB: 8,
        }, // 3
        {
            FirstVertexA: 5, FirstVertexB: 0,
            FaceA: 4, FaceB: 0,
            PrevA: 8, NextA: 3,
            PrevB: 0, NextB: 9,
        }, // 4
        // ring of 5 around the base of the cap
        {
            FirstVertexA: 2, FirstVertexB: 4,
            FaceA: 1, FaceB: 12,
            PrevA: 0, NextA: 1,
            PrevB: 21, NextB: 20,
        }, // 5
        {
            FirstVertexA: 4, FirstVertexB: 8,
            FaceA: 2, FaceB: 14,
            PrevA: 1, NextA: 2,
            PrevB: 23, NextB: 22,
        }, // 6
        {
            FirstVertexA: 8, FirstVertexB: 10,
            FaceA: 3, FaceB: 16,
            PrevA: 2, NextA: 3,
            PrevB: 25, NextB: 24,
        }, // 7
        {
            FirstVertexA: 10, FirstVertexB: 5,
            FaceA: 4, FaceB: 18,
            PrevA: 3, NextA: 4,
            PrevB: 27, NextB: 26,
        }, // 8
        {
            FirstVertexA: 5, FirstVertexB: 2,
            FaceA: 0, FaceB: 10,
            PrevA: 4, NextA: 0,
            PrevB: 29, NextB: 28,
        }, // 9
        // cap 2 around vertex 3
        // 5 spokes starting between face 5 and 6,
        // counter-clockwise from short edge (y-z) rectangle
        {
            FirstVertexA: 3, FirstVertexB: 1,
            FaceA: 5, FaceB: 6,
            PrevA: 14, NextA: 19,
            PrevB: 15, NextB: 11,
        }, // 10
        {
            FirstVertexA: 3, FirstVertexB: 7,
            FaceA: 6, FaceB: 7,
            PrevA: 10, NextA: 15,
            PrevB: 16, NextB: 12,
        }, // 11
        {
            FirstVertexA: 3, FirstVertexB: 11,
            FaceA: 7, FaceB: 8,
            PrevA: 11, NextA: 16,
            PrevB: 17, NextB: 13,
        }, // 12
        {
            FirstVertexA: 3, FirstVertexB: 9,
            FaceA: 8, FaceB: 9,
            PrevA: 12, NextA: 17,
            PrevB: 18, NextB: 14,
        }, // 13
        {
            FirstVertexA: 3, FirstVertexB: 6,
            FaceA: 9, FaceB: 5,
            PrevA: 13, NextA: 18,
            PrevB: 19, NextB: 10,
        }, // 14
        // ring of 5 around the base of cap 2
        {
            FirstVertexA: 7, FirstVertexB: 1,
            FaceA: 6, FaceB: 17,
            PrevA: 11, NextA: 10,
            PrevB: 25, NextB: 26,
        }, // 15
        {
            FirstVertexA: 11, FirstVertexB: 7,
            FaceA: 7, FaceB: 19,
            PrevA: 12, NextA: 11,
            PrevB: 27, NextB: 28,
        }, // 16
        {
            FirstVertexA: 9, FirstVertexB: 11,
            FaceA: 8, FaceB: 11,
            PrevA: 13, NextA: 12,
            PrevB: 29, NextB: 20,
        }, // 17
        {
            FirstVertexA: 6, FirstVertexB: 9,
            FaceA: 9, FaceB: 13,
            PrevA: 14, NextA: 13,
            PrevB: 21, NextB: 22,
        }, // 18
        {
            FirstVertexA: 1, FirstVertexB: 6,
            FaceA: 5, FaceB: 15,
            PrevA: 10, NextA: 14,
            PrevB: 23, NextB: 24,
        }, // 19
        // zig-zag down the middle
        // 10 triangles, 10 new edges
        // starting clockwise from end of edge 0
        {
            FirstVertexA: 9, FirstVertexB: 2,
            FaceA: 11, FaceB: 12,
            PrevA: 17, NextA: 29,
            PrevB: 5, NextB: 21,
        }, // 20
        {
            FirstVertexA: 9, FirstVertexB: 4,
            FaceA: 12, FaceB: 13,
            PrevA: 20, NextA: 5,
            PrevB: 22, NextB: 18,
        }, // 21
        {
            FirstVertexA: 6, FirstVertexB: 4,
            FaceA: 13, FaceB: 14,
            PrevA: 18, NextA: 21,
            PrevB: 6, NextB: 23,
        }, // 22
        {
            FirstVertexA: 6, FirstVertexB: 8,
            FaceA: 14, FaceB: 15,
            PrevA: 22, NextA: 6,
            PrevB: 24, NextB: 19,
        }, // 23
        {
            FirstVertexA: 1, FirstVertexB: 8,
            FaceA: 15, FaceB: 16,
            PrevA: 19, NextA: 23,
            PrevB: 7, NextB: 25,
        }, // 24
        {
            FirstVertexA: 1, FirstVertexB: 10,
            FaceA: 16, FaceB: 17,
            PrevA: 24, NextA: 7,
            PrevB: 26, NextB: 15,
        }, // 25
        {
            FirstVertexA: 7, FirstVertexB: 10,
            FaceA: 17, FaceB: 18,
            PrevA: 15, NextA: 25,
            PrevB: 8, NextB: 27,
        }, // 26
        {
            FirstVertexA: 7, FirstVertexB: 5,
            FaceA: 18, FaceB: 19,
            PrevA: 26, NextA: 8,
            PrevB: 28, NextB: 16,
        }, // 27
        {
            FirstVertexA: 11, FirstVertexB: 5,
            FaceA: 19, FaceB: 10,
            PrevA: 16, NextA: 27,
            PrevB: 9, NextB: 29,
        }, // 28
        {
            FirstVertexA: 11, FirstVertexB: 2,
            FaceA: 10, FaceB: 11,
            PrevA: 28, NextA: 9,
            PrevB: 20, NextB: 17,
        }, // 29 
    }
    return icosahedron, nil
}