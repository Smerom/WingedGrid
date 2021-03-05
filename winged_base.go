package wingedGrid

import (
	"errors"
)

// Contains the basic winged edge data structures and functions for traversing
// the faces, edges, and vertices.
//
// Surfaces are expected to be orientable
// orientability allows verticies on edges to be accociated with a particular
// face such that verticies can be collected for a face without comparison
// between edges (we don't need to test which vertex two edges share.)
//
// Primary information is held in WingedEdge and WingedVertex.Coords
// other information is duplicate for faster traversal of the grid
//
// INDEXES ARE USED INSTEAD OF POINTERS IN ALL STRUCTURES AS THESE STRUCTURES WILL
// BE SENT OVER A NETWORK AS THEIR PRIMARY USE

// represents a face of a tiled surface
type WingedFace struct {
	// three for a triangular tiling, but support others
	// index of edge in wingedGrid
	Edges []int32 // in clockwise order
}

// represents an edge of a tiled surface, internal or boundary
type WingedEdge struct {
	// index of vertices in winged grid
	FirstVertexA, FirstVertexB int32
	// index of faces in winged grid
	FaceA, FaceB int32
	// index of edges in winged grid for face A and face B
	PrevA, NextA, PrevB, NextB int32
}

// represents a vertex point, currently as an embedding in 3-space
type WingedVertex struct {
	Coords          [3]float64 // x, y, z
	Edges           []int32    // in clockwise order, indexed from Grid
	vertexNeighbors []int32
}

// WingedGrid is the structure to be sent across a network as the full
// representation of a tiled surface
// order doesn't matter in the grid
type WingedGrid struct {
	Faces    []WingedFace
	Edges    []WingedEdge
	Vertices []WingedVertex
}

/******************* Winged Face ********************/
// returns faces adjacent to this face
// does not check bounds
func (theGrid WingedGrid) NeighborsForFace(faceIndex int32) ([]int32, error) {
	var err error
	var theFace WingedFace
	theFace = theGrid.Faces[faceIndex]
	// one neighbor for each face
	var neighbors []int32 = make([]int32, len(theFace.Edges))
	for index, edgeIndex := range theFace.Edges {
		var theEdge WingedEdge
		theEdge = theGrid.Edges[edgeIndex]
		var neighborIndex int32
		neighborIndex, err = theEdge.AdjacentForFace(faceIndex)
		// don't check error, pass the last one on to the next function
		neighbors[index] = neighborIndex
	}

	return neighbors, err
}

func (theGrid WingedGrid) FaceCenter(faceIndex int32) ([3]float64, error) {
	var startFace WingedFace = theGrid.Faces[faceIndex]
	var faceCenter [3]float64
	var count float64
	var err error
	for _, edgeIndex := range startFace.Edges {
		var vertex WingedVertex
		var vertexIndex int32
		vertexIndex, err = theGrid.Edges[edgeIndex].FirstVertexForFace(faceIndex)
		if err != nil {
			return faceCenter, err
		}
		vertex = theGrid.Vertices[vertexIndex]
		faceCenter[0] = faceCenter[0] + vertex.Coords[0]
		faceCenter[1] = faceCenter[1] + vertex.Coords[1]
		faceCenter[2] = faceCenter[2] + vertex.Coords[2]
		count = count + 1
	}
	faceCenter[0] = faceCenter[0] / count
	faceCenter[1] = faceCenter[1] / count
	faceCenter[2] = faceCenter[2] / count

	return faceCenter, nil
}

/******************* Winged Edge ********************/
// Returns the index of the next clockwise edge for the given face index, or
// an error if the edge is not associated with the face.
func (theEdge WingedEdge) NextEdgeForFace(faceIndex int32) (int32, error) {
	if theEdge.FaceA == faceIndex {
		return theEdge.NextA, nil
	} else if theEdge.FaceB == faceIndex {
		return theEdge.NextB, nil
	}
	return -1, errors.New("Edge not associated with face.")
}

// Returns the index of the previous clockwise edge for the given face index, or
// an error if the edge is not associated with the face.
func (theEdge WingedEdge) PrevEdgeForFace(faceIndex int32) (int32, error) {
	if theEdge.FaceA == faceIndex {
		return theEdge.PrevA, nil
	} else if theEdge.FaceB == faceIndex {
		return theEdge.PrevB, nil
	}
	return -1, errors.New("Edge not associated with face.")
}

// Returns the index for the first vertex encountered on this edge when traversing
// the face clockwise, or an error if the edge is not associated with the face
func (theEdge WingedEdge) FirstVertexForFace(faceIndex int32) (int32, error) {
	if theEdge.FaceA == faceIndex {
		return theEdge.FirstVertexA, nil
	} else if theEdge.FaceB == faceIndex {
		return theEdge.FirstVertexB, nil
	}
	return -1, errors.New("Edge not associated with face.")
}

// Returns the index for the second vertex encountered on this edge when traversing
// the face clockwise, or an error if the edge is not associated with the face
func (theEdge WingedEdge) SecondVertexForFace(faceIndex int32) (int32, error) {
	if theEdge.FaceA == faceIndex {
		return theEdge.FirstVertexB, nil
	} else if theEdge.FaceB == faceIndex {
		return theEdge.FirstVertexA, nil
	}
	return -1, errors.New("Edge not associated with face.")
}

func (theEdge WingedEdge) NextEdgeForVertex(vertexIndex int32) (int32, error) {
	if theEdge.FirstVertexA == vertexIndex {
		return theEdge.PrevEdgeForFace(theEdge.FaceA)
	} else if theEdge.FirstVertexB == vertexIndex {
		return theEdge.PrevEdgeForFace(theEdge.FaceB)
	}
	return -1, errors.New("Edge not associated with Vertex.")
}

func (theEdge WingedEdge) PrevEdgeForVertex(vertexIndex int32) (int32, error) {
	if theEdge.FirstVertexA == vertexIndex {
		return theEdge.NextEdgeForFace(theEdge.FaceB)
	} else if theEdge.FirstVertexB == vertexIndex {
		return theEdge.NextEdgeForFace(theEdge.FaceA)
	}
	return -1, errors.New("Edge not associated with Vertex.")
}

// returns the other face associated with an edge, or an error
// if the edge is not associated with the face.
func (theEdge WingedEdge) AdjacentForFace(faceIndex int32) (int32, error) {
	if theEdge.FaceA == faceIndex {
		return theEdge.FaceB, nil
	} else if theEdge.FaceB == faceIndex {
		return theEdge.FaceA, nil
	}
	return -1, errors.New("Edge not associated with face.")
}

// returns the other vertex associated with an edge, or an error
// if the edge is not associated with the vertex.
func (theEdge WingedEdge) AdjacentForVertex(vertexIndex int32) (int32, error) {
	if theEdge.FirstVertexA == vertexIndex {
		return theEdge.FirstVertexB, nil
	} else if theEdge.FirstVertexB == vertexIndex {
		return theEdge.FirstVertexA, nil
	}
	return -1, errors.New("Edge not associated with vertex.")
}

/******************* Winged Vertex ********************/
// returns the vertex indices adjacent to this vertex
func (theGrid WingedGrid) NeighborsForVertex(vertexIndex int32) ([]int32, error) {
	var err error
	if vertexIndex > int32(len(theGrid.Vertices))-1 {
		return nil, errors.New("Index out of bounds.")
	}
	// check if the neighbors have already been found
	if len(theGrid.Vertices[vertexIndex].vertexNeighbors) == 0 {
		theGrid.Vertices[vertexIndex].vertexNeighbors = make([]int32, len(theGrid.Vertices[vertexIndex].Edges))
		// one neighbor for each face
		for index, edgeIndex := range theGrid.Vertices[vertexIndex].Edges {
			var theEdge WingedEdge
			theEdge = theGrid.Edges[edgeIndex]
			var neighborIndex int32
			neighborIndex, err = theEdge.AdjacentForVertex(vertexIndex)
			// don't check error, pass the last one on to the next function
			theGrid.Vertices[vertexIndex].vertexNeighbors[index] = neighborIndex
		}
	}

	return theGrid.Vertices[vertexIndex].vertexNeighbors, err
}
