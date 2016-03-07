package wingedGrid

import (
    "errors"
)

// Surfaces are expected to be orientable
//   orientability allows verticies on edges to be accociated with a particular
//   face such that verticies can be collected for a face without comparison
//   between edges (we don't need to test which vertex two edges share.

// primary information is held in WingedEdge and WingedVertex.Coords
// other information is duplicate for faster traversal of the grid

type WingedFace struct {
    // three for a triangular tiling, but support others
    // index of edge in wingedMap
    Edges []int32 // in clockwise order
}

type WingedEdge struct {
    // index of vertices in winged map
    FirstVertexA, FirstVertexB int32
    // index of faces in winged map
    FaceA, FaceB int32
    // index of edges in winged map for face A and face B
    PrevA, NextA, PrevB, NextB int32
}

type WingedVertex struct {
    Coords [3]float64 // x, y, z
    Edges []int32 // in clockwise order, indexed from Grid
}

// order doesn't matter in the grid
type WingedGrid struct {
    Faces []WingedFace
    Edges []WingedEdge
    Vertices []WingedVertex
}


/******************* Winged Edge ********************/
func (theEdge WingedEdge)NextEdgeForFace(faceIndex int32) (int32, error){
    if theEdge.FaceA == faceIndex {
        return theEdge.NextA, nil
    } else if theEdge.FaceB == faceIndex {
        return theEdge.NextB, nil
    }
    return -1, errors.New("Edge not associated with face.")
}
func (theEdge WingedEdge)PrevEdgeForFace(faceIndex int32) (int32, error){
    if theEdge.FaceA == faceIndex {
        return theEdge.PrevA, nil
    } else if theEdge.FaceB == faceIndex {
        return theEdge.PrevB, nil
    }
    return -1, errors.New("Edge not associated with face.")
}
// returns the first vertex encountered on this edge when traversing the face 
//  clockwise
func (theEdge WingedEdge)FirstVertexForFace(faceIndex int32) (int32, error) {
    if theEdge.FaceA == faceIndex {
        return theEdge.FirstVertexA, nil
    } else if theEdge.FaceB == faceIndex {
        return theEdge.FirstVertexB, nil
    }
    return -1, errors.New("Edge not associated with face.")
}
func (theEdge WingedEdge)SecondVertexForFace(faceIndex int32) (int32, error) {
    if theEdge.FaceA == faceIndex {
        return theEdge.FirstVertexB, nil
    } else if theEdge.FaceB == faceIndex {
        return theEdge.FirstVertexA, nil
    }
    return -1, errors.New("Edge not associated with face.")
}