package wingedGrid

import (
    "testing"
)

func TestNextEdgeForFace(t *testing.T) {
    var theEdge WingedEdge
    theEdge.NextA = 1
    theEdge.NextB = 2
    theEdge.FaceA = 3
    theEdge.FaceB = 4
    // should return the correct index for face A and B
    index, err := theEdge.NextEdgeForFace(3)
    if err != nil {
        t.Errorf("Unexpcted error retrieving next edge: %s",err)
    }
    if index != 1 {
        t.Error("Recieved incorrect index.")
    }
    index, err = theEdge.NextEdgeForFace(4)
    if err != nil {
        t.Errorf("Unexpcted error retrieving next edge: %s",err)
    }
    if index != 2 {
        t.Error("Recieved incorrect index.")
    }
    // should return error if given face not associated with face
    _, err = theEdge.NextEdgeForFace(1000)
    if err == nil {
        t.Error("Expected an error for face not associated with edge.")
    }
}

func TestPrevEdgeForFace(t *testing.T) {
    var theEdge WingedEdge
    theEdge.FaceA = 5
    theEdge.FaceB = 6
    theEdge.PrevA = 7
    theEdge.PrevB = 8
    // should return the correct index for face A and B
    index, err := theEdge.PrevEdgeForFace(5)
    if err != nil {
        t.Errorf("Unexpcted error retrieving next edge: %s",err)
    }
    if index != 7 {
        t.Error("Recieved incorrect index.")
    }
    index, err = theEdge.PrevEdgeForFace(6)
    if err != nil {
        t.Errorf("Unexpcted error retrieving next edge: %s",err)
    }
    if index != 8 {
        t.Error("Recieved incorrect index.")
    }
    // should return error if given face not associated with face
    _, err = theEdge.PrevEdgeForFace(1111)
    if err == nil {
        t.Error("Expected an error for face not associated with edge.")
    }
}

func TestFirstVertexForFace(t *testing.T) {
    t.Fail()
}

func TestSecondVertexForFace(t *testing.T) {
    t.Fail()
}