package dxf

import "math"

// Translate moves a Line entity by the given delta values.
// Returns a new Line instance with translated coordinates.
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 100)
//	moved := line.Translate(50, 50) // Line from (50,50) to (150,150)
func (l *Line) Translate(dx, dy float64) *Line {
	return &Line{
		Layer:    l.Layer,
		Color:    l.Color,
		X1:       l.X1 + dx,
		Y1:       l.Y1 + dy,
		X2:       l.X2 + dx,
		Y2:       l.Y2 + dy,
		LineType: l.LineType,
	}
}

// Rotate rotates a Line entity around a center point by the given angle in degrees.
// Returns a new Line instance with rotated coordinates.
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 0)
//	rotated := line.Rotate(90, 0, 0) // Rotate 90° around origin
func (l *Line) Rotate(angleDeg, cx, cy float64) *Line {
	angle := angleDeg * math.Pi / 180.0
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	// Translate to origin
	x1, y1 := l.X1-cx, l.Y1-cy
	x2, y2 := l.X2-cx, l.Y2-cy

	// Rotate
	rx1 := x1*cos - y1*sin
	ry1 := x1*sin + y1*cos
	rx2 := x2*cos - y2*sin
	ry2 := x2*sin + y2*cos

	// Translate back
	return &Line{
		Layer:    l.Layer,
		Color:    l.Color,
		X1:       rx1 + cx,
		Y1:       ry1 + cy,
		X2:       rx2 + cx,
		Y2:       ry2 + cy,
		LineType: l.LineType,
	}
}

// Scale scales a Line entity from a center point by the given factor.
// Returns a new Line instance with scaled coordinates.
//
// Example:
//
//	line := dxf.NewLine(0, 0, 100, 100)
//	scaled := line.Scale(2.0, 0, 0) // Scale 2x from origin
func (l *Line) Scale(factor, cx, cy float64) *Line {
	return &Line{
		Layer:    l.Layer,
		Color:    l.Color,
		X1:       cx + (l.X1-cx)*factor,
		Y1:       cy + (l.Y1-cy)*factor,
		X2:       cx + (l.X2-cx)*factor,
		Y2:       cy + (l.Y2-cy)*factor,
		LineType: l.LineType,
	}
}

// Translate moves a Circle entity by the given delta values.
// Returns a new Circle instance with translated center.
//
// Example:
//
//	circle := dxf.NewCircle(50, 50, 25)
//	moved := circle.Translate(100, 100) // Center at (150,150)
func (c *Circle) Translate(dx, dy float64) *Circle {
	return &Circle{
		Layer:   c.Layer,
		Color:   c.Color,
		CenterX: c.CenterX + dx,
		CenterY: c.CenterY + dy,
		Radius:  c.Radius,
	}
}

// Scale scales a Circle entity's radius by the given factor.
// Returns a new Circle instance with scaled radius.
//
// Example:
//
//	circle := dxf.NewCircle(50, 50, 25)
//	scaled := circle.Scale(2.0) // Radius becomes 50
func (c *Circle) Scale(factor float64) *Circle {
	return &Circle{
		Layer:   c.Layer,
		Color:   c.Color,
		CenterX: c.CenterX,
		CenterY: c.CenterY,
		Radius:  c.Radius * factor,
	}
}

// Translate moves an Arc entity by the given delta values.
// Returns a new Arc instance with translated center.
//
// Example:
//
//	arc := dxf.NewArc(50, 50, 25, 0, 90)
//	moved := arc.Translate(100, 100) // Center at (150,150)
func (a *Arc) Translate(dx, dy float64) *Arc {
	return &Arc{
		Layer:      a.Layer,
		Color:      a.Color,
		CenterX:    a.CenterX + dx,
		CenterY:    a.CenterY + dy,
		Radius:     a.Radius,
		StartAngle: a.StartAngle,
		EndAngle:   a.EndAngle,
	}
}

// Scale scales an Arc entity's radius by the given factor.
// Returns a new Arc instance with scaled radius.
//
// Example:
//
//	arc := dxf.NewArc(50, 50, 25, 0, 90)
//	scaled := arc.Scale(2.0) // Radius becomes 50
func (a *Arc) Scale(factor float64) *Arc {
	return &Arc{
		Layer:      a.Layer,
		Color:      a.Color,
		CenterX:    a.CenterX,
		CenterY:    a.CenterY,
		Radius:     a.Radius * factor,
		StartAngle: a.StartAngle,
		EndAngle:   a.EndAngle,
	}
}

// Translate moves an Ellipse entity by the given delta values.
// Returns a new Ellipse instance with translated center.
//
// Example:
//
//	ellipse := &dxf.Ellipse{CenterX: 50, CenterY: 50, MajorAxisX: 100, MajorAxisY: 0, MinorRatio: 0.5}
//	moved := ellipse.Translate(100, 100) // Center at (150,150)
func (e *Ellipse) Translate(dx, dy float64) *Ellipse {
	return &Ellipse{
		Layer:      e.Layer,
		Color:      e.Color,
		CenterX:    e.CenterX + dx,
		CenterY:    e.CenterY + dy,
		MajorAxisX: e.MajorAxisX,
		MajorAxisY: e.MajorAxisY,
		MinorRatio: e.MinorRatio,
		StartParam: e.StartParam,
		EndParam:   e.EndParam,
	}
}

// Scale scales an Ellipse entity's axes by the given factor.
// Returns a new Ellipse instance with scaled major axis.
//
// Example:
//
//	ellipse := &dxf.Ellipse{MajorAxisX: 100, MajorAxisY: 0, MinorRatio: 0.5}
//	scaled := ellipse.Scale(2.0) // Major axis doubles
func (e *Ellipse) Scale(factor float64) *Ellipse {
	return &Ellipse{
		Layer:      e.Layer,
		Color:      e.Color,
		CenterX:    e.CenterX,
		CenterY:    e.CenterY,
		MajorAxisX: e.MajorAxisX * factor,
		MajorAxisY: e.MajorAxisY * factor,
		MinorRatio: e.MinorRatio,
		StartParam: e.StartParam,
		EndParam:   e.EndParam,
	}
}

// Translate moves a Point entity by the given delta values.
// Returns a new Point instance with translated coordinates.
//
// Example:
//
//	point := dxf.NewPoint(100, 200)
//	moved := point.Translate(50, 50) // Point at (150,250)
func (p *Point) Translate(dx, dy float64) *Point {
	return &Point{
		Layer: p.Layer,
		Color: p.Color,
		X:     p.X + dx,
		Y:     p.Y + dy,
	}
}

// Translate moves a Text entity by the given delta values.
// Returns a new Text instance with translated position.
//
// Example:
//
//	text := dxf.NewText(10, 10, "Hello")
//	moved := text.Translate(50, 50) // Text at (60,60)
func (t *Text) Translate(dx, dy float64) *Text {
	return &Text{
		Layer:    t.Layer,
		Color:    t.Color,
		X:        t.X + dx,
		Y:        t.Y + dy,
		Height:   t.Height,
		Rotation: t.Rotation,
		Content:  t.Content,
		Style:    t.Style,
	}
}

// Rotate rotates a Text entity's rotation angle by the given degrees.
// Returns a new Text instance with updated rotation.
//
// Example:
//
//	text := dxf.NewText(10, 10, "Hello", dxf.WithTextRotation(0))
//	rotated := text.Rotate(45) // Rotation becomes 45°
func (t *Text) Rotate(angleDeg float64) *Text {
	return &Text{
		Layer:    t.Layer,
		Color:    t.Color,
		X:        t.X,
		Y:        t.Y,
		Height:   t.Height,
		Rotation: t.Rotation + angleDeg,
		Content:  t.Content,
		Style:    t.Style,
	}
}

// Scale scales a Text entity's height by the given factor.
// Returns a new Text instance with scaled height.
//
// Example:
//
//	text := dxf.NewText(10, 10, "Hello", dxf.WithTextHeight(5))
//	scaled := text.Scale(2.0) // Height becomes 10
func (t *Text) Scale(factor float64) *Text {
	return &Text{
		Layer:    t.Layer,
		Color:    t.Color,
		X:        t.X,
		Y:        t.Y,
		Height:   t.Height * factor,
		Rotation: t.Rotation,
		Content:  t.Content,
		Style:    t.Style,
	}
}

// Translate moves a Solid entity by the given delta values.
// Returns a new Solid instance with translated vertices.
//
// Example:
//
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
//	moved := solid.Translate(50, 50)
func (s *Solid) Translate(dx, dy float64) *Solid {
	return &Solid{
		Layer: s.Layer,
		Color: s.Color,
		X1:    s.X1 + dx,
		Y1:    s.Y1 + dy,
		X2:    s.X2 + dx,
		Y2:    s.Y2 + dy,
		X3:    s.X3 + dx,
		Y3:    s.Y3 + dy,
		X4:    s.X4 + dx,
		Y4:    s.Y4 + dy,
	}
}

// Rotate rotates a Solid entity around a center point by the given angle in degrees.
// Returns a new Solid instance with rotated vertices.
//
// Example:
//
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
//	rotated := solid.Rotate(45, 50, 50)
func (s *Solid) Rotate(angleDeg, cx, cy float64) *Solid {
	angle := angleDeg * math.Pi / 180.0
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	rotatePoint := func(x, y float64) (float64, float64) {
		dx, dy := x-cx, y-cy
		rx := dx*cos - dy*sin
		ry := dx*sin + dy*cos
		return rx + cx, ry + cy
	}

	rx1, ry1 := rotatePoint(s.X1, s.Y1)
	rx2, ry2 := rotatePoint(s.X2, s.Y2)
	rx3, ry3 := rotatePoint(s.X3, s.Y3)
	rx4, ry4 := rotatePoint(s.X4, s.Y4)

	return &Solid{
		Layer: s.Layer,
		Color: s.Color,
		X1:    rx1,
		Y1:    ry1,
		X2:    rx2,
		Y2:    ry2,
		X3:    rx3,
		Y3:    ry3,
		X4:    rx4,
		Y4:    ry4,
	}
}

// Scale scales a Solid entity from a center point by the given factor.
// Returns a new Solid instance with scaled vertices.
//
// Example:
//
//	solid := dxf.NewSolid(0, 0, 100, 0, 50, 100, 50, 100)
//	scaled := solid.Scale(2.0, 50, 50)
func (s *Solid) Scale(factor, cx, cy float64) *Solid {
	scalePoint := func(x, y float64) (float64, float64) {
		return cx + (x-cx)*factor, cy + (y-cy)*factor
	}

	sx1, sy1 := scalePoint(s.X1, s.Y1)
	sx2, sy2 := scalePoint(s.X2, s.Y2)
	sx3, sy3 := scalePoint(s.X3, s.Y3)
	sx4, sy4 := scalePoint(s.X4, s.Y4)

	return &Solid{
		Layer: s.Layer,
		Color: s.Color,
		X1:    sx1,
		Y1:    sy1,
		X2:    sx2,
		Y2:    sy2,
		X3:    sx3,
		Y3:    sy3,
		X4:    sx4,
		Y4:    sy4,
	}
}

// Translate moves an Insert entity by the given delta values.
// Returns a new Insert instance with translated insertion point.
//
// Example:
//
//	insert := dxf.NewInsert("MyBlock", 100, 100)
//	moved := insert.Translate(50, 50) // Insert at (150,150)
func (i *Insert) Translate(dx, dy float64) *Insert {
	return &Insert{
		Layer:     i.Layer,
		Color:     i.Color,
		BlockName: i.BlockName,
		X:         i.X + dx,
		Y:         i.Y + dy,
		ScaleX:    i.ScaleX,
		ScaleY:    i.ScaleY,
		Rotation:  i.Rotation,
	}
}

// Rotate rotates an Insert entity's rotation angle by the given degrees.
// Returns a new Insert instance with updated rotation.
//
// Example:
//
//	insert := dxf.NewInsert("MyBlock", 100, 100)
//	rotated := insert.Rotate(45) // Rotation becomes 45°
func (i *Insert) Rotate(angleDeg float64) *Insert {
	return &Insert{
		Layer:     i.Layer,
		Color:     i.Color,
		BlockName: i.BlockName,
		X:         i.X,
		Y:         i.Y,
		ScaleX:    i.ScaleX,
		ScaleY:    i.ScaleY,
		Rotation:  i.Rotation + angleDeg,
	}
}

// Scale scales an Insert entity's scale factors by the given factor.
// Returns a new Insert instance with scaled scale factors.
//
// Example:
//
//	insert := dxf.NewInsert("MyBlock", 100, 100)
//	scaled := insert.Scale(2.0) // Scale factors become 2.0
func (i *Insert) Scale(factor float64) *Insert {
	return &Insert{
		Layer:     i.Layer,
		Color:     i.Color,
		BlockName: i.BlockName,
		X:         i.X,
		Y:         i.Y,
		ScaleX:    i.ScaleX * factor,
		ScaleY:    i.ScaleY * factor,
		Rotation:  i.Rotation,
	}
}
