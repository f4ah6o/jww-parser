package dxf

import "math"

// Length calculates the length of a Line entity.
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 100)
//	length := line.Length() // Returns ~141.42
func (l *Line) Length() float64 {
	dx := l.X2 - l.X1
	dy := l.Y2 - l.Y1
	return math.Sqrt(dx*dx + dy*dy)
}

// BoundingBox returns the bounding box of a Line entity.
// Returns (minX, minY, maxX, maxY).
//
// Example:
//
//	line := dxf.NewLine(10, 20, 100, 200)
//	minX, minY, maxX, maxY := line.BoundingBox() // Returns (10, 20, 100, 200)
func (l *Line) BoundingBox() (minX, minY, maxX, maxY float64) {
	minX = math.Min(l.X1, l.X2)
	maxX = math.Max(l.X1, l.X2)
	minY = math.Min(l.Y1, l.Y2)
	maxY = math.Max(l.Y1, l.Y2)
	return
}

// MidPoint returns the middle point of a Line entity.
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 100)
//	x, y := line.MidPoint() // Returns (50, 50)
func (l *Line) MidPoint() (x, y float64) {
	return (l.X1 + l.X2) / 2, (l.Y1 + l.Y2) / 2
}

// Angle returns the angle of the line in degrees (0-360).
// 0 degrees is to the right (positive X axis).
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 100)
//	angle := line.Angle() // Returns 45.0
func (l *Line) Angle() float64 {
	angle := math.Atan2(l.Y2-l.Y1, l.X2-l.X1) * 180.0 / math.Pi
	if angle < 0 {
		angle += 360
	}
	return angle
}

// Area calculates the area of a Circle entity.
//
// Example:
//
//	circle := dxf.NewCircle(50, 50, 25)
//	area := circle.Area() // Returns π * 25²
func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Circumference calculates the circumference of a Circle entity.
//
// Example:
//
//	circle := dxf.NewCircle(50, 50, 25)
//	circ := circle.Circumference() // Returns 2π * 25
func (c *Circle) Circumference() float64 {
	return 2 * math.Pi * c.Radius
}

// BoundingBox returns the bounding box of a Circle entity.
// Returns (minX, minY, maxX, maxY).
//
// Example:
//
//	circle := dxf.NewCircle(50, 50, 25)
//	minX, minY, maxX, maxY := circle.BoundingBox() // Returns (25, 25, 75, 75)
func (c *Circle) BoundingBox() (minX, minY, maxX, maxY float64) {
	return c.CenterX - c.Radius, c.CenterY - c.Radius,
		c.CenterX + c.Radius, c.CenterY + c.Radius
}

// ArcLength calculates the arc length of an Arc entity.
//
// Example:
//
//	arc := dxf.NewArc(50, 50, 25, 0, 90)
//	length := arc.ArcLength() // Returns π * 25 / 2 (quarter circle)
func (a *Arc) ArcLength() float64 {
	// Normalize angle difference
	angleDiff := a.EndAngle - a.StartAngle
	if angleDiff < 0 {
		angleDiff += 360
	}
	// Convert to radians
	angleRad := angleDiff * math.Pi / 180.0
	return a.Radius * angleRad
}

// Area calculates the area of the sector defined by the Arc entity.
//
// Example:
//
//	arc := dxf.NewArc(50, 50, 25, 0, 90)
//	area := arc.Area() // Returns area of 90° sector
func (a *Arc) Area() float64 {
	angleDiff := a.EndAngle - a.StartAngle
	if angleDiff < 0 {
		angleDiff += 360
	}
	angleRad := angleDiff * math.Pi / 180.0
	return 0.5 * a.Radius * a.Radius * angleRad
}

// BoundingBox returns the bounding box of an Arc entity.
// Returns (minX, minY, maxX, maxY).
//
// Example:
//
//	arc := dxf.NewArc(50, 50, 25, 0, 90)
//	minX, minY, maxX, maxY := arc.BoundingBox()
func (a *Arc) BoundingBox() (minX, minY, maxX, maxY float64) {
	// Start with the center point
	minX, maxX = a.CenterX, a.CenterX
	minY, maxY = a.CenterY, a.CenterY

	// Check start and end points
	startRad := a.StartAngle * math.Pi / 180.0
	endRad := a.EndAngle * math.Pi / 180.0

	startX := a.CenterX + a.Radius*math.Cos(startRad)
	startY := a.CenterY + a.Radius*math.Sin(startRad)
	endX := a.CenterX + a.Radius*math.Cos(endRad)
	endY := a.CenterY + a.Radius*math.Sin(endRad)

	minX = math.Min(minX, math.Min(startX, endX))
	maxX = math.Max(maxX, math.Max(startX, endX))
	minY = math.Min(minY, math.Min(startY, endY))
	maxY = math.Max(maxY, math.Max(startY, endY))

	// Check quadrant extrema (0°, 90°, 180°, 270°)
	checkAngle := func(angle float64) {
		if a.containsAngle(angle) {
			x := a.CenterX + a.Radius*math.Cos(angle*math.Pi/180.0)
			y := a.CenterY + a.Radius*math.Sin(angle*math.Pi/180.0)
			minX = math.Min(minX, x)
			maxX = math.Max(maxX, x)
			minY = math.Min(minY, y)
			maxY = math.Max(maxY, y)
		}
	}

	checkAngle(0)   // Right
	checkAngle(90)  // Top
	checkAngle(180) // Left
	checkAngle(270) // Bottom

	return
}

// containsAngle checks if the arc contains a specific angle.
func (a *Arc) containsAngle(angle float64) bool {
	start := a.StartAngle
	end := a.EndAngle

	// Normalize angles to 0-360
	for start < 0 {
		start += 360
	}
	for start >= 360 {
		start -= 360
	}
	for end < 0 {
		end += 360
	}
	for end >= 360 {
		end -= 360
	}
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}

	if start <= end {
		return angle >= start && angle <= end
	}
	// Arc crosses 0°
	return angle >= start || angle <= end
}

// BoundingBox returns the bounding box of an Ellipse entity.
// Returns (minX, minY, maxX, maxY).
//
// Example:
//
//	ellipse := &dxf.Ellipse{CenterX: 50, CenterY: 50, MajorAxisX: 100, MajorAxisY: 0, MinorRatio: 0.5}
//	minX, minY, maxX, maxY := ellipse.BoundingBox()
func (e *Ellipse) BoundingBox() (minX, minY, maxX, maxY float64) {
	// Calculate major axis length
	majorLength := math.Sqrt(e.MajorAxisX*e.MajorAxisX + e.MajorAxisY*e.MajorAxisY)
	minorLength := majorLength * e.MinorRatio

	// Get angle of major axis
	tilt := math.Atan2(e.MajorAxisY, e.MajorAxisX)
	cos := math.Cos(tilt)
	sin := math.Sin(tilt)

	// Calculate bounding box considering rotation
	a := majorLength * cos
	b := minorLength * sin
	c := majorLength * sin
	d := minorLength * cos

	halfWidth := math.Sqrt(a*a + b*b)
	halfHeight := math.Sqrt(c*c + d*d)

	minX = e.CenterX - halfWidth
	maxX = e.CenterX + halfWidth
	minY = e.CenterY - halfHeight
	maxY = e.CenterY + halfHeight
	return
}

// BoundingBox returns the bounding box of a Point entity.
// Returns (x, y, x, y) since it's a single point.
//
// Example:
//
//	point := dxf.NewPoint(100, 200)
//	minX, minY, maxX, maxY := point.BoundingBox() // Returns (100, 200, 100, 200)
func (p *Point) BoundingBox() (minX, minY, maxX, maxY float64) {
	return p.X, p.Y, p.X, p.Y
}

// BoundingBox returns the approximate bounding box of a Text entity.
// Note: This is a simplified calculation that doesn't account for actual font metrics.
// Returns (minX, minY, maxX, maxY).
//
// Example:
//
//	text := dxf.NewText(10, 10, "Hello", dxf.WithTextHeight(5))
//	minX, minY, maxX, maxY := text.BoundingBox()
func (t *Text) BoundingBox() (minX, minY, maxX, maxY float64) {
	// Simplified: estimate width as height * length * 0.6 (typical aspect ratio)
	estimatedWidth := t.Height * float64(len(t.Content)) * 0.6

	if t.Rotation == 0 {
		return t.X, t.Y, t.X + estimatedWidth, t.Y + t.Height
	}

	// For rotated text, calculate the corners and find min/max
	angle := t.Rotation * math.Pi / 180.0
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	// Four corners of the text box
	corners := [][2]float64{
		{0, 0},
		{estimatedWidth, 0},
		{estimatedWidth, t.Height},
		{0, t.Height},
	}

	minX, minY = math.Inf(1), math.Inf(1)
	maxX, maxY = math.Inf(-1), math.Inf(-1)

	for _, corner := range corners {
		x := t.X + corner[0]*cos - corner[1]*sin
		y := t.Y + corner[0]*sin + corner[1]*cos
		minX = math.Min(minX, x)
		maxX = math.Max(maxX, x)
		minY = math.Min(minY, y)
		maxY = math.Max(maxY, y)
	}

	return
}

// BoundingBox returns the bounding box of a Solid entity.
// Returns (minX, minY, maxX, maxY).
//
// Example:
//
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
//	minX, minY, maxX, maxY := solid.BoundingBox()
func (s *Solid) BoundingBox() (minX, minY, maxX, maxY float64) {
	minX = math.Min(math.Min(s.X1, s.X2), math.Min(s.X3, s.X4))
	maxX = math.Max(math.Max(s.X1, s.X2), math.Max(s.X3, s.X4))
	minY = math.Min(math.Min(s.Y1, s.Y2), math.Min(s.Y3, s.Y4))
	maxY = math.Max(math.Max(s.Y1, s.Y2), math.Max(s.Y3, s.Y4))
	return
}

// Area calculates the area of a Solid entity using the Shoelace formula.
//
// Example:
//
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
//	area := solid.Area()
func (s *Solid) Area() float64 {
	// Shoelace formula for a quadrilateral
	// Area = 0.5 * |x1(y2-y4) + x2(y3-y1) + x3(y4-y2) + x4(y1-y3)|
	area := 0.5 * math.Abs(
		s.X1*(s.Y2-s.Y4)+
			s.X2*(s.Y3-s.Y1)+
			s.X3*(s.Y4-s.Y2)+
			s.X4*(s.Y1-s.Y3))
	return area
}

// IsTriangle checks if a Solid entity is a triangle (4th point equals 3rd point).
//
// Example:
//
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
//	isTriangle := solid.IsTriangle() // Returns true
func (s *Solid) IsTriangle() bool {
	return s.X3 == s.X4 && s.Y3 == s.Y4
}

// BoundingBox returns the bounding box of the entire Document.
// Returns (minX, minY, maxX, maxY) encompassing all entities.
//
// Example:
//
//	doc := &dxf.Document{Entities: []dxf.Entity{...}}
//	minX, minY, maxX, maxY := doc.BoundingBox()
func (d *Document) BoundingBox() (minX, minY, maxX, maxY float64) {
	if len(d.Entities) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY = math.Inf(1), math.Inf(1)
	maxX, maxY = math.Inf(-1), math.Inf(-1)

	for _, entity := range d.Entities {
		var eMinX, eMinY, eMaxX, eMaxY float64

		switch e := entity.(type) {
		case *Line:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		case *Circle:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		case *Arc:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		case *Ellipse:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		case *Point:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		case *Text:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		case *Solid:
			eMinX, eMinY, eMaxX, eMaxY = e.BoundingBox()
		default:
			continue
		}

		minX = math.Min(minX, eMinX)
		maxX = math.Max(maxX, eMaxX)
		minY = math.Min(minY, eMinY)
		maxY = math.Max(maxY, eMaxY)
	}

	return
}

// FilterByLayer returns all entities on a specific layer.
//
// Example:
//
//	doc := &dxf.Document{Entities: []dxf.Entity{...}}
//	entities := doc.FilterByLayer("MyLayer")
func (d *Document) FilterByLayer(layerName string) []Entity {
	var filtered []Entity

	for _, entity := range d.Entities {
		var layer string
		switch e := entity.(type) {
		case *Line:
			layer = e.Layer
		case *Circle:
			layer = e.Layer
		case *Arc:
			layer = e.Layer
		case *Ellipse:
			layer = e.Layer
		case *Point:
			layer = e.Layer
		case *Text:
			layer = e.Layer
		case *Solid:
			layer = e.Layer
		case *Insert:
			layer = e.Layer
		default:
			continue
		}

		if layer == layerName {
			filtered = append(filtered, entity)
		}
	}

	return filtered
}

// CountByType returns a map of entity type names to their counts.
//
// Example:
//
//	doc := &dxf.Document{Entities: []dxf.Entity{...}}
//	counts := doc.CountByType() // Returns {"LINE": 10, "CIRCLE": 5, ...}
func (d *Document) CountByType() map[string]int {
	counts := make(map[string]int)

	for _, entity := range d.Entities {
		counts[entity.EntityType()]++
	}

	return counts
}
