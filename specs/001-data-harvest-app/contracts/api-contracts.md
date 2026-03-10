# API Contracts: Data Harvest App

**Feature**: Data Harvest App for Humanoid Robot
**Date**: 2026-03-10
**Phase**: Phase 1 - Design & Contracts

## 1. Agent Backend REST API (Go)

All endpoints use `application/json` and require a valid JWT `Authorization: Bearer <token>`.

### Recording Lifecycle

**POST `/api/v1/sessions`**
Initialize a new recording session.
- Request: `{ "task_id": "uuid", "robot_id": "string" }`
- Response: `201 Created` with `RecordingSession` entity.

**GET `/api/v1/sessions/:id/prerequisites`**
Retrieve checklist items for a session.
- Response: `200 OK` with `Array<PrerequisiteItem>`.

**POST `/api/v1/sessions/:id/start`**
Execute `record.sh`. Session must be `preparing` and all required prerequisites `verified`.
- Response: `202 Accepted`.

**POST `/api/v1/sessions/:id/stop`**
Signal `record.sh` to stop gracefully. Triggers `convert.sh`.
- Response: `202 Accepted`.

### Upload Management

**GET `/api/v1/uploads`**
List all active and historical upload jobs.
- Response: `200 OK` with `Array<UploadJob>`.

**POST `/api/v1/uploads/:id/resume`**
Manually restart a `failed` upload job.
- Response: `202 Accepted`.

## 2. WebSocket Telemetry Protocol

**Endpoint**: `/ws/telemetry`
**Binary Type**: `arraybuffer`

### Message Format: `[1 byte header][protobuf payload]`

| Header (Hex) | Message Type | Protobuf Definition |
|--------------|--------------|---------------------|
| `0x01` | `TelemetryFrame` | `telemetry.proto` |
| `0x02` | `StatusUpdate` | `status.proto` |

## 3. Remote Data Platform API (External)

*Note: The Web UI communicates directly with these endpoints per the Constitution.*

**GET `/api/tasks`**
List available harvest tasks for the operator.
- Auth: JWT (Platform issued).

**POST `/api/datasets/upload`**
Resumable upload endpoint for LeRobot datasets.
- Handled by Go Backend `uploader` service using `go-resty`.

## 4. Toolkit Script Interface

**record.sh**
- Interface: `record.sh <session_id> <output_path> <duration_limit>`
- Output: Logs telemetry to stdout in JSON format for the Go backend to buffer.

**convert.sh**
- Interface: `convert.sh <session_id> <raw_path> <lerobot_path>`
- Exit Code: `0` on success (valid LeRobot V3.0), non-zero on failure.

**upload_local.sh**
- Interface: `upload_local.sh <job_id> <dataset_path> <platform_url>`
- Responsibility: Handles chunked transfer via `curl` or internal python logic.
