package wingedGrid

import (
    "testing"
    "math"
)
// add connectedness test (all edges on a face can be reached from all other edges)
// square tolerance for floating point equality
const tolerance = .00000001

func TestBaseIcosahedronEdges(t *testing.T){
    var baseIcosahedron WingedMap
    baseIcosahedron, _ = BaseIcosahedron()
    // Edge length should be 2 for each edge
    for index, edge := range baseIcosahedron.Edges {
        var dx, dy, dz float64
        dx = baseIcosahedron.Vertices[edge.Vertex1].Coords[0] - 
                baseIcosahedron.Vertices[edge.Vertex2].Coords[0]
                
        dy = baseIcosahedron.Vertices[edge.Vertex1].Coords[1] - 
                baseIcosahedron.Vertices[edge.Vertex2].Coords[1]
                
        dz = baseIcosahedron.Vertices[edge.Vertex1].Coords[2] - 
                baseIcosahedron.Vertices[edge.Vertex2].Coords[2]
                
        length := math.Sqrt(dx*dx + dy*dy + dz*dz)
        // square of error within tolerance
        if (2 - length) * (2 - length) > tolerance {
            t.Errorf("Edge %d out of tolerance, length is: %f", index, length)
        }
    }
    // each face should have 3 edges
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

func TestBaseIcosahedronVertecies(t *testing.T) {
    var baseIcosahedron WingedMap
    baseIcosahedron, _ = BaseIcosahedron()
    if &baseIcosahedron == nil {} // shut up go compiler
    // each vertex should belong to 5 edges
    var count [12]int
    // loop through each edge
    for _, edge := range baseIcosahedron.Edges {
        count[edge.Vertex1] = count[edge.Vertex1] + 1
        count[edge.Vertex2] = count[edge.Vertex2] + 1
    }
    // check edge count for each vertex
    for index, edgeCount := range count {
        if edgeCount != 5 {
            t.Errorf("Vertex %d belongs to %d edges, expected 5.", index, edgeCount)
        }
    }
}

func TestBaseIcosahedronFaces(t *testing.T) {
    var baseIcosahedron WingedMap
    baseIcosahedron, _ = BaseIcosahedron()
    
    // face edges should be a traversable triangle
    //   (ie edge.face(a).next.next.next should == edge)
    for index, face := range baseIcosahedron.Faces {
        // pick an edge, make sure after three unique edges, we are back to the first
        firstIndex := face.Edges[0]
        firstEdge := baseIcosahedron.Edges[firstIndex]
        var nextEdge WingedEdge
        var nextIndex int32
        // get second edge
        if firstEdge.FaceA == int32(index) {
            nextIndex = firstEdge.NextA
        } else if firstEdge.FaceB == int32(index) {
            nextIndex = firstEdge.NextB
        } else {
            t.Errorf("Edge %d does not point to face %d.",face.Edges[0],index)
        }
        // should not be same as first
        if nextIndex == firstIndex {
            t.Error("Second edge should not be same as first.")
        }
        nextEdge = baseIcosahedron.Edges[nextIndex]
        // get third edge
        if nextEdge.FaceA == int32(index) {
            nextIndex = nextEdge.NextA
        } else if nextEdge.FaceB == int32(index) {
            nextIndex = nextEdge.NextB
        } else {
            t.Errorf("Second edge does not point to face %d.", index)
        }
        // should not be same as first
        if nextIndex == firstIndex {
            t.Error("Third edge should not be same as first.")
        }
        nextEdge = baseIcosahedron.Edges[nextIndex]
        // get fourth edge (only need index)
        if nextEdge.FaceA == int32(index) {
            nextIndex = nextEdge.NextA
        } else if nextEdge.FaceB == int32(index) {
            nextIndex = nextEdge.NextB
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
    //  if edgeQ is clockwise from edgeP, the vectors away from their shared vertex,
    //  vectorP and vectorQ should produce a cross product parrallel to the center
    //  of the face (not anti-parrallel)
    for index, face := range baseIcosahedron.Faces {
        edgeP := baseIcosahedron.Edges[face.Edges[0]]
        var edgeQ WingedEdge
        // get next edge
        if edgeP.FaceA == int32(index) {
            edgeQ = baseIcosahedron.Edges[edgeP.NextA]
        } else if edgeP.FaceB == int32(index) {
            edgeQ = baseIcosahedron.Edges[edgeP.NextB]
        } else {
            t.Errorf("Edge %d does not point to face %d.",face.Edges[0],index)
        }
        
        var vectorP, vectorQ, center [3]float64
        // find the shared vertex and make the two vectors
        //  also grab the center, since we've checked that the triangle is equilateral
        //  the centroid will do
        //  could be put into a function
        if edgeP.Vertex1 == edgeQ.Vertex1 {
            vectorP[0] = baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0] -
                            baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0]
            vectorP[1] = baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1] -
                            baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1]
            vectorP[2] = baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2] -
                            baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2]
                            
            vectorQ[0] = baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[0] -
                            baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[0]
            vectorQ[1] = baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[1] -
                            baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[1]
            vectorQ[2] = baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[2] -
                            baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[2]
            // centroid
            center[0] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0] +
                           baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[0]) / 3
            center[1] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1] +
                           baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[1]) / 3
            center[2] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2] +
                           baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[2]) / 3
        } else if edgeP.Vertex1 == edgeQ.Vertex2 {
            vectorP[0] = baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0] -
                            baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0]
            vectorP[1] = baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1] -
                            baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1]
            vectorP[2] = baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2] -
                            baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2]
                            
            vectorQ[0] = baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[0] -
                            baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[0]
            vectorQ[1] = baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[1] -
                            baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[1]
            vectorQ[2] = baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[2] -
                            baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[2]
            // centroid
            center[0] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0] +
                           baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[0]) / 3
            center[1] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1] +
                           baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[1]) / 3
            center[2] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2] +
                           baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[2]) / 3
        } else if edgeP.Vertex2 == edgeQ.Vertex1 {
            vectorP[0] = baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0] -
                            baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0]
            vectorP[1] = baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1] -
                            baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1]
            vectorP[2] = baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2] -
                            baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2]
                            
            vectorQ[0] = baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[0] -
                            baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[0]
            vectorQ[1] = baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[1] -
                            baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[1]
            vectorQ[2] = baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[2] -
                            baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[2]
            // centroid
            center[0] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0] +
                           baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[0]) / 3
            center[1] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1] +
                           baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[1]) / 3
            center[2] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2] +
                           baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[2]) / 3
        } else {
            // edgeP.vertex2 == edgeQ.vertex2
            vectorP[0] = baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0] -
                            baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0]
            vectorP[1] = baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1] -
                            baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1]
            vectorP[2] = baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2] -
                            baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2]
                            
            vectorQ[0] = baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[0] -
                            baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[0]
            vectorQ[1] = baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[1] -
                            baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[1]
            vectorQ[2] = baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[2] -
                            baseIcosahedron.Vertices[edgeQ.Vertex2].Coords[2]
            // centroid
            center[0] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[0] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[0] +
                           baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[0]) / 3
            center[1] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[1] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[1] +
                           baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[1]) / 3
            center[2] = (baseIcosahedron.Vertices[edgeP.Vertex1].Coords[2] + 
                           baseIcosahedron.Vertices[edgeP.Vertex2].Coords[2] +
                           baseIcosahedron.Vertices[edgeQ.Vertex1].Coords[2]) / 3
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
        //  ie, components should subract to zero, since unit vectors
        if !( (normal[0] - center[0])*(normal[0] - center[0])>tolerance ||
              (normal[1] - center[1])*(normal[1] - center[1])>tolerance ||
              (normal[2] - center[2])*(normal[2] - center[2])>tolerance   ) {
             
            t.Errorf("center and normal not parellel for face: %d", index)
        }
    }
}