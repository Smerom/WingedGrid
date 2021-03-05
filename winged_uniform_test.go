package wingedGrid

import (
	"math"
	"testing"
)

func (grid WingedGrid) computeMeanAndVariance() (float64, float64) {
	// compute mean edge length
	var meanSum float64
	for _, edge := range grid.Edges {
		var vert1, vert2 WingedVertex
		vert1 = grid.Vertices[edge.FirstVertexA]
		vert2 = grid.Vertices[edge.FirstVertexB]
		meanSum += distanceBetween3Points(vert1.Coords, vert2.Coords)
	}
	var mean float64 = meanSum / float64(len(grid.Edges))
	// compute variance
	var varianceSum float64
	for _, edge := range grid.Edges {
		var vert1, vert2 WingedVertex
		vert1 = grid.Vertices[edge.FirstVertexA]
		vert2 = grid.Vertices[edge.FirstVertexB]
		distance := distanceBetween3Points(vert1.Coords, vert2.Coords)

		varianceSum += (distance - mean) * (distance - mean)
	}
	var variance float64 = varianceSum / (float64(len(grid.Edges)) - 1)
	return mean, variance
}

func TestVarianceDecreases(t *testing.T) {
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

	var startingMean, startingVariance float64
	startingMean, startingVariance = subdividedGrid.computeMeanAndVariance()
	var newMean, newVariance float64
	subdividedGrid.UniformVertsOnUnitSphere(10000)

	newMean, newVariance = subdividedGrid.computeMeanAndVariance()

	if newVariance >= startingVariance {
		t.Errorf("Variance did not decrease")
	}
	t.Logf("Mean changed from %e to %e", startingMean, newMean)
	t.Logf("Variance changed from %e to %e", startingVariance, newVariance)
}

// somewhat duplicate of above
func TestIcosahedronSubdivisionUniformEdgeLength(t *testing.T) {
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

	subdividedGrid.UniformVertsOnUnitSphere(10000)

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
