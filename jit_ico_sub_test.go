package wingedGrid

import (
	"log"
	"testing"
)

func TestVertsAndNeighbors(t *testing.T) {
	subCount := 10
	// create subdivided
	base, _ := BaseIcosahedron()
	sub, _ := base.SubdivideTriangles(int32(subCount))

	isc := NewIcoSubCalc()

	// test origional
	var idx int
	for ; idx < len(base.Vertices); idx++ {
		vert := sub.Vertices[idx]
		calcVert := isc.VertexAndNeighbors(idx, subCount)

		// check vert is close enough
		dist := distanceBetween3Points(vert.Coords, calcVert.Coords)
		if dist > 0.0000001 {
			log.Printf("distance %f too great on vert %d", dist, idx)
			t.Fail()
		}

		// check neibhros
		vNeighbors, _ := sub.NeighborsForVertex(int32(idx))
		if len(calcVert.vertexNeighbors) != len(vNeighbors) {
			log.Printf("Counts don't match %d %d", len(calcVert.vertexNeighbors), len(vNeighbors))
			t.FailNow()
		}
		startIdx := -1
		for i, edx := range vNeighbors {
			if calcVert.vertexNeighbors[0] == edx {
				startIdx = i
				break
			}
		}

		if startIdx == -1 {
			log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
			t.Fatal("Start Not Found")
		}

		for i, edx := range vNeighbors {
			// find start idx
			calcIdx := (len(vNeighbors) + i - startIdx) % len(vNeighbors)
			if calcVert.vertexNeighbors[calcIdx] != edx {
				log.Printf("Edge %d incorrect", i)
				log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
				t.FailNow()
			}
		}
	}

	// test edge divs
	for ; idx < len(base.Vertices)+len(base.Edges)*subCount; idx++ {
		vert := sub.Vertices[idx]
		calcVert := isc.VertexAndNeighbors(idx, subCount)

		// check vert is close enough
		dist := distanceBetween3Points(vert.Coords, calcVert.Coords)
		if dist > 0.0000001 {
			log.Printf("distance %f too great on vert %d", dist, idx)
			t.FailNow()
		}

		// check neibhros
		vNeighbors, _ := sub.NeighborsForVertex(int32(idx))
		if len(calcVert.vertexNeighbors) != len(vNeighbors) {
			log.Printf("Counts don't match %d %d", len(calcVert.vertexNeighbors), len(vNeighbors))
			t.FailNow()
		}
		startIdx := -1
		for i, edx := range vNeighbors {
			if calcVert.vertexNeighbors[0] == edx {
				startIdx = i
				break
			}
		}

		if startIdx == -1 {
			log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
			t.Fatal("Start Not Found")
		}

		for i, edx := range vNeighbors {
			// find start idx
			calcIdx := (len(vNeighbors) + i - startIdx) % len(vNeighbors)
			if calcVert.vertexNeighbors[calcIdx] != edx {
				log.Printf("Edge %d incorrect", i)
				log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
				t.FailNow()
			}
		}
	}

	// test face divs
	for ; idx < len(sub.Vertices); idx++ {
		vert := sub.Vertices[idx]
		calcVert := isc.VertexAndNeighbors(idx, subCount)

		// check vert is close enough
		dist := distanceBetween3Points(vert.Coords, calcVert.Coords)
		if dist > 0.0000001 {
			log.Printf("distance %f too great on vert %d", dist, idx)
			t.FailNow()
		}

		// check neibhros
		vNeighbors, _ := sub.NeighborsForVertex(int32(idx))
		if len(calcVert.vertexNeighbors) != len(vNeighbors) {
			log.Printf("Counts don't match %d %d", len(calcVert.vertexNeighbors), len(vNeighbors))
			t.FailNow()
		}
		startIdx := -1
		for i, edx := range vNeighbors {
			if calcVert.vertexNeighbors[0] == edx {
				startIdx = i
				break
			}
		}

		if startIdx == -1 {
			log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
			t.Fatal("Start Not Found")
		}

		for i, edx := range vNeighbors {
			// find start idx
			calcIdx := (len(vNeighbors) + i - startIdx) % len(vNeighbors)
			if calcVert.vertexNeighbors[calcIdx] != edx {
				log.Printf("Neighbor %d incorrect", i)
				log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
				t.FailNow()
			}
		}
	}

	// check beyond?
}

func TestSubVertsAndNeighbors(t *testing.T) {
	initialSubCount := 10
	subCount := 10
	// create subdivided
	b, _ := BaseIcosahedron()
	base, _ := b.SubdivideTriangles(int32(initialSubCount))
	sub, _ := base.SubdivideTriangles(int32(subCount))

	isc := NewSubIcoSubCalc(base)

	// test origional
	var idx int
	for ; idx < len(base.Vertices); idx++ {
		vert := sub.Vertices[idx]
		calcVert := isc.VertexAndNeighbors(idx, subCount)

		// check vert is close enough
		dist := distanceBetween3Points(vert.Coords, calcVert.Coords)
		if dist > 0.0000001 {
			log.Printf("distance %f too great on vert %d", dist, idx)
			t.Fail()
		}

		// check neibhros
		vNeighbors, _ := sub.NeighborsForVertex(int32(idx))
		if len(calcVert.vertexNeighbors) != len(vNeighbors) {
			log.Printf("Counts don't match %d %d", len(calcVert.vertexNeighbors), len(vNeighbors))
			t.FailNow()
		}
		startIdx := -1
		for i, edx := range vNeighbors {
			if calcVert.vertexNeighbors[0] == edx {
				startIdx = i
				break
			}
		}

		if startIdx == -1 {
			log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
			t.Fatal("Start Not Found")
		}

		for i, edx := range vNeighbors {
			// find start idx
			calcIdx := (len(vNeighbors) + i - startIdx) % len(vNeighbors)
			if calcVert.vertexNeighbors[calcIdx] != edx {
				log.Printf("Edge %d incorrect", i)
				log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
				t.FailNow()
			}
		}
	}

	// test edge divs
	for ; idx < len(base.Vertices)+len(base.Edges)*subCount; idx++ {
		vert := sub.Vertices[idx]
		calcVert := isc.VertexAndNeighbors(idx, subCount)

		// check vert is close enough
		dist := distanceBetween3Points(vert.Coords, calcVert.Coords)
		if dist > 0.0000001 {
			log.Printf("distance %f too great on vert %d", dist, idx)
			t.FailNow()
		}

		// check neibhros
		vNeighbors, _ := sub.NeighborsForVertex(int32(idx))
		if len(calcVert.vertexNeighbors) != len(vNeighbors) {
			log.Printf("Counts don't match %d %d", len(calcVert.vertexNeighbors), len(vNeighbors))
			t.FailNow()
		}
		startIdx := -1
		for i, edx := range vNeighbors {
			if calcVert.vertexNeighbors[0] == edx {
				startIdx = i
				break
			}
		}

		if startIdx == -1 {
			log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
			t.Fatal("Start Not Found")
		}

		for i, edx := range vNeighbors {
			// find start idx
			calcIdx := (len(vNeighbors) + i - startIdx) % len(vNeighbors)
			if calcVert.vertexNeighbors[calcIdx] != edx {
				log.Printf("Edge %d incorrect", i)
				log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
				t.FailNow()
			}
		}
	}

	// test face divs
	for ; idx < len(sub.Vertices); idx++ {
		vert := sub.Vertices[idx]
		calcVert := isc.VertexAndNeighbors(idx, subCount)

		// check vert is close enough
		dist := distanceBetween3Points(vert.Coords, calcVert.Coords)
		if dist > 0.0000001 {
			log.Printf("distance %f too great on vert %d", dist, idx)
			t.FailNow()
		}

		// check neibhros
		vNeighbors, _ := sub.NeighborsForVertex(int32(idx))
		if len(calcVert.vertexNeighbors) != len(vNeighbors) {
			log.Printf("Counts don't match %d %d", len(calcVert.vertexNeighbors), len(vNeighbors))
			t.FailNow()
		}
		startIdx := -1
		for i, edx := range vNeighbors {
			if calcVert.vertexNeighbors[0] == edx {
				startIdx = i
				break
			}
		}

		if startIdx == -1 {
			log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
			t.Fatal("Start Not Found")
		}

		for i, edx := range vNeighbors {
			// find start idx
			calcIdx := (len(vNeighbors) + i - startIdx) % len(vNeighbors)
			if calcVert.vertexNeighbors[calcIdx] != edx {
				log.Printf("Neighbor %d incorrect", i)
				log.Printf("Calc: %#v\n Sub: %#v", calcVert.vertexNeighbors, vNeighbors)
				t.FailNow()
			}
		}
	}

	// check beyond?
}
