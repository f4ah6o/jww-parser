// Package dxf provides types and generation functions for the DXF (Drawing Exchange Format) file format.
//
// DXF is an ASCII-based CAD data file format developed by Autodesk for enabling
// data interoperability between AutoCAD and other programs.
//
// This package provides:
//   - DXF document structure representation
//   - Entity types (Line, Arc, Circle, Text, etc.)
//   - Layer and block definitions
//   - DXF file writing capabilities
//
// Basic usage:
//
//	doc := &dxf.Document{
//	    Layers: []dxf.Layer{
//	        {Name: "0", Color: 7, LineType: "CONTINUOUS"},
//	    },
//	    Entities: []dxf.Entity{
//	        &dxf.Line{Layer: "0", X1: 0, Y1: 0, X2: 100, Y2: 100},
//	    },
//	}
//
//	w := dxf.NewWriter(outputFile)
//	w.WriteDocument(doc)
package dxf

// Document represents a complete DXF document structure.
// It contains layer definitions, drawing entities, and optional block definitions.
type Document struct {
	// Layers contains the layer definitions used by entities.
	Layers []Layer

	// Entities contains all drawing entities in the document.
	Entities []Entity

	// Blocks contains reusable block definitions.
	Blocks []Block
}

// Layer represents a DXF layer definition.
// Layers are used to organize entities by grouping related objects together.
type Layer struct {
	// Name is the layer name (e.g., "0" for the default layer).
	Name string

	// Color is the AutoCAD Color Index (ACI) value (1-255).
	// Common values: 1=red, 2=yellow, 3=green, 4=cyan, 5=blue, 6=magenta, 7=white/black.
	Color int

	// LineType specifies the line pattern (e.g., "CONTINUOUS", "DASHED").
	LineType string

	// Frozen indicates if the layer is frozen (not visible and not printable).
	Frozen bool

	// Locked indicates if the layer is locked (visible but not editable).
	Locked bool
}

// Entity is the interface implemented by all DXF drawing entities.
// Each entity must provide its type name and group code representation.
type Entity interface {
	// EntityType returns the DXF entity type name (e.g., "LINE", "CIRCLE", "TEXT").
	EntityType() string

	// GroupCodes returns the entity's data as DXF group code/value pairs.
	GroupCodes() []GroupCode
}

// GroupCode represents a DXF group code and its associated value.
// DXF files are structured as pairs of group codes (integers) and values.
// Group codes indicate the type of data element (e.g., 0=entity type, 10=X coordinate, 8=layer name).
type GroupCode struct {
	// Code is the DXF group code integer (0-999).
	Code int

	// Value is the associated value (string, int, or float64).
	Value interface{}
}

// Line represents a DXF LINE entity.
// A line is defined by two points in 2D or 3D space.
type Line struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER, 1-255 = specific colors).
	Color int

	// LineType specifies the line pattern (e.g., "CONTINUOUS", "DASHED").
	LineType string

	// X1, Y1 are the coordinates of the line's start point.
	X1, Y1 float64

	// X2, Y2 are the coordinates of the line's end point.
	X2, Y2 float64
}

// EntityType returns "LINE".
func (l *Line) EntityType() string { return "LINE" }

// GroupCodes returns the DXF group codes for this line entity.
func (l *Line) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "LINE"},
		{8, l.Layer},
		{62, l.Color},
		{6, l.LineType},
		{10, l.X1},
		{20, l.Y1},
		{30, 0.0},
		{11, l.X2},
		{21, l.Y2},
		{31, 0.0},
	}
}

// Circle represents a DXF CIRCLE entity.
// A circle is defined by its center point and radius.
type Circle struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern for the circle outline.
	LineType string

	// CenterX, CenterY are the coordinates of the circle's center point.
	CenterX float64
	CenterY float64

	// Radius is the circle's radius.
	Radius float64
}

// EntityType returns "CIRCLE".
func (c *Circle) EntityType() string { return "CIRCLE" }

// GroupCodes returns the DXF group codes for this circle entity.
func (c *Circle) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "CIRCLE"},
		{8, c.Layer},
		{62, c.Color},
		{6, c.LineType},
		{10, c.CenterX},
		{20, c.CenterY},
		{30, 0.0},
		{40, c.Radius},
	}
}

// Arc represents a DXF ARC entity.
// An arc is a portion of a circle defined by center, radius, and start/end angles.
type Arc struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern for the arc.
	LineType string

	// CenterX, CenterY are the coordinates of the arc's center point.
	CenterX float64
	CenterY float64

	// Radius is the arc's radius.
	Radius float64

	// StartAngle is the starting angle in degrees (0-360).
	StartAngle float64

	// EndAngle is the ending angle in degrees (0-360).
	EndAngle float64
}

// EntityType returns "ARC".
func (a *Arc) EntityType() string { return "ARC" }

func (a *Arc) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "ARC"},
		{8, a.Layer},
		{62, a.Color},
		{6, a.LineType},
		{10, a.CenterX},
		{20, a.CenterY},
		{30, 0.0},
		{40, a.Radius},
		{50, a.StartAngle},
		{51, a.EndAngle},
	}
}

// Ellipse represents a DXF ELLIPSE entity.
// An ellipse is defined by center point, major/minor axes, and optional start/end parameters for partial ellipses.
type Ellipse struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern for the ellipse.
	LineType string

	// CenterX, CenterY are the coordinates of the ellipse's center point.
	CenterX float64
	CenterY float64

	// MajorAxisX, MajorAxisY are the endpoint of the major axis relative to the center.
	MajorAxisX float64
	MajorAxisY float64

	// MinorRatio is the ratio of minor axis to major axis (0.0 to 1.0).
	MinorRatio float64

	// StartParam is the start parameter in radians (0.0 for full ellipse).
	StartParam float64

	// EndParam is the end parameter in radians (2*PI for full ellipse).
	EndParam float64
}

// EntityType returns "ELLIPSE".
func (e *Ellipse) EntityType() string { return "ELLIPSE" }

func (e *Ellipse) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "ELLIPSE"},
		{8, e.Layer},
		{62, e.Color},
		{6, e.LineType},
		{10, e.CenterX},
		{20, e.CenterY},
		{30, 0.0},
		{11, e.MajorAxisX},
		{21, e.MajorAxisY},
		{31, 0.0},
		{40, e.MinorRatio},
		{41, e.StartParam},
		{42, e.EndParam},
	}
}

// Point represents a DXF POINT entity.
// A point is a single location in 2D or 3D space.
type Point struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern for the point marker.
	LineType string

	// X, Y are the coordinates of the point.
	X, Y float64
}

// EntityType returns "POINT".
func (p *Point) EntityType() string { return "POINT" }

// GroupCodes returns the DXF group codes for this point entity.
func (p *Point) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "POINT"},
		{8, p.Layer},
		{62, p.Color},
		{6, p.LineType},
		{10, p.X},
		{20, p.Y},
		{30, 0.0},
	}
}

// Text represents a DXF TEXT entity.
// Text entities display a single line of text at a specified location.
type Text struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern applied to the text entity.
	LineType string

	// X, Y are the coordinates of the text insertion point.
	X, Y float64

	// Height is the text height in drawing units.
	Height float64

	// Rotation is the text rotation angle in degrees.
	Rotation float64

	// Content is the actual text string to display.
	Content string

	// Style is the text style name (e.g., "STANDARD").
	Style string
}

// EntityType returns "TEXT".
func (t *Text) EntityType() string { return "TEXT" }

func (t *Text) GroupCodes() []GroupCode {
	codes := []GroupCode{
		{0, "TEXT"},
		{8, EscapeUnicode(t.Layer)},
		{62, t.Color},
		{6, t.LineType},
		{10, t.X},
		{20, t.Y},
		{30, 0.0},
		{40, t.Height},
		{1, EscapeUnicode(t.Content)},
	}
	if t.Rotation != 0 {
		codes = append(codes, GroupCode{50, t.Rotation})
	}
	if t.Style != "" {
		codes = append(codes, GroupCode{7, t.Style})
	}
	return codes
}

// Solid represents a DXF SOLID entity (filled triangle or quadrilateral).
// Solids are used to create filled areas and hatching patterns.
type Solid struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern applied to the solid's outline.
	LineType string

	// X1, Y1 are the coordinates of the first corner point.
	X1, Y1 float64

	// X2, Y2 are the coordinates of the second corner point.
	X2, Y2 float64

	// X3, Y3 are the coordinates of the third corner point.
	X3, Y3 float64

	// X4, Y4 are the coordinates of the fourth corner point (same as X3, Y3 for triangles).
	X4, Y4 float64
}

// EntityType returns "SOLID".
func (s *Solid) EntityType() string { return "SOLID" }

// GroupCodes returns the DXF group codes for this solid entity.
func (s *Solid) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "SOLID"},
		{8, s.Layer},
		{62, s.Color},
		{6, s.LineType},
		{10, s.X1},
		{20, s.Y1},
		{30, 0.0},
		{11, s.X2},
		{21, s.Y2},
		{31, 0.0},
		{12, s.X3},
		{22, s.Y3},
		{32, 0.0},
		{13, s.X4},
		{23, s.Y4},
		{33, 0.0},
	}
}

// Insert represents a DXF INSERT entity (block reference).
// Inserts allow reusing block definitions with different positions, scales, and rotations.
type Insert struct {
	// Layer is the name of the layer this entity belongs to.
	Layer string

	// Color is the ACI color number (0 = BYLAYER).
	Color int

	// LineType specifies the line pattern applied to the insert reference.
	LineType string

	// BlockName is the name of the block definition to insert.
	BlockName string

	// X, Y are the coordinates of the insertion point.
	X, Y float64

	// ScaleX is the X-axis scale factor.
	ScaleX float64

	// ScaleY is the Y-axis scale factor.
	ScaleY float64

	// Rotation is the rotation angle in degrees.
	Rotation float64
}

// EntityType returns "INSERT".
func (i *Insert) EntityType() string { return "INSERT" }

// GroupCodes returns the DXF group codes for this insert entity.
func (i *Insert) GroupCodes() []GroupCode {
	return []GroupCode{
		{0, "INSERT"},
		{8, i.Layer},
		{62, i.Color},
		{6, i.LineType},
		{2, i.BlockName},
		{10, i.X},
		{20, i.Y},
		{30, 0.0},
		{41, i.ScaleX},
		{42, i.ScaleY},
		{43, 1.0}, // ScaleZ
		{50, i.Rotation},
	}
}

// Block represents a DXF block definition.
// Blocks are reusable collections of entities that can be inserted multiple times
// via Insert entities with different transformations.
type Block struct {
	// Name is the unique block name.
	Name string

	// BaseX, BaseY are the coordinates of the block's base point.
	BaseX float64
	BaseY float64

	// Entities contains the entities that comprise this block.
	Entities []Entity
}
