package jww

// Document represents a complete JWW (Jw_cad) file structure.
// JWW files are binary CAD files used by Jw_cad, a popular Japanese CAD software.
// The document contains layer information, drawing entities, and optional block definitions.
type Document struct {
	// Version indicates the JWW file format version (e.g., 351 for Ver.3.51, 420 for Ver.4.20).
	Version uint32

	// Memo is the file memo/description stored in the JWW header.
	Memo string

	// PaperSize specifies the paper size: 0-4 for A0-A4, 8 for 2A, 9 for 3A, etc.
	PaperSize uint32

	// WriteLayerGroup is the currently active layer group for writing (0-15).
	WriteLayerGroup uint32

	// LayerGroups contains 16 layer groups, each with 16 layers.
	// This provides a total of 256 possible layers organized in a hierarchical structure.
	LayerGroups [16]LayerGroup

	// Entities contains all drawing entities (lines, arcs, text, etc.) in the file.
	Entities []Entity

	// BlockDefs contains block definitions that can be referenced by block insert entities.
	BlockDefs []BlockDef
}

// LayerGroup represents a layer group (レイヤグループ) in a JWW file.
// JWW organizes layers into 16 groups, with each group containing 16 layers.
// Each layer group can have its own display state, scale, and protection settings.
type LayerGroup struct {
	// State indicates the layer group's visibility and editability:
	// 0: hidden, 1: display only, 2: editable, 3: write mode
	State uint32

	// WriteLayer is the currently active layer for writing within this group (0-15).
	WriteLayer uint32

	// Scale is the scale denominator for this layer group (e.g., 100.0 for 1:100).
	Scale float64

	// Protect is the protection flag to prevent accidental modifications.
	Protect uint32

	// Layers contains the 16 layers within this layer group.
	Layers [16]Layer

	// Name is the user-defined name of this layer group.
	Name string
}

// Layer represents an individual layer within a layer group.
// Layers are used to organize drawing entities by type, discipline, or other criteria.
type Layer struct {
	// State indicates the layer's visibility and editability:
	// 0: hidden, 1: display only, 2: editable, 3: write mode
	State uint32

	// Protect is the protection flag to prevent accidental modifications.
	Protect uint32

	// Name is the user-defined name of this layer.
	Name string
}

// EntityBase contains common attributes shared by all JWW drawing entities.
// These attributes control appearance properties like line type, color, and layer assignment.
type EntityBase struct {
	// Group is the curve attribute number (線種グループ).
	Group uint32

	// PenStyle is the line type number (線種).
	PenStyle byte

	// PenColor is the line color number (1-9 for basic colors, extended values for SXF colors).
	PenColor uint16

	// PenWidth is the line width in internal units (available in Ver.3.51 and later).
	PenWidth uint16

	// Layer is the layer number within the layer group (0-15).
	Layer uint16

	// LayerGroup is the layer group number (0-15).
	LayerGroup uint16

	// Flag contains various attribute flags for the entity.
	Flag uint16
}

// Entity is the interface implemented by all JWW drawing entities.
// Each entity type (Line, Arc, Point, Text, etc.) must provide access to its
// base attributes and identify its type.
type Entity interface {
	// Base returns a pointer to the common EntityBase attributes.
	Base() *EntityBase

	// Type returns the entity type name (e.g., "LINE", "ARC", "TEXT").
	Type() string
}

// Line represents a straight line segment entity (JWW class: CDataSen).
// Lines are defined by their start and end points in 2D space.
type Line struct {
	EntityBase

	// StartX is the X coordinate of the line's starting point.
	StartX, StartY float64

	// EndX is the X coordinate of the line's ending point.
	EndX, EndY float64
}

// Base returns the entity's base attributes.
func (l *Line) Base() *EntityBase { return &l.EntityBase }

// Type returns "LINE".
func (l *Line) Type() string { return "LINE" }

// Arc represents an arc or circle entity (JWW class: CDataEnko).
// Arcs can be circular or elliptical, and may represent a full circle or a partial arc.
type Arc struct {
	EntityBase

	// CenterX is the X coordinate of the arc's center point.
	CenterX, CenterY float64

	// Radius is the arc's radius (for the major axis in the case of ellipses).
	Radius float64

	// StartAngle is the starting angle in radians.
	StartAngle float64

	// ArcAngle is the angular extent of the arc in radians.
	ArcAngle float64

	// TiltAngle is the rotation angle of the ellipse major axis in radians (0 for circles).
	TiltAngle float64

	// Flatness is the ratio of minor to major axis (1.0 for circles, <1.0 or >1.0 for ellipses).
	Flatness float64

	// IsFullCircle indicates whether this represents a complete circle/ellipse.
	IsFullCircle bool
}

// Base returns the entity's base attributes.
func (a *Arc) Base() *EntityBase { return &a.EntityBase }

// Type returns "CIRCLE" for full circles, "ARC" otherwise.
func (a *Arc) Type() string {
	if a.IsFullCircle {
		return "CIRCLE"
	}
	return "ARC"
}

// Point represents a point entity (JWW class: CDataTen).
// Points can be temporary construction points or permanent marker points.
type Point struct {
	EntityBase

	// X is the X coordinate of the point.
	X, Y float64

	// IsTemporary indicates if this is a temporary construction point (仮点).
	IsTemporary bool

	// Code specifies the point marker type (arrow, cross, circle, etc.).
	Code uint32

	// Angle is the rotation angle for directional point markers.
	Angle float64

	// Scale is the size scale factor for the point marker.
	Scale float64
}

// Base returns the entity's base attributes.
func (p *Point) Base() *EntityBase { return &p.EntityBase }

// Type returns "POINT".
func (p *Point) Type() string { return "POINT" }

// Text represents a text entity (JWW class: CDataMoji).
// Text can be single-line or multi-line, with support for various fonts and styles.
type Text struct {
	EntityBase

	// StartX is the X coordinate of the text's starting point.
	StartX, StartY float64

	// EndX is the X coordinate of the text's ending point (for text box).
	EndX, EndY float64

	// TextType contains text style flags: +10000 for italic, +20000 for bold.
	TextType uint32

	// SizeX is the character width.
	SizeX, SizeY float64

	// Spacing is the character spacing factor.
	Spacing float64

	// Angle is the text rotation angle in degrees.
	Angle float64

	// FontName is the name of the font to use.
	FontName string

	// Content is the actual text content (Shift-JIS encoded in file, converted to UTF-8).
	Content string
}

// Base returns the entity's base attributes.
func (t *Text) Base() *EntityBase { return &t.EntityBase }

// Type returns "TEXT".
func (t *Text) Type() string { return "TEXT" }

// Solid represents a solid fill entity (JWW class: CDataSolid).
// Solids are filled quadrilaterals or triangles used for hatching and shading.
type Solid struct {
	EntityBase

	// Point1X is the X coordinate of the first corner point.
	Point1X, Point1Y float64

	// Point2X is the X coordinate of the second corner point.
	Point2X, Point2Y float64

	// Point3X is the X coordinate of the third corner point.
	Point3X, Point3Y float64

	// Point4X is the X coordinate of the fourth corner point.
	Point4X, Point4Y float64

	// Color is the RGB color value (used when PenColor == 10).
	Color uint32
}

// Base returns the entity's base attributes.
func (s *Solid) Base() *EntityBase { return &s.EntityBase }

// Type returns "SOLID".
func (s *Solid) Type() string { return "SOLID" }

// Block represents a block insert entity (JWW class: CDataBlock).
// Blocks allow reuse of geometry defined in a BlockDef.
type Block struct {
	EntityBase

	// RefX is the X coordinate of the insertion reference point.
	RefX, RefY float64

	// ScaleX is the X-axis scale factor.
	ScaleX float64

	// ScaleY is the Y-axis scale factor.
	ScaleY float64

	// Rotation is the rotation angle in radians.
	Rotation float64

	// DefNumber is the block definition number to reference.
	DefNumber uint32
}

// Base returns the entity's base attributes.
func (b *Block) Base() *EntityBase { return &b.EntityBase }

// Type returns "BLOCK".
func (b *Block) Type() string { return "BLOCK" }

// BlockDef represents a block definition (JWW class: CDataList).
// Block definitions are reusable collections of entities that can be inserted
// multiple times via Block entities.
type BlockDef struct {
	EntityBase

	// Number is the unique block definition identifier.
	Number uint32

	// IsReferenced indicates whether this block is used by any Block entity.
	IsReferenced bool

	// Name is the user-defined name of this block.
	Name string

	// Entities contains the drawing entities that comprise this block.
	Entities []Entity
}
