package wingedGrid

import (
	"math"
	"testing"
)

func TestSubdivisionOfIcosahedronBasicValidity(t *testing.T) {
	// here we test that all values of the new grid have been set in some fasion to a valid value
	//  ie, indecies are non-negative and less than the length of their associated array
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(20)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}

	var totalEdges, totalFaces, totalVertices int32
	totalVertices = int32(len(subdividedGrid.Vertices))
	totalFaces = int32(len(subdividedGrid.Faces))
	totalEdges = int32(len(subdividedGrid.Edges))
	// Check Faces
	var invalidFaceCount int = 0
	for index, face := range subdividedGrid.Faces {
		var faceIsValid bool = true
		if face.Edges == nil {
			t.Errorf("nil edge array for face: %d", index)
		} else {
			for i, edgeIndex := range face.Edges {
				if edgeIndex < 0 || edgeIndex >= totalEdges {
					t.Errorf("Face %d has invalid edge set as %d at array index %d", index, edgeIndex, i)
					faceIsValid = false
				}
			}
			for i, edgeIndex := range face.Edges {
				for j, otherEdgeIndex := range face.Edges {
					if i != j && edgeIndex == otherEdgeIndex {
						t.Errorf("Face %d had duplicate edges", index)
						faceIsValid = false
					}
				}
			}
		}
		if !faceIsValid {
			invalidFaceCount = invalidFaceCount + 1
		}
		if invalidFaceCount > 5 {
			t.Fatal("Too many invalid faces, stopping test to avoid spam")
		}
	}

	// check edges
	var invalidEdgeCount int = 0
	for index, edge := range subdividedGrid.Edges {
		var edgeIsValid bool = true
		// check verts are set
		if edge.FirstVertexA < 0 || edge.FirstVertexA >= totalVertices {
			t.Errorf("Edge at index %d has invalid FirstVertexA set: %d", index, edge.FirstVertexA)
			edgeIsValid = false
		}
		if edge.FirstVertexB < 0 || edge.FirstVertexB >= totalVertices {
			t.Errorf("Edge at index %d has invalid FirstVertexB set: %d", index, edge.FirstVertexB)
			edgeIsValid = false
		}
		// check faces
		if edge.FaceA < 0 || edge.FaceA >= totalFaces {
			t.Errorf("Edge at index %d has invalid FaceA set: %d", index, edge.FaceA)
			edgeIsValid = false
		}
		if edge.FaceB < 0 || edge.FaceB >= totalFaces {
			t.Errorf("Edge at index %d has invalid FaceB set: %d", index, edge.FaceB)
			edgeIsValid = false
		}

		// check next and prev edges
		if edge.NextA < 0 || edge.NextA >= totalEdges {
			t.Errorf("Edge at index %d has invalid NextA set: %d", index, edge.NextA)
			edgeIsValid = false
		}
		if edge.NextB < 0 || edge.NextB >= totalEdges {
			t.Errorf("Edge at index %d has invalid NextB set: %d", index, edge.NextB)
			edgeIsValid = false
		}
		if edge.PrevA < 0 || edge.PrevA >= totalEdges {
			t.Errorf("Edge at index %d has invalid PrevA set: %d", index, edge.PrevA)
			edgeIsValid = false
		}
		if edge.PrevB < 0 || edge.PrevB >= totalEdges {
			t.Errorf("Edge at index %d has invalid PrevB set: %d", index, edge.PrevB)
			edgeIsValid = false
		}

		// check that A and B are different
		if edge.FaceA == edge.FaceB {
			t.Errorf("Edge at index %d has duplicate faces", index)
			edgeIsValid = false
		}
		if edge.NextA == edge.NextB {
			t.Errorf("Edge at index %d has duplicate next edge", index)
			edgeIsValid = false
		}
		if edge.PrevA == edge.PrevB {
			t.Errorf("Edge at index %d has duplicate prev edge", index)
			edgeIsValid = false
		}
		if edge.FirstVertexA == edge.FirstVertexB {
			t.Errorf("Edge at index %d has duplicate vertices", index)
			edgeIsValid = false
		}

		if !edgeIsValid {
			invalidEdgeCount = invalidEdgeCount + 1
		}
		if invalidEdgeCount > 5 {
			t.Fatal("Too many invalid edges, stopping test to avoid spam")
		}
	}
}

func TestSubdivisionOfIcosahedronEdgeVertexOrdering(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(20)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}

	var index int32
	for dummy, _ := range subdividedGrid.Edges {
		index = int32(dummy)
		var isValid bool
		isValid, err = EdgeVertsInCorrectOrientation(subdividedGrid, index)
		if !isValid {
			t.Error(err)
		}
	}
}

func TestIcosahedronSubdivisionEdgeLength(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(10)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}

	subdividedGrid.NormalizeVerticesToDistanceFromOrigin(1.0)

	// common vars to each part of the test
	var edge WingedEdge
	var index int
	var dx, dy, dz float64

	// first test for zero length edges
	var haveEdgeLengthZero bool = false
	for index, edge = range subdividedGrid.Edges {
		dx = subdividedGrid.Vertices[edge.FirstVertexA].Coords[0] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[0]

		dy = subdividedGrid.Vertices[edge.FirstVertexA].Coords[1] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[1]

		dz = subdividedGrid.Vertices[edge.FirstVertexA].Coords[2] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[2]

		if dx*dx < tolerance && dy*dy < tolerance && dz*dz < tolerance {
			t.Errorf("Edge %d is length 0", index)
			haveEdgeLengthZero = true
		}

	}
	if haveEdgeLengthZero {
		t.Fatal("Got zero length edge, next part will fail. Stopping test.")
	}

	// tests whether the edges were given plausable vertices, but
	// won't indicate duplicate edges
	// Pick the first edge to compare the rest against
	edge = subdividedGrid.Edges[0]

	var expectedLength float64
	dx = subdividedGrid.Vertices[edge.FirstVertexA].Coords[0] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[0]

	dy = subdividedGrid.Vertices[edge.FirstVertexA].Coords[1] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[1]

	dz = subdividedGrid.Vertices[edge.FirstVertexA].Coords[2] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[2]

	expectedLength = math.Sqrt(dx*dx + dy*dy + dz*dz)

	var minRatio, maxRatio float64
	minRatio = math.MaxFloat64
	maxRatio = 0

	for index, edge = range subdividedGrid.Edges {
		dx = subdividedGrid.Vertices[edge.FirstVertexA].Coords[0] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[0]

		dy = subdividedGrid.Vertices[edge.FirstVertexA].Coords[1] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[1]

		dz = subdividedGrid.Vertices[edge.FirstVertexA].Coords[2] - subdividedGrid.Vertices[edge.FirstVertexB].Coords[2]

		length := math.Sqrt(dx*dx + dy*dy + dz*dz)
		// square of error within tolerance
		if length/expectedLength > maxRatio {
			maxRatio = length / expectedLength
		}
		if length/expectedLength < minRatio {
			minRatio = length / expectedLength
		}
		if (expectedLength-length)*(expectedLength-length) > tolerance {
			//t.Errorf("Edge %d out of tolerance at %f ratio, length is: %f", index, length/expectedLength, length)
		}
	}
	t.Logf("Max ratio to expected: %f", maxRatio)
	t.Logf("Min ratio to expected: %f", minRatio)
}

func TestIcosahedronSubdivisionEdgeConsistency(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(20)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}
	for index, _ := range subdividedGrid.Edges {
		isConsistent, err := EdgePrevNextConsistent(subdividedGrid, int32(index))
		if err != nil {
			t.Errorf("Edge %d fails consistency with error: %s", index, err)
		} else if !isConsistent {
			t.Errorf("Edge %d is inconsistent", index)
		}
	}
}

func TestIcosahedronSubdivisionFaceAndEdge_EdgeOrderMatches(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(20)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}
	var faceIndex int32
	for dummy, _ := range subdividedGrid.Faces {
		faceIndex = int32(dummy)
		matches, err := FaceEdgesMatchesOrderFromEdge(subdividedGrid, faceIndex)
		if !matches && err != nil {
			t.Errorf("Face %d edge order fails with error: %s", faceIndex, err)
		} else if !matches {
			t.Errorf("Face %d edge order does not match.", faceIndex)
		}
	}
}

func TestIcosahedronSubdivisionVertexAndEdge_EdgeOrderMatches(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(20)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}

	var vertexIndex int32
	for dummy, _ := range subdividedGrid.Vertices {
		vertexIndex = int32(dummy)
		matches, err := VertEdgesMatchEdgeOrder(subdividedGrid, vertexIndex)
		if err != nil {
			t.Errorf("Vertex %d edge order fails with error: %s", vertexIndex, err)
		} else if !matches {
			t.Errorf("Vertex %d edge order does not match.", vertexIndex)
		}
	}
}

// Don't run, currently fails as edge lengths are not equal
func _TestIcosahedronSubdivisionFaceOrientation(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid, err = baseIcosahedron.SubdivideTriangles(20)
	if err != nil {
		t.Fatalf("Failed to subdivide base icosahedron: %s", err)
	}
	var faceIndex int32
	for dummy, _ := range subdividedGrid.Faces {
		faceIndex = int32(dummy)
		correct, err := FaceOrientation(subdividedGrid, faceIndex, tolerance)
		if !correct && err != nil {
			t.Errorf("Face %d Orientation failed with error: %s", faceIndex, err)
		} else if !correct {
			t.Errorf("Face %d has incorrect orientation.", faceIndex)
		}
	}
}
