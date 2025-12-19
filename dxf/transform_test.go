package dxf

import (
	"math"
	"testing"
)

func TestLineTranslate(t *testing.T) {
	line := NewLine(0, 0, 100, 100)
	moved := line.Translate(50, 50)

	if moved.X1 != 50 || moved.Y1 != 50 {
		t.Errorf("Expected start point (50, 50), got (%f, %f)", moved.X1, moved.Y1)
	}
	if moved.X2 != 150 || moved.Y2 != 150 {
		t.Errorf("Expected end point (150, 150), got (%f, %f)", moved.X2, moved.Y2)
	}
}

func TestLineRotate(t *testing.T) {
	line := NewLine(100, 0, 100, 0)
	rotated := line.Rotate(90, 0, 0)

	// After 90° rotation, (100, 0) should become approximately (0, 100)
	epsilon := 0.0001
	if math.Abs(rotated.X1) > epsilon || math.Abs(rotated.Y1-100) > epsilon {
		t.Errorf("Expected start point near (0, 100), got (%f, %f)", rotated.X1, rotated.Y1)
	}
}

func TestLineScale(t *testing.T) {
	line := NewLine(0, 0, 100, 100)
	scaled := line.Scale(2.0, 0, 0)

	if scaled.X2 != 200 || scaled.Y2 != 200 {
		t.Errorf("Expected end point (200, 200), got (%f, %f)", scaled.X2, scaled.Y2)
	}
}

func TestCircleTranslate(t *testing.T) {
	circle := NewCircle(50, 50, 25)
	moved := circle.Translate(100, 100)

	if moved.CenterX != 150 || moved.CenterY != 150 {
		t.Errorf("Expected center (150, 150), got (%f, %f)", moved.CenterX, moved.CenterY)
	}
	if moved.Radius != 25 {
		t.Errorf("Expected radius 25, got %f", moved.Radius)
	}
}

func TestCircleScale(t *testing.T) {
	circle := NewCircle(50, 50, 25)
	scaled := circle.Scale(2.0)

	if scaled.Radius != 50 {
		t.Errorf("Expected radius 50, got %f", scaled.Radius)
	}
}

func TestArcTranslate(t *testing.T) {
	arc := NewArc(50, 50, 25, 0, 90)
	moved := arc.Translate(100, 100)

	if moved.CenterX != 150 || moved.CenterY != 150 {
		t.Errorf("Expected center (150, 150), got (%f, %f)", moved.CenterX, moved.CenterY)
	}
	if moved.StartAngle != 0 || moved.EndAngle != 90 {
		t.Errorf("Expected angles (0, 90), got (%f, %f)", moved.StartAngle, moved.EndAngle)
	}
}

func TestPointTranslate(t *testing.T) {
	point := NewPoint(100, 200)
	moved := point.Translate(50, 50)

	if moved.X != 150 || moved.Y != 250 {
		t.Errorf("Expected point (150, 250), got (%f, %f)", moved.X, moved.Y)
	}
}

func TestTextTranslate(t *testing.T) {
	text := NewText(10, 10, "Hello")
	moved := text.Translate(50, 50)

	if moved.X != 60 || moved.Y != 60 {
		t.Errorf("Expected position (60, 60), got (%f, %f)", moved.X, moved.Y)
	}
}

func TestTextRotate(t *testing.T) {
	text := NewText(10, 10, "Hello", WithTextRotation(0))
	rotated := text.Rotate(45)

	if rotated.Rotation != 45 {
		t.Errorf("Expected rotation 45, got %f", rotated.Rotation)
	}
}

func TestTextScale(t *testing.T) {
	text := NewText(10, 10, "Hello", WithTextHeight(5))
	scaled := text.Scale(2.0)

	if scaled.Height != 10 {
		t.Errorf("Expected height 10, got %f", scaled.Height)
	}
}

func TestSolidTranslate(t *testing.T) {
	solid := NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
	moved := solid.Translate(50, 50)

	if moved.X1 != 50 || moved.Y1 != 50 {
		t.Errorf("Expected point 1 (50, 50), got (%f, %f)", moved.X1, moved.Y1)
	}
}

func TestSolidScale(t *testing.T) {
	solid := NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
	scaled := solid.Scale(2.0, 0, 0)

	if scaled.X2 != 200 || scaled.Y2 != 0 {
		t.Errorf("Expected point 2 (200, 0), got (%f, %f)", scaled.X2, scaled.Y2)
	}
}

func TestInsertTranslate(t *testing.T) {
	insert := NewInsert("MyBlock", 100, 100)
	moved := insert.Translate(50, 50)

	if moved.X != 150 || moved.Y != 150 {
		t.Errorf("Expected position (150, 150), got (%f, %f)", moved.X, moved.Y)
	}
}

func TestInsertRotate(t *testing.T) {
	insert := NewInsert("MyBlock", 100, 100)
	rotated := insert.Rotate(45)

	if rotated.Rotation != 45 {
		t.Errorf("Expected rotation 45, got %f", rotated.Rotation)
	}
}

func TestInsertScale(t *testing.T) {
	insert := NewInsert("MyBlock", 100, 100, WithInsertScale(1.0, 1.0))
	scaled := insert.Scale(2.0)

	if scaled.ScaleX != 2.0 || scaled.ScaleY != 2.0 {
		t.Errorf("Expected scale (2.0, 2.0), got (%f, %f)", scaled.ScaleX, scaled.ScaleY)
	}
}
