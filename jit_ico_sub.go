package wingedGrid

import (
	"math"
)

type IcoSubCalc struct {
	baseIco WingedGrid
}

func NewIcoSubCalc() *IcoSubCalc {
	ico, _ := BaseIcosahedron()
	return &IcoSubCalc{
		baseIco: ico,
	}
}

func NewSubIcoSubCalc(grid WingedGrid) *IcoSubCalc {
	return &IcoSubCalc{
		baseIco: grid,
	}
}

// VertexAndNeighbors does things
func (isc *IcoSubCalc) VertexAndNeighbors(idx int, subDivs int) WingedVertex {
	return WingedVertex{
		Coords:          isc.Vertex(idx, subDivs).Coords,
		vertexNeighbors: isc.VertexNeighbors(idx, subDivs).vertexNeighbors,
	}
}

func (isc *IcoSubCalc) VertexNeighbors(idx int, subDivs int) WingedVertex {
	baseVerts := len(isc.baseIco.Vertices)
	baseEdges := len(isc.baseIco.Edges)
	baseFaces := len(isc.baseIco.Faces)
	// three cases
	// origional vert
	// vert on origional edge
	// vert inside original face
	if idx < baseVerts {
		// origional vert
		return isc.vertexNeighborsOrigional(idx, subDivs)
	}
	idx -= baseVerts
	if idx < baseEdges*subDivs {
		return isc.vertexNeighborsEdge(idx, subDivs)
	}
	idx -= baseEdges * subDivs
	if idx < (subDivs-1)*subDivs/2*baseFaces {
		return isc.vertexNeighborsFace(idx, subDivs)
	}
	panic("past end")
}

func (isc *IcoSubCalc) vertexNeighborsOrigional(idx int, subDivs int) WingedVertex {
	neighbors := make([]int32, len(isc.baseIco.Vertices[idx].Edges))
	var origVertexCount int32 = int32(len(isc.baseIco.Vertices))

	for i, edx := range isc.baseIco.Vertices[idx].Edges {
		edge := isc.baseIco.Edges[edx]

		if edge.FirstVertexA == int32(idx) {
			neighbors[i] = origVertexCount + edx*int32(subDivs)
		} else {
			neighbors[i] = origVertexCount + edx*int32(subDivs) + int32(subDivs) - 1
		}
	}
	return WingedVertex{
		vertexNeighbors: neighbors,
	}
}

func (isc *IcoSubCalc) vertexNeighborsEdge(idx int, subDivs int) WingedVertex {
	baseVerts := len(isc.baseIco.Vertices)
	// six for all but 12 origional pentagons
	neighbors := make([]int32, 6)

	edgeIdx, div := divmod(idx, subDivs)
	edge := isc.baseIco.Edges[edgeIdx]

	// two cases, right next to corner, and not
	if div == 0 || div == subDivs-1 {
		// on edge ones
		if div == 0 {
			neighbors[0] = edge.FirstVertexA
			neighbors[3] = int32(idx+baseVerts) + 1
		} else {
			neighbors[0] = edge.FirstVertexB
			neighbors[3] = int32(idx+baseVerts) - 1
		}

		// find this edge on vertex
		edges := isc.baseIco.Vertices[neighbors[0]].Edges
		edgeInVertEdx := -1
		for i, edx := range edges {
			if edx == int32(edgeIdx) {
				edgeInVertEdx = i
				break
			}
		}

		// 5 is next, 1 is previous
		nextIdx := edges[(len(edges)+edgeInVertEdx+1)%len(edges)]
		nextEdge := isc.baseIco.Edges[int(nextIdx)]
		if nextEdge.FirstVertexA == neighbors[0] {
			neighbors[5] = int32(baseVerts) + nextIdx*int32(subDivs)
		} else {
			neighbors[5] = int32(baseVerts) + nextIdx*int32(subDivs) + int32(subDivs) - 1
		}
		prevIdx := edges[(len(edges)+edgeInVertEdx-1)%len(edges)]
		prevEdge := isc.baseIco.Edges[int(prevIdx)]
		if prevEdge.FirstVertexA == neighbors[0] {
			neighbors[1] = int32(baseVerts) + prevIdx*int32(subDivs)
		} else {
			neighbors[1] = int32(baseVerts) + prevIdx*int32(subDivs) + int32(subDivs) - 1
		}

		// 2 and 4
		// offset from face options are
		// 0, (subdivs-2)*(subdivs-1), subdivs*(subdivs-1) - 1

		// check face a
		var faceFourIdx, faceTwoIdx int32
		if div == 0 {
			faceFourIdx = edge.FaceA
			faceTwoIdx = edge.FaceB
		} else {
			faceFourIdx = edge.FaceB
			faceTwoIdx = edge.FaceA
		}
		faceFour := isc.baseIco.Faces[int(faceFourIdx)]
		var offset int
		if faceFour.Edges[0] == int32(edgeIdx) {
			// we are next to first vertex, first element of last row
			offset = (subDivs - 2) * (subDivs - 1) / 2
		} else if faceFour.Edges[1] == int32(edgeIdx) {
			// we are next to second vertex, first element
			offset = 0
		} else {
			// we are next to third vertext, last element
			offset = subDivs*(subDivs-1)/2 - 1
		}

		neighbors[4] = int32(baseVerts + len(isc.baseIco.Edges)*subDivs + (subDivs-1)*subDivs/2*int(faceFourIdx) + offset)

		// check face a
		faceTwo := isc.baseIco.Faces[int(faceTwoIdx)]
		if faceTwo.Edges[0] == int32(edgeIdx) {
			// we are next to second vertex, first element
			offset = 0
		} else if faceTwo.Edges[1] == int32(edgeIdx) {
			// we are next to third vertext, last element
			offset = subDivs*(subDivs-1)/2 - 1
		} else {
			// we are next to first vertex, first element of last row
			offset = (subDivs - 2) * (subDivs - 1) / 2
		}
		neighbors[2] = int32(baseVerts + len(isc.baseIco.Edges)*subDivs + (subDivs-1)*subDivs/2*int(faceTwoIdx) + offset)

	} else {
		// not next to corner
		neighbors[0] = int32(idx+baseVerts) - 1
		neighbors[3] = int32(idx+baseVerts) + 1

		faceA := isc.baseIco.Faces[int(edge.FaceA)]
		faceOffset := int32(baseVerts + len(isc.baseIco.Edges)*subDivs + (subDivs-1)*subDivs/2*int(edge.FaceA))

		if faceA.Edges[0] == int32(edgeIdx) {
			neighbors[4] = faceOffset + int32(firstInRow(subDivs-div-1))
			neighbors[5] = faceOffset + int32(firstInRow(subDivs-div))
		} else if faceA.Edges[1] == int32(edgeIdx) {
			neighbors[4] = faceOffset + int32(lastInRow(div+1))
			neighbors[5] = faceOffset + int32(lastInRow(div))
		} else {
			// parallel to rows
			neighbors[4] = faceOffset + int32(lastInRow(subDivs-1)-div)
			neighbors[5] = faceOffset + int32(lastInRow(subDivs-1)-div+1)
		}

		faceB := isc.baseIco.Faces[int(edge.FaceB)]
		faceOffset = int32(baseVerts + len(isc.baseIco.Edges)*subDivs + (subDivs-1)*subDivs/2*int(edge.FaceB))
		if faceB.Edges[0] == int32(edgeIdx) {
			neighbors[1] = faceOffset + int32(firstInRow(div))
			neighbors[2] = faceOffset + int32(firstInRow(div+1))
		} else if faceB.Edges[1] == int32(edgeIdx) {
			neighbors[1] = faceOffset + int32(lastInRow(subDivs-div))
			neighbors[2] = faceOffset + int32(lastInRow(subDivs-div-1))
		} else {
			// parallel to rows
			neighbors[1] = faceOffset + int32(lastInRow(subDivs-2)+div)
			neighbors[2] = faceOffset + int32(lastInRow(subDivs-2)+div+1)
		}
	}

	return WingedVertex{
		vertexNeighbors: neighbors,
	}
}

func (isc *IcoSubCalc) vertexNeighborsFace(idx int, subDivs int) WingedVertex {
	faceIdx, loc := divmod(idx, (subDivs-1)*subDivs/2)
	neighbors := make([]int32, 6)

	// zero indexed row
	row := int(math.Ceil((math.Sqrt(float64(8*(loc+1))+1)-1)*0.5)) - 1

	along := loc - row*(row+1)/2

	faceOffset := len(isc.baseIco.Vertices) + len(isc.baseIco.Edges)*subDivs + (subDivs-1)*subDivs/2*int(faceIdx)

	// write all in, some will be incorrect
	// row above
	neighbors[0] = int32(faceOffset + loc - row - 1)
	neighbors[1] = int32(faceOffset + loc - row)
	// next in row
	neighbors[2] = int32(faceOffset + loc + 1)
	// row below
	neighbors[3] = int32(faceOffset + loc + row + 2)
	neighbors[4] = int32(faceOffset + loc + row + 1)
	// prev in row
	neighbors[5] = int32(faceOffset + loc - 1)

	// check if first in row, 0 and 5
	if along == 0 {
		neighbors[0] = isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(faceIdx), 0, int32(subDivs-row-1), int32(subDivs))
		neighbors[5] = isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(faceIdx), 0, int32(subDivs-row-2), int32(subDivs))
	}

	// check if last in row, 1 and 2
	if along == row {
		neighbors[1] = isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(faceIdx), 1, int32(row), int32(subDivs))
		neighbors[2] = isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(faceIdx), 1, int32(row+1), int32(subDivs))
	}

	// check if in last row, 3 and 4
	if row == subDivs-2 {
		neighbors[3] = isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(faceIdx), 2, int32(subDivs-along-2), int32(subDivs))
		neighbors[4] = isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(faceIdx), 2, int32(subDivs-along-1), int32(subDivs))
	}

	return WingedVertex{
		vertexNeighbors: neighbors,
	}
}

func firstInRow(row int) int {
	return (row - 1) * row / 2
}

func lastInRow(row int) int {
	return firstInRow(row+1) - 1
}

func (isc *IcoSubCalc) Vertex(idx int, subDivs int) WingedVertex {
	baseVerts := len(isc.baseIco.Vertices)
	baseEdges := len(isc.baseIco.Edges)
	baseFaces := len(isc.baseIco.Faces)
	// three cases
	// origional vert
	// vert on origional edge
	// vert inside original face
	if idx < baseVerts {
		// origional vert
		return isc.vertexOrigional(idx, subDivs)
	}
	idx -= baseVerts
	if idx < baseEdges*subDivs {
		return isc.vertexEdge(idx, subDivs)
	}
	idx -= baseEdges * subDivs
	if idx < (subDivs-1)*subDivs/2*baseFaces {
		return isc.vertexFace(idx, subDivs)
	}

	panic("past end")
}

func (isc *IcoSubCalc) vertexOrigional(idx int, subDivs int) WingedVertex {

	return WingedVertex{
		Coords: isc.baseIco.Vertices[idx].Coords,
	}
}

func (isc *IcoSubCalc) vertexEdge(idx int, subDivs int) WingedVertex {
	edgeIdx, div := divmod(idx, subDivs)
	edge := isc.baseIco.Edges[edgeIdx]
	// calcIdx := len(isc.baseIco.Vertices) + edgeIdx*subDivs + div
	// log.Printf("Calculating %d vertex %d divisions along edge %d ", calcIdx, div, edgeIdx)

	var firstVertex WingedVertex = isc.baseIco.Vertices[edge.FirstVertexA]
	var secondVertex WingedVertex = isc.baseIco.Vertices[edge.FirstVertexB]
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

	var vert WingedVertex
	divisionLength := math.Sin(angleToSubdivide*(float64(div+1)/float64(subDivs+1))) * sphereRadius / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(float64(div+1)/float64(subDivs+1)))

	vert.Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
	vert.Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
	vert.Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength

	return vert
}

func (isc *IcoSubCalc) vertexFace(idx int, subDivs int) WingedVertex {
	face, loc := divmod(idx, (subDivs-1)*subDivs/2)

	// log.Printf("Calculating vertex %d locations along division of face %d ", loc, face)

	// zero indexed row
	row := int(math.Ceil((math.Sqrt(float64(8*(loc+1))+1)-1)*0.5)) - 1

	// calculate trinagle number for previous row [row*(row+1)] since we are zero indexed
	along := loc - row*(row+1)/2

	// get edge verts to work from
	baseVerts := len(isc.baseIco.Vertices)
	first := isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(face), 0, int32(subDivs-2-row), int32(subDivs)) - int32(baseVerts)
	firstVertex := isc.vertexEdge(int(first), subDivs)
	second := isc.baseIco.vertexIndexAtClockwiseIndexOnOldFace(int32(face), 1, int32(1+row), int32(subDivs)) - int32(baseVerts)
	secondVertex := isc.vertexEdge(int(second), subDivs)

	// log.Printf("First Vert idx: %d, Second Vert idx: %d", first, second)
	// log.Printf("First %#v, Second: %#v", firstVertex, secondVertex)

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

	var vert WingedVertex

	divisionLength := math.Sin(angleToSubdivide*(float64(along+1)/float64(row+2))) * firstVectorLength / math.Sin(math.Pi-cornerAngle-angleToSubdivide*(float64(along+1)/float64(row+2)))
	vert.Coords[0] = firstVertex.Coords[0] + stepDirection[0]*divisionLength
	vert.Coords[1] = firstVertex.Coords[1] + stepDirection[1]*divisionLength
	vert.Coords[2] = firstVertex.Coords[2] + stepDirection[2]*divisionLength

	return vert
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
