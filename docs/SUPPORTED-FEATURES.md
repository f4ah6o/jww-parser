# Supported JWW Features

This document lists all JWW (Jw_cad) features supported by the jww-parser library.

## File Format Versions

| Version Range | Support Status |
|---------------|----------------|
| 200-299 | Partial (legacy) |
| 300-399 | Full Support |
| 400-499 | Full Support |
| 500-599 | Full Support |
| 600-699 | Full Support |
| 700-810 | Full Support |

## Entity Types

### Line (Sen)

| Feature | JWW | DXF | Notes |
|---------|-----|-----|-------|
| Basic line | ✅ | LINE | |
| Construction line | ✅ | LINE | |
| Line color | ✅ | ✅ | Mapped to ACI |
| Line type | ✅ | ⚠️ | Basic types only |
| Line width | ✅ | ⚠️ | Converted to weight |

### Arc/Circle (Enko)

| Feature | JWW | DXF | Notes |
|---------|-----|-----|-------|
| Full circle | ✅ | CIRCLE | |
| Arc | ✅ | ARC | |
| Ellipse | ✅ | ELLIPSE | |
| Elliptical arc | ✅ | ELLIPSE | |
| Color | ✅ | ✅ | |
| Flatness ratio | ✅ | ✅ | Converted to minor ratio |

### Point (Ten)

| Feature | JWW | DXF | Notes |
|---------|-----|-----|-------|
| Standard point | ✅ | POINT | |
| Temporary point | ⚠️ | - | Skipped by default |
| Point code | ✅ | - | Available in JWW JSON |

### Text (Moji)

| Feature | JWW | DXF | Notes |
|---------|-----|-----|-------|
| Single line text | ✅ | TEXT | |
| Text height | ✅ | ✅ | |
| Text width | ✅ | ⚠️ | Width factor approximation |
| Rotation angle | ✅ | ✅ | |
| Font name | ✅ | - | Not converted to DXF |
| Japanese text | ✅ | ✅ | Shift-JIS to UTF-8 |
| Special characters | ✅ | ⚠️ | Unicode escape in DXF |

### Solid Fill (Soryomen)

| Feature | JWW | DXF | Notes |
|---------|-----|-----|-------|
| Triangle | ✅ | SOLID | |
| Quadrilateral | ✅ | SOLID | |
| Polygon (>4 points) | ⚠️ | ⚠️ | Triangulated |
| Solid color | ✅ | ✅ | |

### Block (Buzoku)

| Feature | JWW | DXF | Notes |
|---------|-----|-----|-------|
| Block definition | ✅ | BLOCK | |
| Block reference | ✅ | INSERT | |
| Scale X/Y | ✅ | ✅ | |
| Rotation | ✅ | ✅ | |
| Nested blocks | ⚠️ | ⚠️ | Limited depth |

## Layer Structure

### Layer Groups

- 16 layer groups (0-F)
- Group names preserved
- Visibility state converted

### Layers

- 16 layers per group (0-F)
- Layer names preserved
- Visibility state → DXF frozen
- Lock state → DXF locked

### Layer Naming

JWW layers are converted to DXF layers using the format:
```
{GroupIndex in Hex}-{LayerIndex in Hex}
```

Examples:
- Group 0, Layer 0 → "0-0"
- Group 5, Layer 10 → "5-A"
- Group 15, Layer 15 → "F-F"

## Colors

### Standard Colors

| JWW Code | Name | DXF ACI |
|----------|------|---------|
| 1 | Black (on white bg) / White (on black bg) | 7 |
| 2 | Blue | 5 |
| 3 | Red | 1 |
| 4 | Magenta | 6 |
| 5 | Green | 3 |
| 6 | Cyan | 4 |
| 7 | Yellow | 2 |
| 8 | White (on white bg) / Black (on black bg) | 7 |
| 9 | Gray | 8 |

### Extended Colors (100+)

Extended colors are mapped to the nearest ACI color. The mapping algorithm:

1. Extract RGB components from JWW color code
2. Find closest match in ACI palette
3. Use that ACI code

## Line Types

| JWW Type | Name | DXF Type |
|----------|------|----------|
| 0 | Solid | CONTINUOUS |
| 1 | Dashed | DASHED |
| 2 | Dash-dot | DASHDOT |
| 3 | Dotted | DOT |
| 4 | Long dash | DASHED2 |
| 5-9 | Custom | BYLAYER |

## Unsupported Features

The following JWW features are NOT currently supported:

### Entities
- ❌ Dimensions (Sunpou)
- ❌ Hatching patterns
- ❌ Splines/Bezier curves
- ❌ Images/raster graphics
- ❌ OLE objects

### Attributes
- ❌ Extended entity data (XDATA)
- ❌ Hyperlinks
- ❌ Custom properties

### Drawing Features
- ❌ Paper space / model space separation
- ❌ Viewports
- ❌ Named views
- ❌ Printing settings

## Known Limitations

### Precision

- Coordinates are stored as 64-bit floating point
- Maximum precision: 15-16 significant digits
- Recommended working precision: 6 decimal places

### Large Files

- Files > 10MB may require streaming mode
- Entity limit: No hard limit (memory dependent)
- Recommended: Use `maxEntities` option for previews

### Compatibility

- Output DXF uses AutoCAD R12/R14 format
- Some CAD applications may have limited support
- ODA FileConverter may show compatibility warnings

### Text

- Font substitution may occur in target applications
- Special JWW font features not preserved
- Vertical text converted to rotated horizontal

### Blocks

- Block attributes not fully supported
- Dynamic blocks converted to static
- Block table limited to 255 blocks

## Feature Matrix

| Feature | parse() | toDxf() | toDxfString() |
|---------|---------|---------|---------------|
| Line | ✅ | ✅ | ✅ |
| Circle | ✅ | ✅ | ✅ |
| Arc | ✅ | ✅ | ✅ |
| Ellipse | ✅ | ✅ | ✅ |
| Point | ✅ | ⚠️ | ⚠️ |
| Text | ✅ | ✅ | ✅ |
| Solid | ✅ | ⚠️ | ⚠️ |
| Block Def | ✅ | ✅ | ✅ |
| Block Ref | ✅ | ✅ | ✅ |
| Layers | ✅ | ✅ | ✅ |
| Colors | ✅ | ✅ | ✅ |
| Line Types | ✅ | ⚠️ | ⚠️ |

Legend:
- ✅ Full support
- ⚠️ Partial support / limitations
- ❌ Not supported

## Requesting Features

If you need support for additional JWW features, please:

1. Open an issue on GitHub
2. Provide sample JWW files if possible
3. Describe the use case

We prioritize features based on community demand and feasibility.
