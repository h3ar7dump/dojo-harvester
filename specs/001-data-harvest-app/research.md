# Research: Data Harvest App Technical Decisions

**Feature**: Data Harvest App for Humanoid Robot
**Date**: 2026-03-10
**Phase**: Phase 0 - Research & Technical Decisions

## Research Areas

### 1. WebSocket Streaming Patterns for Real-time Telemetry

**Context**: High-frequency robot telemetry (joint states, sensor readings) must be streamed to the UI with <100ms latency for 20-30 concurrent robots.

**Decision**: Use binary WebSocket frames with Protocol Buffers (Protobuf) serialization.

**Rationale**: 
- Protobuf offers significantly smaller payload sizes and faster serialization/deserialization than JSON, critical for meeting the 100ms latency target.
- Binary framing reduces the overhead of base64 encoding/decoding.
- Implementation will use a type-prefix byte `[1 byte type][payload bytes]` to allow future multiplexing of different message types over the same socket.

**Alternatives Considered**:
- **JSON over WebSocket**: Rejected due to payload size and string parsing overhead.
- **gRPC-Web**: Rejected to avoid the complexity of an additional proxy (Envoy) given the simple dojo environment.

### 2. LeRobot V3.0 Format Compliance & Validation

**Context**: Datasets must strictly adhere to the LeRobot V3.0 format to ensure compatibility with training pipelines.

**Decision**: Implement a two-stage validation: 
1. **Toolkit Level**: `convert.sh` (Python) performs schema check using the `lerobot` library.
2. **Backend Level**: The Go backend performs a "structural integrity check" (existence of `meta/info.json`, `data/*.parquet`, `videos/*.mp4`) before enqueuing for upload.

**Rationale**: Multi-layer validation prevents corrupt data from reaching the platform and provides early feedback to the operator.

### 3. Shell Script Execution & Resource Management

**Context**: The Go backend must orchestrate external scripts while maintaining system stability for up to 30 concurrent sessions.

**Decision**: Use a centralized `Executor` service in Go leveraging `os/exec` with dedicated `context.Context` for each process.

**Rationale**: 
- Provides deterministic lifecycle management (start, stop, cancel).
- Allows setting `syscall.SysProcAttr` for process grouping and priority (niceness) to prevent a single recording session from starving the OS.
- Timeouts and resource limits (via context) prevent zombie processes.

### 4. BadgerDB Persistence Strategy

**Context**: Local state (upload jobs, session status) must survive app restarts and handle concurrent writes from 30 robots.

**Decision**: Use BadgerDB with an LSM tree configuration optimized for write throughput.

**Rationale**: 
- Embedded KV store avoids external service dependencies (Redis/Postgres).
- BadgerDB supports high-concurrency writes and provides ACID transactions.
- Queue management will use lexicographical keys (e.g., `queue:pending:<timestamp>:<id>`) for efficient iteration.

### 5. 3D Rendering Performance in React

**Context**: 60 FPS visualization of humanoid joints.

**Decision**: Use `@react-three/fiber` with manual mesh manipulation in the `useFrame` loop.

**Rationale**: 
- Avoids React reconciliation overhead by directly updating Three.js object properties from the Zustand store.
- `urdf-loader` will be used to parse the robot model once and cache the geometry.
