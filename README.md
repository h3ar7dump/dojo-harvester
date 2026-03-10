# Dojo Harvester

A comprehensive full-stack application for capturing, converting, and uploading high-fidelity robotic telemetry and sensor data in LeRobot V3.0 format.

## Architecture

The system is built with 3 main components:
- **Go Backend (`/backend`)**: Handles orchestrating jobs (recording, conversion, upload), state management with BadgerDB, and WebSocket telemetry streaming.
- **React Frontend (`/frontend`)**: Web UI for operators to manage tasks, perform prerequisite checks, and monitor real-time 3D visualizations.
- **Python Toolkit (`/toolkit`)**: Low-level shell scripts and Python utilities for hardware interfacing and PyArrow parquet conversions.

## Prerequisites

- Go 1.26+
- Node.js 20+ & npm
- Python 3.13.12+
- `uv` Python package manager

## Quick Start

### 1. Python Toolkit Setup
```bash
cd toolkit
uv init -p 3.13.12 .
uv add lerobot mcap rosbags pyarrow click grpcio-tools
uv run python -m grpc_tools.protoc -I../shared/proto --python_out=src/proto --grpc_python_out=src/proto ../shared/proto/telemetry.proto ../shared/proto/command.proto ../shared/proto/status.proto
```

### 2. Backend Setup
```bash
cd backend
go mod tidy
go build ./cmd/server/main.go
# Run the server
./main.exe
```

### 3. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

## System Configuration
Global configuration can be found in `backend/config/config.yaml`.
Metrics and states are stored locally in the `data/badger` directory created by the backend.
