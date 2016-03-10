package wingedGrid

import (
    
)

func (startGrid WingedGrid)CreateDual() (WingedGrid, error) {
    var dualGrid WingedGrid
    // same number of edges
    dualGrid.Edges = make([]WingedEdge, len(startGrid.Edges))
    // faces and vertices swap
    dualGrid.Faces = make([]WingedFace, len(startGrid.Vertices))
    dualGrid.Vertices = make([]WingedVertex, len(startGrid.Faces))
    
    return dualGrid, nil
}