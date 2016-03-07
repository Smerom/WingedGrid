package wingedGrid

import (
    "testing"
    "math"
)
// add connectedness test (all edges on a face can be reached from all other edges)
// square tolerance for floating point equality
const tolerance = .00000001

func TestBaseIcosahedronEdges(t *testing.T){
    var baseIcosahedron WingedGrid
    baseIcosahedron, _ = BaseIcosahedron()
    // Edge length should be 2 for each edge
    for index, edge := range baseIcosahedron.Edges {
        var dx, dy, dz float64
        dx = baseIcosahedron.Vertices[edge.FirstVertexA].Coords[0] - 
                baseIcosahedron.Vertices[edge.FirstVertexB].Coords[0]
                
        dy = baseIcosahedron.Vertices[edge.FirstVertexA].Coords[1] - 
                baseIcosahedron.Vertices[edge.FirstVertexB].Coords[1]
                
        dz = baseIcosahedron.Vertices[edge.FirstVertexA].Coords[2] - 
                baseIcosahedron.Vertices[edge.FirstVertexB].Coords[2]
                
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
    
    // verticies should be in correct order
    for index, _ := range baseIcosahedron.Edges {
        correct, err := EdgeVertsInCorrectOrientation(baseIcosahedron, int32(index))
        if err != nil {
            t.Errorf("Vertex ordering error for edge %d: %s", index, err)
        } else if !correct {
            t.Errorf("Vertices not in correct order somewhere near edge: %d", index)
        }
    }
}

func TestBaseIcosahedronVertecies(t *testing.T) {
    var baseIcosahedron WingedGrid
    baseIcosahedron, _ = BaseIcosahedron()
    if &baseIcosahedron == nil {} // shut up go compiler
    // each vertex should belong to 5 edges
    var count [12]int
    // loop through each edge
    for _, edge := range baseIcosahedron.Edges {
        count[edge.FirstVertexA] = count[edge.FirstVertexA] + 1
        count[edge.FirstVertexB] = count[edge.FirstVertexB] + 1
    }
    // check edge count for each vertex
    for index, edgeCount := range count {
        if edgeCount != 5 {
            t.Errorf("Vertex %d belongs to %d edges, expected 5.", index, edgeCount)
        }
    }
}

func TestBaseIcosahedronFaces(t *testing.T) {
    var baseIcosahedron WingedGrid
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
    for index, _ := range baseIcosahedron.Faces {
        correct, err := FaceOrientation(baseIcosahedron, int32(index), tolerance)
        if err != nil {
            t.Errorf("Unexpected error determining face orientation: %s", err)
        } else if !correct {
            t.Errorf("Incorrect orientation: center and normal not parellel for face: %d", index)
        }
    }
}

func TestBaseIcosahedronFaceEdgeOrdering(t *testing.T) {
    var baseIcosahedron WingedGrid
    baseIcosahedron, _ = BaseIcosahedron()
    
    for index, _ := range baseIcosahedron.Faces {
        correct, err := FaceEdgesMatchesOrderFromEdge(baseIcosahedron, int32(index))
        if err != nil {
            t.Errorf("Error testing order: %s", err)
        }
        if !correct {
            t.Errorf("Face %d does not have edges in correct order.", index)
        }
    }
}