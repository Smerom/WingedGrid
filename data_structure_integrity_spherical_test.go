package wingedGrid

import (
    "errors"
    "math"
)
// for winged grids expected to be spherical with tiles facing away from the origin

// Returns whether the edge order from the face index matches that of the edges
//  taken from edge.next
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
    return false, nil
}
// Returns whether a the vertices of a face are within square tolerance distance
//  from the plane of the first three vertices
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
//  position of its center, within the square tolerance
func FaceOrientation(theGrid WingedGrid, faceIndex int32, tolerance float64) (bool, error) {
    var err error
    var theFace WingedFace 
    theFace = theGrid.Faces[faceIndex]
    // find center of face
    var vectorP, vectorQ, center [3]float64
    var count float64
    for _, edgeIndex := range theFace.Edges {
        var vertex WingedVertex
        var vertexIndex int32
        vertexIndex, err = theGrid.Edges[edgeIndex].FirstVertexForFace(faceIndex)
        if err != nil {
            return false, nil
        }
        vertex = theGrid.Vertices[vertexIndex]
        center[0] = center[0] + vertex.Coords[0]
        center[1] = center[1] + vertex.Coords[1]
        center[2] = center[2] + vertex.Coords[2]
        count = count + 1
    }
    center[0] = center[0] / count
    center[1] = center[1] / count
    center[1] = center[1] / count
    
    // simply use the first two edges, vectors away from their shared vertex
    var edgeP, edgeQ WingedEdge
    edgeP = theGrid.Edges[theFace.Edges[0]]
    edgeQ = theGrid.Edges[theFace.Edges[1]]
    var startVertexIndex, endVertexIndex int32
    endVertexIndex, err = edgeP.FirstVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }
    startVertexIndex, err = edgeP.SecondVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }
    vectorP[0] = theGrid.Vertices[endVertexIndex].Coords[0] -
                    theGrid.Vertices[startVertexIndex].Coords[0]
    vectorP[1] = theGrid.Vertices[endVertexIndex].Coords[1] -
                    theGrid.Vertices[startVertexIndex].Coords[1]
    vectorP[2] = theGrid.Vertices[endVertexIndex].Coords[2] -
                    theGrid.Vertices[startVertexIndex].Coords[2]
    
    endVertexIndex, err = edgeQ.SecondVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }
    startVertexIndex, err = edgeQ.FirstVertexForFace(faceIndex)
    if err != nil {
        return false, err
    }                        
    vectorQ[0] = theGrid.Vertices[endVertexIndex].Coords[0] -
                    theGrid.Vertices[startVertexIndex].Coords[0]
    vectorQ[1] = theGrid.Vertices[endVertexIndex].Coords[1] -
                    theGrid.Vertices[startVertexIndex].Coords[1]
    vectorQ[2] = theGrid.Vertices[endVertexIndex].Coords[2] -
                    theGrid.Vertices[startVertexIndex].Coords[2]
                    
    // replace with cross product at some point, copying here from original
    //   icosahedron test
    var scaleFactor float64
    // normalize the center to unit vector
    scaleFactor = 1/math.Sqrt(center[0]*center[0] + 
                              center[1]*center[1] +
                              center[2]*center[2])
    
    center[0] = center[0] * scaleFactor
    center[1] = center[1] * scaleFactor
    center[2] = center[2] * scaleFactor
    // find normal to face, cross product!
    var normal [3]float64
    normal[0] = vectorP[1]*vectorQ[2] - vectorP[2]*vectorQ[1]
    normal[1] = -(vectorP[0]*vectorQ[2] - vectorP[2]*vectorQ[0])
    normal[2] = vectorP[0]*vectorQ[1] - vectorP[1]*vectorQ[0]
    // normalize it!
    scaleFactor = 1/math.Sqrt(normal[0]*normal[0] + 
                              normal[1]*normal[1] +
                              normal[2]*normal[2])
                            
    normal[0] = normal[0] * scaleFactor
    normal[1] = normal[1] * scaleFactor
    normal[2] = normal[2] * scaleFactor
    // they should be parallel (not antiparrallel!)
    //  ie, components should subract to zero, since unit vectors
    if !( (normal[0] - center[0])*(normal[0] - center[0])>tolerance ||
          (normal[1] - center[1])*(normal[1] - center[1])>tolerance ||
          (normal[2] - center[2])*(normal[2] - center[2])>tolerance   ) {
         
        return false, nil
    }
    
    return true, nil
}