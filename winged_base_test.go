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
    // should return error if given face not associated with edge
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
        t.Errorf("Unexpcted error retrieving prev edge: %s",err)
    }
    if index != 7 {
        t.Error("Recieved incorrect index.")
    }
    index, err = theEdge.PrevEdgeForFace(6)
    if err != nil {
        t.Errorf("Unexpcted error retrieving prev edge: %s",err)
    }
    if index != 8 {
        t.Error("Recieved incorrect index.")
    }
    // should return error if given face not associated with edge
    _, err = theEdge.PrevEdgeForFace(1111)
    if err == nil {
        t.Error("Expected an error for face not associated with edge.")
    }
}

func TestFirstVertexForFace(t *testing.T) {
    var theEdge WingedEdge
    theEdge.FaceA = 9
    theEdge.FaceB = 10
    theEdge.FirstVertexA = 11
    theEdge.FirstVertexB = 12
    // should return the correct index for face A and B
    index, err := theEdge.FirstVertexForFace(9)
    if err != nil {
        t.Errorf("Unexpcted error retrieving first vertex: %s",err)
    }
    if index != 11 {
        t.Error("Recieved incorrect index.")
    }
    index, err = theEdge.FirstVertexForFace(10)
    if err != nil {
        t.Errorf("Unexpcted error retrieving first vertex: %s",err)
    }
    if index != 12 {
        t.Error("Recieved incorrect index.")
    }
    // should return error if given face not associated with edge
    _, err = theEdge.FirstVertexForFace(1141)
    if err == nil {
        t.Error("Expected an error for face not associated with edge.")
    }
}

func TestSecondVertexForFace(t *testing.T) {
    var theEdge WingedEdge
    theEdge.FaceA = 13
    theEdge.FaceB = 14
    theEdge.FirstVertexA = 15
    theEdge.FirstVertexB = 16
    // should return the correct index for face A and B
    index, err := theEdge.SecondVertexForFace(13)
    if err != nil {
        t.Errorf("Unexpcted error retrieving second vertex: %s",err)
    }
    if index != 16 {
        t.Error("Recieved incorrect index.")
    }
    index, err = theEdge.SecondVertexForFace(14)
    if err != nil {
        t.Errorf("Unexpcted error retrieving second vertex: %s",err)
    }
    if index != 15 {
        t.Error("Recieved incorrect index.")
    }
    // should return error if given face not associated with edge
    _, err = theEdge.SecondVertexForFace(1121)
    if err == nil {
        t.Error("Expected an error for face not associated with edge.")
    }
}