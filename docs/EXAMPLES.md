# Usage Examples

This document provides practical examples for using jww-parser in various scenarios.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Node.js Examples](#nodejs-examples)
- [Browser Examples](#browser-examples)
- [Three.js Integration](#threejs-integration)
- [React Integration](#react-integration)
- [Error Handling](#error-handling)
- [Progress Tracking](#progress-tracking)
- [Performance Optimization](#performance-optimization)
- [Batch Processing](#batch-processing)

## Basic Usage

### Parsing a JWW File

```typescript
import { createParser } from 'jww-parser';
import { readFileSync } from 'fs';

async function parseJwwFile(filePath: string) {
  // Read the file
  const buffer = readFileSync(filePath);
  const data = new Uint8Array(buffer);

  // Create parser and parse
  const parser = await createParser();
  const doc = parser.parse(data);

  console.log(`Version: ${doc.Version}`);
  console.log(`Entities: ${doc.Entities.length}`);
  console.log(`Blocks: ${doc.Blocks.length}`);

  // Cleanup
  parser.dispose();

  return doc;
}
```

### Converting to DXF

```typescript
import { createParser } from 'jww-parser';
import { readFileSync, writeFileSync } from 'fs';

async function convertToDxf(inputPath: string, outputPath: string) {
  const buffer = readFileSync(inputPath);
  const data = new Uint8Array(buffer);

  const parser = await createParser();

  // Get DXF as string (ready to save)
  const dxfContent = parser.toDxfString(data);
  writeFileSync(outputPath, dxfContent);

  // Or get as JSON for further processing
  const dxfDoc = parser.toDxf(data);
  console.log(`Converted ${dxfDoc.Entities.length} entities`);

  parser.dispose();
}
```

## Node.js Examples

### CLI Tool

```typescript
#!/usr/bin/env node
import { createParser, isJwwFile, quickValidate } from 'jww-parser';
import { readFileSync, writeFileSync } from 'fs';
import { basename, extname } from 'path';

async function main() {
  const args = process.argv.slice(2);

  if (args.length < 1) {
    console.log('Usage: jww-convert <input.jww> [output.dxf]');
    process.exit(1);
  }

  const inputPath = args[0];
  const outputPath = args[1] || inputPath.replace(/\.jww$/i, '.dxf');

  // Read and validate
  const data = new Uint8Array(readFileSync(inputPath));

  if (!isJwwFile(data)) {
    console.error('Error: Not a valid JWW file');
    process.exit(1);
  }

  const validation = quickValidate(data);
  if (!validation.valid) {
    console.error('Validation errors:');
    validation.issues.forEach(issue => {
      console.error(`  [${issue.severity}] ${issue.message}`);
    });
    process.exit(1);
  }

  console.log(`JWW Version: ${validation.version}`);
  console.log(`Estimated entities: ~${validation.estimatedEntityCount}`);

  // Convert
  const parser = await createParser();
  const dxfContent = parser.toDxfString(data);
  writeFileSync(outputPath, dxfContent);

  console.log(`Converted to: ${outputPath}`);

  parser.dispose();
}

main().catch(console.error);
```

### Batch Processing

```typescript
import { createParser } from 'jww-parser';
import { readFileSync, writeFileSync, readdirSync } from 'fs';
import { join, extname } from 'path';

async function batchConvert(inputDir: string, outputDir: string) {
  const parser = await createParser();
  const files = readdirSync(inputDir).filter(f =>
    extname(f).toLowerCase() === '.jww'
  );

  console.log(`Found ${files.length} JWW files`);

  let converted = 0;
  let failed = 0;

  for (const file of files) {
    const inputPath = join(inputDir, file);
    const outputPath = join(outputDir, file.replace(/\.jww$/i, '.dxf'));

    try {
      const data = new Uint8Array(readFileSync(inputPath));
      const dxfContent = parser.toDxfString(data);
      writeFileSync(outputPath, dxfContent);
      converted++;
      console.log(`✓ ${file}`);
    } catch (error) {
      failed++;
      console.error(`✗ ${file}: ${error.message}`);
    }
  }

  console.log(`\nCompleted: ${converted} converted, ${failed} failed`);

  const stats = parser.getStats();
  console.log(`Average parse time: ${stats.averageParseTimeMs.toFixed(2)}ms`);

  parser.dispose();
}
```

## Browser Examples

### HTML Setup

```html
<!DOCTYPE html>
<html>
<head>
  <title>JWW Viewer</title>
  <!-- Required: Include wasm_exec.js before using the parser -->
  <script src="node_modules/jww-parser/wasm/wasm_exec.js"></script>
</head>
<body>
  <input type="file" id="fileInput" accept=".jww">
  <div id="output"></div>

  <script type="module">
    import { createParser, isJwwFile } from 'jww-parser';

    const fileInput = document.getElementById('fileInput');
    const output = document.getElementById('output');

    fileInput.addEventListener('change', async (e) => {
      const file = e.target.files[0];
      if (!file) return;

      const arrayBuffer = await file.arrayBuffer();
      const data = new Uint8Array(arrayBuffer);

      if (!isJwwFile(data)) {
        output.textContent = 'Error: Not a valid JWW file';
        return;
      }

      try {
        const parser = await createParser('jww-parser.wasm');
        const doc = parser.parse(data);

        output.innerHTML = `
          <h3>File Info</h3>
          <p>Version: ${doc.Version}</p>
          <p>Entities: ${doc.Entities.length}</p>
          <p>Blocks: ${doc.Blocks.length}</p>
        `;

        parser.dispose();
      } catch (error) {
        output.textContent = `Error: ${error.message}`;
      }
    });
  </script>
</body>
</html>
```

### Download DXF

```typescript
async function downloadAsDxf(file: File) {
  const data = new Uint8Array(await file.arrayBuffer());

  const parser = await createParser();
  const dxfContent = parser.toDxfString(data);
  parser.dispose();

  // Create download
  const blob = new Blob([dxfContent], { type: 'application/dxf' });
  const url = URL.createObjectURL(blob);

  const a = document.createElement('a');
  a.href = url;
  a.download = file.name.replace(/\.jww$/i, '.dxf');
  a.click();

  URL.revokeObjectURL(url);
}
```

## Three.js Integration

### Basic Setup

```typescript
import * as THREE from 'three';
import { createParser, DxfDocument, DxfEntity } from 'jww-parser';

class JwwThreeViewer {
  private scene: THREE.Scene;
  private camera: THREE.OrthographicCamera;
  private renderer: THREE.WebGLRenderer;

  constructor(container: HTMLElement) {
    this.scene = new THREE.Scene();
    this.scene.background = new THREE.Color(0x1a1a2e);

    // Orthographic camera for 2D CAD viewing
    const aspect = container.clientWidth / container.clientHeight;
    this.camera = new THREE.OrthographicCamera(
      -100 * aspect, 100 * aspect, 100, -100, 0.1, 1000
    );
    this.camera.position.z = 10;

    this.renderer = new THREE.WebGLRenderer({ antialias: true });
    this.renderer.setSize(container.clientWidth, container.clientHeight);
    container.appendChild(this.renderer.domElement);
  }

  async loadJww(data: Uint8Array) {
    const parser = await createParser();
    const dxfDoc = parser.toDxf(data);
    parser.dispose();

    this.clearScene();
    this.addEntities(dxfDoc);
    this.fitToView();
    this.render();
  }

  private clearScene() {
    while (this.scene.children.length > 0) {
      this.scene.remove(this.scene.children[0]);
    }
  }

  private addEntities(doc: DxfDocument) {
    for (const entity of doc.Entities) {
      const object = this.createObject(entity);
      if (object) {
        this.scene.add(object);
      }
    }
  }

  private createObject(entity: DxfEntity): THREE.Object3D | null {
    const color = this.getColor(entity.Color);

    switch (entity.Type) {
      case 'LINE':
        return this.createLine(entity, color);
      case 'CIRCLE':
        return this.createCircle(entity, color);
      case 'ARC':
        return this.createArc(entity, color);
      case 'ELLIPSE':
        return this.createEllipse(entity, color);
      case 'TEXT':
        return this.createText(entity, color);
      default:
        return null;
    }
  }

  private createLine(entity: any, color: number): THREE.Line {
    const geometry = new THREE.BufferGeometry();
    const points = [
      new THREE.Vector3(entity.X1, entity.Y1, 0),
      new THREE.Vector3(entity.X2, entity.Y2, 0)
    ];
    geometry.setFromPoints(points);

    const material = new THREE.LineBasicMaterial({ color });
    return new THREE.Line(geometry, material);
  }

  private createCircle(entity: any, color: number): THREE.Line {
    const geometry = new THREE.BufferGeometry();
    const points: THREE.Vector3[] = [];
    const segments = 64;

    for (let i = 0; i <= segments; i++) {
      const angle = (i / segments) * Math.PI * 2;
      points.push(new THREE.Vector3(
        entity.CenterX + Math.cos(angle) * entity.Radius,
        entity.CenterY + Math.sin(angle) * entity.Radius,
        0
      ));
    }

    geometry.setFromPoints(points);
    const material = new THREE.LineBasicMaterial({ color });
    return new THREE.Line(geometry, material);
  }

  private createArc(entity: any, color: number): THREE.Line {
    const geometry = new THREE.BufferGeometry();
    const points: THREE.Vector3[] = [];
    const segments = 64;

    const startRad = (entity.StartAngle * Math.PI) / 180;
    const endRad = (entity.EndAngle * Math.PI) / 180;
    let sweep = endRad - startRad;
    if (sweep < 0) sweep += Math.PI * 2;

    for (let i = 0; i <= segments; i++) {
      const angle = startRad + (i / segments) * sweep;
      points.push(new THREE.Vector3(
        entity.CenterX + Math.cos(angle) * entity.Radius,
        entity.CenterY + Math.sin(angle) * entity.Radius,
        0
      ));
    }

    geometry.setFromPoints(points);
    const material = new THREE.LineBasicMaterial({ color });
    return new THREE.Line(geometry, material);
  }

  private createEllipse(entity: any, color: number): THREE.Line {
    const geometry = new THREE.BufferGeometry();
    const points: THREE.Vector3[] = [];
    const segments = 64;

    const majorLength = Math.sqrt(
      entity.MajorAxisX ** 2 + entity.MajorAxisY ** 2
    );
    const minorLength = majorLength * entity.MinorRatio;
    const rotation = Math.atan2(entity.MajorAxisY, entity.MajorAxisX);

    const startParam = entity.StartParam;
    const endParam = entity.EndParam;
    let sweep = endParam - startParam;
    if (sweep < 0) sweep += Math.PI * 2;

    for (let i = 0; i <= segments; i++) {
      const t = startParam + (i / segments) * sweep;
      const x = majorLength * Math.cos(t);
      const y = minorLength * Math.sin(t);

      points.push(new THREE.Vector3(
        entity.CenterX + x * Math.cos(rotation) - y * Math.sin(rotation),
        entity.CenterY + x * Math.sin(rotation) + y * Math.cos(rotation),
        0
      ));
    }

    geometry.setFromPoints(points);
    const material = new THREE.LineBasicMaterial({ color });
    return new THREE.Line(geometry, material);
  }

  private createText(entity: any, color: number): THREE.Object3D | null {
    // For simple text, use sprite or TextGeometry
    // This is a placeholder - real implementation would use
    // THREE.TextGeometry or troika-three-text
    return null;
  }

  private getColor(aciColor?: number): number {
    const colorMap: Record<number, number> = {
      1: 0xff0000,  // Red
      2: 0xffff00,  // Yellow
      3: 0x00ff00,  // Green
      4: 0x00ffff,  // Cyan
      5: 0x0000ff,  // Blue
      6: 0xff00ff,  // Magenta
      7: 0xffffff,  // White
      8: 0x808080,  // Gray
    };
    return colorMap[aciColor || 7] || 0xffffff;
  }

  private fitToView() {
    const box = new THREE.Box3().setFromObject(this.scene);
    const size = box.getSize(new THREE.Vector3());
    const center = box.getCenter(new THREE.Vector3());

    const maxDim = Math.max(size.x, size.y);
    const aspect = this.renderer.domElement.width /
                   this.renderer.domElement.height;

    const padding = 1.1;
    this.camera.left = -maxDim * padding * aspect / 2;
    this.camera.right = maxDim * padding * aspect / 2;
    this.camera.top = maxDim * padding / 2;
    this.camera.bottom = -maxDim * padding / 2;

    this.camera.position.set(center.x, center.y, 10);
    this.camera.lookAt(center);
    this.camera.updateProjectionMatrix();
  }

  render() {
    this.renderer.render(this.scene, this.camera);
  }
}

// Usage
const container = document.getElementById('viewer')!;
const viewer = new JwwThreeViewer(container);

fileInput.addEventListener('change', async (e) => {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (file) {
    const data = new Uint8Array(await file.arrayBuffer());
    await viewer.loadJww(data);
  }
});
```

### Advanced Three.js with OrbitControls

```typescript
import * as THREE from 'three';
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls';

class AdvancedJwwViewer extends JwwThreeViewer {
  private controls: OrbitControls;

  constructor(container: HTMLElement) {
    super(container);

    // Add orbit controls for pan/zoom
    this.controls = new OrbitControls(
      this.camera,
      this.renderer.domElement
    );
    this.controls.enableRotate = false; // 2D only
    this.controls.screenSpacePanning = true;

    this.controls.addEventListener('change', () => this.render());

    // Zoom with mouse wheel
    this.controls.enableZoom = true;
    this.controls.zoomSpeed = 1.2;
  }

  // Add layer visibility control
  setLayerVisibility(layerName: string, visible: boolean) {
    this.scene.traverse((object) => {
      if (object.userData.layer === layerName) {
        object.visible = visible;
      }
    });
    this.render();
  }

  // Color override for a layer
  setLayerColor(layerName: string, color: number) {
    this.scene.traverse((object) => {
      if (object.userData.layer === layerName &&
          object instanceof THREE.Line) {
        (object.material as THREE.LineBasicMaterial).color.setHex(color);
      }
    });
    this.render();
  }
}
```

## React Integration

```tsx
import React, { useEffect, useRef, useState } from 'react';
import { createParser, JwwParser, JwwDocument, isJwwFile } from 'jww-parser';

interface JwwViewerProps {
  file: File | null;
}

export function JwwViewer({ file }: JwwViewerProps) {
  const [parser, setParser] = useState<JwwParser | null>(null);
  const [document, setDocument] = useState<JwwDocument | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  // Initialize parser once
  useEffect(() => {
    let mounted = true;

    createParser().then(p => {
      if (mounted) setParser(p);
    }).catch(err => {
      if (mounted) setError(`Failed to initialize: ${err.message}`);
    });

    return () => {
      mounted = false;
      parser?.dispose();
    };
  }, []);

  // Parse file when it changes
  useEffect(() => {
    if (!file || !parser) return;

    setLoading(true);
    setError(null);

    file.arrayBuffer().then(buffer => {
      const data = new Uint8Array(buffer);

      if (!isJwwFile(data)) {
        setError('Not a valid JWW file');
        setLoading(false);
        return;
      }

      try {
        const doc = parser.parse(data);
        setDocument(doc);
      } catch (err: any) {
        setError(err.message);
      }
      setLoading(false);
    });
  }, [file, parser]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!document) return <div>Select a JWW file</div>;

  return (
    <div className="jww-viewer">
      <h3>Document Info</h3>
      <p>Version: {document.Version}</p>
      <p>Entities: {document.Entities.length}</p>

      <h4>Entity Types</h4>
      <ul>
        {Object.entries(
          document.Entities.reduce((acc, e) => {
            acc[e.Type] = (acc[e.Type] || 0) + 1;
            return acc;
          }, {} as Record<string, number>)
        ).map(([type, count]) => (
          <li key={type}>{type}: {count}</li>
        ))}
      </ul>
    </div>
  );
}
```

## Error Handling

```typescript
import {
  createParser,
  JwwParserError,
  NotInitializedError,
  WasmLoadError,
  ValidationError,
  ParseError,
  JwwErrorCode
} from 'jww-parser';

async function safeConvert(data: Uint8Array) {
  try {
    const parser = await createParser();

    try {
      parser.validateOrThrow(data);
      return parser.toDxf(data);
    } finally {
      parser.dispose();
    }

  } catch (error) {
    if (error instanceof NotInitializedError) {
      console.error('Parser not initialized');
    } else if (error instanceof WasmLoadError) {
      console.error('WASM loading failed:', error.message);
    } else if (error instanceof ValidationError) {
      console.error('Validation failed:');
      error.issues.forEach(issue => {
        console.error(`  [${issue.severity}] ${issue.code}: ${issue.message}`);
      });
    } else if (error instanceof ParseError) {
      console.error(`Parse error at offset ${error.offset}:`, error.message);
    } else if (error instanceof JwwParserError) {
      console.error(`Error [${error.code}]:`, error.toDetailedString());
    } else {
      throw error;
    }
    return null;
  }
}
```

## Progress Tracking

```typescript
import { createParser, ProgressInfo } from 'jww-parser';

async function parseWithProgress(data: Uint8Array) {
  const parser = await createParser();

  const progressBar = document.querySelector('.progress-bar') as HTMLElement;
  const statusText = document.querySelector('.status-text') as HTMLElement;

  const doc = parser.parse(data, {
    onProgress: (info: ProgressInfo) => {
      // Update progress bar
      progressBar.style.width = `${info.overallProgress}%`;

      // Update status text
      statusText.textContent = info.message || `Stage: ${info.stage}`;

      // Log timing
      console.log(`[${info.elapsedMs}ms] ${info.stage}: ${info.progress}%`);

      // Estimated time remaining
      if (info.estimatedRemainingMs) {
        console.log(`ETA: ${(info.estimatedRemainingMs / 1000).toFixed(1)}s`);
      }
    }
  });

  progressBar.style.width = '100%';
  statusText.textContent = 'Complete!';

  parser.dispose();
  return doc;
}
```

## Performance Optimization

```typescript
import { createParser, ParseOptions } from 'jww-parser';

// Preview mode - fast loading with limited entities
async function previewFile(data: Uint8Array) {
  const parser = await createParser();

  const options: ParseOptions = {
    maxEntities: 500,        // Limit entities
    includeBlocks: false,    // Skip block definitions
    skipEntityTypes: ['Point', 'Solid']  // Skip heavy types
  };

  const doc = parser.parse(data, options);
  parser.dispose();
  return doc;
}

// Layer-specific loading
async function loadSpecificLayers(data: Uint8Array, layers: number[]) {
  const parser = await createParser();

  const doc = parser.parse(data, {
    layerGroupFilter: [0, 1],  // Only groups 0 and 1
    layerFilter: {
      0: [0, 1, 2],   // Layers 0-2 from group 0
      1: [5, 6]       // Layers 5-6 from group 1
    }
  });

  parser.dispose();
  return doc;
}

// With abort controller
async function parseWithTimeout(data: Uint8Array, timeoutMs: number) {
  const parser = await createParser();
  const controller = new AbortController();

  const timeout = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const doc = parser.parse(data, {
      signal: controller.signal
    });
    clearTimeout(timeout);
    return doc;
  } catch (error) {
    if (error.message.includes('aborted')) {
      console.log('Parse operation timed out');
    }
    throw error;
  } finally {
    parser.dispose();
  }
}
```
