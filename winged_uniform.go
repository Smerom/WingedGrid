package wingedGrid

import (
	_ "log"
	"math"
)

func normalizedForceVectorWithDistance(vectorFrom [3]float64, vectorTo [3]float64) ([3]float64, float64) {
	var forceVector [3]float64
	forceVector[0] = vectorTo[0] - vectorFrom[0]
	forceVector[1] = vectorTo[1] - vectorFrom[1]
	forceVector[2] = vectorTo[2] - vectorFrom[2]

	return normalize3VectorWithScale(forceVector)
}

func normalize3VectorWithScale(vector [3]float64) ([3]float64, float64) {
	var scale float64
	var result [3]float64
	scale = math.Sqrt(vector[0]*vector[0] + vector[1]*vector[1] + vector[2]*vector[2])
	if scale == 0 {
		panic("no scale!")
	}

	result[0] = vector[0] / scale
	result[1] = vector[1] / scale
	result[2] = vector[2] / scale
	return result, scale
}

func distanceBetween3Points(point1, point2 [3]float64) float64 {
	var x, y, z float64
	x = point1[0] - point2[0]
	y = point1[1] - point2[1]
	z = point1[2] - point2[2]

	return math.Sqrt(x*x + y*y + z*z)
}

func int32InSlice(a int32, list []int32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (modifiedGrid *WingedGrid) UniformVertsOnUnitSphere(steps int) {
	var newCoords [][3]float64
	newCoords = make([][3]float64, len(modifiedGrid.Vertices))

	//var smallDistance float64 = distanceBetween3Points(modifiedGrid.Vertices[modifiedGrid.Edges[0].FirstVertexA].Coords, modifiedGrid.Vertices[modifiedGrid.Edges[0].FirstVertexB].Coords)
	modifiedGrid.NormalizeVerticesToDistanceFromOrigin(1.0)
	for i := 0; i < steps; i++ {
		for index, vertex := range modifiedGrid.Vertices {
			var centerVertex [3]float64
			var outerCenterVertex [3]float64
			var addedNeighbors []int32
			// append self to added neighbors
			addedNeighbors = append(addedNeighbors, int32(index))
			neighbors, _ := modifiedGrid.NeighborsForVertex(int32(index))
			for _, neighborIndex := range neighbors {

				var neighborVert WingedVertex
				neighborVert = modifiedGrid.Vertices[neighborIndex]

				centerVertex[0] += neighborVert.Coords[0]
				centerVertex[1] += neighborVert.Coords[1]
				centerVertex[2] += neighborVert.Coords[2]
				// append to added
				addedNeighbors = append(addedNeighbors, neighborIndex)
			}

			centerVertex, _ = normalize3VectorWithScale(centerVertex)

			forceVector, distance := normalizedForceVectorWithDistance(vertex.Coords, centerVertex)

			newCoords[index][0] = vertex.Coords[0] + forceVector[0]*distance*0.01
			newCoords[index][1] = vertex.Coords[1] + forceVector[1]*distance*0.01
			newCoords[index][2] = vertex.Coords[2] + forceVector[2]*distance*0.01

			// create outer center vertex
			for _, neighborIndex := range neighbors {
				// loop outer neighbors
				outerNeighbors, _ := modifiedGrid.NeighborsForVertex(neighborIndex)
				for _, outerNeighborIndex := range outerNeighbors {
					if !int32InSlice(outerNeighborIndex, addedNeighbors) {
						addedNeighbors = append(addedNeighbors, outerNeighborIndex)

						var neighborVert WingedVertex
						neighborVert = modifiedGrid.Vertices[outerNeighborIndex]

						outerCenterVertex[0] += neighborVert.Coords[0]
						outerCenterVertex[1] += neighborVert.Coords[1]
						outerCenterVertex[2] += neighborVert.Coords[2]
					}
				}
			}
			outerCenterVertex, _ = normalize3VectorWithScale(outerCenterVertex)

			forceVector, distance = normalizedForceVectorWithDistance(vertex.Coords, outerCenterVertex)

			newCoords[index][0] += forceVector[0] * distance * 0.01
			newCoords[index][1] += forceVector[1] * distance * 0.01
			newCoords[index][2] += forceVector[2] * distance * 0.01

		}

		// set coords after all new ones are computed
		for index, _ := range modifiedGrid.Vertices {
			modifiedGrid.Vertices[index].Coords = newCoords[index]
		}
		modifiedGrid.NormalizeVerticesToDistanceFromOrigin(1.0)
	}
}
