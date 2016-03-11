package wingedGrid

import (
    "errors"
    //"fmt"
    "math"
    //"log"
)
// This file contains functions for testing the expected data structure of
// any wingedGrid

// Returns whether the edge order from the vertex index matches that of the edges
// NOT YET IMPLEMENTED
func VertEdgesMatchEdgeOrder(theGrid WingedGrid, vertIndex int32) (bool, error) {
    
    return false, nil
}
// Returns whether the edge has its vertices ordered correctly, 
// such that FirstVertexA it the first vertex encounter when clockwise traversing
// FaceA
func EdgeVertsInCorrectOrientation(theGrid WingedGrid, edgeIndex int32) (bool, error){
    var theEdge WingedEdge = theGrid.Edges[edgeIndex]
    var nextEdge WingedEdge = theGrid.Edges[theEdge.NextA]
    // Second clockwise vertex on theEdge should be the same as the first for 
    // nextEdgeA for faceA, only need to test for one pair.
    // Check with second vertex for next edge incase the next edge is the wrong one
    if theEdge.FirstVertexB == nextEdge.FirstVertexA ||
       theEdge.FirstVertexB == nextEdge.FirstVertexB {
        return true, nil
    }
    return false, errors.New("This Edge has wrong order")
}
// Returns whether the edge order from the face index matches that of the edges
// taken from edge.next
func FaceEdgesMatchesOrderFromEdge(theGrid WingedGrid, faceIndex int32) (bool, error) {
    var theFace WingedFace = theGrid.Faces[faceIndex]
    if len(theFace.Edges) <= 2 {
        return false, errors.New("Face has too few edges.")
    }
    var currentEdge WingedEdge
    currentEdge = theGrid.Edges[theFace.Edges[len(theFace.Edges) - 1]]
    for i := 0; i < len(theFace.Edges); i++ {
        next, err := currentEdge.NextEdgeForFace(faceIndex)
        if err != nil {
            return false, err
        }
        if next == theFace.Edges[i] {
            currentEdge = theGrid.Edges[next]
        } else {
            return false, nil
        }
    }
    return true, nil
}
// Returns whether a the vertices of a face are within square tolerance distance
// from the plane of the first three vertices
// NOT FULLY IMPLEMENTED
func FaceVerticesPlanar(theGrid WingedGrid, faceIndex int32, tolerance float64) (bool) {
    var theFace WingedFace
    theFace = theGrid.Faces[faceIndex]
    // always true for three or fewer points
    if len(theFace.Edges) <= 3 {
        return true
    }
    
    return false
}
// Returns whether a face has an orientation in the same direction as the
// position of its center, within the square tolerance.
// This test is only applicable to WingedGrids that are expected to aproximate a 
// sphere.
func FaceOrientation(theGrid WingedGrid, faceIndex int32, tolerance float64) (bool, error) {
    var err error
    var theFace WingedFace = theGrid.Faces[faceIndex]
    // find center of face
    var vectorP, vectorQ, faceCenter [3]float64
    var count float64
    for _, edgeIndex := range theFace.Edges {
        var vertex WingedVertex
        var vertexIndex int32
        vertexIndex, err = theGrid.Edges[edgeIndex].FirstVertexForFace(faceIndex)
        if err != nil {
            return false, nil
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
    
    // simply use the first two edges, vectors away from their shared vertex
    // from these we get the face normal vector
    var edgeP, edgeQ WingedEdge
    edgeP = theGrid.Edges[theFace.Edges[1]]
    edgeQ = theGrid.Edges[theFace.Edges[0]]
    
    // get and check we find valid vertices
    var startVertexIndex, endVertexIndex int32
    endVertexIndex, err = edgeP.FirstVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }
    startVertexIndex, err = edgeP.SecondVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }
    // get vector along edge P
    vectorP[0] = theGrid.Vertices[endVertexIndex].Coords[0] -
                    theGrid.Vertices[startVertexIndex].Coords[0]
    vectorP[1] = theGrid.Vertices[endVertexIndex].Coords[1] -
                    theGrid.Vertices[startVertexIndex].Coords[1]
    vectorP[2] = theGrid.Vertices[endVertexIndex].Coords[2] -
                    theGrid.Vertices[startVertexIndex].Coords[2]
    
    // get and check we find valid vertices
    endVertexIndex, err = edgeQ.SecondVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }
    startVertexIndex, err = edgeQ.FirstVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }                        
    // get vector along edge Q
    vectorQ[0] = theGrid.Vertices[endVertexIndex].Coords[0] -
                    theGrid.Vertices[startVertexIndex].Coords[0]
    vectorQ[1] = theGrid.Vertices[endVertexIndex].Coords[1] -
                    theGrid.Vertices[startVertexIndex].Coords[1]
    vectorQ[2] = theGrid.Vertices[endVertexIndex].Coords[2] -
                    theGrid.Vertices[startVertexIndex].Coords[2]
    
    // replace with cross product at some point, copying here from original
    // icosahedron test where I simply subtracted normalized components
    var scaleFactor float64
    // normalize the center to unit vector
    scaleFactor = 1/math.Sqrt(faceCenter[0]*faceCenter[0] + 
                              faceCenter[1]*faceCenter[1] +
                              faceCenter[2]*faceCenter[2])
    faceCenter[0] = faceCenter[0] * scaleFactor
    faceCenter[1] = faceCenter[1] * scaleFactor
    faceCenter[2] = faceCenter[2] * scaleFactor
    
    // find normal to face, cross product!
    var faceNormal [3]float64
    faceNormal[0] = vectorP[1]*vectorQ[2] - vectorP[2]*vectorQ[1]
    faceNormal[1] = -(vectorP[0]*vectorQ[2] - vectorP[2]*vectorQ[0])
    faceNormal[2] = vectorP[0]*vectorQ[1] - vectorP[1]*vectorQ[0]
    // normalize it!
    scaleFactor = 1/math.Sqrt(faceNormal[0]*faceNormal[0] + 
                              faceNormal[1]*faceNormal[1] +
                              faceNormal[2]*faceNormal[2])
    faceNormal[0] = faceNormal[0] * scaleFactor
    faceNormal[1] = faceNormal[1] * scaleFactor
    faceNormal[2] = faceNormal[2] * scaleFactor
    //log.Printf("Center: %v Normal: %v",faceCenter,faceNormal)
    // they should be parallel (not antiparrallel!)
    //  ie, components should subract to zero, since unit vectors
    if ( (faceNormal[0] - faceCenter[0])*(faceNormal[0] - faceCenter[0])>tolerance ||
         (faceNormal[1] - faceCenter[1])*(faceNormal[1] - faceCenter[1])>tolerance ||
         (faceNormal[2] - faceCenter[2])*(faceNormal[2] - faceCenter[2])>tolerance   ) {
         
        return false, nil
    }
    return true, nil
}