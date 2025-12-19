# JWW Parser API Reference

This document provides comprehensive API documentation for the jww-parser library.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [JwwParser Class](#jwwparser-class)
- [Factory Functions](#factory-functions)
- [Type Definitions](#type-definitions)
- [Error Classes](#error-classes)
- [Options](#options)

## Installation

```bash
npm install jww-parser
```

## Quick Start

```typescript
import { createParser, isJwwFile } from 'jww-parser';

// Read file as Uint8Array
const fileData = new Uint8Array(await file.arrayBuffer());

// Check if it's a JWW file
if (!isJwwFile(fileData)) {
  console.error('Not a valid JWW file');
  return;
}

// Create and initialize parser
const parser = await createParser();

// Parse to JWW document
const jwwDoc = parser.parse(fileData);

// Or convert directly to DXF
const dxfDoc = parser.toDxf(fileData);
const dxfString = parser.toDxfString(fileData);
```

## JwwParser Class

The main class for parsing JWW files.

### Constructor

```typescript
constructor(wasmPath?: string)
```

Creates a new parser instance.

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| wasmPath | string | auto-detected | Path to jww-parser.wasm file |

### Methods

#### `init(): Promise<void>`

Initializes the WebAssembly module. Must be called before parsing.

```typescript
const parser = new JwwParser();
await parser.init();
```

#### `isInitialized(): boolean`

Returns whether the parser has been initialized.

#### `parse(data: Uint8Array, options?: ParseOptions): JwwDocument`

Parses a JWW file and returns the document structure.

```typescript
const doc = parser.parse(fileData, {
  maxEntities: 1000,
  skipEntityTypes: ['Point'],
  onProgress: (info) => console.log(`${info.overallProgress}%`)
});
```

#### `toDxf(data: Uint8Array, options?: ConvertOptions): DxfDocument`

Parses a JWW file and converts to DXF document structure (JSON).

```typescript
const dxfDoc = parser.toDxf(fileData);
console.log(`Entities: ${dxfDoc.Entities.length}`);
```

#### `toDxfString(data: Uint8Array, options?: ConvertOptions): string`

Parses and converts to DXF file content (ready to save as .dxf).

```typescript
const dxfContent = parser.toDxfString(fileData);
fs.writeFileSync('output.dxf', dxfContent);
```

#### `validate(data: Uint8Array): ValidationResult`

Validates a JWW file without fully parsing it. Does not require initialization.

```typescript
const result = parser.validate(fileData);
if (!result.valid) {
  console.error('Validation failed:', result.issues);
}
```

#### `validateOrThrow(data: Uint8Array): void`

Validates and throws `ValidationError` if invalid.

```typescript
try {
  parser.validateOrThrow(fileData);
} catch (error) {
  if (error instanceof ValidationError) {
    console.error('Issues:', error.issues);
  }
}
```

#### `setDebug(options: DebugOptions | boolean): void`

Enables or configures debug mode.

```typescript
// Simple enable
parser.setDebug(true);

// With options
parser.setDebug({
  enabled: true,
  logLevel: 'debug',
  logToConsole: true,
  includeMemoryUsage: true
});
```

#### `getDebugLogs(level?: LogLevel): DebugLogEntry[]`

Gets debug log entries.

```typescript
const errors = parser.getDebugLogs('error');
const allLogs = parser.getDebugLogs();
```

#### `clearDebugLogs(): void`

Clears stored debug logs.

#### `getMemoryStats(): MemoryStats`

Gets current memory usage statistics.

```typescript
const stats = parser.getMemoryStats();
console.log(`WASM memory: ${stats.totalFormatted}`);
```

#### `getStats(): ParserStats`

Gets parser performance statistics.

```typescript
const stats = parser.getStats();
console.log(`Average parse time: ${stats.averageParseTimeMs}ms`);
console.log(`Total bytes processed: ${stats.totalBytesProcessed}`);
```

#### `resetStats(): void`

Resets parser statistics.

#### `dispose(): void`

Cleans up resources and releases memory.

```typescript
parser.dispose();
// Parser is no longer usable after dispose
```

#### `getVersion(): string`

Gets the WASM module version.

## Factory Functions

### `createParser(wasmPath?: string, options?: { debug?: DebugOptions }): Promise<JwwParser>`

Creates and initializes a parser instance.

```typescript
const parser = await createParser(undefined, {
  debug: { enabled: true, logToConsole: true }
});
```

### `quickValidate(data: Uint8Array): ValidationResult`

Validates a file without initializing a full parser.

```typescript
const result = quickValidate(fileData);
if (result.valid) {
  console.log(`Version: ${result.version}`);
}
```

### `isJwwFile(data: Uint8Array): boolean`

Quick check if data appears to be a JWW file.

```typescript
if (isJwwFile(fileData)) {
  // Proceed with parsing
}
```

## Type Definitions

### JWW Document Types

#### `JwwDocument`

```typescript
interface JwwDocument {
  Version: number;        // JWW file version
  Memo: string;           // Document memo
  PaperSize: number;      // Paper size code
  LayerGroups: JwwLayerGroup[];
  Entities: JwwEntity[];
  Blocks: JwwBlock[];
}
```

#### `JwwEntity` (Union Type)

```typescript
type JwwEntity =
  | JwwLine
  | JwwArc
  | JwwPoint
  | JwwText
  | JwwSolid
  | JwwBlockRef;
```

#### Entity Base Properties

All entities share these base properties:

```typescript
interface JwwEntityBase {
  Type: string;       // Entity type discriminator
  Group: number;      // Curve attribute number
  PenStyle: number;   // Line type (0-9)
  PenColor: number;   // Color code (1-256+)
  PenWidth: number;   // Line width
  Layer: number;      // Layer index (0-15)
  LayerGroup: number; // Layer group index (0-15)
}
```

### DXF Document Types

#### `DxfDocument`

```typescript
interface DxfDocument {
  Layers: DxfLayer[];
  Entities: DxfEntity[];
  Blocks: DxfBlock[];
}
```

#### `DxfEntity` (Union Type)

```typescript
type DxfEntity =
  | DxfLine
  | DxfCircle
  | DxfArc
  | DxfEllipse
  | DxfPoint
  | DxfText
  | DxfMText
  | DxfSolid
  | DxfInsert
  | DxfLwPolyline
  | (DxfEntityBase & Record<string, unknown>);
```

### Validation Types

#### `ValidationResult`

```typescript
interface ValidationResult {
  valid: boolean;
  version?: number;
  sizeCategory: 'small' | 'medium' | 'large' | 'very_large';
  estimatedEntityCount?: number;
  issues: ValidationIssue[];
  validationTimeMs: number;
}
```

#### `ValidationIssue`

```typescript
interface ValidationIssue {
  severity: 'error' | 'warning' | 'info';
  code: string;
  message: string;
  offset?: number;
  details?: Record<string, unknown>;
}
```

### Progress Types

#### `ProgressInfo`

```typescript
interface ProgressInfo {
  stage: ProgressStage;
  progress: number;           // 0-100 within stage
  overallProgress: number;    // 0-100 overall
  message?: string;
  entitiesProcessed?: number;
  totalEntities?: number;
  elapsedMs: number;
  estimatedRemainingMs?: number;
}
```

#### `ProgressStage`

```typescript
type ProgressStage =
  | 'loading'
  | 'parsing_header'
  | 'parsing_layers'
  | 'parsing_entities'
  | 'parsing_blocks'
  | 'converting'
  | 'complete';
```

## Error Classes

### `JwwParserError`

Base error class for all parser errors.

```typescript
class JwwParserError extends Error {
  readonly code: JwwErrorCode;
  readonly cause?: Error;
  readonly context?: Record<string, unknown>;
  readonly timestamp: Date;

  toDetailedString(): string;
  toJSON(): Record<string, unknown>;
}
```

### Error Codes

```typescript
enum JwwErrorCode {
  NOT_INITIALIZED = 'NOT_INITIALIZED',
  WASM_LOAD_FAILED = 'WASM_LOAD_FAILED',
  WASM_TIMEOUT = 'WASM_TIMEOUT',
  INVALID_SIGNATURE = 'INVALID_SIGNATURE',
  UNSUPPORTED_VERSION = 'UNSUPPORTED_VERSION',
  PARSE_ERROR = 'PARSE_ERROR',
  CONVERSION_ERROR = 'CONVERSION_ERROR',
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  MEMORY_ERROR = 'MEMORY_ERROR',
  INVALID_ARGUMENT = 'INVALID_ARGUMENT',
}
```

### Specialized Error Classes

- `NotInitializedError` - Parser not initialized
- `WasmLoadError` - WASM module loading failed
- `ValidationError` - File validation failed
- `ParseError` - Parsing failed

## Options

### `ParseOptions`

```typescript
interface ParseOptions {
  streamingMode?: boolean;              // Low memory mode
  skipEntityTypes?: JwwEntity['Type'][]; // Skip specific types
  layerGroupFilter?: number[];          // Filter by layer groups
  layerFilter?: Record<number, number[]>; // Filter by layers
  maxEntities?: number;                 // Limit entity count
  includeBlocks?: boolean;              // Include blocks (default: true)
  expandBlockReferences?: boolean;      // Expand INSERT to entities
  onProgress?: ProgressCallback;        // Progress callback
  signal?: AbortSignal;                 // For cancellation
}
```

### `ConvertOptions`

Extends `ParseOptions` with:

```typescript
interface ConvertOptions extends ParseOptions {
  includeTemporaryPoints?: boolean;     // Convert temp points
  colorMapping?: Record<number, number>; // Color override
  layerNamePattern?: string;            // Layer naming pattern
  precision?: number;                   // Coordinate precision
}
```

### `DebugOptions`

```typescript
interface DebugOptions {
  enabled?: boolean;              // Enable debug mode
  logLevel?: LogLevel;            // Minimum log level
  onDebug?: DebugCallback;        // Custom callback
  logToConsole?: boolean;         // Log to console
  maxLogEntries?: number;         // Max stored entries
  includeTiming?: boolean;        // Include timing info
  includeMemoryUsage?: boolean;   // Include memory stats
}
```
