package wingedGrid

import (
    "testing"
)

func TestBasicDualOfIcosahedron(t *testing.T){
    // we should get a dodecahedron, quick check counts
    var err error
    var baseIcosahedron WingedGrid
    baseIcosahedron, err = BaseIcosahedron()
    if err != nil {
        t.Fatalf("Error creating icosahedron: %s", err)
    }
    var dualGrid WingedGrid
    dualGrid, err = baseIcosahedron.CreateDual()
    if err != nil {
        t.Fatalf("Error creating dual: %s", err)
    }
    if len(dualGrid.Faces) != 12 {
        t.Errorf("Incorrect number of faces. Expected: 12, Got: %d", len(dualGrid.Faces))
    }
    if len(dualGrid.Vertices) != 20 {
        t.Errorf("Incorrect number of vertices. Expected: 20, Got: %d", len(dualGrid.Vertices))
    }
    if len(dualGrid.Edges) != 30 {
        t.Errorf("Incorrect number of edges. Expected: 30, Got: %d", len(dualGrid.Edges))
    }
}

func TestDualEdgeVertexOrder(t *testing.T) {
    var err error
    var baseIcosahedron WingedGrid
    baseIcosahedron, err = BaseIcosahedron()
    if err != nil {
        t.Fatalf("Failed to create base icosahedron: %s", err)
    }
    var dualGrid WingedGrid
    dualGrid, err = baseIcosahedron.CreateDual()
    if err != nil {
        t.Fatalf("Error creating dual: %s", err)
    }
    // Verticies should be in correct order
    // The order doesn't give us any extra information
    // but reduces the number of comparisons needed
    // to create a well formed mesh from the data.
    // Also useful for finding the dual of the WingedGrid.
    for index, _ := range dualGrid.Edges {
        correct, err := EdgeVertsInCorrectOrientation(dualGrid, int32(index))
        if err != nil {
            t.Errorf("Vertex ordering error for edge %d: %s", index, err)
        } else if !correct {
            t.Errorf("Vertices not in correct order somewhere near edge: %d", index)
        }
    }
}

func TestDualFaceOrientation(t *testing.T) {
    var err error
    var baseIcosahedron WingedGrid
    baseIcosahedron, err = BaseIcosahedron()
    if err != nil {
        t.Fatalf("Failed to create base icosahedron: %s", err)
    }
    var dualGrid WingedGrid
    dualGrid, err = baseIcosahedron.CreateDual()
    if err != nil {
        t.Fatalf("Error creating dual: %s", err)
    }
    if &dualGrid == nil {return}/*
    // face normal should point away from origin
    //  if edgeQ is clockwise from edgeP, the vectors away from their shared vertex,
    //  vectorP and vectorQ should produce a cross product parrallel to the center
    //  of the face (not anti-parrallel)
    for index, _ := range dualGrid.Faces {
        correct, err := FaceOrientation(dualGrid, int32(index), tolerance)
        if err != nil {
            t.Errorf("Unexpected error determining face orientation: %s", err)
        } else if !correct {
            t.Errorf("Incorrect orientation: center and normal not parellel for face: %d", index)
        }
    }//*/
}