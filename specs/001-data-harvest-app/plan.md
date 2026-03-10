# Implementation Plan: Data Harvest App for Humanoid Robot

**Branch**: `001-data-harvest-app` | **Date**: 2026-03-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-data-harvest-app/spec.md`

## Summary

Build a high-performance data harvest application for humanoid robots. The system uses a Go-based edge agent to orchestrate local recording and conversion workflows, a Python toolkit for hardware-specific ROS processing, and a React frontend for operator control. All telemetry is streamed via binary Protobuf over WebSockets to ensure <100ms latency, while datasets are validated against the LeRobot V3.0 schema before being uploaded to the remote data platform.

## Technical Context

**Language/Version**: Go 1.26, React 19 (TS), Python 3.13.12  
**Primary Dependencies**: gin, gorilla/websocket, lerobot, three.js, Zustand  
**Storage**: BadgerDB (embedded KV), Local Filesystem (Datasets)  
**Testing**: Go Test, Vitest, Pytest  
**Target Platform**: Linux Edge Devices, Modern Web Browsers  
**Project Type**: Full-stack hardware orchestration service  
**Performance Goals**: <100ms Telemetry Latency, 60 FPS 3D rendering  
**Constraints**: 20-30 concurrent robot sessions, offline-first recording  
**Scale/Scope**: Dojo-wide data collection (30 robots), 10GB+ daily data volume

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| Serialization Protocol | ✅ PASS | Binary Protobuf for telemetry, JSON for Management APIs. |
| Task Management Isolation | ✅ PASS | React UI calls Remote Platform API directly for task claiming. |
| Execution Abstraction | ✅ PASS | Go Backend calls shell scripts in `toolkit/` using Executor service. |
| 3D Asset Delivery | ✅ PASS | Go Agent serves URDF/meshes via HTTP endpoints. |
| Operational Reliability | ✅ PASS | Exponential backoff for uploads, structural validation post-conversion. |

## Project Structure

### Documentation (this feature)

```text
specs/001-data-harvest-app/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
shared/
└── proto/              # Shared schemas

backend/                # Go Agent
├── internal/
│   ├── api/            # Gin handlers
│   ├── executor/       # Process lifecycle
│   ├── storage/        # BadgerDB wrappers
│   └── websocket/      # Telemetry streams
└── pkg/
    └── proto/          # Generated Go types

frontend/               # React UI
├── src/
│   ├── stores/         # Zustand state
│   ├── three/          # 3D visuals
│   └── websocket/      # Binary client
└── public/             # Static assets (robot models)

toolkit/                # Python Tools
├── scripts/            # Shell entry points
└── src/
    └── converter/      # LeRobot V3.0 logic
```

**Structure Decision**: Quad-project layout (Shared, Backend, Frontend, Toolkit) to strictly isolate hardware-level Python logic from dojo-level Go orchestration.

## Complexity Tracking

> No constitutional violations. Adheres to all requested boundaries.
