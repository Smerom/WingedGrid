package wingedGrid

func (startGrid WingedGrid) CreateDual() (WingedGrid, error) {
	var dualGrid WingedGrid
	// same number of edges
	dualGrid.Edges = make([]WingedEdge, len(startGrid.Edges))
	// faces and vertices swap
	dualGrid.Faces = make([]WingedFace, len(startGrid.Vertices))
	dualGrid.Vertices = make([]WingedVertex, len(startGrid.Faces))

	// create faces
	for index, startVertex := range startGrid.Vertices {
		dualGrid.Faces[index].Edges = make([]int32, len(startVertex.Edges))
		for otherIndex, edgeIndex := range startVertex.Edges {
			dualGrid.Faces[index].Edges[otherIndex] = edgeIndex
		}
	}

	// create vertices
	for index, startFace := range startGrid.Faces {
		// set edges
		dualGrid.Vertices[index].Edges = make([]int32, len(startFace.Edges))
		for otherIndex, edgeIndex := range startFace.Edges {
			dualGrid.Vertices[index].Edges[otherIndex] = edgeIndex
		}
		// set coords from center of old face
		var faceCenter [3]float64
		var count float64
		var err error
		for _, edgeIndex := range startFace.Edges {
			var vertex WingedVertex
			var vertexIndex int32
			vertexIndex, err = startGrid.Edges[edgeIndex].FirstVertexForFace(int32(index))
			if err != nil {
				return dualGrid, err
			}
			vertex = startGrid.Vertices[vertexIndex]
			faceCenter[0] = faceCenter[0] + vertex.Coords[0]
			faceCenter[1] = faceCenter[1] + vertex.Coords[1]
			faceCenter[2] = faceCenter[2] + vertex.Coords[2]
			count = count + 1
		}
		dualGrid.Vertices[index].Coords[0] = faceCenter[0] / count
		dualGrid.Vertices[index].Coords[1] = faceCenter[1] / count
		dualGrid.Vertices[index].Coords[2] = faceCenter[2] / count
	}

	// set edges
	for index, startEdge := range startGrid.Edges {
		dualGrid.Edges[index].FaceA = startEdge.FirstVertexA
		dualGrid.Edges[index].FaceB = startEdge.FirstVertexB
		// faces are swapped from vertices
		dualGrid.Edges[index].FirstVertexA = startEdge.FaceB
		dualGrid.Edges[index].FirstVertexB = startEdge.FaceA

		// set prev and next
		for faceIndex, face := range dualGrid.Faces {
			for index, edgeIndex := range face.Edges {
				var theEdge WingedEdge = dualGrid.Edges[edgeIndex]
				if theEdge.FaceA == int32(faceIndex) {
					if index == 0 {
						dualGrid.Edges[edgeIndex].PrevA = face.Edges[len(face.Edges)-1]
						dualGrid.Edges[edgeIndex].NextA = face.Edges[1]
					} else if index == len(face.Edges)-1 {
						dualGrid.Edges[edgeIndex].PrevA = face.Edges[len(face.Edges)-2]
						dualGrid.Edges[edgeIndex].NextA = face.Edges[0]
					} else {
						dualGrid.Edges[edgeIndex].PrevA = face.Edges[index-1]
						dualGrid.Edges[edgeIndex].NextA = face.Edges[index+1]
					}
				}
				if theEdge.FaceB == int32(faceIndex) {
					if index == 0 {
						dualGrid.Edges[edgeIndex].PrevB = face.Edges[len(face.Edges)-1]
						dualGrid.Edges[edgeIndex].NextB = face.Edges[1]
					} else if index == len(face.Edges)-1 {
						dualGrid.Edges[edgeIndex].PrevB = face.Edges[len(face.Edges)-2]
						dualGrid.Edges[edgeIndex].NextB = face.Edges[0]
					} else {
						dualGrid.Edges[edgeIndex].PrevB = face.Edges[index-1]
						dualGrid.Edges[edgeIndex].NextB = face.Edges[index+1]
					}
				}
			}
		}
	}

	return dualGrid, nil
}
