/**
 * JWW (Jw_cad) file parser and DXF converter
 *
 * This module provides functionality to parse JWW binary files
 * and convert them to DXF format using WebAssembly.
 *
 * @packageDocumentation
 */

// =============================================================================
// Error Classes
// =============================================================================

/**
 * Error codes for JWW parser operations
 */
export enum JwwErrorCode {
  /** Parser has not been initialized */
  NOT_INITIALIZED = "NOT_INITIALIZED",
  /** WASM module failed to load */
  WASM_LOAD_FAILED = "WASM_LOAD_FAILED",
  /** WASM functions not available after timeout */
  WASM_TIMEOUT = "WASM_TIMEOUT",
  /** Invalid JWW file signature */
  INVALID_SIGNATURE = "INVALID_SIGNATURE",
  /** Unsupported JWW version */
  UNSUPPORTED_VERSION = "UNSUPPORTED_VERSION",
  /** General parse error */
  PARSE_ERROR = "PARSE_ERROR",
  /** DXF conversion error */
  CONVERSION_ERROR = "CONVERSION_ERROR",
  /** Validation error */
  VALIDATION_ERROR = "VALIDATION_ERROR",
  /** Memory allocation error */
  MEMORY_ERROR = "MEMORY_ERROR",
  /** Invalid argument provided */
  INVALID_ARGUMENT = "INVALID_ARGUMENT",
}

/**
 * Base error class for JWW parser errors
 */
export class JwwParserError extends Error {
  /** Error code identifying the error type */
  readonly code: JwwErrorCode;
  /** Original error that caused this error, if any */
  readonly cause?: Error;
  /** Additional context about the error */
  readonly context?: Record<string, unknown>;
  /** Timestamp when the error occurred */
  readonly timestamp: Date;

  constructor(
    code: JwwErrorCode,
    message: string,
    options?: { cause?: Error; context?: Record<string, unknown> }
  ) {
    super(message);
    this.name = "JwwParserError";
    this.code = code;
    this.cause = options?.cause;
    this.context = options?.context;
    this.timestamp = new Date();

    // Maintain proper stack trace for where our error was thrown
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, JwwParserError);
    }
  }

  /**
   * Returns a detailed string representation of the error
   */
  toDetailedString(): string {
    let details = `[${this.code}] ${this.message}`;
    if (this.context) {
      details += `\nContext: ${JSON.stringify(this.context, null, 2)}`;
    }
    if (this.cause) {
      details += `\nCaused by: ${this.cause.message}`;
    }
    return details;
  }

  /**
   * Converts the error to a plain object for logging/serialization
   */
  toJSON(): Record<string, unknown> {
    return {
      name: this.name,
      code: this.code,
      message: this.message,
      context: this.context,
      timestamp: this.timestamp.toISOString(),
      cause: this.cause?.message,
      stack: this.stack,
    };
  }
}

/**
 * Error thrown when the parser is not initialized
 */
export class NotInitializedError extends JwwParserError {
  constructor() {
    super(
      JwwErrorCode.NOT_INITIALIZED,
      "Parser not initialized. Call init() first or use createParser() factory function."
    );
    this.name = "NotInitializedError";
  }
}

/**
 * Error thrown when WASM module fails to load
 */
export class WasmLoadError extends JwwParserError {
  constructor(message: string, cause?: Error) {
    super(JwwErrorCode.WASM_LOAD_FAILED, message, { cause });
    this.name = "WasmLoadError";
  }
}

/**
 * Error thrown when file validation fails
 */
export class ValidationError extends JwwParserError {
  /** Specific validation issues found */
  readonly issues: ValidationIssue[];

  constructor(message: string, issues: ValidationIssue[]) {
    super(JwwErrorCode.VALIDATION_ERROR, message, {
      context: { issues },
    });
    this.name = "ValidationError";
    this.issues = issues;
  }
}

/**
 * Error thrown during parsing
 */
export class ParseError extends JwwParserError {
  /** Byte offset where the error occurred, if available */
  readonly offset?: number;
  /** Section being parsed when the error occurred */
  readonly section?: string;

  constructor(
    message: string,
    options?: { cause?: Error; offset?: number; section?: string }
  ) {
    super(JwwErrorCode.PARSE_ERROR, message, {
      cause: options?.cause,
      context: {
        offset: options?.offset,
        section: options?.section,
      },
    });
    this.name = "ParseError";
    this.offset = options?.offset;
    this.section = options?.section;
  }
}

// =============================================================================
// JWW Document Types
// =============================================================================

/**
 * Complete JWW document structure
 */
export interface JwwDocument {
  /** JWW file version number */
  Version: number;
  /** Document memo/comments */
  Memo: string;
  /** Paper size code */
  PaperSize: number;
  /** Layer groups (16 groups) */
  LayerGroups: JwwLayerGroup[];
  /** All entities in the document */
  Entities: JwwEntity[];
  /** Block definitions */
  Blocks: JwwBlock[];
}

/**
 * JWW layer group containing 16 layers
 */
export interface JwwLayerGroup {
  /** Layer group name */
  Name: string;
  /** Layers within this group (16 layers) */
  Layers: JwwLayer[];
}

/**
 * JWW layer within a layer group
 */
export interface JwwLayer {
  /** Layer name */
  Name: string;
  /** Whether the layer is visible */
  Visible: boolean;
  /** Whether the layer is locked */
  Locked: boolean;
}

/**
 * Base interface for all JWW entities
 */
export interface JwwEntityBase {
  /** Entity type discriminator */
  Type: string;
  /** Curve attribute number */
  Group: number;
  /** Line type/style */
  PenStyle: number;
  /** Color code */
  PenColor: number;
  /** Line width */
  PenWidth: number;
  /** Layer index within the group (0-15) */
  Layer: number;
  /** Layer group index (0-15) */
  LayerGroup: number;
}

/**
 * JWW Line entity
 */
export interface JwwLine extends JwwEntityBase {
  Type: "Line";
  /** Start X coordinate */
  X1: number;
  /** Start Y coordinate */
  Y1: number;
  /** End X coordinate */
  X2: number;
  /** End Y coordinate */
  Y2: number;
}

/**
 * JWW Arc entity (includes circles and ellipses)
 */
export interface JwwArc extends JwwEntityBase {
  Type: "Arc";
  /** Center X coordinate */
  CenterX: number;
  /** Center Y coordinate */
  CenterY: number;
  /** Arc radius */
  Radius: number;
  /** Start angle in degrees */
  StartAngle: number;
  /** End angle in degrees */
  EndAngle: number;
  /** Flatness ratio (1.0 for circles, other values for ellipses) */
  Flatness: number;
}

/**
 * JWW Point entity
 */
export interface JwwPoint extends JwwEntityBase {
  Type: "Point";
  /** X coordinate */
  X: number;
  /** Y coordinate */
  Y: number;
  /** Point code/type */
  Code: number;
}

/**
 * JWW Text entity
 */
export interface JwwText extends JwwEntityBase {
  Type: "Text";
  /** Text insertion X coordinate */
  X: number;
  /** Text insertion Y coordinate */
  Y: number;
  /** Text content */
  Text: string;
  /** Font name */
  FontName: string;
  /** Text height */
  Height: number;
  /** Text width (character width) */
  Width: number;
  /** Rotation angle in degrees */
  Angle: number;
}

/**
 * JWW Solid (filled polygon) entity
 */
export interface JwwSolid extends JwwEntityBase {
  Type: "Solid";
  /** Array of [x, y] coordinate pairs */
  Points: [number, number][];
}

/**
 * JWW Block reference entity
 */
export interface JwwBlockRef extends JwwEntityBase {
  Type: "Block";
  /** Insertion X coordinate */
  X: number;
  /** Insertion Y coordinate */
  Y: number;
  /** X scale factor */
  ScaleX: number;
  /** Y scale factor */
  ScaleY: number;
  /** Rotation angle in degrees */
  Angle: number;
  /** Referenced block definition number */
  BlockNumber: number;
}

/**
 * Union type of all JWW entity types
 */
export type JwwEntity =
  | JwwLine
  | JwwArc
  | JwwPoint
  | JwwText
  | JwwSolid
  | JwwBlockRef;

/**
 * JWW block definition
 */
export interface JwwBlock {
  /** Block name */
  Name: string;
  /** Entities within the block */
  Entities: JwwEntity[];
}

// =============================================================================
// DXF Document Types
// =============================================================================

/**
 * DXF document structure
 */
export interface DxfDocument {
  /** DXF layers */
  Layers: DxfLayer[];
  /** All entities in the document */
  Entities: DxfEntity[];
  /** Block definitions */
  Blocks: DxfBlock[];
}

/**
 * DXF layer definition
 */
export interface DxfLayer {
  /** Layer name */
  Name: string;
  /** ACI color code (1-255) */
  Color: number;
  /** Whether the layer is frozen */
  Frozen: boolean;
  /** Whether the layer is locked */
  Locked: boolean;
}

/**
 * Base interface for all DXF entities
 */
export interface DxfEntityBase {
  /** Entity type (LINE, CIRCLE, ARC, etc.) */
  Type: string;
  /** Layer name */
  Layer: string;
  /** Optional ACI color code */
  Color?: number;
  /** Optional line type name */
  LineType?: string;
}

/**
 * DXF LINE entity
 */
export interface DxfLine extends DxfEntityBase {
  Type: "LINE";
  /** Start X coordinate */
  X1: number;
  /** Start Y coordinate */
  Y1: number;
  /** Start Z coordinate */
  Z1?: number;
  /** End X coordinate */
  X2: number;
  /** End Y coordinate */
  Y2: number;
  /** End Z coordinate */
  Z2?: number;
}

/**
 * DXF CIRCLE entity
 */
export interface DxfCircle extends DxfEntityBase {
  Type: "CIRCLE";
  /** Center X coordinate */
  CenterX: number;
  /** Center Y coordinate */
  CenterY: number;
  /** Center Z coordinate */
  CenterZ?: number;
  /** Circle radius */
  Radius: number;
}

/**
 * DXF ARC entity
 */
export interface DxfArc extends DxfEntityBase {
  Type: "ARC";
  /** Center X coordinate */
  CenterX: number;
  /** Center Y coordinate */
  CenterY: number;
  /** Center Z coordinate */
  CenterZ?: number;
  /** Arc radius */
  Radius: number;
  /** Start angle in degrees */
  StartAngle: number;
  /** End angle in degrees */
  EndAngle: number;
}

/**
 * DXF ELLIPSE entity
 */
export interface DxfEllipse extends DxfEntityBase {
  Type: "ELLIPSE";
  /** Center X coordinate */
  CenterX: number;
  /** Center Y coordinate */
  CenterY: number;
  /** Center Z coordinate */
  CenterZ?: number;
  /** Major axis X component (relative to center) */
  MajorAxisX: number;
  /** Major axis Y component (relative to center) */
  MajorAxisY: number;
  /** Major axis Z component (relative to center) */
  MajorAxisZ?: number;
  /** Minor to major axis ratio (0 to 1) */
  MinorRatio: number;
  /** Start parameter (0 to 2*PI) */
  StartParam: number;
  /** End parameter (0 to 2*PI) */
  EndParam: number;
}

/**
 * DXF POINT entity
 */
export interface DxfPoint extends DxfEntityBase {
  Type: "POINT";
  /** X coordinate */
  X: number;
  /** Y coordinate */
  Y: number;
  /** Z coordinate */
  Z?: number;
}

/**
 * DXF TEXT entity
 */
export interface DxfText extends DxfEntityBase {
  Type: "TEXT";
  /** Insertion X coordinate */
  X: number;
  /** Insertion Y coordinate */
  Y: number;
  /** Insertion Z coordinate */
  Z?: number;
  /** Text height */
  Height: number;
  /** Text content */
  Content: string;
  /** Rotation angle in degrees */
  Rotation?: number;
  /** Text style name */
  Style?: string;
}

/**
 * DXF MTEXT (multiline text) entity
 */
export interface DxfMText extends DxfEntityBase {
  Type: "MTEXT";
  /** Insertion X coordinate */
  X: number;
  /** Insertion Y coordinate */
  Y: number;
  /** Insertion Z coordinate */
  Z?: number;
  /** Text height */
  Height: number;
  /** Text content (with formatting codes) */
  Content: string;
  /** Rotation angle in degrees */
  Rotation?: number;
  /** Reference rectangle width */
  Width?: number;
}

/**
 * DXF SOLID entity (filled triangle or quadrilateral)
 */
export interface DxfSolid extends DxfEntityBase {
  Type: "SOLID";
  /** First corner X */
  X1: number;
  /** First corner Y */
  Y1: number;
  /** Second corner X */
  X2: number;
  /** Second corner Y */
  Y2: number;
  /** Third corner X */
  X3: number;
  /** Third corner Y */
  Y3: number;
  /** Fourth corner X (same as third for triangles) */
  X4: number;
  /** Fourth corner Y (same as third for triangles) */
  Y4: number;
}

/**
 * DXF INSERT entity (block reference)
 */
export interface DxfInsert extends DxfEntityBase {
  Type: "INSERT";
  /** Block name */
  BlockName: string;
  /** Insertion X coordinate */
  X: number;
  /** Insertion Y coordinate */
  Y: number;
  /** Insertion Z coordinate */
  Z?: number;
  /** X scale factor */
  ScaleX?: number;
  /** Y scale factor */
  ScaleY?: number;
  /** Z scale factor */
  ScaleZ?: number;
  /** Rotation angle in degrees */
  Rotation?: number;
}

/**
 * DXF POLYLINE vertex
 */
export interface DxfVertex {
  /** X coordinate */
  X: number;
  /** Y coordinate */
  Y: number;
  /** Z coordinate */
  Z?: number;
  /** Bulge factor (for curved segments) */
  Bulge?: number;
}

/**
 * DXF LWPOLYLINE entity
 */
export interface DxfLwPolyline extends DxfEntityBase {
  Type: "LWPOLYLINE";
  /** Whether the polyline is closed */
  Closed: boolean;
  /** Polyline vertices */
  Vertices: DxfVertex[];
}

/**
 * Union type of all DXF entity types
 */
export type DxfEntity =
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
  | (DxfEntityBase & Record<string, unknown>); // Allow for unknown entity types

/**
 * DXF block definition
 */
export interface DxfBlock {
  /** Block name */
  Name: string;
  /** Base point X coordinate */
  BaseX?: number;
  /** Base point Y coordinate */
  BaseY?: number;
  /** Base point Z coordinate */
  BaseZ?: number;
  /** Entities within the block */
  Entities: DxfEntity[];
}

// =============================================================================
// Progress and Callback Types
// =============================================================================

/**
 * Progress stages during parsing
 */
export type ProgressStage =
  | "loading"
  | "parsing_header"
  | "parsing_layers"
  | "parsing_entities"
  | "parsing_blocks"
  | "converting"
  | "complete";

/**
 * Progress information during parsing
 */
export interface ProgressInfo {
  /** Current stage of parsing */
  stage: ProgressStage;
  /** Progress within the current stage (0-100) */
  progress: number;
  /** Overall progress (0-100) */
  overallProgress: number;
  /** Optional message about current operation */
  message?: string;
  /** Number of entities processed so far */
  entitiesProcessed?: number;
  /** Total number of entities (if known) */
  totalEntities?: number;
  /** Elapsed time in milliseconds */
  elapsedMs: number;
  /** Estimated remaining time in milliseconds (if available) */
  estimatedRemainingMs?: number;
}

/**
 * Callback function for reporting progress
 */
export type ProgressCallback = (progress: ProgressInfo) => void;

// =============================================================================
// Parse Options
// =============================================================================

/**
 * Options for parsing JWW files
 */
export interface ParseOptions {
  /**
   * Enable streaming mode for lower memory usage
   * When enabled, entities are processed in chunks
   * @default false
   */
  streamingMode?: boolean;

  /**
   * Skip parsing of specific entity types
   * Useful for improving performance when not all entities are needed
   */
  skipEntityTypes?: Array<JwwEntity["Type"]>;

  /**
   * Only parse entities on specific layer groups
   * Array of layer group indices (0-15)
   */
  layerGroupFilter?: number[];

  /**
   * Only parse entities on specific layers
   * Format: { layerGroup: layerIndex[] }
   */
  layerFilter?: Record<number, number[]>;

  /**
   * Limit the number of entities to parse
   * Useful for previewing large files
   */
  maxEntities?: number;

  /**
   * Include block definitions in the output
   * @default true
   */
  includeBlocks?: boolean;

  /**
   * Expand block references inline (replaces INSERT with actual entities)
   * Warning: May significantly increase output size
   * @default false
   */
  expandBlockReferences?: boolean;

  /**
   * Progress callback function
   * Called periodically during parsing with progress information
   */
  onProgress?: ProgressCallback;

  /**
   * Abort signal for cancelling the parse operation
   */
  signal?: AbortSignal;
}

/**
 * Options for DXF conversion
 */
export interface ConvertOptions extends ParseOptions {
  /**
   * Convert temporary points (Code > 0)
   * @default false
   */
  includeTemporaryPoints?: boolean;

  /**
   * Color mapping override
   * Maps JWW color codes to DXF ACI colors
   */
  colorMapping?: Record<number, number>;

  /**
   * Layer naming pattern
   * Use {group} and {layer} placeholders
   * @default "{group}-{layer}"
   */
  layerNamePattern?: string;

  /**
   * Precision for coordinate values (decimal places)
   * @default 6
   */
  precision?: number;
}

// =============================================================================
// Validation Types
// =============================================================================

/**
 * Severity level for validation issues
 */
export type ValidationSeverity = "error" | "warning" | "info";

/**
 * A single validation issue
 */
export interface ValidationIssue {
  /** Severity of the issue */
  severity: ValidationSeverity;
  /** Issue code */
  code: string;
  /** Human-readable message */
  message: string;
  /** Byte offset where the issue was found, if applicable */
  offset?: number;
  /** Additional details */
  details?: Record<string, unknown>;
}

/**
 * Result of file validation
 */
export interface ValidationResult {
  /** Whether the file is valid (no errors) */
  valid: boolean;
  /** File format version, if detected */
  version?: number;
  /** Estimated file size category */
  sizeCategory: "small" | "medium" | "large" | "very_large";
  /** Estimated entity count, if detectable */
  estimatedEntityCount?: number;
  /** List of validation issues found */
  issues: ValidationIssue[];
  /** Validation took this many milliseconds */
  validationTimeMs: number;
}

// =============================================================================
// Debug Types
// =============================================================================

/**
 * Debug log levels
 */
export type LogLevel = "debug" | "info" | "warn" | "error";

/**
 * Debug log entry
 */
export interface DebugLogEntry {
  /** Timestamp of the log entry */
  timestamp: Date;
  /** Log level */
  level: LogLevel;
  /** Log message */
  message: string;
  /** Additional data */
  data?: Record<string, unknown>;
}

/**
 * Debug callback function
 */
export type DebugCallback = (entry: DebugLogEntry) => void;

/**
 * Debug options
 */
export interface DebugOptions {
  /**
   * Enable debug mode
   * @default false
   */
  enabled?: boolean;

  /**
   * Minimum log level to capture
   * @default "info"
   */
  logLevel?: LogLevel;

  /**
   * Debug callback function
   * If not provided, logs will be stored internally
   */
  onDebug?: DebugCallback;

  /**
   * Log to console
   * @default false
   */
  logToConsole?: boolean;

  /**
   * Maximum number of log entries to store
   * @default 1000
   */
  maxLogEntries?: number;

  /**
   * Include timing information
   * @default true
   */
  includeTiming?: boolean;

  /**
   * Include memory usage information
   * @default false
   */
  includeMemoryUsage?: boolean;
}

// =============================================================================
// Memory and Statistics Types
// =============================================================================

/**
 * Memory usage statistics
 */
export interface MemoryStats {
  /** WASM memory buffer size in bytes */
  wasmMemoryBytes: number;
  /** Estimated JS heap usage in bytes (if available) */
  jsHeapBytes?: number;
  /** Total estimated memory usage */
  totalBytes: number;
  /** Human-readable total */
  totalFormatted: string;
}

/**
 * Parser statistics
 */
export interface ParserStats {
  /** Number of parse operations performed */
  parseCount: number;
  /** Total bytes processed */
  totalBytesProcessed: number;
  /** Average parse time in milliseconds */
  averageParseTimeMs: number;
  /** Fastest parse time in milliseconds */
  fastestParseTimeMs: number;
  /** Slowest parse time in milliseconds */
  slowestParseTimeMs: number;
  /** Number of errors encountered */
  errorCount: number;
  /** Current memory usage */
  memoryStats: MemoryStats;
}

// =============================================================================
// WASM Result Types
// =============================================================================

/**
 * Result type from WASM operations
 * @internal
 */
interface WasmResult {
  ok: boolean;
  data?: string;
  error?: string;
  errorCode?: string;
  offset?: number;
  section?: string;
}

/**
 * Validation result from WASM
 * @internal
 */
interface WasmValidationResult {
  ok: boolean;
  valid: boolean;
  version?: number;
  estimatedEntities?: number;
  issues?: Array<{
    severity: string;
    code: string;
    message: string;
    offset?: number;
  }>;
  error?: string;
}

// =============================================================================
// Global Declarations
// =============================================================================

declare global {
  var Go: new () => GoInstance;
  var jwwParse: ((data: Uint8Array) => WasmResult) | undefined;
  var jwwToDxf: ((data: Uint8Array) => WasmResult) | undefined;
  var jwwToDxfString: ((data: Uint8Array) => WasmResult) | undefined;
  var jwwValidate: ((data: Uint8Array) => WasmValidationResult) | undefined;
  var jwwGetVersion: (() => string) | undefined;
  var jwwSetDebug: ((enabled: boolean) => void) | undefined;
}

interface GoInstance {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
  _inst?: WebAssembly.Instance;
}

// =============================================================================
// Helper Functions
// =============================================================================

/**
 * Format bytes to human-readable string
 */
function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`;
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}

/**
 * Get file size category based on byte count
 */
function getSizeCategory(
  bytes: number
): "small" | "medium" | "large" | "very_large" {
  if (bytes < 100 * 1024) return "small"; // < 100KB
  if (bytes < 1024 * 1024) return "medium"; // < 1MB
  if (bytes < 10 * 1024 * 1024) return "large"; // < 10MB
  return "very_large"; // >= 10MB
}

/**
 * Create a parse error from WASM result
 */
function createParseError(result: WasmResult): ParseError {
  const message = result.error || "Parse failed";

  // Try to extract more specific error type
  if (message.includes("invalid JWW signature")) {
    return new ParseError("Invalid JWW file format: missing or corrupt signature", {
      section: "header",
      offset: 0,
    });
  }

  if (message.includes("unsupported") && message.includes("version")) {
    return new ParseError(message, {
      section: "header",
    });
  }

  return new ParseError(message, {
    offset: result.offset,
    section: result.section,
  });
}

// =============================================================================
// JwwParser Class
// =============================================================================

/**
 * JWW file parser with WebAssembly backend
 *
 * @example
 * ```typescript
 * // Using the factory function (recommended)
 * const parser = await createParser();
 * const doc = parser.parse(fileData);
 *
 * // Manual initialization
 * const parser = new JwwParser();
 * await parser.init();
 * const doc = parser.parse(fileData);
 * ```
 */
export class JwwParser {
  private initialized = false;
  private initPromise: Promise<void> | null = null;
  private wasmPath: string;
  private goInstance: GoInstance | null = null;
  private wasmInstance: WebAssembly.Instance | null = null;

  // Debug and logging
  private debugOptions: DebugOptions = { enabled: false };
  private debugLogs: DebugLogEntry[] = [];

  // Statistics
  private stats: {
    parseCount: number;
    totalBytesProcessed: number;
    parseTimes: number[];
    errorCount: number;
  } = {
    parseCount: 0,
    totalBytesProcessed: 0,
    parseTimes: [],
    errorCount: 0,
  };

  /**
   * Create a new JWW parser instance
   * @param wasmPath - Path to the jww-parser.wasm file
   */
  constructor(wasmPath?: string) {
    this.wasmPath = wasmPath || this.getDefaultWasmPath();
  }

  private getDefaultWasmPath(): string {
    if (typeof process !== "undefined" && process.versions?.node) {
      return new URL("../wasm/jww-parser.wasm", import.meta.url).pathname;
    }
    return "jww-parser.wasm";
  }

  // ===========================================================================
  // Debug Methods
  // ===========================================================================

  /**
   * Enable or configure debug mode
   * @param options - Debug configuration options
   */
  setDebug(options: DebugOptions | boolean): void {
    if (typeof options === "boolean") {
      this.debugOptions = { enabled: options };
    } else {
      this.debugOptions = { ...this.debugOptions, ...options };
    }

    // Also set debug mode in WASM if available
    if (this.initialized && typeof globalThis.jwwSetDebug === "function") {
      globalThis.jwwSetDebug(this.debugOptions.enabled ?? false);
    }

    this.log("info", "Debug mode " + (this.debugOptions.enabled ? "enabled" : "disabled"));
  }

  /**
   * Get debug logs
   * @param level - Optional minimum log level to filter
   * @returns Array of debug log entries
   */
  getDebugLogs(level?: LogLevel): DebugLogEntry[] {
    if (!level) return [...this.debugLogs];

    const levels: LogLevel[] = ["debug", "info", "warn", "error"];
    const minIndex = levels.indexOf(level);

    return this.debugLogs.filter(
      (entry) => levels.indexOf(entry.level) >= minIndex
    );
  }

  /**
   * Clear debug logs
   */
  clearDebugLogs(): void {
    this.debugLogs = [];
  }

  private log(
    level: LogLevel,
    message: string,
    data?: Record<string, unknown>
  ): void {
    if (!this.debugOptions.enabled) return;

    const levels: LogLevel[] = ["debug", "info", "warn", "error"];
    const minLevel = this.debugOptions.logLevel || "info";
    if (levels.indexOf(level) < levels.indexOf(minLevel)) return;

    const entry: DebugLogEntry = {
      timestamp: new Date(),
      level,
      message,
      data: this.debugOptions.includeMemoryUsage
        ? { ...data, memory: this.getMemoryStats() }
        : data,
    };

    // Store log entry
    this.debugLogs.push(entry);
    if (this.debugLogs.length > (this.debugOptions.maxLogEntries || 1000)) {
      this.debugLogs.shift();
    }

    // Call debug callback if provided
    if (this.debugOptions.onDebug) {
      this.debugOptions.onDebug(entry);
    }

    // Log to console if enabled
    if (this.debugOptions.logToConsole) {
      const consoleFn =
        level === "error"
          ? console.error
          : level === "warn"
            ? console.warn
            : level === "debug"
              ? console.debug
              : console.log;
      consoleFn(`[JwwParser:${level}] ${message}`, data || "");
    }
  }

  // ===========================================================================
  // Initialization Methods
  // ===========================================================================

  /**
   * Check if the parser is initialized
   */
  isInitialized(): boolean {
    return this.initialized;
  }

  /**
   * Initialize the WASM module
   * Must be called before using parse methods
   */
  async init(): Promise<void> {
    if (this.initialized) return;
    if (this.initPromise) return this.initPromise;

    this.log("info", "Initializing WASM module");
    const startTime = Date.now();

    this.initPromise = this.loadWasm();
    try {
      await this.initPromise;
      this.initialized = true;
      this.log("info", "WASM module initialized", {
        elapsedMs: Date.now() - startTime,
      });
    } catch (error) {
      this.initPromise = null;
      this.log("error", "WASM initialization failed", {
        error: error instanceof Error ? error.message : String(error),
      });
      throw error;
    }
  }

  private async loadWasm(): Promise<void> {
    // Load wasm_exec.js if Go is not defined
    if (typeof Go === "undefined") {
      await this.loadWasmExec();
    }

    this.goInstance = new Go();

    try {
      if (typeof process !== "undefined" && process.versions?.node) {
        // Node.js environment
        const fs = await import("fs");
        const wasmBuffer = fs.readFileSync(this.wasmPath);
        const wasmModule = await WebAssembly.compile(wasmBuffer);
        this.wasmInstance = await WebAssembly.instantiate(
          wasmModule,
          this.goInstance.importObject
        );
      } else {
        // Browser environment
        const result = await WebAssembly.instantiateStreaming(
          fetch(this.wasmPath),
          this.goInstance.importObject
        ).catch(async () => {
          const response = await fetch(this.wasmPath);
          if (!response.ok) {
            throw new WasmLoadError(
              `Failed to fetch WASM file: ${response.status} ${response.statusText}`
            );
          }
          const bytes = await response.arrayBuffer();
          return WebAssembly.instantiate(bytes, this.goInstance!.importObject);
        });
        this.wasmInstance = result.instance;
      }
    } catch (error) {
      throw new WasmLoadError(
        `Failed to load WASM module: ${error instanceof Error ? error.message : String(error)}`,
        error instanceof Error ? error : undefined
      );
    }

    // Don't await - Go.run() blocks until the program exits
    this.goInstance.run(this.wasmInstance);

    // Wait for functions to be available
    await this.waitForWasmFunctions();

    // Set debug mode if already enabled
    if (
      this.debugOptions.enabled &&
      typeof globalThis.jwwSetDebug === "function"
    ) {
      globalThis.jwwSetDebug(true);
    }
  }

  private async loadWasmExec(): Promise<void> {
    if (typeof process !== "undefined" && process.versions?.node) {
      const wasmExecPath = new URL("../wasm/wasm_exec.js", import.meta.url)
        .pathname;
      await import(wasmExecPath);
    } else {
      throw new WasmLoadError(
        "Go runtime not loaded. Please include wasm_exec.js in your HTML before using the parser."
      );
    }
  }

  private async waitForWasmFunctions(
    timeout = 5000,
    interval = 50
  ): Promise<void> {
    const start = Date.now();
    while (Date.now() - start < timeout) {
      if (
        typeof globalThis.jwwParse === "function" &&
        typeof globalThis.jwwToDxf === "function" &&
        typeof globalThis.jwwToDxfString === "function"
      ) {
        return;
      }
      await new Promise((resolve) => setTimeout(resolve, interval));
    }
    throw new JwwParserError(
      JwwErrorCode.WASM_TIMEOUT,
      "WASM functions not available after timeout. The WASM module may have failed to initialize."
    );
  }

  private ensureInitialized(): void {
    if (!this.initialized) {
      throw new NotInitializedError();
    }
  }

  // ===========================================================================
  // Validation Methods
  // ===========================================================================

  /**
   * Validate a JWW file without fully parsing it
   * Useful for checking file validity before processing
   *
   * @param data - JWW file content as Uint8Array
   * @returns Validation result with any issues found
   */
  validate(data: Uint8Array): ValidationResult {
    const startTime = Date.now();

    // Basic validation without WASM (always available)
    const issues: ValidationIssue[] = [];

    // Check minimum size
    if (data.length < 16) {
      issues.push({
        severity: "error",
        code: "FILE_TOO_SMALL",
        message: "File is too small to be a valid JWW file",
        details: { size: data.length, minimumRequired: 16 },
      });
    }

    // Check JWW signature "JwwData."
    const signature = new TextDecoder().decode(data.slice(0, 8));
    if (signature !== "JwwData.") {
      issues.push({
        severity: "error",
        code: "INVALID_SIGNATURE",
        message: `Invalid file signature: expected "JwwData.", got "${signature.replace(/[^\x20-\x7E]/g, "?")}"`,
        offset: 0,
      });
    }

    // Try to read version (bytes 8-11, little-endian)
    let version: number | undefined;
    if (data.length >= 12) {
      const view = new DataView(data.buffer, data.byteOffset, data.byteLength);
      version = view.getInt32(8, true);

      if (version < 200 || version > 1000) {
        issues.push({
          severity: "warning",
          code: "UNUSUAL_VERSION",
          message: `Unusual version number: ${version}. Expected between 200-1000.`,
          details: { version },
        });
      }
    }

    // Estimate entity count based on file size (rough heuristic)
    // Average entity is approximately 50-100 bytes
    const estimatedEntityCount = Math.floor(data.length / 75);

    // Size warnings
    const sizeCategory = getSizeCategory(data.length);
    if (sizeCategory === "very_large") {
      issues.push({
        severity: "warning",
        code: "LARGE_FILE",
        message: `File is very large (${formatBytes(data.length)}). Consider using streaming mode for better performance.`,
        details: { size: data.length, sizeFormatted: formatBytes(data.length) },
      });
    }

    const hasErrors = issues.some((i) => i.severity === "error");

    return {
      valid: !hasErrors,
      version,
      sizeCategory,
      estimatedEntityCount,
      issues,
      validationTimeMs: Date.now() - startTime,
    };
  }

  /**
   * Validate and throw if invalid
   * Convenience method that throws a ValidationError if the file is invalid
   *
   * @param data - JWW file content as Uint8Array
   * @throws {ValidationError} If the file is invalid
   */
  validateOrThrow(data: Uint8Array): void {
    const result = this.validate(data);
    if (!result.valid) {
      throw new ValidationError(
        "File validation failed: " +
          result.issues
            .filter((i) => i.severity === "error")
            .map((i) => i.message)
            .join("; "),
        result.issues
      );
    }
  }

  // ===========================================================================
  // Parsing Methods
  // ===========================================================================

  /**
   * Parse a JWW file and return the document structure
   *
   * @param data - JWW file content as Uint8Array
   * @param options - Optional parsing options
   * @returns Parsed JWW document
   * @throws {NotInitializedError} If parser is not initialized
   * @throws {ParseError} If parsing fails
   */
  parse(data: Uint8Array, options?: ParseOptions): JwwDocument {
    this.ensureInitialized();

    // Check for abort signal
    if (options?.signal?.aborted) {
      throw new JwwParserError(
        JwwErrorCode.PARSE_ERROR,
        "Parse operation was aborted"
      );
    }

    const startTime = Date.now();
    this.log("info", "Starting parse operation", {
      dataSize: data.length,
      options: options ? { ...options, onProgress: undefined } : undefined,
    });

    // Report initial progress
    if (options?.onProgress) {
      options.onProgress({
        stage: "loading",
        progress: 0,
        overallProgress: 0,
        message: "Loading file data",
        elapsedMs: 0,
      });
    }

    try {
      const result = globalThis.jwwParse!(data);

      if (!result.ok) {
        this.stats.errorCount++;
        throw createParseError(result);
      }

      // Report parsing complete
      if (options?.onProgress) {
        options.onProgress({
          stage: "complete",
          progress: 100,
          overallProgress: 100,
          message: "Parsing complete",
          elapsedMs: Date.now() - startTime,
        });
      }

      const doc = JSON.parse(result.data!) as JwwDocument;

      // Apply filtering options
      let filteredDoc = doc;
      if (options) {
        filteredDoc = this.applyParseOptions(doc, options);
      }

      // Update stats
      const parseTime = Date.now() - startTime;
      this.stats.parseCount++;
      this.stats.totalBytesProcessed += data.length;
      this.stats.parseTimes.push(parseTime);
      if (this.stats.parseTimes.length > 100) {
        this.stats.parseTimes.shift();
      }

      this.log("info", "Parse complete", {
        parseTimeMs: parseTime,
        entityCount: filteredDoc.Entities.length,
        blockCount: filteredDoc.Blocks.length,
      });

      return filteredDoc;
    } catch (error) {
      this.stats.errorCount++;
      this.log("error", "Parse failed", {
        error: error instanceof Error ? error.message : String(error),
      });

      if (error instanceof JwwParserError) {
        throw error;
      }
      throw new ParseError(
        error instanceof Error ? error.message : String(error),
        { cause: error instanceof Error ? error : undefined }
      );
    }
  }

  private applyParseOptions(doc: JwwDocument, options: ParseOptions): JwwDocument {
    let entities = doc.Entities;
    let blocks = doc.Blocks;

    // Filter by entity type
    if (options.skipEntityTypes && options.skipEntityTypes.length > 0) {
      const skipTypes = new Set(options.skipEntityTypes);
      entities = entities.filter((e) => !skipTypes.has(e.Type));
    }

    // Filter by layer group
    if (options.layerGroupFilter && options.layerGroupFilter.length > 0) {
      const groups = new Set(options.layerGroupFilter);
      entities = entities.filter((e) => groups.has(e.LayerGroup));
    }

    // Filter by specific layers
    if (options.layerFilter) {
      const filter = options.layerFilter;
      entities = entities.filter((e) => {
        const layers = filter[e.LayerGroup];
        return !layers || layers.includes(e.Layer);
      });
    }

    // Limit entity count
    if (options.maxEntities && entities.length > options.maxEntities) {
      entities = entities.slice(0, options.maxEntities);
    }

    // Handle blocks
    if (options.includeBlocks === false) {
      blocks = [];
    }

    return {
      ...doc,
      Entities: entities,
      Blocks: blocks,
    };
  }

  /**
   * Parse a JWW file and convert to DXF document structure
   *
   * @param data - JWW file content as Uint8Array
   * @param options - Optional conversion options
   * @returns DXF document object
   */
  toDxf(data: Uint8Array, options?: ConvertOptions): DxfDocument {
    this.ensureInitialized();

    const startTime = Date.now();
    this.log("info", "Starting DXF conversion", { dataSize: data.length });

    if (options?.onProgress) {
      options.onProgress({
        stage: "converting",
        progress: 0,
        overallProgress: 0,
        message: "Converting to DXF",
        elapsedMs: 0,
      });
    }

    try {
      const result = globalThis.jwwToDxf!(data);

      if (!result.ok) {
        this.stats.errorCount++;
        throw createParseError(result);
      }

      if (options?.onProgress) {
        options.onProgress({
          stage: "complete",
          progress: 100,
          overallProgress: 100,
          message: "Conversion complete",
          elapsedMs: Date.now() - startTime,
        });
      }

      const doc = JSON.parse(result.data!) as DxfDocument;

      // Update stats
      const parseTime = Date.now() - startTime;
      this.stats.parseCount++;
      this.stats.totalBytesProcessed += data.length;
      this.stats.parseTimes.push(parseTime);

      this.log("info", "DXF conversion complete", {
        parseTimeMs: parseTime,
        entityCount: doc.Entities.length,
      });

      return doc;
    } catch (error) {
      this.stats.errorCount++;
      this.log("error", "DXF conversion failed", {
        error: error instanceof Error ? error.message : String(error),
      });

      if (error instanceof JwwParserError) {
        throw error;
      }
      throw new ParseError(
        error instanceof Error ? error.message : String(error),
        { cause: error instanceof Error ? error : undefined }
      );
    }
  }

  /**
   * Parse a JWW file and convert to DXF file content string
   *
   * @param data - JWW file content as Uint8Array
   * @param options - Optional conversion options
   * @returns DXF file content as string (ready to save as .dxf file)
   */
  toDxfString(data: Uint8Array, options?: ConvertOptions): string {
    this.ensureInitialized();

    const startTime = Date.now();
    this.log("info", "Starting DXF string generation", { dataSize: data.length });

    try {
      const result = globalThis.jwwToDxfString!(data);

      if (!result.ok) {
        this.stats.errorCount++;
        throw createParseError(result);
      }

      const parseTime = Date.now() - startTime;
      this.stats.parseCount++;
      this.stats.totalBytesProcessed += data.length;
      this.stats.parseTimes.push(parseTime);

      this.log("info", "DXF string generation complete", {
        parseTimeMs: parseTime,
        outputLength: result.data!.length,
      });

      return result.data!;
    } catch (error) {
      this.stats.errorCount++;

      if (error instanceof JwwParserError) {
        throw error;
      }
      throw new ParseError(
        error instanceof Error ? error.message : String(error),
        { cause: error instanceof Error ? error : undefined }
      );
    }
  }

  // ===========================================================================
  // Memory Management
  // ===========================================================================

  /**
   * Get current memory usage statistics
   */
  getMemoryStats(): MemoryStats {
    let wasmMemoryBytes = 0;

    // Try to get WASM memory size
    if (this.wasmInstance?.exports.memory) {
      const memory = this.wasmInstance.exports.memory as WebAssembly.Memory;
      wasmMemoryBytes = memory.buffer.byteLength;
    }

    // Try to get JS heap size (Node.js only)
    let jsHeapBytes: number | undefined;
    if (
      typeof process !== "undefined" &&
      typeof process.memoryUsage === "function"
    ) {
      jsHeapBytes = process.memoryUsage().heapUsed;
    }

    const totalBytes = wasmMemoryBytes + (jsHeapBytes || 0);

    return {
      wasmMemoryBytes,
      jsHeapBytes,
      totalBytes,
      totalFormatted: formatBytes(totalBytes),
    };
  }

  /**
   * Get parser statistics
   */
  getStats(): ParserStats {
    const parseTimes = this.stats.parseTimes;

    return {
      parseCount: this.stats.parseCount,
      totalBytesProcessed: this.stats.totalBytesProcessed,
      averageParseTimeMs:
        parseTimes.length > 0
          ? parseTimes.reduce((a, b) => a + b, 0) / parseTimes.length
          : 0,
      fastestParseTimeMs: parseTimes.length > 0 ? Math.min(...parseTimes) : 0,
      slowestParseTimeMs: parseTimes.length > 0 ? Math.max(...parseTimes) : 0,
      errorCount: this.stats.errorCount,
      memoryStats: this.getMemoryStats(),
    };
  }

  /**
   * Reset parser statistics
   */
  resetStats(): void {
    this.stats = {
      parseCount: 0,
      totalBytesProcessed: 0,
      parseTimes: [],
      errorCount: 0,
    };
    this.log("debug", "Statistics reset");
  }

  /**
   * Clean up resources and release memory
   * Call this when you're done using the parser to free WASM memory
   */
  dispose(): void {
    this.log("info", "Disposing parser resources");

    // Clear references
    this.goInstance = null;
    this.wasmInstance = null;
    this.initialized = false;
    this.initPromise = null;

    // Clear stats and logs
    this.stats = {
      parseCount: 0,
      totalBytesProcessed: 0,
      parseTimes: [],
      errorCount: 0,
    };
    this.debugLogs = [];

    // Note: We cannot truly "free" WASM memory, but clearing references
    // allows garbage collection to reclaim the memory
  }

  /**
   * Get the WASM module version
   */
  getVersion(): string {
    if (typeof globalThis.jwwGetVersion === "function") {
      return globalThis.jwwGetVersion();
    }
    return "1.0.0"; // Fallback version
  }
}

// =============================================================================
// Factory Functions
// =============================================================================

/**
 * Create and initialize a JWW parser instance
 *
 * @param wasmPath - Optional path to the jww-parser.wasm file
 * @param options - Optional debug options
 * @returns Initialized JwwParser instance
 *
 * @example
 * ```typescript
 * const parser = await createParser();
 * const doc = parser.parse(fileData);
 * ```
 */
export async function createParser(
  wasmPath?: string,
  options?: { debug?: DebugOptions }
): Promise<JwwParser> {
  const parser = new JwwParser(wasmPath);

  if (options?.debug) {
    parser.setDebug(options.debug);
  }

  await parser.init();
  return parser;
}

/**
 * Quick validate a JWW file without initializing a full parser
 * Performs basic validation only (signature, version check)
 *
 * @param data - JWW file content as Uint8Array
 * @returns Validation result
 */
export function quickValidate(data: Uint8Array): ValidationResult {
  const parser = new JwwParser();
  return parser.validate(data);
}

/**
 * Check if a Uint8Array looks like a JWW file
 *
 * @param data - File content as Uint8Array
 * @returns true if the file appears to be a JWW file
 */
export function isJwwFile(data: Uint8Array): boolean {
  if (data.length < 8) return false;
  const signature = new TextDecoder().decode(data.slice(0, 8));
  return signature === "JwwData.";
}

// =============================================================================
// Default Export
// =============================================================================

export default {
  JwwParser,
  createParser,
  quickValidate,
  isJwwFile,
  JwwParserError,
  NotInitializedError,
  WasmLoadError,
  ValidationError,
  ParseError,
  JwwErrorCode,
};
