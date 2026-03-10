# Data Model: Data Harvest App

**Feature**: Data Harvest App for Humanoid Robot
**Date**: 2026-03-10
**Phase**: Phase 1 - Design & Contracts

## Entity Definitions

### RecordingSession

Represents a single data collection event.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| session_id | string | UUID, PK | Unique session identifier |
| task_id | string | FK → HarvestTask | Associated collection task |
| robot_id | string | Required | ID of the robot used |
| operator_id | string | Required | ID of the operator |
| status | enum | [preparing, recording, converting, uploading, completed, failed] | Current session state |
| raw_path | string | Local Path | Path to raw (MCAP/ROS) data |
| storage_path | string | Local Path | Target LeRobot V3.0 path |
| duration_sec | float | - | Recording duration |
| error_msg | string | Nullable | Last failure reason |
| created_at | datetime | ISO 8601 | Initialization time |
| updated_at | datetime | ISO 8601 | Last state change |

---

### PrerequisiteItem

Checklist items required for a session to transition from `preparing` to `recording`.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| item_id | string | PK | e.g., "storage_check" |
| session_id | string | FK → RecordingSession | Parent session |
| name | string | Required | Display name |
| type | enum | [manual, auto] | Verification method |
| status | enum | [pending, verified, failed] | Verification result |
| metadata | json | Nullable | e.g., free space remaining |

---

### UploadJob

Tracks the background upload process to the data platform.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| job_id | string | PK | Unique job ID |
| session_id | string | FK → RecordingSession | Session being uploaded |
| status | enum | [pending, in_progress, completed, failed] | Current status |
| progress | float | 0.0 - 1.0 | Transfer percentage |
| bytes_total | int64 | - | Total size |
| bytes_sent | int64 | - | Transferred size |
| retries | int | Max 5 | Current attempt count |
| last_retry | datetime | ISO 8601 | Timestamp of last attempt |

---

### HarvestTask

Cached metadata for tasks assigned from the remote platform.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| task_id | string | PK | Platform task ID |
| title | string | Required | Task name |
| description | string | Required | Task instructions |
| robot_config | json | Required | Target robot joint/sensor config |
| required_eps | int | - | Target episode count |
| completed_eps | int | - | Local completion count |
| priority | enum | [P1, P2, P3] | Task priority |

---

### RobotTelemetry (Protobuf)

Real-time state message streamed via WebSocket.

| Field | Type | Description |
|-------|------|-------------|
| timestamp_ns | int64 | Unix timestamp in nanoseconds |
| joint_positions | float[] | Array of current joint angles |
| joint_velocities | float[] | Array of current joint velocities |
| end_effector_pose | Pose | XYZ + Quaternion |
| battery_level | float | Percentage |
| status_flags | string[] | Active hardware errors |

## Relationships

- A **HarvestTask** can have multiple **RecordingSessions**.
- A **RecordingSession** has many **PrerequisiteItems**.
- A **RecordingSession** has at most one **UploadJob** active at a time.
- **RobotTelemetry** is transient and associated with an active `recording` session.

## State Transitions

### RecordingSession
`preparing` -> `recording` (via verify prerequisites)
`recording` -> `converting` (via script exit 0)
`converting` -> `uploading` (via validation pass)
`uploading` -> `completed` (via upload success)
`*` -> `failed` (any error or abort)

### UploadJob
`pending` -> `in_progress` -> `completed`
`in_progress` -> `failed` (after 5 retries)
`failed` -> `in_progress` (manual resume)
