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
	Coords [3]float64 // x, y, z
	Edges  []int32    // in clockwise order, indexed from Grid
}

// WingedGrid is the structure to be sent across a network as the full
// representation of a tiled surface
// order doesn't matter in the grid
type WingedGrid struct {
	Faces    []WingedFace
	Edges    []WingedEdge
	Vertices []WingedVertex
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
