package wingedGrid

import (
	"math"
	"testing"
)

func TestIcosahedronSubdivisionEdgeLength(t *testing.T) {
	var err error
	var baseIcosahedron WingedGrid
	baseIcosahedron, err = BaseIcosahedron()
	if err != nil {
		t.Fatalf("Failed to create base icosahedron: %s", err)
	}
	var subdividedGrid WingedGrid
	subdividedGrid = baseIcosahedron.SubdivideTriangles(1)
	// tests whether the edges were given plausable vertices, but
	// won't indicate duplicate edges
	// Pick the first edge to compare the rest against
	var edge WingedEdge = subdividedGrid.Edges[0]
	var dx, dy, dz float64
	var expectedLength float64
	dx = subdividedGrid.Vertices[edge.FirstVertexA].Coords[0] -
		subdividedGrid.Vertices[edge.FirstVertexB].Coords[0]

	dy = subdividedGrid.Vertices[edge.FirstVertexA].Coords[1] -
		subdividedGrid.Vertices[edge.FirstVertexB].Coords[1]

	dz = subdividedGrid.Vertices[edge.FirstVertexA].Coords[2] -
		subdividedGrid.Vertices[edge.FirstVertexB].Coords[2]

	expectedLength = math.Sqrt(dx*dx + dy*dy + dz*dz)

	for index, edge := range subdividedGrid.Edges {
		dx = subdividedGrid.Vertices[edge.FirstVertexA].Coords[0] -
			subdividedGrid.Vertices[edge.FirstVertexB].Coords[0]

		dy = subdividedGrid.Vertices[edge.FirstVertexA].Coords[1] -
			subdividedGrid.Vertices[edge.FirstVertexB].Coords[1]

		dz = subdividedGrid.Vertices[edge.FirstVertexA].Coords[2] -
			subdividedGrid.Vertices[edge.FirstVertexB].Coords[2]

		length := math.Sqrt(dx*dx + dy*dy + dz*dz)
		// square of error within tolerance
		if (expectedLength-length)*(expectedLength-length) > tolerance {
			t.Errorf("Edge %d out of tolerance, length is: %f", index, length)
		}
	}
}
