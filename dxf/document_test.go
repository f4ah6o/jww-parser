package dxf

import "testing"

func TestNewDocument(t *testing.T) {
	doc := NewDocument()

	if len(doc.Layers) != 1 {
		t.Errorf("Expected 1 default layer, got %d", len(doc.Layers))
	}

	if doc.Layers[0].Name != "0" {
		t.Errorf("Expected default layer '0', got '%s'", doc.Layers[0].Name)
	}

	if len(doc.Entities) != 0 {
		t.Errorf("Expected 0 entities, got %d", len(doc.Entities))
	}
}

func TestDocumentAddLayer(t *testing.T) {
	doc := NewDocument().
		AddLayer("Layer1", 1, "CONTINUOUS").
		AddLayer("Layer2", 2, "DASHED")

	if doc.LayerCount() != 3 {
		t.Errorf("Expected 3 layers (including default), got %d", doc.LayerCount())
	}

	layer := doc.GetLayer("Layer1")
	if layer == nil {
		t.Fatal("Layer1 not found")
	}
	if layer.Color != 1 {
		t.Errorf("Expected color 1, got %d", layer.Color)
	}
}

func TestDocumentAddEntity(t *testing.T) {
	line := NewLine(0, 0, 100, 100)
	doc := NewDocument().AddEntity(line)

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", doc.EntityCount())
	}
}

func TestDocumentAddLine(t *testing.T) {
	doc := NewDocument().
		AddLine(0, 0, 100, 100).
		AddLine(10, 10, 50, 50)

	if doc.EntityCount() != 2 {
		t.Errorf("Expected 2 entities, got %d", doc.EntityCount())
	}

	line, ok := doc.Entities[0].(*Line)
	if !ok {
		t.Fatal("First entity is not a Line")
	}
	if line.X1 != 0 || line.Y1 != 0 {
		t.Errorf("First line start point mismatch")
	}
}

func TestDocumentAddCircle(t *testing.T) {
	doc := NewDocument().
		AddCircle(50, 50, 25).
		AddCircle(100, 100, 50)

	if doc.EntityCount() != 2 {
		t.Errorf("Expected 2 entities, got %d", doc.EntityCount())
	}

	circle, ok := doc.Entities[0].(*Circle)
	if !ok {
		t.Fatal("First entity is not a Circle")
	}
	if circle.Radius != 25 {
		t.Errorf("Expected radius 25, got %f", circle.Radius)
	}
}

func TestDocumentAddArc(t *testing.T) {
	doc := NewDocument().AddArc(50, 50, 25, 0, 90)

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", doc.EntityCount())
	}

	arc, ok := doc.Entities[0].(*Arc)
	if !ok {
		t.Fatal("Entity is not an Arc")
	}
	if arc.StartAngle != 0 || arc.EndAngle != 90 {
		t.Errorf("Arc angles mismatch")
	}
}

func TestDocumentAddPoint(t *testing.T) {
	doc := NewDocument().AddPoint(100, 200)

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", doc.EntityCount())
	}

	point, ok := doc.Entities[0].(*Point)
	if !ok {
		t.Fatal("Entity is not a Point")
	}
	if point.X != 100 || point.Y != 200 {
		t.Errorf("Point coordinates mismatch")
	}
}

func TestDocumentAddText(t *testing.T) {
	doc := NewDocument().AddText(10, 10, "Hello World")

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", doc.EntityCount())
	}

	text, ok := doc.Entities[0].(*Text)
	if !ok {
		t.Fatal("Entity is not a Text")
	}
	if text.Content != "Hello World" {
		t.Errorf("Text content mismatch")
	}
}

func TestDocumentAddSolid(t *testing.T) {
	doc := NewDocument().AddSolid(0, 0, 100, 0, 50, 100, 50, 100)

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", doc.EntityCount())
	}

	solid, ok := doc.Entities[0].(*Solid)
	if !ok {
		t.Fatal("Entity is not a Solid")
	}
	if solid.X1 != 0 || solid.Y1 != 0 {
		t.Errorf("Solid point 1 mismatch")
	}
}

func TestDocumentAddInsert(t *testing.T) {
	doc := NewDocument().AddInsert("MyBlock", 100, 100)

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", doc.EntityCount())
	}

	insert, ok := doc.Entities[0].(*Insert)
	if !ok {
		t.Fatal("Entity is not an Insert")
	}
	if insert.BlockName != "MyBlock" {
		t.Errorf("Block name mismatch")
	}
}

func TestDocumentAddBlock(t *testing.T) {
	block := Block{
		Name:  "MyBlock",
		BaseX: 0,
		BaseY: 0,
		Entities: []Entity{
			NewLine(0, 0, 100, 100),
		},
	}

	doc := NewDocument().AddBlock(block)

	if doc.BlockCount() != 1 {
		t.Errorf("Expected 1 block, got %d", doc.BlockCount())
	}

	found := doc.GetBlock("MyBlock")
	if found == nil {
		t.Fatal("Block not found")
	}
	if found.Name != "MyBlock" {
		t.Errorf("Block name mismatch")
	}
}

func TestDocumentRemoveEntity(t *testing.T) {
	doc := NewDocument().
		AddLine(0, 0, 100, 100).
		AddCircle(50, 50, 25).
		RemoveEntity(0)

	if doc.EntityCount() != 1 {
		t.Errorf("Expected 1 entity after removal, got %d", doc.EntityCount())
	}

	// First entity should now be the circle
	_, ok := doc.Entities[0].(*Circle)
	if !ok {
		t.Error("Remaining entity should be a Circle")
	}
}

func TestDocumentClearEntities(t *testing.T) {
	doc := NewDocument().
		AddLine(0, 0, 100, 100).
		AddCircle(50, 50, 25).
		ClearEntities()

	if doc.EntityCount() != 0 {
		t.Errorf("Expected 0 entities after clear, got %d", doc.EntityCount())
	}
}

func TestDocumentGetLayer(t *testing.T) {
	doc := NewDocument().AddLayer("MyLayer", 1, "CONTINUOUS")

	layer := doc.GetLayer("MyLayer")
	if layer == nil {
		t.Fatal("Layer not found")
	}
	if layer.Name != "MyLayer" {
		t.Errorf("Layer name mismatch")
	}

	notFound := doc.GetLayer("NonExistent")
	if notFound != nil {
		t.Error("Expected nil for non-existent layer")
	}
}

func TestDocumentHasLayer(t *testing.T) {
	doc := NewDocument().AddLayer("MyLayer", 1, "CONTINUOUS")

	if !doc.HasLayer("MyLayer") {
		t.Error("Expected HasLayer to return true for MyLayer")
	}

	if doc.HasLayer("NonExistent") {
		t.Error("Expected HasLayer to return false for non-existent layer")
	}
}

func TestDocumentGetBlock(t *testing.T) {
	block := Block{Name: "MyBlock", Entities: []Entity{}}
	doc := NewDocument().AddBlock(block)

	found := doc.GetBlock("MyBlock")
	if found == nil {
		t.Fatal("Block not found")
	}

	notFound := doc.GetBlock("NonExistent")
	if notFound != nil {
		t.Error("Expected nil for non-existent block")
	}
}

func TestDocumentHasBlock(t *testing.T) {
	block := Block{Name: "MyBlock", Entities: []Entity{}}
	doc := NewDocument().AddBlock(block)

	if !doc.HasBlock("MyBlock") {
		t.Error("Expected HasBlock to return true for MyBlock")
	}

	if doc.HasBlock("NonExistent") {
		t.Error("Expected HasBlock to return false for non-existent block")
	}
}

func TestDocumentFluentAPI(t *testing.T) {
	// Test method chaining
	doc := NewDocument().
		AddLayer("Layer1", 1, "CONTINUOUS").
		AddLayer("Layer2", 2, "DASHED").
		AddLine(0, 0, 100, 100, WithLineLayer("Layer1")).
		AddCircle(50, 50, 25, WithCircleLayer("Layer2")).
		AddPoint(100, 100).
		AddText(10, 10, "Hello", WithTextHeight(5.0))

	if doc.LayerCount() != 3 { // Including default "0"
		t.Errorf("Expected 3 layers, got %d", doc.LayerCount())
	}

	if doc.EntityCount() != 4 {
		t.Errorf("Expected 4 entities, got %d", doc.EntityCount())
	}
}
