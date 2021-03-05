package wingedGrid

import (
	"errors"
	"log"
	"math"
)

// assuming triangular tiling of a surface homeomorphic to S2
func (oldGrid WingedGrid) SubdivideTriangles(edgeSubdivisions int32) (WingedGrid, error) {
	var err error
	var dividedGrid WingedGrid
	if edgeSubdivisions < 1 {
		err = errors.New("Invalid number of subdivisions")
		return dividedGrid, err
	}

	var faceCount int32
	// subdividing each edge n times produces a number of faces equal
	//  to the base face count multiplied by (1/2(n+2)(n+1) + 1/2(n+1)(n))
	faceCount = int32(len(oldGrid.Faces)) * ((edgeSubdivisions+2)*(edgeSubdivisions+1)/2 + (edgeSubdivisions+1)*edgeSubdivisions/2)
	dividedGrid.Faces = make([]WingedFace, faceCount)
	// since each face 'owns' 1/2 of three edges, there are 1.5 times as
	//  many edges as faces
	dividedGrid.Edges = make([]WingedEdge, 3*faceCount/2)
	// Euler-somebody or other gives us the vertex count of
	//  1/2*(face count) + 2
	dividedGrid.Vertices = make([]WingedVertex, faceCount/2+2)

	// Invalidate all values
	var i int32
	for i = 0; i < faceCount; i++ {
		dividedGrid.Faces[i].Edges = make([]int32, 3)
		dividedGrid.Faces[i].Edges[0] = -1
		dividedGrid.Faces[i].Edges[1] = -1
		dividedGrid.Faces[i].Edges[2] = -1
	}
	for i = 0; i < 3*faceCount/2; i++ {
		dividedGrid.Edges[i].FaceA = -1
		dividedGrid.Edges[i].FaceB = -1
		dividedGrid.Edges[i].FirstVertexA = -1
		dividedGrid.Edges[i].FirstVertexB = -1
		dividedGrid.Edges[i].NextA = -1
		dividedGrid.Edges[i].NextB = -1
		dividedGrid.Edges[i].PrevA = -1
		dividedGrid.Edges[i].PrevB = -1
	}
	// verticies corrisponding with old ones will have the same number of
	//  associated edges to preserve the Euler characteristic ( =2 for S2)
	for i = 0; i < int32(len(oldGrid.Vertices)); i++ {
		dividedGrid.Vertices[i].Edges = make([]int32, len(oldGrid.Vertices[i].Edges))
		for j := 0; j < len(oldGrid.Vertices[i].Edges); j++ {
			dividedGrid.Vertices[i].Edges[j] = -1
		}
		// set the coords
		dividedGrid.Vertices[i].Coords[0] = math.MaxInt32
		dividedGrid.Vertices[i].Coords[1] = math.MaxInt32
		dividedGrid.Vertices[i].Coords[2] = math.MaxInt32
	}
	// the way we divide the faces creates six accociated edges for the remaining
	//  vertecies
	for i = int32(len(oldGrid.Vertices)); i < faceCount/2+2; i++ {
		dividedGrid.Vertices[i].Edges = make([]int32, 6)
		dividedGrid.Vertices[i].Edges[0] = -1
		dividedGrid.Vertices[i].Edges[1] = -1
		dividedGrid.Vertices[i].Edges[2] = -1
		dividedGrid.Vertices[i].Edges[3] = -1
		dividedGrid.Vertices[i].Edges[4] = -1
		dividedGrid.Vertices[i].Edges[5] = -1
		// set the coords
		dividedGrid.Vertices[i].Coords[0] = math.MaxInt32
		dividedGrid.Vertices[i].Coords[1] = math.MaxInt32
		dividedGrid.Vertices[i].Coords[2] = math.MaxInt32
	}

	/***************** Subdivide the grid ****************/

	// create the edges
	oldGrid.setSubdivisionEdgeVertices(edgeSubdivisions, dividedGrid)

	// create the faces, update the edges
	oldGrid.setSubdivisionFaceEdges(edgeSubdivisions, dividedGrid)

	// create the verticies
	oldGrid.subdivideVertices(edgeSubdivisions, dividedGrid)

	// set vertex edge array
	dividedGrid.setEdgesForVerticesIfInvalid()

	return dividedGrid, err
}

/******************* EDGE SUBDIVISION ***********************/

func (oldGrid WingedGrid) setSubdivisionEdgeVertices(edgeSubdivisions int32, dividedGrid WingedGrid) error {

	/****** Old Edge Subdivision goes in the first section of the new array, ordered by edge *******/
	var origVertexCount int32 = int32(len(oldGrid.Vertices))
	var i, j int32
	for derp, edge := range oldGrid.Edges {
		i = int32(derp)

		// first edge has origional vertex
		dividedGrid.Edges[i*(edgeSubdivisions+1)].FirstVertexA = edge.FirstVertexA
		dividedGrid.Edges[i*(edgeSubdivisions+1)].FirstVertexB = origVertexCount + i*edgeSubdivisions

		for j = 1; j < edgeSubdivisions; j++ {
			// set the edge vertex indecies
			dividedGrid.Edges[i*(edgeSubdivisions+1)+j].FirstVertexA = origVertexCount + i*edgeSubdivisions + j - 1
			dividedGrid.Edges[i*(edgeSubdivisions+1)+j].FirstVertexB = origVertexCount + i*edgeSubdivisions + j
		}

		// connect last new edge
		dividedGrid.Edges[i*(edgeSubdivisions+1)+edgeSubdivisions].FirstVertexA = origVertexCount + i*edgeSubdivisions + edgeSubdivisions - 1
		dividedGrid.Edges[i*(edgeSubdivisions+1)+edgeSubdivisions].FirstVertexB = edge.FirstVertexB
	}

	/********* Edges created interior to old faces go in the second section, ordered by face. ****/
	var faceIndex int32
	for i, oldFace := range oldGrid.Faces {
		faceIndex = int32(i)
		// vertex offset for vertices interior to the face
		var vertexOffset int32 = int32(len(oldGrid.Vertices)) + int32(len(oldGrid.Edges))*edgeSubdivisions + (edgeSubdivisions*(edgeSubdivisions-1)/2)*faceIndex

		var edgeOffset int32 = (edgeSubdivisions+1)*int32(len(oldGrid.Edges)) + 3*(edgeSubdivisions)*(edgeSubdivisions+1)/2*faceIndex
		if edgeSubdivisions == 1 {
			dividedGrid.Edges[edgeOffset+0].FirstVertexA = int32(len(oldGrid.Vertices)) + oldFace.Edges[0]
			dividedGrid.Edges[edgeOffset+0].FirstVertexB = int32(len(oldGrid.Vertices)) + oldFace.Edges[1]

			dividedGrid.Edges[edgeOffset+1].FirstVertexA = int32(len(oldGrid.Vertices)) + oldFace.Edges[0]
			dividedGrid.Edges[edgeOffset+1].FirstVertexB = int32(len(oldGrid.Vertices)) + oldFace.Edges[2]

			dividedGrid.Edges[edgeOffset+2].FirstVertexA = int32(len(oldGrid.Vertices)) + oldFace.Edges[1]
			dividedGrid.Edges[edgeOffset+2].FirstVertexB = int32(len(oldGrid.Vertices)) + oldFace.Edges[2]
		} else {
			// first row
			dividedGrid.Edges[edgeOffset+0].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1, edgeSubdivisions)
			dividedGrid.Edges[edgeOffset+0].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, 0, edgeSubdivisions)

			dividedGrid.Edges[edgeOffset+1].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1, edgeSubdivisions)
			dividedGrid.Edges[edgeOffset+1].FirstVertexB = vertexOffset + 0

			dividedGrid.Edges[edgeOffset+2].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, 0, edgeSubdivisions)
			dividedGrid.Edges[edgeOffset+2].FirstVertexB = vertexOffset + 0

			// middle rows
			var i, j int32
			var rowOffset int32
			for i = 1; i < edgeSubdivisions-1; i++ {
				rowOffset = i * (i + 1) * 3 / 2
				// first border
				dividedGrid.Edges[edgeOffset+rowOffset+0].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1-i, edgeSubdivisions)
				dividedGrid.Edges[edgeOffset+rowOffset+0].FirstVertexB = vertexOffset + (i * (i - 1) / 2)

				dividedGrid.Edges[edgeOffset+rowOffset+1].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-1-i, edgeSubdivisions)
				dividedGrid.Edges[edgeOffset+rowOffset+1].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + i

				dividedGrid.Edges[edgeOffset+rowOffset+2].FirstVertexA = vertexOffset + (i * (i - 1) / 2)
				dividedGrid.Edges[edgeOffset+rowOffset+2].FirstVertexB = vertexOffset + (i*(i-1))/2 + i
				// interior of face
				for j = 1; j < i; j++ {
					dividedGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + j - 1
					dividedGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + j

					dividedGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + j - 1
					dividedGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + i + j

					dividedGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + j
					dividedGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + i + j
				}

				// second border
				dividedGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + i - 1
				dividedGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)

				dividedGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + i - 1
				dividedGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + 2*i

				dividedGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)
				dividedGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + 2*i
			}

			// last row
			i = edgeSubdivisions - 1 // should already be set, but just incase
			rowOffset = i * (i + 1) * 3 / 2
			// border 1
			dividedGrid.Edges[edgeOffset+rowOffset+0].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
			dividedGrid.Edges[edgeOffset+rowOffset+0].FirstVertexB = vertexOffset + (i * (i - 1) / 2)

			dividedGrid.Edges[edgeOffset+rowOffset+1].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
			dividedGrid.Edges[edgeOffset+rowOffset+1].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-1, edgeSubdivisions)

			dividedGrid.Edges[edgeOffset+rowOffset+2].FirstVertexA = vertexOffset + (i * (i - 1) / 2)
			dividedGrid.Edges[edgeOffset+rowOffset+2].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-1, edgeSubdivisions)

			// middle
			for j = 1; j < i; j++ {
				dividedGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + j - 1
				dividedGrid.Edges[edgeOffset+rowOffset+j*3+0].FirstVertexB = vertexOffset + (i * (i - 1) / 2) + j

				dividedGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + j - 1
				dividedGrid.Edges[edgeOffset+rowOffset+j*3+1].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-j-1, edgeSubdivisions)

				dividedGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + j
				dividedGrid.Edges[edgeOffset+rowOffset+j*3+2].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-j-1, edgeSubdivisions)
			}

			// border 2
			dividedGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + i - 1
			dividedGrid.Edges[edgeOffset+rowOffset+i*3+0].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)

			dividedGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexA = vertexOffset + (i * (i - 1) / 2) + i - 1
			dividedGrid.Edges[edgeOffset+rowOffset+i*3+1].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, 0, edgeSubdivisions)

			dividedGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexA = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)
			dividedGrid.Edges[edgeOffset+rowOffset+i*3+2].FirstVertexB = oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 2, 0, edgeSubdivisions)
		}
	}

	return nil
}

/******************* FACE SUBDIVISION ***********************/
func (oldGrid WingedGrid) setSubdivisionFaceEdges(edgeSubdivisions int32, dividedGrid WingedGrid) error {
	var faceIndex int32
	var i, j int32
	for dummy, _ := range oldGrid.Faces {
		faceIndex = int32(dummy)

		// get the number of faces we divide this one into
		var subFaceCount int32
		subFaceCount = (edgeSubdivisions+2)*(edgeSubdivisions+1)/2 + (edgeSubdivisions+1)*edgeSubdivisions/2

		// number of internal edges created on this face
		//  equal to ( (edgeSubdivisions + 1) choose 2 ) * 3
		var subEdgeCount int32
		subEdgeCount = 3 * edgeSubdivisions * (edgeSubdivisions + 1) / 2

		// each new face set starts at
		var indexStart int32 = subFaceCount * faceIndex

		var edgeOffset int32 = (edgeSubdivisions+1)*int32(len(oldGrid.Edges)) + subEdgeCount*faceIndex

		// first corner
		// set edges in face
		var faceEdges []int32
		faceEdges = dividedGrid.Faces[indexStart+0].Edges
		faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions, edgeSubdivisions)
		faceEdges[1] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 1, 0, edgeSubdivisions)
		faceEdges[2] = edgeOffset

		// loop through the middle section of faces
		// edges grow by 1/2*i*(i-1)*3
		for i = 1; i < edgeSubdivisions; i++ {
			var rowIndexStart int32 = i*(i+1)/2 + i*(i-1)/2

			// edge
			faceEdges = dividedGrid.Faces[indexStart+rowIndexStart+0].Edges
			faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-i, edgeSubdivisions)
			faceEdges[1] = edgeOffset + i*(i-1)*3/2 + 1
			faceEdges[2] = edgeOffset + i*(i+1)*3/2

			// middle
			// up to one less than next index start
			for j = 1; j < (i+1)*(i+2)/2+(i+1)*i/2-rowIndexStart-1; j++ {
				if j%2 == 1 {
					faceEdges = dividedGrid.Faces[indexStart+rowIndexStart+j].Edges
					faceEdges[0] = edgeOffset + i*(i-1)*3/2 + (j-1)*3/2
					faceEdges[1] = edgeOffset + i*(i-1)*3/2 + (j-1)*3/2 + 2
					faceEdges[2] = edgeOffset + i*(i-1)*3/2 + (j-1)*3/2 + 1
				} else {
					faceEdges = dividedGrid.Faces[indexStart+rowIndexStart+j].Edges
					faceEdges[0] = edgeOffset + i*(i+1)*3/2 + j*3/2
					faceEdges[1] = edgeOffset + i*(i-1)*3/2 + (j-2)*3/2 + 2
					faceEdges[2] = edgeOffset + i*(i-1)*3/2 + j*3/2 + 1
				}
			}

			// edge
			faceEdges = dividedGrid.Faces[indexStart+(i+1)*(i+2)/2+(i+1)*i/2-1].Edges
			faceEdges[0] = edgeOffset + i*(i+1)*3/2 + ((i+1)*(i+2)/2+(i+1)*i/2-rowIndexStart-1)*3/2
			faceEdges[1] = edgeOffset + i*(i-1)*3/2 + ((i+1)*(i+2)/2+(i+1)*i/2-rowIndexStart-3)*3/2 + 2
			faceEdges[2] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 1, i, edgeSubdivisions)
		}

		var rowIndexStart int32 = edgeSubdivisions*(edgeSubdivisions+1)/2 + edgeSubdivisions*(edgeSubdivisions-1)/2
		// bottom corner 1
		faceEdges = dividedGrid.Faces[indexStart+rowIndexStart+0].Edges
		faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 0, 0, edgeSubdivisions)
		faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + 1
		faceEdges[2] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions, edgeSubdivisions)
		// bottom edge
		for j = 1; j < subFaceCount-rowIndexStart-1; j++ {
			if j%2 == 1 {
				faceEdges = dividedGrid.Faces[indexStart+rowIndexStart+j].Edges
				faceEdges[0] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-1)*3/2
				faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-1)*3/2 + 2
				faceEdges[2] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-1)*3/2 + 1
			} else {
				faceEdges = dividedGrid.Faces[indexStart+rowIndexStart+j].Edges
				faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 2, edgeSubdivisions-j/2, edgeSubdivisions)
				faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (j-2)*3/2 + 2
				faceEdges[2] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + j*3/2 + 1
			}
		}
		// bottom corner 2
		faceEdges = dividedGrid.Faces[indexStart+subFaceCount-1].Edges
		faceEdges[0] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 2, 0, edgeSubdivisions)
		faceEdges[1] = edgeOffset + edgeSubdivisions*(edgeSubdivisions-1)*3/2 + (subFaceCount-rowIndexStart-3)*3/2 + 2
		faceEdges[2] = oldGrid.edgeIndexAtClockwiseIndexOnOldFace(faceIndex, 1, edgeSubdivisions, edgeSubdivisions)
	}

	// set edge faces from the previously build edge arrays
	for index, _ := range dividedGrid.Faces {
		dividedGrid.updateEdgesFromFace(int32(index))
	}

	return nil
}

func (grid WingedGrid) updateEdgesFromFace(faceIndex int32) {
	var edges []int32 = grid.Faces[faceIndex].Edges
	var thisEdge, nextEdge, prevEdge WingedEdge
	if edges[len(edges)-1] >= int32(len(grid.Edges)) || edges[len(edges)-1] < 0 {
		// breakpoint
		log.Printf("break")
	}
	prevEdge = grid.Edges[edges[len(edges)-1]]
	thisEdge = grid.Edges[edges[0]]
	nextEdge = grid.Edges[edges[1]]
	// test verticies
	if thisEdge.FirstVertexA == prevEdge.FirstVertexA || thisEdge.FirstVertexA == prevEdge.FirstVertexB {
		// check the next edge also matches and face A is not set
		if thisEdge.FirstVertexB == nextEdge.FirstVertexA || thisEdge.FirstVertexB == nextEdge.FirstVertexB {
			if thisEdge.FaceA == -1 {
				grid.Edges[edges[0]].FaceA = faceIndex
				grid.Edges[edges[0]].PrevA = edges[len(edges)-1]
				grid.Edges[edges[0]].NextA = edges[1]
			} else {
				log.Printf("For face %d. Face A has already been set for edge: %v With edge set: %v", faceIndex, thisEdge, edges)
			}

		} else {
			log.Printf("For face %d. Previous edge matches, but next edge doesn't share correct vertex!", faceIndex)
		}
	} else if thisEdge.FirstVertexB == prevEdge.FirstVertexA || thisEdge.FirstVertexB == prevEdge.FirstVertexB {
		// check the next edge also matches and face B is not set
		if thisEdge.FirstVertexA == nextEdge.FirstVertexA || thisEdge.FirstVertexA == nextEdge.FirstVertexB {
			if thisEdge.FaceB == -1 {
				grid.Edges[edges[0]].FaceB = faceIndex
				grid.Edges[edges[0]].PrevB = edges[len(edges)-1]
				grid.Edges[edges[0]].NextB = edges[1]
			} else {
				log.Printf("For face %d. Face B has already been set for edge: %v With edge set: %v", faceIndex, thisEdge, edges)
			}
		} else {
			log.Printf("For face %d. Previous edge matches, but next edge doesn't share correct vertex!", faceIndex)
		}
	} else {
		log.Printf("For face %d. Edges Don't share a vertex!", faceIndex)
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
					grid.Edges[edges[i]].FaceA = faceIndex
					grid.Edges[edges[i]].PrevA = edges[i-1]
					grid.Edges[edges[i]].NextA = edges[i+1]
				} else {
					log.Printf("For face %d. Face A has already been set for edge: %v With edge set: %v", faceIndex, thisEdge, edges)
				}

			} else {
				log.Printf("For face %d. Previous edge matches, but next edge doesn't share correct vertex!", faceIndex)
			}
		} else if thisEdge.FirstVertexB == prevEdge.FirstVertexA || thisEdge.FirstVertexB == prevEdge.FirstVertexB {
			// check the next edge also matches
			if thisEdge.FirstVertexA == nextEdge.FirstVertexA || thisEdge.FirstVertexA == nextEdge.FirstVertexB {
				if thisEdge.FaceB == -1 {
					grid.Edges[edges[i]].FaceB = faceIndex
					grid.Edges[edges[i]].PrevB = edges[i-1]
					grid.Edges[edges[i]].NextB = edges[i+1]
				} else {
					log.Printf("Face B has already been set for edge: %v With edge set: %v", thisEdge, edges)
				}

			} else {
				log.Printf("For face %d. Previous edge matches, but next edge doesn't share correct vertex!", faceIndex)
			}
		} else {
			log.Printf("For face %d. Edges Don't share a vertex!", faceIndex)
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
				grid.Edges[edges[len(edges)-1]].FaceA = faceIndex
				grid.Edges[edges[len(edges)-1]].PrevA = edges[len(edges)-2]
				grid.Edges[edges[len(edges)-1]].NextA = edges[0]
			} else {
				log.Printf("For face %d. Face A has already been set for edge: %v With edge set: %v", faceIndex, thisEdge, edges)
			}
		} else {
			log.Printf("For face %d. Previous edge matches, but next edge doesn't share correct vertex!", faceIndex)
		}
	} else if thisEdge.FirstVertexB == prevEdge.FirstVertexA || thisEdge.FirstVertexB == prevEdge.FirstVertexB {
		// check the next edge also matches
		if thisEdge.FirstVertexA == nextEdge.FirstVertexA || thisEdge.FirstVertexA == nextEdge.FirstVertexB {
			if thisEdge.FaceB == -1 {
				grid.Edges[edges[len(edges)-1]].FaceB = faceIndex
				grid.Edges[edges[len(edges)-1]].PrevB = edges[len(edges)-2]
				grid.Edges[edges[len(edges)-1]].NextB = edges[0]
			} else {
				log.Printf("For face %d. Face B has already been set for edge: %v With edge set: %v", faceIndex, thisEdge, edges)
			}
		} else {
			log.Printf("For face %d. Previous edge matches, but next edge doesn't share correct vertex!", faceIndex)
		}
	} else {
		log.Printf("For face %d. Edges Don't share a vertex!", faceIndex)
	}

}

/******************* VERTEX SUBDIVISION ***********************/
func (oldGrid WingedGrid) subdivideVertices(edgeSubdivisions int32, dividedGrid WingedGrid) {
	// set coords for the origional verts
	for index, vertex := range oldGrid.Vertices {
		dividedGrid.Vertices[index].Coords[0] = vertex.Coords[0]
		dividedGrid.Vertices[index].Coords[1] = vertex.Coords[1]
		dividedGrid.Vertices[index].Coords[2] = vertex.Coords[2]
	}

	// subdivide along each edge
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

		var divisionLength float64
		for j = 0; j < edgeSubdivisions; j++ {
			// find the new vertex position and create the vertex
			// but don't correct it's length yet
			divisionLength = math.Sin(angleToSubdivide*(float64(j+1)/float64(edgeSubdivisions+1))) * sphereRadius / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(float64(j+1)/float64(edgeSubdivisions+1)))

			dividedGrid.Vertices[origVertexCount+i*edgeSubdivisions+j].Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
			dividedGrid.Vertices[origVertexCount+i*edgeSubdivisions+j].Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
			dividedGrid.Vertices[origVertexCount+i*edgeSubdivisions+j].Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength
		}
	}

	// subdivide face interior
	var faceIndex int32
	for dummy, _ := range oldGrid.Faces {
		faceIndex = int32(dummy)
		var vertexOffset int32 = int32(len(oldGrid.Vertices)) + int32(len(oldGrid.Edges))*edgeSubdivisions + edgeSubdivisions*(edgeSubdivisions-1)/2*faceIndex
		// only if we have more than one division
		if edgeSubdivisions > 1 {
			for i = 0; i < edgeSubdivisions-1; i++ {
				var firstVertex WingedVertex = dividedGrid.Vertices[oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 0, edgeSubdivisions-2-i, edgeSubdivisions)]
				var secondVertex WingedVertex = dividedGrid.Vertices[oldGrid.vertexIndexAtClockwiseIndexOnOldFace(faceIndex, 1, 1+i, edgeSubdivisions)]
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
					divisionLength = math.Sin(angleToSubdivide*(float64(j+1)/float64(i+2))) * firstVectorLength / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(float64(j+1)/float64(i+2)))
					if vertexOffset+(i*(i+1)/2)+j >= int32(len(dividedGrid.Vertices)) {
						log.Printf("breakpoint")
					}
					dividedGrid.Vertices[vertexOffset+(i*(i+1)/2)+j].Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
					dividedGrid.Vertices[vertexOffset+(i*(i+1)/2)+j].Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
					dividedGrid.Vertices[vertexOffset+(i*(i+1)/2)+j].Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength
				}
			}
		}
	}
}

func (grid WingedGrid) normalizeVerticesToSphere(baseVertexIndex int) {
	var wantedLength float64
	wantedLength = vectorLength(grid.Vertices[baseVertexIndex].Coords)
	for i, vertex := range grid.Vertices {
		var currentLength float64
		currentLength = vectorLength(vertex.Coords)
		grid.Vertices[i].Coords[0] = vertex.Coords[0] * wantedLength / currentLength
		grid.Vertices[i].Coords[1] = vertex.Coords[1] * wantedLength / currentLength
		grid.Vertices[i].Coords[2] = vertex.Coords[2] * wantedLength / currentLength
	}
}

func (grid WingedGrid) NormalizeVerticesToDistanceFromOrigin(wantedLength float64) {
	for i, vertex := range grid.Vertices {
		var currentLength float64
		currentLength = vectorLength(vertex.Coords)
		grid.Vertices[i].Coords[0] = vertex.Coords[0] * wantedLength / currentLength
		grid.Vertices[i].Coords[1] = vertex.Coords[1] * wantedLength / currentLength
		grid.Vertices[i].Coords[2] = vertex.Coords[2] * wantedLength / currentLength
	}
}

func (grid WingedGrid) setEdgesForVerticesIfInvalid() {
	// loop through edges so we only have to touch each one once
	for index, edge := range grid.Edges {
		var nextEdgeIndex int32 = -1
		var nextEdge WingedEdge
		var theVertexIndex int32 = edge.FirstVertexA
		var theVertex WingedVertex = grid.Vertices[theVertexIndex]
		if theVertex.Edges[0] == -1 {
			nextEdgeIndex, _ = edge.NextEdgeForVertex(theVertexIndex)
			nextEdge = grid.Edges[nextEdgeIndex]
			theVertex.Edges[0] = int32(index)
			var i int = 1
			for int32(index) != nextEdgeIndex {
				theVertex.Edges[i] = nextEdgeIndex
				nextEdgeIndex, _ = nextEdge.NextEdgeForVertex(theVertexIndex)
				nextEdge = grid.Edges[nextEdgeIndex]
				i = i + 1
			}
		}
		theVertexIndex = edge.FirstVertexB
		theVertex = grid.Vertices[theVertexIndex]
		if theVertex.Edges[0] == -1 {
			nextEdgeIndex, _ = edge.NextEdgeForVertex(theVertexIndex)
			nextEdge = grid.Edges[nextEdgeIndex]
			theVertex.Edges[0] = int32(index)
			var i int = 1
			for int32(index) != nextEdgeIndex {
				theVertex.Edges[i] = nextEdgeIndex
				nextEdgeIndex, _ = nextEdge.NextEdgeForVertex(theVertexIndex)
				nextEdge = grid.Edges[nextEdgeIndex]
				i = i + 1
			}
		}
	}
}

/******************* Helper Functions ***********************/

func vectorAngle(first, second [3]float64) float64 {
	return math.Acos((first[0]*second[0] + first[1]*second[1] + first[2]*second[2]) / (math.Sqrt(first[0]*first[0]+first[1]*first[1]+first[2]*first[2]) * math.Sqrt(second[0]*second[0]+second[1]*second[1]+second[2]*second[2])))
}
func vectorLength(vector [3]float64) float64 {
	return math.Sqrt(vector[0]*vector[0] + vector[1]*vector[1] + vector[2]*vector[2])
}

func (grid WingedGrid) vertexIndexAtClockwiseIndexOnOldFace(faceIndex, edgeInFaceIndex, clockwiseVertexIndex, edgeSubdivisions int32) int32 {
	var edgeIndex int32 = grid.Faces[faceIndex].Edges[edgeInFaceIndex]
	var edge WingedEdge = grid.Edges[edgeIndex]

	if edge.FaceA == faceIndex {
		return int32(len(grid.Vertices)) + edgeIndex*edgeSubdivisions + clockwiseVertexIndex
	}
	if edge.FaceB == faceIndex {
		return int32(len(grid.Vertices)) + edgeIndex*edgeSubdivisions + edgeSubdivisions - 1 - clockwiseVertexIndex
	}
	return -1
}

func (grid WingedGrid) edgeIndexAtClockwiseIndexOnOldFace(faceIndex, edgeInFaceIndex, clockwiseEdgeIndex, edgeSubdivisions int32) int32 {
	var oldEdgeIndex int32 = grid.Faces[faceIndex].Edges[edgeInFaceIndex]
	var oldEdge WingedEdge = grid.Edges[oldEdgeIndex]
	if oldEdge.FaceA == faceIndex {
		return oldEdgeIndex*(edgeSubdivisions+1) + clockwiseEdgeIndex
	}
	if oldEdge.FaceB == faceIndex {
		return oldEdgeIndex*(edgeSubdivisions+1) + edgeSubdivisions - clockwiseEdgeIndex
	}
	return -1
}
