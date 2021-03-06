package wingedGrid

import (
	"math"
	"testing"
)

// NEED:
// add connectedness test (all edges on a face can be reached from all other edges)

// square tolerance for floating point equality
const tolerance = .00000001

func TestBaseIcosahedronEdgeLength(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// tests whether the edges were given plausable vertices, but
	// won't indicate duplicate edges
	// Edge length should be 2 for each edge
	for index, edge := range baseIcosahedron.Edges {
		var dx, dy, dz float64
		dx = baseIcosahedron.Vertices[edge.FirstVertexA].Coords[0] -
			baseIcosahedron.Vertices[edge.FirstVertexB].Coords[0]

		dy = baseIcosahedron.Vertices[edge.FirstVertexA].Coords[1] -
			baseIcosahedron.Vertices[edge.FirstVertexB].Coords[1]

		dz = baseIcosahedron.Vertices[edge.FirstVertexA].Coords[2] -
			baseIcosahedron.Vertices[edge.FirstVertexB].Coords[2]

		length := math.Sqrt(dx*dx + dy*dy + dz*dz)
		// square of error within tolerance
		if (2-length)*(2-length) > tolerance {
			t.Errorf("Edge %d out of tolerance, length is: %f", index, length)
		}
	}
}
func TestBaseIcosahedronEdgesPerFace(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// Expecting triangles, each face should have 3 edges
	// For each face with loop through all edges and count
	// the associations.
	for index, _ := range baseIcosahedron.Faces {
		var count int32 = 0
		for _, edge := range baseIcosahedron.Edges {
			if edge.FaceA == int32(index) {
				count = count + 1
			}
			if edge.FaceB == int32(index) {
				count = count + 1
			}
		}
		if count != 3 {
			t.Errorf("Face %d has %d edges, expected 3.", index, count)
		}
	}
}

func TestBaseIcosahedronEdgeVertexOrder(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// Verticies should be in correct order
	// The order doesn't give us any extra information
	// but reduces the number of comparisons needed
	// to create a well formed mesh from the data.
	// Also useful for finding the dual of the WingedGrid.
	for index, _ := range baseIcosahedron.Edges {
		correct, err := EdgeVertsInCorrectOrientation(baseIcosahedron, int32(index))
		if err != nil {
			t.Errorf("Vertex ordering error for edge %d: %s", index, err)
		} else if !correct {
			t.Errorf("Vertices not in correct order somewhere near edge: %d", index)
		}
	}
}

func TestBaseIcosahedronVertecies(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// Each vertex should belong to 5 edges
	// Along with the edges per face and edge length,
	// this test determines whether we hooked up the edges
	// correctly from the 3 golden retangles.
	var count [12]int
	// loop through each edge
	for _, edge := range baseIcosahedron.Edges {
		// add one to the count for each vertex referenced by the edge
		count[edge.FirstVertexA] = count[edge.FirstVertexA] + 1
		count[edge.FirstVertexB] = count[edge.FirstVertexB] + 1
	}
	// check edge count for each vertex
	for index, edgeCount := range count {
		if edgeCount != 5 {
			t.Errorf("Vertex %d belongs to %d edges, expected 5.", index, edgeCount)
		}
	}
}

func TestBaseIcosahedronFaceTraversability(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// face edges should be a traversable triangle
	// (ie edge.face(a).next.next.next should == edge)
	for index, face := range baseIcosahedron.Faces {
		// pick an edge, make sure after three unique edges, we are back to the first
		firstEdgeIndex := face.Edges[0]
		firstEdge := baseIcosahedron.Edges[firstEdgeIndex]
		var nextEdge WingedEdge
		var nextEdgeIndex int32
		// get second edge
		if firstEdge.FaceA == int32(index) {
			nextEdgeIndex = firstEdge.NextA
		} else if firstEdge.FaceB == int32(index) {
			nextEdgeIndex = firstEdge.NextB
		} else {
			t.Errorf("Edge %d does not point to face %d.", face.Edges[0], index)
		}
		// should not be same as first
		if nextEdgeIndex == firstEdgeIndex {
			t.Error("Second edge should not be same as first.")
		}
		nextEdge = baseIcosahedron.Edges[nextEdgeIndex]
		// get third edge
		if nextEdge.FaceA == int32(index) {
			nextEdgeIndex = nextEdge.NextA
		} else if nextEdge.FaceB == int32(index) {
			nextEdgeIndex = nextEdge.NextB
		} else {
			t.Errorf("Second edge does not point to face %d.", index)
		}
		// should not be same as first
		if nextEdgeIndex == firstEdgeIndex {
			t.Error("Third edge should not be same as first.")
		}
		nextEdge = baseIcosahedron.Edges[nextEdgeIndex]
		// get fourth edge (only need index)
		if nextEdge.FaceA == int32(index) {
			nextEdgeIndex = nextEdge.NextA
		} else if nextEdge.FaceB == int32(index) {
			nextEdgeIndex = nextEdge.NextB
		} else {
			t.Errorf("Third edge does not point to face %d.", index)
		}
		// forth should be the same as the first
		if nextEdgeIndex != firstEdgeIndex {
			t.Errorf("Face %d: forth edge does not match the first! Expected a triangle.", index)
			t.Logf("Expected %d == %d", nextEdgeIndex, firstEdgeIndex)
		}
	}
}

func TestBaseIcosahedronFaceOrientation(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// face normal should point away from origin
	//  if edgeQ is clockwise from edgeP, the vectors away from their shared vertex,
	//  vectorP and vectorQ should produce a cross product parrallel to the center
	//  of the face (not anti-parrallel)
	for index, _ := range baseIcosahedron.Faces {
		correct, err := FaceOrientation(baseIcosahedron, int32(index), tolerance)
		if err != nil {
			t.Errorf("Unexpected error determining face orientation: %s", err)
		} else if !correct {
			t.Errorf("Incorrect orientation: center and normal not parellel for face: %d", index)
		}
	}
}

func TestBaseIcosahedronEdgeConsistency(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	for index, _ := range baseIcosahedron.Edges {
		isConsistent, err := EdgePrevNextConsistent(baseIcosahedron, int32(index))
		if err != nil {
			t.Errorf("Edge %d fails consistency with error: %s", index, err)
		} else if !isConsistent {
			t.Errorf("Edge %d is inconsistent", index)
		}
	}
}

func TestBaseIcosahedronFaceEdgeOrdering(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	// theFace.Edges should have the same order as given by
	// edge.NextEdgeForFace(theFace)
	for index, _ := range baseIcosahedron.Faces {
		correct, err := FaceEdgesMatchesOrderFromEdge(baseIcosahedron, int32(index))
		if err != nil {
			t.Errorf("Error testing order: %s", err)
		}
		if !correct {
			t.Errorf("Face %d does not have edges in correct order.", index)
		}
	}
}

func TestBaseIcosahedronVertexEdgeOrdering(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}

	var vertexIndex int32
	for dummy, _ := range baseIcosahedron.Vertices {
		vertexIndex = int32(dummy)
		matches, err := VertEdgesMatchEdgeOrder(baseIcosahedron, vertexIndex)
		if err != nil {
			t.Errorf("Vertex %d edge order fails with error: %s", vertexIndex, err)
		} else if !matches {
			t.Errorf("Vertex %d edge order does not match.", vertexIndex)
		}
	}
}
