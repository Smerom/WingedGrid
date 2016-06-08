package wingedGrid

import (
	"log"
	"math"
)

// assuming triangular tiling of a surface homeomorphic to S2
func (oldGrid WingedGrid) SubdivideTriangles(edgeSubdivisions int32) (WingedGrid, error) {
	return oldGrid.subdivideTriangles_new(edgeSubdivisions)
}
func (oldGrid WingedGrid) subdivideTriangles_old(edgeSubdivisions int32) WingedGrid {
	var newGrid WingedGrid
	var faceCount int32
	if edgeSubdivisions < 1 {
		return newGrid
	} else {
		log.Printf("Dividing edges %d times", edgeSubdivisions)
	}

	// subdividing each edge n times produces a number of faces equal
	//  to the base face count multiplied by (1/2(n+2)(n+1) + 1/2(n+1)(n))
	faceCount = int32(len(oldGrid.Faces)) * ((edgeSubdivisions+2)*(edgeSubdivisions+1)/2 + (edgeSubdivisions+1)*edgeSubdivisions/2)
	newGrid.Faces = make([]WingedFace, faceCount)
	// since each face 'owns' 1/2 of three edges, there are 1.5 times as
	//  many edges as faces
	newGrid.Edges = make([]WingedEdge, 3*faceCount/2)
	// Euler-somebody or other gives us the vertex count of
	//  1/2*(face count) + 2
	newGrid.Vertices = make([]WingedVertex, faceCount/2+2)

	// set all indecies to -1 so we know they are invalid
	var i int32
	for i = 0; i < faceCount; i++ {
		newGrid.Faces[i].Edges = make([]int32, 3)
		newGrid.Faces[i].Edges[0] = -1
		newGrid.Faces[i].Edges[1] = -1
		newGrid.Faces[i].Edges[2] = -1
	}
	for i = 0; i < 3*faceCount/2; i++ {
		newGrid.Edges[i].FaceA = -1
		newGrid.Edges[i].FaceB = -1
		newGrid.Edges[i].FirstVertexA = -1
		newGrid.Edges[i].FirstVertexB = -1
		newGrid.Edges[i].NextA = -1
		newGrid.Edges[i].NextB = -1
		newGrid.Edges[i].PrevA = -1
		newGrid.Edges[i].PrevB = -1
	}
	// verticies corrisponding with old ones will have the same number of
	//  associated edges to preserve the Euler characteristic ( =2 for S2)
	for i = 0; i < int32(len(oldGrid.Vertices)); i++ {
		newGrid.Vertices[i].Edges = make([]int32, len(oldGrid.Vertices[i].Edges))
		for j := 0; j < len(oldGrid.Vertices[i].Edges); j++ {
			newGrid.Vertices[i].Edges[j] = -1
		}
	}
	// the way we divide the faces creates six accociated edges for the remaining
	//  vertecies
	for i = int32(len(oldGrid.Vertices)); i < faceCount/2+2; i++ {
		newGrid.Vertices[i].Edges = make([]int32, 6)
		newGrid.Vertices[i].Edges[0] = -1
		newGrid.Vertices[i].Edges[1] = -1
		newGrid.Vertices[i].Edges[2] = -1
		newGrid.Vertices[i].Edges[3] = -1
		newGrid.Vertices[i].Edges[4] = -1
		newGrid.Vertices[i].Edges[5] = -1
	}

	// divide edges
	oldGrid.subdivideEdges(edgeSubdivisions, newGrid)

	// divide each face
	for i = 0; i < int32(len(oldGrid.Faces)); i++ {
		oldGrid.subdivideFaceAtIndex(i, edgeSubdivisions, newGrid)
	}

	// normalize vertecies
	newGrid.normalizeVerticesToSphere(0)

	return newGrid
}

func (oldGrid WingedGrid) subdivideEdges(edgeSubdivisions int32, newGrid WingedGrid) {
	var origVertexCount int32 = int32(len(oldGrid.Vertices))
	var i, j int32
	for derp, edge := range oldGrid.Edges {
		i = int32(derp)
		var firstVertex WingedVertex = oldGrid.Vertices[edge.FirstVertexA]
		var secondVertex WingedVertex = oldGrid.Vertices[edge.FirstVertexB]
		// angle to subdivide
		var angleToSubdivide float64
		angleToSubdivide = vectorAngle(firstVertex.Coords, secondVertex.Coords)

		// angle between origin, first vertex, and second vertex
		var vectorA, vectorB [3]float64
		vectorA[0] = -1 * firstVertex.Coords[0]
		vectorA[1] = -1 * firstVertex.Coords[1]
		vectorA[2] = -1 * firstVertex.Coords[2]

		vectorB[0] = secondVertex.Coords[0] - firstVertex.Coords[0]
		vectorB[1] = secondVertex.Coords[1] - firstVertex.Coords[1]
		vectorB[2] = secondVertex.Coords[2] - firstVertex.Coords[2]

		var cornerAngle float64
		cornerAngle = vectorAngle(vectorA, vectorB)

		// unit vector from first to second vertex
		var stepDirection [3]float64
		stepDirection[0] = vectorB[0] / vectorLength(vectorB)
		stepDirection[1] = vectorB[1] / vectorLength(vectorB)
		stepDirection[2] = vectorB[2] / vectorLength(vectorB)

		// origional radius of the
		var sphereRadius float64 = vectorLength(firstVertex.Coords)

		// first edge has origional vertex
		newGrid.Edges[i*(edgeSubdivisions+1)].FirstVertexA = edge.FirstVertexA
		newGrid.Edges[i*(edgeSubdivisions+1)].FirstVertexB = origVertexCount + i*edgeSubdivisions

		var divisionLength float64
		divisionLength = math.Sin(angleToSubdivide*(1/float64(edgeSubdivisions))) * sphereRadius / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(1)/float64(edgeSubdivisions))

		newGrid.Vertices[origVertexCount+i*edgeSubdivisions].Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
		newGrid.Vertices[origVertexCount+i*edgeSubdivisions].Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
		newGrid.Vertices[origVertexCount+i*edgeSubdivisions].Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength

		for j = 1; j < edgeSubdivisions; j++ {
			// set the edge vertex indecies
			newGrid.Edges[i*(edgeSubdivisions+1)+j].FirstVertexA = origVertexCount + i*edgeSubdivisions + j - 1
			newGrid.Edges[i*(edgeSubdivisions+1)+j].FirstVertexB = origVertexCount + i*edgeSubdivisions + j

			// find the new vertex position and create the vertex
			// but correct it's length yet
			divisionLength = math.Sin(angleToSubdivide*(float64(j+1)/float64(edgeSubdivisions))) * sphereRadius / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(float64(j+1)/float64(edgeSubdivisions)))

			newGrid.Vertices[origVertexCount+i*edgeSubdivisions+j].Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
			newGrid.Vertices[origVertexCount+i*edgeSubdivisions+j].Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
			newGrid.Vertices[origVertexCount+i*edgeSubdivisions+j].Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength
		}

		// connect last new edge
		newGrid.Edges[i*(edgeSubdivisions+1)+edgeSubdivisions].FirstVertexA = origVertexCount + i*edgeSubdivisions + edgeSubdivisions - 1
		newGrid.Edges[i*(edgeSubdivisions+1)+edgeSubdivisions].FirstVertexB = edge.FirstVertexB
	}
}

func (grid WingedGrid) setFaceEdges(faceIndex int32, edges []int32) {
	var thisEdge, nextEdge, prevEdge WingedEdge
	//log.Printf("With edges %v, of total edges: %d", edges, len(grid.Edges))
	prevEdge = grid.Edges[edges[len(edges)-1]]
	thisEdge = grid.Edges[edges[0]]
	nextEdge = grid.Edges[edges[1]]
	// test verticies
	if thisEdge.FirstVertexA == prevEdge.FirstVertexA || thisEdge.FirstVertexA == prevEdge.FirstVertexB {
		// check the next edge also matches and face A is not set
		if thisEdge.FirstVertexB == nextEdge.FirstVertexA || thisEdge.FirstVertexB == nextEdge.FirstVertexB {
			if thisEdge.FaceA == -1 {
				thisEdge.FaceA = faceIndex
				thisEdge.PrevA = edges[len(edges)-1]
				thisEdge.NextA = edges[1]
			} else {
				log.Printf("Face A has already been set for edge: %v With edge set: %v", thisEdge, edges)
			}

		} else {
			log.Printf("Edges don't match face!")
		}
	} else if thisEdge.FirstVertexB == prevEdge.FirstVertexA || thisEdge.FirstVertexB == prevEdge.FirstVertexB {
		// check the next edge also matches and face B is not set
		if thisEdge.FirstVertexA == nextEdge.FirstVertexA || thisEdge.FirstVertexA == nextEdge.FirstVertexB {
			if thisEdge.FaceB == -1 {
				thisEdge.FaceB = faceIndex
				thisEdge.PrevB = edges[len(edges)-1]
				thisEdge.NextB = edges[1]
			} else {
				log.Printf("Face B has already been set for edge: %v With edge set: %v", thisEdge, edges)
			}
		} else {
			log.Printf("Edges don't match face!")
		}
	} else {
		log.Printf("Edges Don't share a vertex!")
	}

	// loop through middle edges
	var i int32
	for i = 1; i < int32(len(edges))-1; i++ {
		prevEdge = grid.Edges[edges[i-1]]
		thisEdge = grid.Edges[edges[i]]
		nextEdge = grid.Edges[edges[i+1]]
		// test verticies
		if thisEdge.FirstVertexA == prevEdge.FirstVertexA || thisEdge.FirstVertexA == prevEdge.FirstVertexB {
			// check the next edge also matches
			if thisEdge.FirstVertexB == nextEdge.FirstVertexA || thisEdge.FirstVertexB == nextEdge.FirstVertexB {
				if thisEdge.FaceA == -1 {
					thisEdge.FaceA = faceIndex
					thisEdge.PrevA = edges[i-1]
					thisEdge.NextA = edges[i+1]
				} else {
					log.Printf("Face A has already been set for edge: %v With edge set: %v", thisEdge, edges)
				}

			} else {
				log.Printf("Edges don't match face!")
			}
		} else if thisEdge.FirstVertexB == prevEdge.FirstVertexA || thisEdge.FirstVertexB == prevEdge.FirstVertexB {
			// check the next edge also matches
			if thisEdge.FirstVertexA == nextEdge.FirstVertexA || thisEdge.FirstVertexA == nextEdge.FirstVertexB {
				if thisEdge.FaceB == -1 {
					thisEdge.FaceB = faceIndex
					thisEdge.PrevB = edges[i-1]
					thisEdge.NextB = edges[i+1]
				} else {
					log.Printf("Face B has already been set for edge: %v With edge set: %v", thisEdge, edges)
				}

			} else {
				log.Printf("Edges don't match face!")
			}
		} else {
			log.Printf("Edges Don't share a vertex!")
		}
	}

	// last edge
	prevEdge = grid.Edges[edges[len(edges)-2]]
	thisEdge = grid.Edges[edges[len(edges)-1]]
	nextEdge = grid.Edges[edges[0]]
	// test verticies
	if thisEdge.FirstVertexA == prevEdge.FirstVertexA || thisEdge.FirstVertexA == prevEdge.FirstVertexB {
		// check the next edge also matches
		if thisEdge.FirstVertexB == nextEdge.FirstVertexA || thisEdge.FirstVertexB == nextEdge.FirstVertexB {
			if thisEdge.FaceA == -1 {
				thisEdge.FaceA = faceIndex
				thisEdge.PrevA = edges[len(edges)-2]
				thisEdge.NextA = edges[0]
			} else {
				log.Printf("Face A has already been set for edge: %v With edge set: %v", thisEdge, edges)
			}
		} else {
			log.Printf("Edges don't match face!")
		}
	} else if thisEdge.FirstVertexB == prevEdge.FirstVertexA || thisEdge.FirstVertexB == prevEdge.FirstVertexB {
		// check the next edge also matches
		if thisEdge.FirstVertexA == nextEdge.FirstVertexA || thisEdge.FirstVertexA == nextEdge.FirstVertexB {
			if thisEdge.FaceB == -1 {
				thisEdge.FaceB = faceIndex
				thisEdge.PrevB = edges[len(edges)-2]
				thisEdge.NextB = edges[0]
			} else {
				log.Printf("Face B has already been set for edge: %v With edge set: %v", thisEdge, edges)
			}
		} else {
			log.Printf("Edges don't match face!")
		}
	} else {
		log.Printf("Edges Don't share a vertex!")
	}

}

func (oldGrid WingedGrid) subdivideFaceAtIndex(faceIndex, edgeSubdivisions int32, newGrid WingedGrid) {
	log.Printf("deviding face: %d", faceIndex)
	// get the number of faces we divide this one into
	var subFaceCount int32
	subFaceCount = ((edgeSubdivisions+2)*(edgeSubdivisions+1)/2 + (edgeSubdivisions+1)*edgeSubdivisions/2)

	// number of internal edges created on this face
	//  equal to ( (edgeSubdivisions + 1) choose 2 ) * 3
	var subEdgeCount int32
	subEdgeCount = 3 * edgeSubdivisions * (edgeSubdivisions + 1) / 2

	// number of internal vertecies created on this face
	//  equal to ( edgeSubdivisions choose 2 )
	var subVertexCount int32
	if edgeSubdivisions > 1 {
		subVertexCount = edgeSubdivisions * (edgeSubdivisions - 1) / 2
	} else {
		subVertexCount = 0
	}

	var vertexOffset int32 = int32(len(oldGrid.Vertices)) + int32(len(oldGrid.Edges))*edgeSubdivisions + subVertexCount*faceIndex

	var bigFace WingedFace = oldGrid.Faces[faceIndex]

	/****** VERTS *******/
	// only if we have more than one division
	var i, j int32
	if edgeSubdivisions > 1 {
		for i = 0; i < edgeSubdivisions-1; i++ {
			var firstVertex WingedVertex = newGrid.Vertices[oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-2-i, edgeSubdivisions)]
			var secondVertex WingedVertex = newGrid.Vertices[oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, 1+i, edgeSubdivisions)]
			// angle to subdivide
			var angleToSubdivide float64
			angleToSubdivide = vectorAngle(firstVertex.Coords, secondVertex.Coords)

			// angle between origin, first vertex, and second vertex
			var vectorA, vectorB [3]float64
			vectorA[0] = -1 * firstVertex.Coords[0]
			vectorA[1] = -1 * firstVertex.Coords[1]
			vectorA[2] = -1 * firstVertex.Coords[2]

			vectorB[0] = secondVertex.Coords[0] - firstVertex.Coords[0]
			vectorB[1] = secondVertex.Coords[1] - firstVertex.Coords[1]
			vectorB[2] = secondVertex.Coords[2] - firstVertex.Coords[2]

			var cornerAngle float64
			cornerAngle = vectorAngle(vectorA, vectorB)

			// unit vector from first to second vertex
			var stepDirection [3]float64
			stepDirection[0] = vectorB[0] / vectorLength(vectorB)
			stepDirection[1] = vectorB[1] / vectorLength(vectorB)
			stepDirection[2] = vectorB[2] / vectorLength(vectorB)

			var firstVectorLength float64 = vectorLength(firstVertex.Coords)
			var divisionLength float64 = 0
			for j = 0; j < i+1; j++ {
				divisionLength = math.Sin(angleToSubdivide*(float64(j+1)/float64(i+1))) * firstVectorLength / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(float64(j+1)/float64(i+1)))

				newGrid.Vertices[vertexOffset+((i+1)*(i+2)/2)+j].Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
				newGrid.Vertices[vertexOffset+((i+1)*(i+2)/2)+j].Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
				newGrid.Vertices[vertexOffset+((i+1)*(i+2)/2)+j].Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength
			}
		}
	}

	/****** EDGES *******/
	var edgeOffset int32 = (edgeSubdivisions+1)*int32(len(oldGrid.Edges)) + 3*(edgeSubdivisions)*(edgeSubdivisions+1)/2*faceIndex
	if edgeSubdivisions == 1 {
		newGrid.Edges[edgeOffset+0].FirstVertexA = int32(len(oldGrid.Vertices)) + bigFace.Edges[0]
		newGrid.Edges[edgeOffset+0].FirstVertexB = int32(len(oldGrid.Vertices)) + bigFace.Edges[1]

		newGrid.Edges[edgeOffset+1].FirstVertexA = int32(len(oldGrid.Vertices)) + bigFace.Edges[0]
		newGrid.Edges[edgeOffset+1].FirstVertexB = int32(len(oldGrid.Vertices)) + bigFace.Edges[2]

		newGrid.Edges[edgeOffset+2].FirstVertexA = int32(len(oldGrid.Vertices)) + bigFace.Edges[1]
		newGrid.Edges[edgeOffset+2].FirstVertexB = int32(len(oldGrid.Vertices)) + bigFace.Edges[2]
	} else {
		// first row
		newGrid.Edges[edgeOffset+0].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1, edgeSubdivisions)
		newGrid.Edges[edgeOffset+0].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)

		newGrid.Edges[edgeOffset+1].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1, edgeSubdivisions)
		newGrid.Edges[edgeOffset+1].FirstVertexB = vertexOffset + 0

		newGrid.Edges[edgeOffset+2].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
		newGrid.Edges[edgeOffset+2].FirstVertexB = vertexOffset + 0

		// middle rows
		var i, j int32
		var rowOffset int32
		for i = 1; i < edgeSubdivisions-1; i++ {
			rowOffset = i * (i + 1) * 3 / 2
			// first border
			newGrid.Edges[edgeOffset+rowOffset+0].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1-i, edgeSubdivisions)
			newGrid.Edges[edgeOffset+rowOffset+0].FirstVertexB = vertexOffset + (i * (i + 1) / 2)

			newGrid.Edges[edgeOffset+rowOffset+1].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1-i, edgeSubdivisions)
			newGrid.Edges[edgeOffset+rowOffset+1].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + i

			newGrid.Edges[edgeOffset+rowOffset+2].FirstVertexA = vertexOffset + (i * (i + 1) / 2)
			newGrid.Edges[edgeOffset+rowOffset+2].FirstVertexB = vertexOffset + (1 / 2 * i * (i + 1)) + i
			// interior of face
			for j = 1; j < i; j++ {
				newGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + j - 1
				newGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + j

				newGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + j - 1
				newGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + i + j

				newGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + j
				newGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + i + j
			}

			// second border
			newGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + i - 1
			newGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)

			newGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + i - 1
			newGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + 2*i

			newGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)
			newGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + 2*i
		}

		// last row
		i = edgeSubdivisions - 1 // should already be set, but just incase
		rowOffset = i * (i + 1) * 3 / 2
		// border 1
		newGrid.Edges[edgeOffset+rowOffset+0].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
		newGrid.Edges[edgeOffset+rowOffset+0].FirstVertexB = vertexOffset + (i * (i + 1) / 2)

		newGrid.Edges[edgeOffset+rowOffset+1].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
		newGrid.Edges[edgeOffset+rowOffset+1].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-1, edgeSubdivisions)

		newGrid.Edges[edgeOffset+rowOffset+2].FirstVertexA = vertexOffset + (i * (i + 1) / 2)
		newGrid.Edges[edgeOffset+rowOffset+2].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-1, edgeSubdivisions)

		// middle
		for j = 1; j < i; j++ {
			newGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + j - 1
			newGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexB = vertexOffset + (i * (i + 1) / 2) + j

			newGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + j - 1
			newGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-j-1, edgeSubdivisions)

			newGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + j
			newGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-j-1, edgeSubdivisions)
		}

		// border 2
		newGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + i - 1
		newGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)

		newGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexA = vertexOffset + (i * (i + 1) / 2) + i - 1
		newGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, 0, edgeSubdivisions)

		newGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)
		newGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, 0, edgeSubdivisions)
	}

	/******* FACES *******/
	// each new face set starts at
	var indexStart int32 = subFaceCount * faceIndex

	// first corner
	// set edges in face
	var faceEdges []int32
	faceEdges = newGrid.Faces[indexStart+0].Edges
	faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions, edgeSubdivisions)
	faceEdges[1] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
	faceEdges[2] = edgeOffset + faceIndex*subEdgeCount
	newGrid.setFaceEdges(indexStart, faceEdges)

	// loop through the middle section of faces
	// edges grow by 1/2*i*(i-1)*3
	for i = 1; i < edgeSubdivisions; i++ {
		var rowIndexStart int32 = i*(i+1)/2 + i*(i-1)/2

		// edge
		faceEdges = newGrid.Faces[indexStart+rowIndexStart+0].Edges
		faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-i, edgeSubdivisions)
		faceEdges[1] = edgeOffset + i*(i-1)*3/2 + 1
		faceEdges[2] = edgeOffset + i*(i+1)*3/2
		newGrid.setFaceEdges(indexStart+rowIndexStart+0, faceEdges)

		// middle
		// up to one less than next index start
		for j = 1; j < (i+1)*(i+2)/2+1/2*(i+1)*i/2-rowIndexStart-1; j++ {
			if j%2 == 1 {
				faceEdges = newGrid.Faces[indexStart+rowIndexStart+j].Edges
				faceEdges[0] = edgeOffset + i*(i-1)*3/2 + (j-1)/2
				faceEdges[1] = edgeOffset + i*(i-1)*3/2 + (j-1)/2 + 2
				faceEdges[2] = edgeOffset + i*(i-1)*3/2 + (j-1)/2 + 1
				newGrid.setFaceEdges(indexStart+rowIndexStart+j, faceEdges)
			} else {
				faceEdges = newGrid.Faces[indexStart+rowIndexStart+j].Edges
				faceEdges[0] = edgeOffset + i*(i+1)*3/2 + j/2
				faceEdges[1] = edgeOffset + i*(i-1)*3/2 + (j-2)/2 + 2
				faceEdges[2] = edgeOffset + i*(i-1)*3/2 + j/2 + 1
				newGrid.setFaceEdges(indexStart+rowIndexStart+j, faceEdges)
			}
		}

		// edge
		faceEdges = newGrid.Faces[indexStart+(i+1)*(i+2)/2+(i+1)*i/2-1].Edges
		faceEdges[0] = edgeOffset + i*(i+1)*3/2 + ((i+1)*(i+2)/2+(i+1)*i/2-rowIndexStart-1)/2
		faceEdges[1] = edgeOffset + i*(i-1)*3/2 + ((i+1)*(i+2)/2+(i+1)*i/2-rowIndexStart-3)/2 + 2
		faceEdges[2] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)
		newGrid.setFaceEdges(indexStart+(i+1)*(i+2)/2+(i+1)*i/2-1, faceEdges)
	}

	var rowIndexStart int32 = edgeSubdivisions*(edgeSubdivisions+1)/2 + edgeSubdivisions*(edgeSubdivisions-1)/2
	// bottom corner 1
	faceEdges = newGrid.Faces[indexStart+rowIndexStart+0].Edges
	faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
	faceEdges[1] = edgeOffset + i*(edgeSubdivisions-1)*3/2 + 1
	faceEdges[2] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions, edgeSubdivisions)
	newGrid.setFaceEdges(indexStart+rowIndexStart+0, faceEdges)
	// bottom edge
	for j = 1; j < subFaceCount-rowIndexStart-1; j++ {
		if j%2 == 1 {
			faceEdges = newGrid.Faces[indexStart+rowIndexStart+j].Edges
			faceEdges[0] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-1)/2
			faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-1)/2 + 2
			faceEdges[2] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-1)/2 + 1
			newGrid.setFaceEdges(indexStart+rowIndexStart+j, faceEdges)
		} else {
			faceEdges = newGrid.Faces[indexStart+rowIndexStart+j].Edges
			faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-j/2, edgeSubdivisions)
			faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-2)/2 + 2
			faceEdges[2] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + j/2 + 1
			newGrid.setFaceEdges(indexStart+rowIndexStart+j, faceEdges)
		}
	}
	// bottom corner 2
	faceEdges = newGrid.Faces[indexStart+subFaceCount-1].Edges
	faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 2, 0, edgeSubdivisions)
	faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (subFaceCount-rowIndexStart-3)/2 + 2
	faceEdges[2] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 1, edgeSubdivisions, edgeSubdivisions)
	newGrid.setFaceEdges(indexStart+subFaceCount-1, faceEdges)
}
