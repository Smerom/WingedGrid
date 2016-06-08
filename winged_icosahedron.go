package wingedGrid

const goldenRatio = 1.61803398875

// Sets up our base icosahedron from three golden retangles, resulting in edges
// or length 2
func BaseIcosahedron() (WingedGrid, error) {
	var icosahedron WingedGrid

	icosahedron.Faces = []WingedFace{
		// Duplicate info could be done pragmatically after contructing edges
		// must be in counter-clockwise order
		// cap 1
		{Edges: []int32{9, 4, 0}}, // 0
		{Edges: []int32{1, 5, 0}}, // 1
		{Edges: []int32{2, 6, 1}}, // 2
		{Edges: []int32{3, 7, 2}}, // 3
		{Edges: []int32{4, 8, 3}}, // 4
		// cap 2
		{Edges: []int32{14, 19, 10}}, // 5
		{Edges: []int32{15, 11, 10}}, // 6
		{Edges: []int32{16, 12, 11}}, // 7
		{Edges: []int32{17, 13, 12}}, // 8
		{Edges: []int32{18, 14, 13}}, // 9
		// Ring of 10 between caps, clockwise around cap 1
		//   starting with the first face adjacent to face 0
		{Edges: []int32{29, 28, 9}},  // 10
		{Edges: []int32{29, 20, 17}}, // 11
		{Edges: []int32{21, 20, 5}},  // 12
		{Edges: []int32{21, 22, 18}}, // 13
		{Edges: []int32{23, 22, 6}},  // 14
		{Edges: []int32{23, 24, 19}}, // 15
		{Edges: []int32{25, 24, 7}},  // 16
		{Edges: []int32{25, 26, 15}}, // 17
		{Edges: []int32{27, 26, 8}},  // 18
		{Edges: []int32{27, 28, 16}}, // 19
	}
	icosahedron.Vertices = []WingedVertex{
		// Edges are duplicate info, could be done pragmatically after constucting
		// edges
		// Must be in counter-clockwise order
		// y-z plane rectangle
		{Coords: [3]float64{0, 1, goldenRatio},
			Edges: []int32{4, 3, 2, 1, 0}}, // 0
		{Coords: [3]float64{0, 1, -goldenRatio},
			Edges: []int32{19, 24, 25, 15, 10}}, // 1
		{Coords: [3]float64{0, -1, goldenRatio},
			Edges: []int32{5, 20, 29, 9, 0}}, // 2
		{Coords: [3]float64{0, -1, -goldenRatio},
			Edges: []int32{11, 12, 13, 14, 10}}, // 3
		// x-z plane rectangle
		{Coords: [3]float64{-goldenRatio, 0, 1},
			Edges: []int32{6, 22, 21, 5, 1}}, // 4
		{Coords: [3]float64{goldenRatio, 0, 1},
			Edges: []int32{9, 28, 27, 8, 4}}, // 5
		{Coords: [3]float64{-goldenRatio, 0, -1},
			Edges: []int32{18, 22, 23, 19, 14}}, // 6
		{Coords: [3]float64{goldenRatio, 0, -1},
			Edges: []int32{15, 26, 27, 16, 11}}, // 7
		// x-y plane rectangle
		{Coords: [3]float64{-1, goldenRatio, 0},
			Edges: []int32{7, 24, 23, 6, 2}}, // 8
		{Coords: [3]float64{-1, -goldenRatio, 0},
			Edges: []int32{17, 20, 21, 18, 13}}, // 9
		{Coords: [3]float64{1, goldenRatio, 0},
			Edges: []int32{8, 26, 25, 7, 3}}, // 10
		{Coords: [3]float64{1, -goldenRatio, 0},
			Edges: []int32{16, 28, 29, 17, 12}}, // 11
	}
	icosahedron.Edges = []WingedEdge{
		// cap 1 around vertex 0
		// 5 spokes starting between face 0 and 1,
		// clockwise around starting with short edge
		// of the rectangle (y-z)
		{
			FirstVertexA: 0, FirstVertexB: 2,
			FaceA: 0, FaceB: 1,
			PrevA: 4, NextA: 9,
			PrevB: 5, NextB: 1,
		}, // 0
		{
			FirstVertexA: 0, FirstVertexB: 4,
			FaceA: 1, FaceB: 2,
			PrevA: 0, NextA: 5,
			PrevB: 6, NextB: 2,
		}, // 1
		{
			FirstVertexA: 0, FirstVertexB: 8,
			FaceA: 2, FaceB: 3,
			PrevA: 1, NextA: 6,
			PrevB: 7, NextB: 3,
		}, // 2
		{
			FirstVertexA: 0, FirstVertexB: 10,
			FaceA: 3, FaceB: 4,
			PrevA: 2, NextA: 7,
			PrevB: 8, NextB: 4,
		}, // 3
		{
			FirstVertexA: 0, FirstVertexB: 5,
			FaceA: 4, FaceB: 0,
			PrevA: 3, NextA: 8,
			PrevB: 9, NextB: 0,
		}, // 4
		// ring of 5 around the base of the cap
		{
			FirstVertexA: 4, FirstVertexB: 2,
			FaceA: 1, FaceB: 12,
			PrevA: 1, NextA: 0,
			PrevB: 20, NextB: 21,
		}, // 5
		{
			FirstVertexA: 8, FirstVertexB: 4,
			FaceA: 2, FaceB: 14,
			PrevA: 2, NextA: 1,
			PrevB: 22, NextB: 23,
		}, // 6
		{
			FirstVertexA: 10, FirstVertexB: 8,
			FaceA: 3, FaceB: 16,
			PrevA: 3, NextA: 2,
			PrevB: 24, NextB: 25,
		}, // 7
		{
			FirstVertexA: 5, FirstVertexB: 10,
			FaceA: 4, FaceB: 18,
			PrevA: 4, NextA: 3,
			PrevB: 26, NextB: 27,
		}, // 8
		{
			FirstVertexA: 2, FirstVertexB: 5,
			FaceA: 0, FaceB: 10,
			PrevA: 0, NextA: 4,
			PrevB: 28, NextB: 29,
		}, // 9
		// cap 2 around vertex 3
		// 5 spokes starting between face 5 and 6,
		// counter-clockwise from short edge (y-z) rectangle
		{
			FirstVertexA: 1, FirstVertexB: 3,
			FaceA: 5, FaceB: 6,
			PrevA: 19, NextA: 14,
			PrevB: 11, NextB: 15,
		}, // 10
		{
			FirstVertexA: 7, FirstVertexB: 3,
			FaceA: 6, FaceB: 7,
			PrevA: 15, NextA: 10,
			PrevB: 12, NextB: 16,
		}, // 11
		{
			FirstVertexA: 11, FirstVertexB: 3,
			FaceA: 7, FaceB: 8,
			PrevA: 16, NextA: 11,
			PrevB: 13, NextB: 17,
		}, // 12
		{
			FirstVertexA: 9, FirstVertexB: 3,
			FaceA: 8, FaceB: 9,
			PrevA: 17, NextA: 12,
			PrevB: 14, NextB: 18,
		}, // 13
		{
			FirstVertexA: 6, FirstVertexB: 3,
			FaceA: 9, FaceB: 5,
			PrevA: 18, NextA: 13,
			PrevB: 10, NextB: 19,
		}, // 14
		// ring of 5 around the base of cap 2
		{
			FirstVertexA: 1, FirstVertexB: 7,
			FaceA: 6, FaceB: 17,
			PrevA: 10, NextA: 11,
			PrevB: 26, NextB: 25,
		}, // 15
		{
			FirstVertexA: 7, FirstVertexB: 11,
			FaceA: 7, FaceB: 19,
			PrevA: 11, NextA: 12,
			PrevB: 28, NextB: 27,
		}, // 16
		{
			FirstVertexA: 11, FirstVertexB: 9,
			FaceA: 8, FaceB: 11,
			PrevA: 12, NextA: 13,
			PrevB: 20, NextB: 29,
		}, // 17
		{
			FirstVertexA: 9, FirstVertexB: 6,
			FaceA: 9, FaceB: 13,
			PrevA: 13, NextA: 14,
			PrevB: 22, NextB: 21,
		}, // 18
		{
			FirstVertexA: 6, FirstVertexB: 1,
			FaceA: 5, FaceB: 15,
			PrevA: 14, NextA: 10,
			PrevB: 24, NextB: 23,
		}, // 19
		// zig-zag down the middle
		// 10 triangles, 10 new edges
		// starting clockwise from end of edge 0
		{
			FirstVertexA: 2, FirstVertexB: 9,
			FaceA: 11, FaceB: 12,
			PrevA: 29, NextA: 17,
			PrevB: 21, NextB: 5,
		}, // 20
		{
			FirstVertexA: 4, FirstVertexB: 9,
			FaceA: 12, FaceB: 13,
			PrevA: 5, NextA: 20,
			PrevB: 18, NextB: 22,
		}, // 21
		{
			FirstVertexA: 4, FirstVertexB: 6,
			FaceA: 13, FaceB: 14,
			PrevA: 21, NextA: 18,
			PrevB: 23, NextB: 6,
		}, // 22
		{
			FirstVertexA: 8, FirstVertexB: 6,
			FaceA: 14, FaceB: 15,
			PrevA: 6, NextA: 22,
			PrevB: 19, NextB: 24,
		}, // 23
		{
			FirstVertexA: 8, FirstVertexB: 1,
			FaceA: 15, FaceB: 16,
			PrevA: 23, NextA: 19,
			PrevB: 25, NextB: 7,
		}, // 24
		{
			FirstVertexA: 10, FirstVertexB: 1,
			FaceA: 16, FaceB: 17,
			PrevA: 7, NextA: 24,
			PrevB: 15, NextB: 26,
		}, // 25
		{
			FirstVertexA: 10, FirstVertexB: 7,
			FaceA: 17, FaceB: 18,
			PrevA: 25, NextA: 15,
			PrevB: 27, NextB: 8,
		}, // 26
		{
			FirstVertexA: 5, FirstVertexB: 7,
			FaceA: 18, FaceB: 19,
			PrevA: 8, NextA: 26,
			PrevB: 16, NextB: 28,
		}, // 27
		{
			FirstVertexA: 5, FirstVertexB: 11,
			FaceA: 19, FaceB: 10,
			PrevA: 27, NextA: 16,
			PrevB: 29, NextB: 9,
		}, // 28
		{
			FirstVertexA: 2, FirstVertexB: 11,
			FaceA: 10, FaceB: 11,
			PrevA: 9, NextA: 28,
			PrevB: 17, NextB: 20,
		}, // 29
	}
	return icosahedron, nil
}
