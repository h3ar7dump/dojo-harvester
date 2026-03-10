# Quick Start: Data Harvest App

**Feature**: Data Harvest App for Humanoid Robot
**Date**: 2026-03-10

## Environment Initialization

### 1. Python Toolkit (uv)
We use `uv` for high-performance dependency management.
```bash
cd toolkit
uv init -p 3.13.12 .
uv add lerobot mcap rosbags pyarrow click grpcio-tools
```

### 2. Protobuf Generation
Generate Go and Python types from `shared/proto`.
```bash
# Python
uv run python -m grpc_tools.protoc -I../shared/proto --python_out=src/proto --grpc_python_out=src/proto ../shared/proto/*.proto

# Go
protoc -I=../shared/proto --go_out=pkg/proto ../shared/proto/*.proto
```

### 3. Backend (Go)
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 4. Frontend (React)
```bash
cd frontend
npm install
npm run dev
```

## Mandatory Architecture Check

Before committing, ensure:
1. **No direct robot connection** in React. All telemetry must come from the Go Backend WebSocket.
2. **No direct Python execution** in Go. All toolkit operations must go through `.sh` scripts in `toolkit/scripts/`.
3. **No mixing of protocols**. Use binary Protobuf for telemetry streams and JSON for all HTTP management APIs.
4. **LeRobot V3.0 Compliance**. Verify the generated dataset contains `meta/info.json` and `.parquet` chunks before triggering an upload.

## Local Testing

1. Use the **Task Dashboard** to claim a mock task (seeds provided in `backend/internal/api/tasks.go` for dev mode).
2. Complete the **Prerequisite Checklist**. Manual items require clicking "Verify"; auto items (storage/connect) trigger backend checks.
3. Launch **Recording**. Verify the 3D visualizer reflects the simulated joint movements from the `record.sh` stdout.
4. Click **Stop Recording** and monitor the **Conversion Progress**.
5. Initiate **Upload** and verify the job status moves to `completed` in the BadgerDB explorer or UI dashboard.
