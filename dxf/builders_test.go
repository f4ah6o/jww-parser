package dxf

import (
	"testing"
)

func TestNewLine(t *testing.T) {
	line := NewLine(0, 0, 100, 100)
	if line.X1 != 0 || line.Y1 != 0 || line.X2 != 100 || line.Y2 != 100 {
		t.Errorf("NewLine coordinates mismatch")
	}
	if line.Layer != "0" {
		t.Errorf("Expected default layer '0', got '%s'", line.Layer)
	}
	if line.Color != 0 {
		t.Errorf("Expected default color 0, got %d", line.Color)
	}
}

func TestNewLineWithOptions(t *testing.T) {
	line := NewLine(0, 0, 100, 100,
		WithLineLayer("MyLayer"),
		WithLineColor(5),
		WithLineType("DASHED"))

	if line.Layer != "MyLayer" {
		t.Errorf("Expected layer 'MyLayer', got '%s'", line.Layer)
	}
	if line.Color != 5 {
		t.Errorf("Expected color 5, got %d", line.Color)
	}
	if line.LineType != "DASHED" {
		t.Errorf("Expected line type 'DASHED', got '%s'", line.LineType)
	}
}

func TestNewCircle(t *testing.T) {
	circle := NewCircle(50, 50, 25)
	if circle.CenterX != 50 || circle.CenterY != 50 || circle.Radius != 25 {
		t.Errorf("NewCircle parameters mismatch")
	}
	if circle.Layer != "0" {
		t.Errorf("Expected default layer '0', got '%s'", circle.Layer)
	}
}

func TestNewCircleWithOptions(t *testing.T) {
	circle := NewCircle(50, 50, 25,
		WithCircleLayer("MyLayer"),
		WithCircleColor(3))

	if circle.Layer != "MyLayer" {
		t.Errorf("Expected layer 'MyLayer', got '%s'", circle.Layer)
	}
	if circle.Color != 3 {
		t.Errorf("Expected color 3, got %d", circle.Color)
	}
}

func TestNewArc(t *testing.T) {
	arc := NewArc(50, 50, 25, 0, 90)
	if arc.CenterX != 50 || arc.CenterY != 50 || arc.Radius != 25 {
		t.Errorf("NewArc parameters mismatch")
	}
	if arc.StartAngle != 0 || arc.EndAngle != 90 {
		t.Errorf("NewArc angles mismatch")
	}
}

func TestNewPoint(t *testing.T) {
	point := NewPoint(100, 200)
	if point.X != 100 || point.Y != 200 {
		t.Errorf("NewPoint coordinates mismatch")
	}
}

func TestNewText(t *testing.T) {
	text := NewText(10, 10, "Hello World")
	if text.Content != "Hello World" {
		t.Errorf("Expected content 'Hello World', got '%s'", text.Content)
	}
	if text.X != 10 || text.Y != 10 {
		t.Errorf("NewText coordinates mismatch")
	}
	if text.Height != 2.5 {
		t.Errorf("Expected default height 2.5, got %f", text.Height)
	}
}

func TestNewTextWithOptions(t *testing.T) {
	text := NewText(10, 10, "Hello",
		WithTextHeight(5.0),
		WithTextRotation(45),
		WithTextStyle("MyStyle"))

	if text.Height != 5.0 {
		t.Errorf("Expected height 5.0, got %f", text.Height)
	}
	if text.Rotation != 45 {
		t.Errorf("Expected rotation 45, got %f", text.Rotation)
	}
	if text.Style != "MyStyle" {
		t.Errorf("Expected style 'MyStyle', got '%s'", text.Style)
	}
}

func TestNewSolid(t *testing.T) {
	solid := NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
	if solid.X1 != 0 || solid.Y1 != 0 {
		t.Errorf("NewSolid point 1 mismatch")
	}
	if solid.X2 != 100 || solid.Y2 != 0 {
		t.Errorf("NewSolid point 2 mismatch")
	}
}

func TestNewInsert(t *testing.T) {
	insert := NewInsert("MyBlock", 100, 100)
	if insert.BlockName != "MyBlock" {
		t.Errorf("Expected block name 'MyBlock', got '%s'", insert.BlockName)
	}
	if insert.X != 100 || insert.Y != 100 {
		t.Errorf("NewInsert coordinates mismatch")
	}
	if insert.ScaleX != 1.0 || insert.ScaleY != 1.0 {
		t.Errorf("Expected default scale 1.0, got (%f, %f)", insert.ScaleX, insert.ScaleY)
	}
}

func TestNewInsertWithOptions(t *testing.T) {
	insert := NewInsert("MyBlock", 100, 100,
		WithInsertScale(2.0, 3.0),
		WithInsertRotation(45))

	if insert.ScaleX != 2.0 || insert.ScaleY != 3.0 {
		t.Errorf("Expected scale (2.0, 3.0), got (%f, %f)", insert.ScaleX, insert.ScaleY)
	}
	if insert.Rotation != 45 {
		t.Errorf("Expected rotation 45, got %f", insert.Rotation)
	}
}
