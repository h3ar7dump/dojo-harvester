<!--
Sync Impact Report:
- Version change: 1.0.0 → 1.1.0
- List of modified principles: Toolkit & Python Scripts (added uv package management requirements)
- Added sections: None
- Removed sections: None
- Templates requiring updates:
  - ✅ .specify/templates/plan-template.md (No changes needed, compatible)
  - ✅ .specify/templates/spec-template.md (No changes needed, compatible)
  - ✅ .specify/templates/tasks-template.md (No changes needed, compatible)
- Follow-up TODOs: None
-->
# Data Harvester App for Humanoid Robot Constitution

## Core Principles

### I. Architectural Boundaries

**Serialization Protocol Mandate**
- Protobuf MUST be used strictly for data pipeline communication (SensorData, VideoFrame, JointState, Metadata)
- JSON MUST be used for HTTP/REST API communication
- No mixing of serialization formats within the same communication channel

**Task Management Isolation**
- Client Web UI MUST communicate directly via HTTP REST to Remote Data Platform
- Task management MUST NOT route through Agent Backend
- This isolation ensures task queue integrity and prevents circular dependencies

**Execution Abstraction**
- Go Agent Backend MUST call shell scripts (`record.sh`, `convert.sh`, `upload_local.sh`) via executor service
- Direct Python execution from Go Agent Backend is PROHIBITED
- All toolkit operations MUST be encapsulated in script boundaries with explicit arguments

**3D Asset Delivery**
- Go Agent Backend MUST serve URDF/meshes over HTTP endpoints
- React Client MUST fetch 3D assets dynamically at runtime
- No compile-time embedding of 3D assets in the client

### II. Technology Stack Mandates

**Agent Backend (Go)**
The following libraries are REQUIRED:
- `gin-gonic/gin` - HTTP REST framework
- `gorilla/websocket` - Real-time streaming communication
- `fsnotify/fsnotify` - Dataset scanner and filesystem watcher
- `go-resty/resty/v2` - HTTP upload client
- `shirou/gopsutil/v4` - System status and telemetry
- `spf13/viper` - Configuration management
- `apache/arrow-go/v18` - Parquet file reading
- `go.uber.org/zap` - Structured logging
- `dgraph-io/badger/v4` - Embedded key-value store

**Client Web UI (React 19 + TypeScript + Vite)**
The following libraries are REQUIRED:
- `three` + `@react-three/fiber` - 3D rendering engine
- `urdf-loader` - URDF robot model parsing
- `shadcn/ui` + `Radix UI` + `TailwindCSS` - 2D UI components and styling
- `Zustand` - State management
- Native WebSockets - Communication with Agent Backend

**Toolkit & Python Scripts**
The following tools and libraries are REQUIRED:
- `uv` MUST be used to manage the Python project and dependencies
- Python projects MUST be initialized using `uv init -p 3.13.12 .`
- Packages MUST be installed using `uv add <package>`
- Python scripts MUST be executed using `uv run`
- `lerobot` (V3.0) - Robot learning framework compliance
- `mcap` - MCAP format handling
- `rosbags` - ROS bag conversion
- `pyarrow` - Parquet/data conversion utilities
- CLI execution with explicit arguments ONLY

### III. Data Storage & Format Compliance

**Raw Storage Structure**
Raw data MUST be stored at:
```
/data/<YYYY-MM-DD>/<session_id(task_id)>/raw/<episode_id>/
```

**LeRobot Storage Structure**
LeRobot-formatted data MUST be stored at:
```
/data/<YYYY-MM-DD>/<session_id(task_id)>/lerobot/<episode_id>/
```

**LeRobot V3.0 Compliance**
Every LeRobot dataset MUST contain:
- `meta/` directory with:
  - `info.json` - Dataset metadata
  - `stats.json` - Statistics metadata
  - `tasks.json` - Task definitions
- `data/` directory with:
  - `.parquet` chunk files
- `videos/` directory with:
  - `.mp4` video files

Non-compliance with LeRobot V3.0 format is NOT permitted.

### IV. Operational Reliability

**Upload Resilience**
- Failed uploads MUST be retried with exponential backoff
- Upload state MUST be tracked in BadgerDB
- No silent failures - all upload outcomes MUST be logged

**Upload Validation**
- Scanner MUST verify `meta/info.json` exists before adding to upload queue
- Datasets missing required metadata MUST be rejected from upload queue
- Validation errors MUST be surfaced to the Web UI

**Real-time Feedback**
- Web UI MUST reflect upload progress in real-time
- Robot telemetry MUST be streamed via WebSocket
- Live 3D joint movements MUST be rendered as they arrive
- Latency between robot state and UI display SHOULD be < 100ms

## Governance

This Constitution supersedes all other project practices and guidelines. All implementation decisions MUST align with these principles.

**Amendment Process:**
1. Any constitutional change requires a Plan document justifying the amendment
2. Changes MUST be reviewed for impact across all project components
3. Migration plans MUST be provided for breaking changes
4. Version numbers MUST follow semantic versioning (MAJOR.MINOR.BUILD)

**Compliance Verification:**
- All PRs MUST verify compliance with constitutional principles
- Technology stack deviations require explicit constitutional amendment
- Architecture boundary violations are BLOCKING issues

**Version**: 1.1.0 | **Ratified**: 2026-03-09 | **Last Amended**: 2026-03-10
