package wingedGrid

import (
    "testing"
    "math"
)

// square tolerance for floating point equality
const tolerance = .00000001

func TestBaseIcosahedronEdges(t *testing.T){
    var baseIcosahedron WingedMap
    baseIcosahedron, _ = BaseIcosahedron()
    // Edge length should be 2 for each edge
    for index, edge := range baseIcosahedron.edges {
        var dx, dy, dz float64
        dx = baseIcosahedron.vertices[edge.vertex1].coords[0] - 
                baseIcosahedron.vertices[edge.vertex2].coords[0]
                
        dy = baseIcosahedron.vertices[edge.vertex1].coords[1] - 
                baseIcosahedron.vertices[edge.vertex2].coords[1]
                
        dz = baseIcosahedron.vertices[edge.vertex1].coords[2] - 
                baseIcosahedron.vertices[edge.vertex2].coords[2]
                
        length := math.Sqrt(dx*dx + dy*dy + dz*dz)
        // square of error within tolerance
        if (2 - length) * (2 - length) > tolerance {
            t.Errorf("Edge %d out of tolerance, length is: %f", index, length)
        }
    }
    // each face should have 3 edges
    for index, _ := range baseIcosahedron.faces {
        var count int32 = 0
        for _, edge := range baseIcosahedron.edges {
            if edge.faceA == int32(index) {
                count = count + 1
            }
            if edge.faceB == int32(index) {
                count = count + 1
            }
        }
        if count != 3 {
            t.Errorf("Face %d has %d edges, expected 3.", index, count)
        }
    }
}
func TestBaseIcosahedronVertecies(t *testing.T) {
    var baseIcosahedron WingedMap
    baseIcosahedron, _ = BaseIcosahedron()
    if &baseIcosahedron == nil {} // shut up go compiler
    t.Error("PENDING")
}


func TestBaseIcosahedronFaces(t *testing.T) {
    var baseIcosahedron WingedMap
    baseIcosahedron, _ = BaseIcosahedron()
    
    // face edges should be traversable (ie edge.face(a).next.next.next should be edge)
    for index, face := range baseIcosahedron.faces {
        // pick an edge, make sure after three unique edges, we are back to the first
        firstIndex := face.edges[0]
        firstEdge := baseIcosahedron.edges[firstIndex]
        var nextEdge WingedEdge
        var nextIndex int32
        // get second edge
        if firstEdge.faceA == int32(index) {
            nextIndex = firstEdge.nextA
        } else if firstEdge.faceB == int32(index) {
            nextIndex = firstEdge.nextB
        } else {
            t.Errorf("Edge %d does not point to face %d.",face.edges[0],index)
        }
        // should not be same as first
        if nextIndex == firstIndex {
            t.Error("Second edge should not be same as first.")
        }
        nextEdge = baseIcosahedron.edges[nextIndex]
        // get third edge
        if nextEdge.faceA == int32(index) {
            nextIndex = nextEdge.nextA
        } else if nextEdge.faceB == int32(index) {
            nextIndex = nextEdge.nextB
        } else {
            t.Errorf("Second edge does not point to face %d.", index)
        }
        // should not be same as first
        if nextIndex == firstIndex {
            t.Error("Third edge should not be same as first.")
        }
        nextEdge = baseIcosahedron.edges[nextIndex]
        // get fourth edge (only need index)
        if nextEdge.faceA == int32(index) {
            nextIndex = nextEdge.nextA
        } else if nextEdge.faceB == int32(index) {
            nextIndex = nextEdge.nextB
        } else {
            t.Errorf("Third edge does not point to face %d.", index)
        }
        // forth should be the same as the first
        if nextIndex != firstIndex {
            t.Errorf("Face %d: forth edge does not match the first! Expected a triangle.", index)
            t.Logf("Expected %d == %d", nextIndex, firstIndex)
        }
    }
    // face normal should point away from origin
    //  if edgeQ is clockwise from edgeP, the vectors away from their shared vertex
    //  vertexP and vertexQ should produce a cross product parrallel to the center
    //  of the face (not anti-parrallel)
    for index, face := range baseIcosahedron.faces {
        edgeP := baseIcosahedron.edges[face.edges[0]]
        var edgeQ WingedEdge
        // get next edge
        if edgeP.faceA == int32(index) {
            edgeQ = baseIcosahedron.edges[edgeP.nextA]
        } else if edgeP.faceB == int32(index) {
            edgeQ = baseIcosahedron.edges[edgeP.nextB]
        } else {
            t.Errorf("Edge %d does not point to face %d.",face.edges[0],index)
        }
        
        var vectorP, vectorQ, center [3]float64
        // find the shared vertex and make the two vectors
        //  also grab the center, since we've checked that the triangle is equilateral
        //  the centroid will do
        if edgeP.vertex1 == edgeQ.vertex1 {
            vectorP[0] = baseIcosahedron.vertices[edgeP.vertex2].coords[0] -
                            baseIcosahedron.vertices[edgeP.vertex1].coords[0]
            vectorP[1] = baseIcosahedron.vertices[edgeP.vertex2].coords[1] -
                            baseIcosahedron.vertices[edgeP.vertex1].coords[1]
            vectorP[2] = baseIcosahedron.vertices[edgeP.vertex2].coords[2] -
                            baseIcosahedron.vertices[edgeP.vertex1].coords[2]
                            
            vectorQ[0] = baseIcosahedron.vertices[edgeQ.vertex2].coords[0] -
                            baseIcosahedron.vertices[edgeQ.vertex1].coords[0]
            vectorQ[1] = baseIcosahedron.vertices[edgeQ.vertex2].coords[1] -
                            baseIcosahedron.vertices[edgeQ.vertex1].coords[1]
            vectorQ[2] = baseIcosahedron.vertices[edgeQ.vertex2].coords[2] -
                            baseIcosahedron.vertices[edgeQ.vertex1].coords[2]
            // centroid
            center[0] = (baseIcosahedron.vertices[edgeP.vertex1].coords[0] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[0] +
                           baseIcosahedron.vertices[edgeQ.vertex2].coords[0]) / 3
            center[1] = (baseIcosahedron.vertices[edgeP.vertex1].coords[1] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[1] +
                           baseIcosahedron.vertices[edgeQ.vertex2].coords[1]) / 3
            center[2] = (baseIcosahedron.vertices[edgeP.vertex1].coords[2] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[2] +
                           baseIcosahedron.vertices[edgeQ.vertex2].coords[2]) / 3
        } else if edgeP.vertex1 == edgeQ.vertex2 {
            vectorP[0] = baseIcosahedron.vertices[edgeP.vertex2].coords[0] -
                            baseIcosahedron.vertices[edgeP.vertex1].coords[0]
            vectorP[1] = baseIcosahedron.vertices[edgeP.vertex2].coords[1] -
                            baseIcosahedron.vertices[edgeP.vertex1].coords[1]
            vectorP[2] = baseIcosahedron.vertices[edgeP.vertex2].coords[2] -
                            baseIcosahedron.vertices[edgeP.vertex1].coords[2]
                            
            vectorQ[0] = baseIcosahedron.vertices[edgeQ.vertex1].coords[0] -
                            baseIcosahedron.vertices[edgeQ.vertex2].coords[0]
            vectorQ[1] = baseIcosahedron.vertices[edgeQ.vertex1].coords[1] -
                            baseIcosahedron.vertices[edgeQ.vertex2].coords[1]
            vectorQ[2] = baseIcosahedron.vertices[edgeQ.vertex1].coords[2] -
                            baseIcosahedron.vertices[edgeQ.vertex2].coords[2]
            // centroid
            center[0] = (baseIcosahedron.vertices[edgeP.vertex1].coords[0] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[0] +
                           baseIcosahedron.vertices[edgeQ.vertex1].coords[0]) / 3
            center[1] = (baseIcosahedron.vertices[edgeP.vertex1].coords[1] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[1] +
                           baseIcosahedron.vertices[edgeQ.vertex1].coords[1]) / 3
            center[2] = (baseIcosahedron.vertices[edgeP.vertex1].coords[2] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[2] +
                           baseIcosahedron.vertices[edgeQ.vertex1].coords[2]) / 3
        } else if edgeP.vertex2 == edgeQ.vertex1 {
            vectorP[0] = baseIcosahedron.vertices[edgeP.vertex1].coords[0] -
                            baseIcosahedron.vertices[edgeP.vertex2].coords[0]
            vectorP[1] = baseIcosahedron.vertices[edgeP.vertex1].coords[1] -
                            baseIcosahedron.vertices[edgeP.vertex2].coords[1]
            vectorP[2] = baseIcosahedron.vertices[edgeP.vertex1].coords[2] -
                            baseIcosahedron.vertices[edgeP.vertex2].coords[2]
                            
            vectorQ[0] = baseIcosahedron.vertices[edgeQ.vertex2].coords[0] -
                            baseIcosahedron.vertices[edgeQ.vertex1].coords[0]
            vectorQ[1] = baseIcosahedron.vertices[edgeQ.vertex2].coords[1] -
                            baseIcosahedron.vertices[edgeQ.vertex1].coords[1]
            vectorQ[2] = baseIcosahedron.vertices[edgeQ.vertex2].coords[2] -
                            baseIcosahedron.vertices[edgeQ.vertex1].coords[2]
            // centroid
            center[0] = (baseIcosahedron.vertices[edgeP.vertex1].coords[0] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[0] +
                           baseIcosahedron.vertices[edgeQ.vertex2].coords[0]) / 3
            center[1] = (baseIcosahedron.vertices[edgeP.vertex1].coords[1] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[1] +
                           baseIcosahedron.vertices[edgeQ.vertex2].coords[1]) / 3
            center[2] = (baseIcosahedron.vertices[edgeP.vertex1].coords[2] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[2] +
                           baseIcosahedron.vertices[edgeQ.vertex2].coords[2]) / 3
        } else {
            // edgeP.vertex2 == edgeQ.vertex2
            vectorP[0] = baseIcosahedron.vertices[edgeP.vertex1].coords[0] -
                            baseIcosahedron.vertices[edgeP.vertex2].coords[0]
            vectorP[1] = baseIcosahedron.vertices[edgeP.vertex1].coords[1] -
                            baseIcosahedron.vertices[edgeP.vertex2].coords[1]
            vectorP[2] = baseIcosahedron.vertices[edgeP.vertex1].coords[2] -
                            baseIcosahedron.vertices[edgeP.vertex2].coords[2]
                            
            vectorQ[0] = baseIcosahedron.vertices[edgeQ.vertex1].coords[0] -
                            baseIcosahedron.vertices[edgeQ.vertex2].coords[0]
            vectorQ[1] = baseIcosahedron.vertices[edgeQ.vertex1].coords[1] -
                            baseIcosahedron.vertices[edgeQ.vertex2].coords[1]
            vectorQ[2] = baseIcosahedron.vertices[edgeQ.vertex1].coords[2] -
                            baseIcosahedron.vertices[edgeQ.vertex2].coords[2]
            // centroid
            center[0] = (baseIcosahedron.vertices[edgeP.vertex1].coords[0] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[0] +
                           baseIcosahedron.vertices[edgeQ.vertex1].coords[0]) / 3
            center[1] = (baseIcosahedron.vertices[edgeP.vertex1].coords[1] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[1] +
                           baseIcosahedron.vertices[edgeQ.vertex1].coords[1]) / 3
            center[2] = (baseIcosahedron.vertices[edgeP.vertex1].coords[2] + 
                           baseIcosahedron.vertices[edgeP.vertex2].coords[2] +
                           baseIcosahedron.vertices[edgeQ.vertex1].coords[2]) / 3
        }
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
        //  ie, components should subract to zero, cince unit vectors
        if !( (normal[0] - center[0])*(normal[0] - center[0])>tolerance ||
              (normal[1] - center[1])*(normal[1] - center[1])>tolerance ||
              (normal[2] - center[2])*(normal[2] - center[2])>tolerance   ) {
             
            t.Errorf("center and normal not parellel for face: %d", index)
        }
    }//*/
}