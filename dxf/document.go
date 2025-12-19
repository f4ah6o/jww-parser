package dxf

// NewDocument creates a new empty DXF document with a default layer "0".
//
// Example:
//
//	doc := dxf.NewDocument()
//	doc.AddLine(0, 0, 100, 100)
func NewDocument() *Document {
	return &Document{
		Layers: []Layer{
			{
				Name:     "0",
				Color:    7,
				LineType: "CONTINUOUS",
				Frozen:   false,
				Locked:   false,
			},
		},
		Entities: []Entity{},
		Blocks:   []Block{},
	}
}

// AddLayer adds a new layer to the document and returns the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddLayer("MyLayer", 1, "CONTINUOUS").
//		AddLayer("AnotherLayer", 2, "DASHED")
func (d *Document) AddLayer(name string, color int, lineType string) *Document {
	d.Layers = append(d.Layers, Layer{
		Name:     name,
		Color:    color,
		LineType: lineType,
		Frozen:   false,
		Locked:   false,
	})
	return d
}

// AddEntity adds an entity to the document and returns the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddEntity(dxf.NewLine(0, 0, 100, 100)).
//		AddEntity(dxf.NewCircle(50, 50, 25))
func (d *Document) AddEntity(entity Entity) *Document {
	d.Entities = append(d.Entities, entity)
	return d
}

// AddLine creates and adds a Line entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddLine(0, 0, 100, 100, dxf.WithLineLayer("MyLayer"), dxf.WithLineColor(1))
func (d *Document) AddLine(x1, y1, x2, y2 float64, opts ...LineOption) *Document {
	d.Entities = append(d.Entities, NewLine(x1, y1, x2, y2, opts...))
	return d
}

// AddCircle creates and adds a Circle entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddCircle(50, 50, 25, dxf.WithCircleLayer("MyLayer"), dxf.WithCircleColor(2))
func (d *Document) AddCircle(centerX, centerY, radius float64, opts ...CircleOption) *Document {
	d.Entities = append(d.Entities, NewCircle(centerX, centerY, radius, opts...))
	return d
}

// AddArc creates and adds an Arc entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddArc(50, 50, 25, 0, 90, dxf.WithArcLayer("MyLayer"), dxf.WithArcColor(3))
func (d *Document) AddArc(centerX, centerY, radius, startAngle, endAngle float64, opts ...ArcOption) *Document {
	d.Entities = append(d.Entities, NewArc(centerX, centerY, radius, startAngle, endAngle, opts...))
	return d
}

// AddPoint creates and adds a Point entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddPoint(100, 200, dxf.WithPointLayer("MyLayer"), dxf.WithPointColor(4))
func (d *Document) AddPoint(x, y float64, opts ...PointOption) *Document {
	d.Entities = append(d.Entities, NewPoint(x, y, opts...))
	return d
}

// AddText creates and adds a Text entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddText(10, 10, "Hello World",
//			dxf.WithTextLayer("MyLayer"),
//			dxf.WithTextHeight(5.0),
//			dxf.WithTextRotation(45))
func (d *Document) AddText(x, y float64, content string, opts ...TextOption) *Document {
	d.Entities = append(d.Entities, NewText(x, y, content, opts...))
	return d
}

// AddSolid creates and adds a Solid entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddSolid(0, 0, 100, 0, 50, 100, 50, 100,
//			dxf.WithSolidLayer("MyLayer"),
//			dxf.WithSolidColor(5))
func (d *Document) AddSolid(p1x, p1y, p2x, p2y, p3x, p3y, p4x, p4y float64, opts ...SolidOption) *Document {
	d.Entities = append(d.Entities, NewSolid(p1x, p1y, p2x, p2y, p3x, p3y, p4x, p4y, opts...))
	return d
}

// AddInsert creates and adds an Insert entity to the document, returning the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddInsert("MyBlock", 100, 100,
//			dxf.WithInsertLayer("MyLayer"),
//			dxf.WithInsertScale(2.0, 2.0),
//			dxf.WithInsertRotation(45))
func (d *Document) AddInsert(blockName string, x, y float64, opts ...InsertOption) *Document {
	d.Entities = append(d.Entities, NewInsert(blockName, x, y, opts...))
	return d
}

// AddBlock adds a block definition to the document and returns the document for chaining.
//
// Example:
//
//	block := dxf.Block{
//		Name: "MyBlock",
//		BaseX: 0,
//		BaseY: 0,
//		Entities: []dxf.Entity{
//			dxf.NewLine(0, 0, 100, 100),
//		},
//	}
//	doc := dxf.NewDocument().AddBlock(block)
func (d *Document) AddBlock(block Block) *Document {
	d.Blocks = append(d.Blocks, block)
	return d
}

// RemoveEntity removes the entity at the specified index from the document.
// Returns the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddLine(0, 0, 100, 100).
//		AddCircle(50, 50, 25).
//		RemoveEntity(0) // Removes the line
func (d *Document) RemoveEntity(index int) *Document {
	if index >= 0 && index < len(d.Entities) {
		d.Entities = append(d.Entities[:index], d.Entities[index+1:]...)
	}
	return d
}

// ClearEntities removes all entities from the document.
// Returns the document for chaining.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddLine(0, 0, 100, 100).
//		AddCircle(50, 50, 25).
//		ClearEntities() // Removes all entities
func (d *Document) ClearEntities() *Document {
	d.Entities = []Entity{}
	return d
}

// GetLayer returns a layer by name, or nil if not found.
//
// Example:
//
//	doc := dxf.NewDocument().AddLayer("MyLayer", 1, "CONTINUOUS")
//	layer := doc.GetLayer("MyLayer")
func (d *Document) GetLayer(name string) *Layer {
	for i := range d.Layers {
		if d.Layers[i].Name == name {
			return &d.Layers[i]
		}
	}
	return nil
}

// HasLayer checks if a layer with the given name exists.
//
// Example:
//
//	doc := dxf.NewDocument().AddLayer("MyLayer", 1, "CONTINUOUS")
//	exists := doc.HasLayer("MyLayer") // Returns true
func (d *Document) HasLayer(name string) bool {
	return d.GetLayer(name) != nil
}

// GetBlock returns a block by name, or nil if not found.
//
// Example:
//
//	block := dxf.Block{Name: "MyBlock", Entities: []dxf.Entity{}}
//	doc := dxf.NewDocument().AddBlock(block)
//	found := doc.GetBlock("MyBlock")
func (d *Document) GetBlock(name string) *Block {
	for i := range d.Blocks {
		if d.Blocks[i].Name == name {
			return &d.Blocks[i]
		}
	}
	return nil
}

// HasBlock checks if a block with the given name exists.
//
// Example:
//
//	block := dxf.Block{Name: "MyBlock", Entities: []dxf.Entity{}}
//	doc := dxf.NewDocument().AddBlock(block)
//	exists := doc.HasBlock("MyBlock") // Returns true
func (d *Document) HasBlock(name string) bool {
	return d.GetBlock(name) != nil
}

// EntityCount returns the number of entities in the document.
//
// Example:
//
//	doc := dxf.NewDocument().
//		AddLine(0, 0, 100, 100).
//		AddCircle(50, 50, 25)
//	count := doc.EntityCount() // Returns 2
func (d *Document) EntityCount() int {
	return len(d.Entities)
}

// LayerCount returns the number of layers in the document.
//
// Example:
//
//	doc := dxf.NewDocument().AddLayer("MyLayer", 1, "CONTINUOUS")
//	count := doc.LayerCount() // Returns 2 (includes default "0" layer)
func (d *Document) LayerCount() int {
	return len(d.Layers)
}

// BlockCount returns the number of blocks in the document.
//
// Example:
//
//	block := dxf.Block{Name: "MyBlock", Entities: []dxf.Entity{}}
//	doc := dxf.NewDocument().AddBlock(block)
//	count := doc.BlockCount() // Returns 1
func (d *Document) BlockCount() int {
	return len(d.Blocks)
}
