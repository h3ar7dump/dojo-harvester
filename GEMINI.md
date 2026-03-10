---
description: Custom instructions for the Gemini CLI agent. 
---

# Project Context: Dojo Harvester

## Technology Stack

- **Language/Version**: Go 1.26 (Backend), Node.js 20+ with React 19 / TypeScript 5.x (Frontend), Python 3.13.12 (Toolkit)
- **Primary Dependencies**: `gin-gonic/gin`, `gorilla/websocket`, `go-resty/resty/v2`, `spf13/viper`, `go.uber.org/zap`, `three`, `@react-three/fiber`, `urdf-loader`, `shadcn/ui`, `Zustand`, `lerobot` (V3.0), `mcap`, `rosbags`, `pyarrow`
- **Database/Storage**: BadgerDB (embedded KV for state/metadata), Filesystem for raw/converted data
- **Project Type**: Full-stack web application with local hardware orchestrator

## Current Feature

- **Branch**: `001-data-harvest-app`
- **Specification**: `specs/001-data-harvest-app/spec.md`
- **Implementation Plan**: `specs/001-data-harvest-app/plan.md`
- **Tasks**: `specs/001-data-harvest-app/tasks.md`

## Architecture & Boundaries

1. **Serialization**: Protobuf for data pipeline (telemetry); JSON for HTTP/REST APIs.
2. **Task Isolation**: Web UI communicates directly with Remote Data Platform; does not route through Agent Backend.
3. **Execution**: Go Agent Backend executes toolkit Python via explicit shell scripts (`record.sh`, `convert.sh`, `upload_local.sh`).
4. **Storage Format**: LeRobot V3.0 (`meta/info.json`, `data/chunk-*.parquet`, `videos/*.mp4`).

<!-- MANUALLY_ADDED_CONTEXT_START -->
<!-- Add any manual context notes here. This section is preserved during updates. -->
<!-- MANUALLY_ADDED_CONTEXT_END -->
