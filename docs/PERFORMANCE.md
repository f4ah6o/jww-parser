# Performance Tuning Guide

This guide provides best practices and optimization techniques for using jww-parser with large files and in performance-critical applications.

## Table of Contents

- [Understanding Performance Characteristics](#understanding-performance-characteristics)
- [Memory Management](#memory-management)
- [Parse Options for Performance](#parse-options-for-performance)
- [Batch Processing](#batch-processing)
- [Browser Optimization](#browser-optimization)
- [Node.js Optimization](#nodejs-optimization)
- [Benchmarking](#benchmarking)

## Understanding Performance Characteristics

### File Size Categories

| Category | Size | Typical Entities | Parse Time | Memory |
|----------|------|------------------|------------|--------|
| Small | < 100KB | < 1,000 | < 50ms | < 10MB |
| Medium | 100KB - 1MB | 1,000 - 10,000 | 50-200ms | 10-50MB |
| Large | 1MB - 10MB | 10,000 - 100,000 | 200ms-2s | 50-200MB |
| Very Large | > 10MB | > 100,000 | > 2s | > 200MB |

### Performance Bottlenecks

1. **WASM Initialization** (~100-500ms cold start)
2. **Binary Parsing** (O(n) where n = file size)
3. **JSON Serialization** (O(m) where m = entity count)
4. **JavaScript Object Creation** (memory-intensive)

## Memory Management

### Single Parser Instance

Reuse a single parser instance for multiple operations:

```typescript
// BAD: Creating new parser for each file
for (const file of files) {
  const parser = await createParser();  // Expensive!
  const doc = parser.parse(file);
  parser.dispose();
}

// GOOD: Reuse parser instance
const parser = await createParser();
for (const file of files) {
  const doc = parser.parse(file);
  // Process doc...
}
parser.dispose();
```

### Monitor Memory Usage

```typescript
const parser = await createParser();

// Check memory periodically
setInterval(() => {
  const stats = parser.getMemoryStats();
  console.log(`WASM Memory: ${stats.totalFormatted}`);

  if (stats.totalBytes > 500 * 1024 * 1024) { // 500MB threshold
    console.warn('High memory usage detected');
    // Consider disposing and reinitializing
  }
}, 10000);
```

### Dispose When Done

```typescript
async function processFiles(files: File[]) {
  const parser = await createParser();

  try {
    for (const file of files) {
      const data = new Uint8Array(await file.arrayBuffer());
      const doc = parser.parse(data);
      await processDocument(doc);

      // Allow garbage collection
      // doc = null not needed in function scope
    }
  } finally {
    parser.dispose();  // Always cleanup
  }
}
```

### Chunked Processing for Large Files

```typescript
async function processLargeFile(data: Uint8Array) {
  const parser = await createParser();

  // First, get basic info without full parse
  const validation = parser.validate(data);

  if (validation.sizeCategory === 'very_large') {
    console.log('Large file detected, using chunked approach');

    // Parse with entity limit for preview
    const preview = parser.parse(data, {
      maxEntities: 1000,
      includeBlocks: false
    });

    // Display preview to user...

    // Then process full document if needed
    const fullDoc = parser.parse(data);
    // ...
  }

  parser.dispose();
}
```

## Parse Options for Performance

### Skip Unnecessary Entity Types

```typescript
// If you only need lines and arcs
const doc = parser.parse(data, {
  skipEntityTypes: ['Point', 'Text', 'Solid', 'Block']
});
```

### Filter by Layer

```typescript
// Only parse visible layers
const doc = parser.parse(data, {
  layerGroupFilter: [0, 1],  // Groups 0 and 1 only
  layerFilter: {
    0: [0, 1, 2, 3],  // Specific layers in group 0
  }
});
```

### Limit Entity Count for Previews

```typescript
// Quick preview
const preview = parser.parse(data, {
  maxEntities: 500,
  includeBlocks: false
});

// Full parse only if user confirms
if (userConfirms) {
  const full = parser.parse(data);
}
```

### Skip Blocks When Not Needed

```typescript
// If you don't need block definitions
const doc = parser.parse(data, {
  includeBlocks: false
});
```

## Batch Processing

### Optimal Batch Size

```typescript
async function processBatch(files: Uint8Array[], batchSize = 10) {
  const parser = await createParser();
  const results: JwwDocument[] = [];

  // Process in batches
  for (let i = 0; i < files.length; i += batchSize) {
    const batch = files.slice(i, i + batchSize);

    // Allow GC between batches
    await new Promise(resolve => setTimeout(resolve, 0));

    for (const file of batch) {
      results.push(parser.parse(file));
    }
  }

  parser.dispose();
  return results;
}
```

### Parallel Processing with Worker Threads (Node.js)

```typescript
// worker.ts
import { parentPort, workerData } from 'worker_threads';
import { createParser } from 'jww-parser';

async function process() {
  const parser = await createParser();
  const result = parser.toDxfString(workerData.data);
  parser.dispose();
  parentPort?.postMessage(result);
}

process();

// main.ts
import { Worker } from 'worker_threads';

async function parallelConvert(files: Uint8Array[]) {
  const workers = files.map(data => {
    return new Promise<string>((resolve, reject) => {
      const worker = new Worker('./worker.ts', {
        workerData: { data }
      });
      worker.on('message', resolve);
      worker.on('error', reject);
    });
  });

  return Promise.all(workers);
}
```

## Browser Optimization

### Lazy Load WASM

```typescript
let parserPromise: Promise<JwwParser> | null = null;

function getParser(): Promise<JwwParser> {
  if (!parserPromise) {
    parserPromise = createParser();
  }
  return parserPromise;
}

// Use when needed
async function handleFile(file: File) {
  const parser = await getParser();
  // ...
}
```

### Web Worker for Heavy Processing

```typescript
// jww-worker.ts
self.onmessage = async (e) => {
  const { createParser } = await import('jww-parser');
  const parser = await createParser();

  const result = parser.parse(e.data);
  self.postMessage(result);

  parser.dispose();
};

// main.ts
const worker = new Worker(new URL('./jww-worker.ts', import.meta.url));

function parseInWorker(data: Uint8Array): Promise<JwwDocument> {
  return new Promise((resolve, reject) => {
    worker.onmessage = (e) => resolve(e.data);
    worker.onerror = reject;
    worker.postMessage(data, [data.buffer]);
  });
}
```

### Efficient File Reading

```typescript
// Use streaming for large files
async function* streamFile(file: File, chunkSize = 1024 * 1024) {
  let offset = 0;
  while (offset < file.size) {
    const chunk = file.slice(offset, offset + chunkSize);
    yield new Uint8Array(await chunk.arrayBuffer());
    offset += chunkSize;
  }
}

// For validation, you might only need the first chunk
async function quickValidateFile(file: File) {
  const firstChunk = file.slice(0, 1024);
  const data = new Uint8Array(await firstChunk.arrayBuffer());
  return quickValidate(data);
}
```

### Render Optimization for Three.js

```typescript
// Use BufferGeometry merge for many lines
function mergeLines(entities: DxfLine[]): THREE.BufferGeometry {
  const positions: number[] = [];

  for (const line of entities) {
    positions.push(line.X1, line.Y1, 0);
    positions.push(line.X2, line.Y2, 0);
  }

  const geometry = new THREE.BufferGeometry();
  geometry.setAttribute('position',
    new THREE.Float32BufferAttribute(positions, 3));

  return geometry;
}

// Use instanced meshes for repeated blocks
function createBlockInstances(inserts: DxfInsert[], blockGeometry: THREE.BufferGeometry) {
  const mesh = new THREE.InstancedMesh(
    blockGeometry,
    new THREE.MeshBasicMaterial(),
    inserts.length
  );

  const matrix = new THREE.Matrix4();

  inserts.forEach((insert, i) => {
    matrix.compose(
      new THREE.Vector3(insert.X, insert.Y, 0),
      new THREE.Quaternion().setFromAxisAngle(
        new THREE.Vector3(0, 0, 1),
        (insert.Rotation || 0) * Math.PI / 180
      ),
      new THREE.Vector3(insert.ScaleX || 1, insert.ScaleY || 1, 1)
    );
    mesh.setMatrixAt(i, matrix);
  });

  return mesh;
}
```

## Node.js Optimization

### Use Buffer Directly

```typescript
import { readFileSync } from 'fs';

// Direct buffer usage
const buffer = readFileSync('file.jww');
const data = new Uint8Array(buffer.buffer, buffer.byteOffset, buffer.byteLength);
```

### Memory-Mapped Files (for very large files)

```typescript
import { open } from 'fs/promises';

async function processLargeFile(path: string) {
  const file = await open(path, 'r');
  const stats = await file.stat();

  // Read header only first
  const header = Buffer.alloc(1024);
  await file.read(header, 0, 1024, 0);

  // Validate before loading full file
  const validation = quickValidate(new Uint8Array(header));
  if (!validation.valid) {
    await file.close();
    throw new Error('Invalid file');
  }

  // Now read full file
  const data = Buffer.alloc(stats.size);
  await file.read(data, 0, stats.size, 0);
  await file.close();

  const parser = await createParser();
  const result = parser.parse(new Uint8Array(data));
  parser.dispose();

  return result;
}
```

### Stream Output

```typescript
import { createWriteStream } from 'fs';

async function convertToStream(data: Uint8Array, outputPath: string) {
  const parser = await createParser();
  const dxfString = parser.toDxfString(data);
  parser.dispose();

  // Stream large output
  const stream = createWriteStream(outputPath);

  // Write in chunks
  const chunkSize = 64 * 1024;
  for (let i = 0; i < dxfString.length; i += chunkSize) {
    const chunk = dxfString.slice(i, i + chunkSize);
    stream.write(chunk);
  }

  return new Promise<void>((resolve, reject) => {
    stream.on('finish', resolve);
    stream.on('error', reject);
    stream.end();
  });
}
```

## Benchmarking

### Built-in Stats

```typescript
const parser = await createParser();

// Enable timing in debug mode
parser.setDebug({
  enabled: true,
  includeTiming: true,
  logToConsole: true
});

// Parse multiple files
for (const file of files) {
  parser.parse(file);
}

// Get statistics
const stats = parser.getStats();
console.log('Parse Statistics:');
console.log(`  Total files: ${stats.parseCount}`);
console.log(`  Total bytes: ${stats.totalBytesProcessed}`);
console.log(`  Average time: ${stats.averageParseTimeMs.toFixed(2)}ms`);
console.log(`  Fastest: ${stats.fastestParseTimeMs}ms`);
console.log(`  Slowest: ${stats.slowestParseTimeMs}ms`);
console.log(`  Errors: ${stats.errorCount}`);
console.log(`  Memory: ${stats.memoryStats.totalFormatted}`);
```

### Custom Benchmarking

```typescript
interface BenchmarkResult {
  name: string;
  iterations: number;
  totalMs: number;
  avgMs: number;
  minMs: number;
  maxMs: number;
  opsPerSec: number;
}

async function benchmark(
  name: string,
  fn: () => void | Promise<void>,
  iterations = 100
): Promise<BenchmarkResult> {
  const times: number[] = [];

  // Warmup
  for (let i = 0; i < 5; i++) {
    await fn();
  }

  // Benchmark
  for (let i = 0; i < iterations; i++) {
    const start = performance.now();
    await fn();
    times.push(performance.now() - start);
  }

  const totalMs = times.reduce((a, b) => a + b, 0);
  return {
    name,
    iterations,
    totalMs,
    avgMs: totalMs / iterations,
    minMs: Math.min(...times),
    maxMs: Math.max(...times),
    opsPerSec: 1000 / (totalMs / iterations)
  };
}

// Usage
const parser = await createParser();
const testData = readFileSync('test.jww');

const result = await benchmark('parse', () => {
  parser.parse(new Uint8Array(testData));
}, 50);

console.log(`${result.name}: ${result.avgMs.toFixed(2)}ms avg, ${result.opsPerSec.toFixed(1)} ops/s`);
```

## Performance Tips Summary

1. **Reuse parser instances** - Avoid repeated WASM initialization
2. **Use parse options** - Skip unnecessary data
3. **Monitor memory** - Dispose when memory grows too large
4. **Use Web Workers** - Keep main thread responsive
5. **Batch process** - Allow GC between batches
6. **Profile first** - Use getStats() to identify bottlenecks
7. **Preview before full load** - Use maxEntities for large files
8. **Match options to use case** - Don't parse what you don't need
