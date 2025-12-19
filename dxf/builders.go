package dxf

// EntityOption is a function that configures entity properties.
// This pattern allows for flexible, readable entity construction.
type EntityOption func(interface{})

// LineOption configures Line entity properties.
type LineOption func(*Line)

// WithLineLayer sets the layer for a Line entity.
func WithLineLayer(layer string) LineOption {
	return func(l *Line) {
		l.Layer = layer
	}
}

// WithLineColor sets the color for a Line entity.
func WithLineColor(color int) LineOption {
	return func(l *Line) {
		l.Color = color
	}
}

// WithLineType sets the line type for a Line entity.
func WithLineType(lineType string) LineOption {
	return func(l *Line) {
		l.LineType = lineType
	}
}

// NewLine creates a new Line entity with the given coordinates.
// Optional LineOption functions can customize the line properties.
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 100,
//		dxf.WithLineLayer("MyLayer"),
//		dxf.WithLineColor(1))
func NewLine(x1, y1, x2, y2 float64, opts ...LineOption) *Line {
	line := &Line{
		Layer:    "0",
		Color:    0, // BYLAYER
		X1:       x1,
		Y1:       y1,
		X2:       x2,
		Y2:       y2,
		LineType: "CONTINUOUS",
	}
	for _, opt := range opts {
		opt(line)
	}
	return line
}

// CircleOption configures Circle entity properties.
type CircleOption func(*Circle)

// WithCircleLayer sets the layer for a Circle entity.
func WithCircleLayer(layer string) CircleOption {
	return func(c *Circle) {
		c.Layer = layer
	}
}

// WithCircleColor sets the color for a Circle entity.
func WithCircleColor(color int) CircleOption {
	return func(c *Circle) {
		c.Color = color
	}
}

// NewCircle creates a new Circle entity with the given center and radius.
// Optional CircleOption functions can customize the circle properties.
//
// Example:
//
//	circle := dxf.NewCircle(50, 50, 25,
//		dxf.WithCircleLayer("MyLayer"),
//		dxf.WithCircleColor(2))
func NewCircle(centerX, centerY, radius float64, opts ...CircleOption) *Circle {
	circle := &Circle{
		Layer:   "0",
		Color:   0, // BYLAYER
		CenterX: centerX,
		CenterY: centerY,
		Radius:  radius,
	}
	for _, opt := range opts {
		opt(circle)
	}
	return circle
}

// ArcOption configures Arc entity properties.
type ArcOption func(*Arc)

// WithArcLayer sets the layer for an Arc entity.
func WithArcLayer(layer string) ArcOption {
	return func(a *Arc) {
		a.Layer = layer
	}
}

// WithArcColor sets the color for an Arc entity.
func WithArcColor(color int) ArcOption {
	return func(a *Arc) {
		a.Color = color
	}
}

// NewArc creates a new Arc entity with the given center, radius, and angles.
// Angles are in degrees. Optional ArcOption functions can customize the arc properties.
//
// Example:
//
//	arc := dxf.NewArc(50, 50, 25, 0, 90,
//		dxf.WithArcLayer("MyLayer"),
//		dxf.WithArcColor(3))
func NewArc(centerX, centerY, radius, startAngle, endAngle float64, opts ...ArcOption) *Arc {
	arc := &Arc{
		Layer:      "0",
		Color:      0, // BYLAYER
		CenterX:    centerX,
		CenterY:    centerY,
		Radius:     radius,
		StartAngle: startAngle,
		EndAngle:   endAngle,
	}
	for _, opt := range opts {
		opt(arc)
	}
	return arc
}

// PointOption configures Point entity properties.
type PointOption func(*Point)

// WithPointLayer sets the layer for a Point entity.
func WithPointLayer(layer string) PointOption {
	return func(p *Point) {
		p.Layer = layer
	}
}

// WithPointColor sets the color for a Point entity.
func WithPointColor(color int) PointOption {
	return func(p *Point) {
		p.Color = color
	}
}

// NewPoint creates a new Point entity with the given coordinates.
// Optional PointOption functions can customize the point properties.
//
// Example:
//
//	point := dxf.NewPoint(100, 200,
//		dxf.WithPointLayer("MyLayer"),
//		dxf.WithPointColor(4))
func NewPoint(x, y float64, opts ...PointOption) *Point {
	point := &Point{
		Layer: "0",
		Color: 0, // BYLAYER
		X:     x,
		Y:     y,
	}
	for _, opt := range opts {
		opt(point)
	}
	return point
}

// TextOption configures Text entity properties.
type TextOption func(*Text)

// WithTextLayer sets the layer for a Text entity.
func WithTextLayer(layer string) TextOption {
	return func(t *Text) {
		t.Layer = layer
	}
}

// WithTextColor sets the color for a Text entity.
func WithTextColor(color int) TextOption {
	return func(t *Text) {
		t.Color = color
	}
}

// WithTextHeight sets the height for a Text entity.
func WithTextHeight(height float64) TextOption {
	return func(t *Text) {
		t.Height = height
	}
}

// WithTextRotation sets the rotation angle (in degrees) for a Text entity.
func WithTextRotation(rotation float64) TextOption {
	return func(t *Text) {
		t.Rotation = rotation
	}
}

// WithTextStyle sets the text style for a Text entity.
func WithTextStyle(style string) TextOption {
	return func(t *Text) {
		t.Style = style
	}
}

// NewText creates a new Text entity with the given position and content.
// Optional TextOption functions can customize the text properties.
//
// Example:
//
//	text := dxf.NewText(10, 10, "Hello World",
//		dxf.WithTextLayer("MyLayer"),
//		dxf.WithTextHeight(5.0),
//		dxf.WithTextRotation(45))
func NewText(x, y float64, content string, opts ...TextOption) *Text {
	text := &Text{
		Layer:    "0",
		Color:    0, // BYLAYER
		X:        x,
		Y:        y,
		Height:   2.5, // Default height
		Rotation: 0,
		Content:  content,
		Style:    "STANDARD",
	}
	for _, opt := range opts {
		opt(text)
	}
	return text
}

// SolidOption configures Solid entity properties.
type SolidOption func(*Solid)

// WithSolidLayer sets the layer for a Solid entity.
func WithSolidLayer(layer string) SolidOption {
	return func(s *Solid) {
		s.Layer = layer
	}
}

// WithSolidColor sets the color for a Solid entity.
func WithSolidColor(color int) SolidOption {
	return func(s *Solid) {
		s.Color = color
	}
}

// NewSolid creates a new Solid entity (filled polygon) with the given corner points.
// For triangles, set p4x and p4y equal to p3x and p3y.
// Optional SolidOption functions can customize the solid properties.
//
// Example:
//
//	// Triangle
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100,
//		dxf.WithSolidLayer("MyLayer"),
//		dxf.WithSolidColor(5))
func NewSolid(p1x, p1y, p2x, p2y, p3x, p3y, p4x, p4y float64, opts ...SolidOption) *Solid {
	solid := &Solid{
		Layer: "0",
		Color: 0, // BYLAYER
		X1:    p1x,
		Y1:    p1y,
		X2:    p2x,
		Y2:    p2y,
		X3:    p3x,
		Y3:    p3y,
		X4:    p4x,
		Y4:    p4y,
	}
	for _, opt := range opts {
		opt(solid)
	}
	return solid
}

// InsertOption configures Insert entity properties.
type InsertOption func(*Insert)

// WithInsertLayer sets the layer for an Insert entity.
func WithInsertLayer(layer string) InsertOption {
	return func(i *Insert) {
		i.Layer = layer
	}
}

// WithInsertColor sets the color for an Insert entity.
func WithInsertColor(color int) InsertOption {
	return func(i *Insert) {
		i.Color = color
	}
}

// WithInsertScale sets the scale factors for an Insert entity.
func WithInsertScale(scaleX, scaleY float64) InsertOption {
	return func(i *Insert) {
		i.ScaleX = scaleX
		i.ScaleY = scaleY
	}
}

// WithInsertRotation sets the rotation angle (in degrees) for an Insert entity.
func WithInsertRotation(rotation float64) InsertOption {
	return func(i *Insert) {
		i.Rotation = rotation
	}
}

// NewInsert creates a new Insert entity (block reference) with the given block name and position.
// Optional InsertOption functions can customize the insert properties.
//
// Example:
//
//	insert := dxf.NewInsert("MyBlock", 100, 100,
//		dxf.WithInsertLayer("MyLayer"),
//		dxf.WithInsertScale(2.0, 2.0),
//		dxf.WithInsertRotation(45))
func NewInsert(blockName string, x, y float64, opts ...InsertOption) *Insert {
	insert := &Insert{
		Layer:     "0",
		Color:     0, // BYLAYER
		BlockName: blockName,
		X:         x,
		Y:         y,
		ScaleX:    1.0,
		ScaleY:    1.0,
		Rotation:  0,
	}
	for _, opt := range opts {
		opt(insert)
	}
	return insert
}
