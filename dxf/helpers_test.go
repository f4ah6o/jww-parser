package dxf

import (
	"math"
	"testing"
)

func TestLineLength(t *testing.T) {
	line := NewLine(0, 0, 100, 0)
	length := line.Length()

	if length != 100 {
		t.Errorf("Expected length 100, got %f", length)
	}

	// Test diagonal line
	diagonal := NewLine(0, 0, 100, 100)
	diagonalLength := diagonal.Length()
	expected := math.Sqrt(100*100 + 100*100)

	if math.Abs(diagonalLength-expected) > 0.0001 {
		t.Errorf("Expected diagonal length %f, got %f", expected, diagonalLength)
	}
}

func TestLineBoundingBox(t *testing.T) {
	line := NewLine(10, 20, 100, 200)
	minX, minY, maxX, maxY := line.BoundingBox()

	if minX != 10 || minY != 20 || maxX != 100 || maxY != 200 {
		t.Errorf("Expected bounding box (10, 20, 100, 200), got (%f, %f, %f, %f)",
			minX, minY, maxX, maxY)
	}
}

func TestLineMidPoint(t *testing.T) {
	line := NewLine(0, 0, 100, 100)
	x, y := line.MidPoint()

	if x != 50 || y != 50 {
		t.Errorf("Expected midpoint (50, 50), got (%f, %f)", x, y)
	}
}

func TestLineAngle(t *testing.T) {
	line := NewLine(0, 0, 100, 0)
	angle := line.Angle()

	if angle != 0 {
		t.Errorf("Expected angle 0, got %f", angle)
	}

	// Test 45° line
	diagonal := NewLine(0, 0, 100, 100)
	diagAngle := diagonal.Angle()

	if math.Abs(diagAngle-45) > 0.0001 {
		t.Errorf("Expected angle 45, got %f", diagAngle)
	}
}

func TestCircleArea(t *testing.T) {
	circle := NewCircle(50, 50, 10)
	area := circle.Area()
	expected := math.Pi * 100

	if math.Abs(area-expected) > 0.0001 {
		t.Errorf("Expected area %f, got %f", expected, area)
	}
}

func TestCircleCircumference(t *testing.T) {
	circle := NewCircle(50, 50, 10)
	circ := circle.Circumference()
	expected := 2 * math.Pi * 10

	if math.Abs(circ-expected) > 0.0001 {
		t.Errorf("Expected circumference %f, got %f", expected, circ)
	}
}

func TestCircleBoundingBox(t *testing.T) {
	circle := NewCircle(50, 50, 25)
	minX, minY, maxX, maxY := circle.BoundingBox()

	if minX != 25 || minY != 25 || maxX != 75 || maxY != 75 {
		t.Errorf("Expected bounding box (25, 25, 75, 75), got (%f, %f, %f, %f)",
			minX, minY, maxX, maxY)
	}
}

func TestArcArcLength(t *testing.T) {
	// Quarter circle (90°)
	arc := NewArc(50, 50, 10, 0, 90)
	length := arc.ArcLength()
	expected := math.Pi * 10 / 2 // Quarter of circumference

	if math.Abs(length-expected) > 0.0001 {
		t.Errorf("Expected arc length %f, got %f", expected, length)
	}
}

func TestArcArea(t *testing.T) {
	// Quarter circle (90°)
	arc := NewArc(50, 50, 10, 0, 90)
	area := arc.Area()
	expected := 0.5 * 10 * 10 * (math.Pi / 2) // Sector area

	if math.Abs(area-expected) > 0.0001 {
		t.Errorf("Expected sector area %f, got %f", expected, area)
	}
}

func TestPointBoundingBox(t *testing.T) {
	point := NewPoint(100, 200)
	minX, minY, maxX, maxY := point.BoundingBox()

	if minX != 100 || minY != 200 || maxX != 100 || maxY != 200 {
		t.Errorf("Expected bounding box (100, 200, 100, 200), got (%f, %f, %f, %f)",
			minX, minY, maxX, maxY)
	}
}

func TestSolidBoundingBox(t *testing.T) {
	solid := NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
	minX, minY, maxX, maxY := solid.BoundingBox()

	if minX != 0 || minY != 0 || maxX != 100 || maxY != 100 {
		t.Errorf("Expected bounding box (0, 0, 100, 100), got (%f, %f, %f, %f)",
			minX, minY, maxX, maxY)
	}
}

func TestSolidArea(t *testing.T) {
	// Simple triangle
	solid := NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
	area := solid.Area()

	// Area should be positive
	if area <= 0 {
		t.Errorf("Expected positive area, got %f", area)
	}
}

func TestSolidIsTriangle(t *testing.T) {
	triangle := NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
	if !triangle.IsTriangle() {
		t.Errorf("Expected solid to be a triangle")
	}

	quad := NewSolid(0, 0, 100, 0, 100, 100, 0, 100)
	if quad.IsTriangle() {
		t.Errorf("Expected solid to not be a triangle")
	}
}

func TestDocumentBoundingBox(t *testing.T) {
	doc := NewDocument().
		AddLine(0, 0, 100, 100).
		AddCircle(200, 200, 50)

	minX, minY, maxX, maxY := doc.BoundingBox()

	if minX != 0 || minY != 0 {
		t.Errorf("Expected min corner (0, 0), got (%f, %f)", minX, minY)
	}
	if maxX != 250 || maxY != 250 {
		t.Errorf("Expected max corner (250, 250), got (%f, %f)", maxX, maxY)
	}
}

func TestDocumentFilterByLayer(t *testing.T) {
	doc := NewDocument().
		AddLine(0, 0, 100, 100, WithLineLayer("Layer1")).
		AddLine(0, 0, 50, 50, WithLineLayer("Layer2")).
		AddCircle(50, 50, 25, WithCircleLayer("Layer1"))

	layer1Entities := doc.FilterByLayer("Layer1")
	if len(layer1Entities) != 2 {
		t.Errorf("Expected 2 entities on Layer1, got %d", len(layer1Entities))
	}

	layer2Entities := doc.FilterByLayer("Layer2")
	if len(layer2Entities) != 1 {
		t.Errorf("Expected 1 entity on Layer2, got %d", len(layer2Entities))
	}
}

func TestDocumentCountByType(t *testing.T) {
	doc := NewDocument().
		AddLine(0, 0, 100, 100).
		AddLine(0, 0, 50, 50).
		AddCircle(50, 50, 25).
		AddPoint(100, 100)

	counts := doc.CountByType()

	if counts["LINE"] != 2 {
		t.Errorf("Expected 2 lines, got %d", counts["LINE"])
	}
	if counts["CIRCLE"] != 1 {
		t.Errorf("Expected 1 circle, got %d", counts["CIRCLE"])
	}
	if counts["POINT"] != 1 {
		t.Errorf("Expected 1 point, got %d", counts["POINT"])
	}
}
