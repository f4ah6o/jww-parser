# DXF JSON Format Specification

This document describes the JSON structure returned by the `toDxf()` method.

## Overview

The DXF JSON format is a JavaScript-friendly representation of DXF (Drawing Exchange Format) data. It provides structured access to CAD drawing data that can be easily manipulated in JavaScript/TypeScript applications.

## Document Structure

```typescript
interface DxfDocument {
  Layers: DxfLayer[];
  Entities: DxfEntity[];
  Blocks: DxfBlock[];
}
```

### Layers

```typescript
interface DxfLayer {
  Name: string;     // Layer name (e.g., "0-0", "A-5")
  Color: number;    // ACI color code (1-255)
  Frozen: boolean;  // Layer frozen state
  Locked: boolean;  // Layer locked state
}
```

**Layer Naming Convention:**

JWW uses 16 layer groups × 16 layers = 256 layers total. These are converted to DXF layer names using the pattern `{group}-{layer}` where both are hexadecimal digits (0-F).

Examples:
- Layer 0 in Group 0 → `"0-0"`
- Layer 10 in Group 15 → `"F-A"`

## Entity Types

### LINE

Represents a straight line segment.

```typescript
interface DxfLine {
  Type: "LINE";
  Layer: string;
  Color?: number;
  LineType?: string;
  X1: number;      // Start X
  Y1: number;      // Start Y
  Z1?: number;     // Start Z (optional, default 0)
  X2: number;      // End X
  Y2: number;      // End Y
  Z2?: number;     // End Z (optional, default 0)
}
```

**Example:**
```json
{
  "Type": "LINE",
  "Layer": "0-0",
  "Color": 7,
  "X1": 0.0,
  "Y1": 0.0,
  "X2": 100.0,
  "Y2": 50.0
}
```

### CIRCLE

Represents a complete circle.

```typescript
interface DxfCircle {
  Type: "CIRCLE";
  Layer: string;
  Color?: number;
  CenterX: number;
  CenterY: number;
  CenterZ?: number;
  Radius: number;
}
```

**Example:**
```json
{
  "Type": "CIRCLE",
  "Layer": "0-1",
  "Color": 1,
  "CenterX": 50.0,
  "CenterY": 50.0,
  "Radius": 25.0
}
```

### ARC

Represents a circular arc.

```typescript
interface DxfArc {
  Type: "ARC";
  Layer: string;
  Color?: number;
  CenterX: number;
  CenterY: number;
  CenterZ?: number;
  Radius: number;
  StartAngle: number;  // In degrees (0-360)
  EndAngle: number;    // In degrees (0-360)
}
```

**Angle Convention:**
- Angles are measured counter-clockwise from the positive X-axis
- 0° = East, 90° = North, 180° = West, 270° = South
- Arc is drawn from StartAngle to EndAngle in counter-clockwise direction

**Example:**
```json
{
  "Type": "ARC",
  "Layer": "0-0",
  "Color": 3,
  "CenterX": 100.0,
  "CenterY": 100.0,
  "Radius": 50.0,
  "StartAngle": 0.0,
  "EndAngle": 90.0
}
```

### ELLIPSE

Represents an ellipse or elliptical arc.

```typescript
interface DxfEllipse {
  Type: "ELLIPSE";
  Layer: string;
  Color?: number;
  CenterX: number;
  CenterY: number;
  CenterZ?: number;
  MajorAxisX: number;   // Major axis endpoint relative to center
  MajorAxisY: number;
  MajorAxisZ?: number;
  MinorRatio: number;   // Minor/Major axis ratio (0 < ratio ≤ 1)
  StartParam: number;   // Start parameter (0 to 2π)
  EndParam: number;     // End parameter (0 to 2π)
}
```

**Ellipse Calculation:**

For a point on the ellipse at parameter `t`:
```javascript
const majorLength = Math.sqrt(MajorAxisX**2 + MajorAxisY**2);
const minorLength = majorLength * MinorRatio;
const rotation = Math.atan2(MajorAxisY, MajorAxisX);

// Point at parameter t
const x = CenterX + majorLength * Math.cos(t) * Math.cos(rotation)
                  - minorLength * Math.sin(t) * Math.sin(rotation);
const y = CenterY + majorLength * Math.cos(t) * Math.sin(rotation)
                  + minorLength * Math.sin(t) * Math.cos(rotation);
```

**Example:**
```json
{
  "Type": "ELLIPSE",
  "Layer": "0-0",
  "Color": 2,
  "CenterX": 100.0,
  "CenterY": 100.0,
  "MajorAxisX": 50.0,
  "MajorAxisY": 0.0,
  "MinorRatio": 0.5,
  "StartParam": 0.0,
  "EndParam": 6.283185307
}
```

### POINT

Represents a point marker.

```typescript
interface DxfPoint {
  Type: "POINT";
  Layer: string;
  Color?: number;
  X: number;
  Y: number;
  Z?: number;
}
```

**Example:**
```json
{
  "Type": "POINT",
  "Layer": "0-0",
  "X": 25.0,
  "Y": 75.0
}
```

### TEXT

Represents single-line text.

```typescript
interface DxfText {
  Type: "TEXT";
  Layer: string;
  Color?: number;
  X: number;         // Insertion point X
  Y: number;         // Insertion point Y
  Z?: number;
  Height: number;    // Text height
  Content: string;   // Text content (Unicode)
  Rotation?: number; // Rotation angle in degrees
  Style?: string;    // Text style name
}
```

**Character Encoding:**
- Text is converted from JWW's Shift-JIS to UTF-8
- Non-ASCII characters are represented using Unicode escape sequences (`\U+XXXX`) in the DXF string output

**Example:**
```json
{
  "Type": "TEXT",
  "Layer": "0-5",
  "Color": 7,
  "X": 10.0,
  "Y": 20.0,
  "Height": 3.5,
  "Content": "Sample Text",
  "Rotation": 45.0
}
```

### MTEXT

Represents multi-line text with formatting.

```typescript
interface DxfMText {
  Type: "MTEXT";
  Layer: string;
  Color?: number;
  X: number;
  Y: number;
  Z?: number;
  Height: number;
  Content: string;   // May contain formatting codes
  Rotation?: number;
  Width?: number;    // Reference rectangle width
}
```

### SOLID

Represents a filled triangular or quadrilateral region.

```typescript
interface DxfSolid {
  Type: "SOLID";
  Layer: string;
  Color?: number;
  X1: number; Y1: number;  // First corner
  X2: number; Y2: number;  // Second corner
  X3: number; Y3: number;  // Third corner
  X4: number; Y4: number;  // Fourth corner (= third for triangles)
}
```

**Note:** DXF SOLID uses a specific vertex order (1-2-4-3) for quadrilaterals:
```
1 -------- 2
|          |
4 -------- 3
```

**Example:**
```json
{
  "Type": "SOLID",
  "Layer": "0-0",
  "Color": 5,
  "X1": 0.0, "Y1": 0.0,
  "X2": 10.0, "Y2": 0.0,
  "X3": 5.0, "Y3": 10.0,
  "X4": 5.0, "Y4": 10.0
}
```

### INSERT

Represents a block reference (instance).

```typescript
interface DxfInsert {
  Type: "INSERT";
  Layer: string;
  Color?: number;
  BlockName: string;   // Name of referenced block
  X: number;           // Insertion X
  Y: number;           // Insertion Y
  Z?: number;
  ScaleX?: number;     // X scale factor (default 1)
  ScaleY?: number;     // Y scale factor (default 1)
  ScaleZ?: number;     // Z scale factor (default 1)
  Rotation?: number;   // Rotation angle in degrees
}
```

**Example:**
```json
{
  "Type": "INSERT",
  "Layer": "0-0",
  "BlockName": "BLOCK_1",
  "X": 100.0,
  "Y": 50.0,
  "ScaleX": 1.0,
  "ScaleY": 1.0,
  "Rotation": 0.0
}
```

### LWPOLYLINE

Represents a lightweight polyline (2D).

```typescript
interface DxfLwPolyline {
  Type: "LWPOLYLINE";
  Layer: string;
  Color?: number;
  Closed: boolean;
  Vertices: DxfVertex[];
}

interface DxfVertex {
  X: number;
  Y: number;
  Z?: number;
  Bulge?: number;  // Bulge factor for curved segments
}
```

**Bulge Factor:**
- 0 = straight line to next vertex
- Positive = counter-clockwise arc
- Negative = clockwise arc
- |bulge| = tan(arc_angle / 4)
- |bulge| = 1 means a semicircle

## Blocks

Block definitions contain reusable geometry.

```typescript
interface DxfBlock {
  Name: string;
  BaseX?: number;  // Base point X
  BaseY?: number;  // Base point Y
  BaseZ?: number;  // Base point Z
  Entities: DxfEntity[];
}
```

**Example:**
```json
{
  "Name": "SYMBOL_A",
  "BaseX": 0.0,
  "BaseY": 0.0,
  "Entities": [
    {
      "Type": "LINE",
      "Layer": "0",
      "X1": -5.0, "Y1": 0.0,
      "X2": 5.0, "Y2": 0.0
    },
    {
      "Type": "CIRCLE",
      "Layer": "0",
      "CenterX": 0.0, "CenterY": 0.0,
      "Radius": 3.0
    }
  ]
}
```

## Color Mapping

JWW colors are converted to DXF ACI (AutoCAD Color Index) colors:

| JWW Color | Description | DXF ACI |
|-----------|-------------|---------|
| 1 | Black/White | 7 |
| 2 | Blue | 5 |
| 3 | Red | 1 |
| 4 | Magenta | 6 |
| 5 | Green | 3 |
| 6 | Cyan | 4 |
| 7 | Yellow | 2 |
| 8 | White/Black | 7 |
| 9 | Gray | 8 |
| 100+ | Custom colors | Mapped to closest ACI |

## Coordinate System

- All coordinates use millimeters as the unit
- The coordinate system is right-handed:
  - X increases to the right
  - Y increases upward
  - Z increases toward the viewer (out of the screen)
- Origin (0, 0) is typically at the drawing's reference point

## Complete Example

```json
{
  "Layers": [
    {"Name": "0-0", "Color": 7, "Frozen": false, "Locked": false},
    {"Name": "0-1", "Color": 1, "Frozen": false, "Locked": false}
  ],
  "Entities": [
    {
      "Type": "LINE",
      "Layer": "0-0",
      "Color": 7,
      "X1": 0.0, "Y1": 0.0,
      "X2": 100.0, "Y2": 0.0
    },
    {
      "Type": "LINE",
      "Layer": "0-0",
      "Color": 7,
      "X1": 100.0, "Y1": 0.0,
      "X2": 100.0, "Y2": 100.0
    },
    {
      "Type": "CIRCLE",
      "Layer": "0-1",
      "Color": 1,
      "CenterX": 50.0, "CenterY": 50.0,
      "Radius": 30.0
    },
    {
      "Type": "TEXT",
      "Layer": "0-0",
      "Color": 7,
      "X": 10.0, "Y": 110.0,
      "Height": 5.0,
      "Content": "Sample Drawing"
    }
  ],
  "Blocks": []
}
```
