---
description: "Actionable, dependency-ordered tasks for Data Harvest App implementation"
---

# Tasks: Data Harvest App for Humanoid Robot

**Input**: Design documents from `/specs/001-data-harvest-app/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

---

## Phase 1: Setup & Project Initialization

> **Goal**: Initialize all projects, set up shared protobuf definitions, and configure development environment
> **Prerequisites**: None (first phase)
> **Independent Test**: All three projects build successfully and protobuf types compile

### Project Structure Setup

- [x] T001 Create shared/proto directory structure with telemetry.proto, command.proto, status.proto
- [x] T002 Initialize Go backend project with go.mod at github.com/dojo-harvester/backend
- [x] T003 Initialize React frontend project with Vite + TypeScript template
- [x] T004 Initialize Python toolkit project with pyproject.toml

### Protobuf Code Generation

- [x] T005 [P] Install protoc-gen-go and generate Go protobuf types in backend/pkg/proto
- [x] T006 [P] Install grpcio-tools and generate Python protobuf types in toolkit/src/proto

### Dependency Installation

- [x] T007 Install Go dependencies: gin, gorilla/websocket, fsnotify, go-resty, gopsutil, viper, arrow-go, zap, badger
- [x] T008 [P] Install React dependencies: three, @react-three/fiber, urdf-loader, shadcn/ui, Radix UI, TailwindCSS, Zustand
- [x] T009 [P] Install Python dependencies: lerobot, mcap, rosbags, pyarrow, click

### Development Tooling

- [x] T010 Configure Go development tools (air for live reload, golangci-lint)
- [x] T011 [P] Configure React development tools (ESLint, Prettier, Vitest)
- [x] T012 [P] Configure Python development tools (pytest, black, ruff)

---

## Phase 2: Foundational Components (Blocking Prerequisites)

> **Goal**: Build foundational components required by all user stories
> **Prerequisites**: Phase 1 complete
> **Independent Test**: Backend server starts and serves health endpoint; WebSocket connects

### Backend Core Infrastructure

- [x] T013 Create backend configuration management (Viper) with config.yaml structure
- [x] T014 Set up Zap structured logging with configurable levels
- [x] T015 Initialize BadgerDB storage layer with connection pooling
- [x] T016 Create HTTP server bootstrap with Gin framework
- [x] T017 Implement health check endpoint GET /status

### WebSocket Infrastructure

- [x] T018 Implement WebSocket connection manager with concurrent connection support
- [x] T019 Create binary message frame handler for Protobuf [1 byte type][payload bytes]
- [x] T020 Implement telemetry message broadcaster for 20-30 concurrent robots

### Shell Script Executor

- [x] T021 Create executor service for shell script invocation via os/exec
- [x] T022 Implement process lifecycle management (start, stop, cancel via context)
- [x] T023 Add resource limits (CPU, memory, timeout) for script execution

### Frontend Core Infrastructure

- [x] T024 Set up Zustand store structure for global state management
- [x] T025 Create WebSocket client with binary message support
- [x] T026 Implement JWT authentication context and token storage
- [x] T027 Create API client service with axios/fetch for REST endpoints

### UI Component Foundation

- [x] T028 [P] Set up TailwindCSS configuration with custom theme
- [x] T029 [P] Configure shadcn/ui and Radix UI component primitives
- [x] T030 [P] Create base layout components (Header, Sidebar, Main content area)

---

## Phase 3: User Story 1 - Guided Recording Preparation (P1)

> **Goal**: Implement prerequisite checklist workflow that must complete before recording can start
> **Prerequisites**: Phase 2 complete
> **Independent Test**: User can walk through checklist and record button only enables when all prerequisites verified

### Backend - Session & Prerequisite APIs

- [x] T031 [US1] Create RecordingSession model and BadgerDB storage operations
- [x] T032 [US1] Create PrerequisiteItem model and BadgerDB storage operations
- [x] T033 [US1] Implement POST /sessions endpoint for creating new recording sessions
- [x] T034 [US1] Implement GET /sessions/{session_id}/prerequisites endpoint
- [x] T035 [US1] Implement POST /prerequisites/{item_id}/verify endpoint for manual verification
- [x] T036 [US1] Implement automatic prerequisite verification for system-level checks (storage, connectivity)

### Frontend - Recording Preparation UI

- [x] T037 [US1] [P] Create RecordingPreparationPage component
- [x] T038 [US1] [P] Create PrerequisiteChecklist component with checkboxes
- [x] T039 [US1] Create StartRecordingButton component with disabled state logic
- [x] T040 [US1] [P] Create prerequisite status indicators (pending/verified/failed)
- [x] T041 [US1] Implement prerequisite verification API integration
- [x] T042 [US1] Add real-time storage space monitoring display

---

## Phase 4: User Story 2 - Recording Execution and Monitoring (P1)

> **Goal**: Enable recording launch via external scripts with real-time telemetry monitoring
> **Prerequisites**: Phase 3 (US1) complete - needs sessions from US1
> **Independent Test**: User can launch recording, see live telemetry, and stop recording

### Backend - Recording Execution

- [x] T043 [US2] Create record.sh shell script entry point
- [x] T044 [US2] Implement POST /sessions/{session_id}/start endpoint
- [x] T045 [US2] Integrate executor service to launch recording script with session_id, output_path arguments
- [x] T046 [US2] Implement POST /sessions/{session_id}/stop endpoint
- [x] T047 [US2] Create robot disconnection detection and abort logic
- [x] T048 [US2] Implement storage space monitoring during recording
- [x] T049 [US2] Add recording state transition handling (preparing → recording → converting)

### Backend - Telemetry Streaming

- [x] T050 [US2] Create RobotTelemetry model with Protobuf serialization
- [x] T051 [US2] Implement telemetry ingestion from recording script
- [x] T052 [US2] Create WebSocket telemetry broadcaster with <100ms latency
- [x] T053 [US2] Add error handling and robot disconnection detection

### Frontend - Recording Dashboard

- [x] T054 [US2] [P] Create RecordingDashboardPage component
- [x] T055 [US2] [P] Create LiveTelemetryPanel component for joint positions, camera feeds, sensor data
- [x] T056 [US2] [P] Implement three.js + @react-three/fiber 3D robot visualization
- [x] T057 [US2] Create recording duration timer and progress display
- [x] T058 [US2] [P] Create StopRecordingButton component
- [x] T059 [US2] Add recording error display and troubleshooting guidance
- [x] T060 [US2] Integrate WebSocket telemetry consumer with Zustand store

---

## Phase 5: User Story 3 - Data Conversion and Formatting (P2)

> **Goal**: Convert raw recordings to LeRobot V3.0 format
> **Prerequisites**: Phase 4 (US2) complete - needs completed recordings
> **Independent Test**: Raw recording converts to valid LeRobot V3.0 dataset structure

### Backend - Conversion Pipeline

- [x] T061 [US3] Create Dataset model and BadgerDB storage operations
- [x] T062 [US3] Create convert.sh shell script entry point
- [x] T063 [US3] Implement POST /datasets/{dataset_id}/convert endpoint
- [x] T064 [US3] Integrate executor service to launch conversion script
- [x] T065 [US3] Create LeRobot V3.0 validation logic (meta/, data/, videos/ structure check)
- [x] T066 [US3] Implement GET /datasets/{dataset_id}/status endpoint with progress tracking
- [x] T067 [US3] Add conversion failure handling with retry capability

### Python Toolkit - Conversion Scripts

- [x] T068 [US3] [P] Create Python conversion module using lerobot V3.0 library
- [x] T069 [US3] [P] Implement raw data to Parquet conversion using pyarrow
- [x] T070 [US3] Generate LeRobot V3.0 metadata files (info.json, stats.json, tasks.json)
- [x] T071 [US3] Add video file handling and organization

### Frontend - Conversion UI

- [x] T072 [US3] [P] Create ConversionStatusPage component
- [x] T073 [US3] [P] Create conversion progress bar with percentage display
- [x] T074 [US3] Add conversion logs viewer component
- [x] T075 [US3] Implement retry and abort functionality for failed conversions
- [x] T076 [US3] Display LeRobot V3.0 validation results

---

## Phase 6: User Story 4 - Dataset Upload to Data Platform (P2)

> **Goal**: Upload converted datasets to remote platform with resilience
> **Prerequisites**: Phase 5 (US3) complete - needs converted datasets
> **Independent Test**: Dataset uploads successfully with retry on transient failures

### Backend - Upload Management

- [x] T077 [US4] Create UploadJob model and BadgerDB storage operations
- [x] T078 [US4] Implement upload queue management (pending, in-progress, completed, failed)
- [x] T079 [US4] Create upload.sh shell script entry point for local upload operations
- [x] T080 [US4] Implement POST /datasets/{dataset_id}/upload endpoint
- [x] T081 [US4] Create resilient upload client with exponential backoff retry (1s, 2s, 4s, 8s, 16s)
- [x] T082 [US4] Implement resume capability from last successful chunk
- [x] T083 [US4] Create GET /uploads/{job_id} endpoint for progress tracking
- [x] T084 [US4] Add upload state persistence for crash recovery
- [x] T085 [US4] Implement GET /uploads/queue endpoint for queue status

### Frontend - Upload UI

- [x] T086 [US4] [P] Create UploadProgressPage component
- [x] T087 [US4] [P] Create upload progress display with bytes transferred, speed, ETA
- [x] T088 [US4] Implement upload retry and cancel functionality
- [x] T089 [US4] Add upload queue visualization
- [x] T090 [US4] Display upload completion and platform dataset URL

---

## Phase 7: User Story 5 - Data Platform Task Management (P3)

> **Goal**: Integrate with data platform for task viewing and claiming
> **Prerequisites**: Phase 2 complete (foundational) - can be developed in parallel
> **Independent Test**: User can log in, view tasks, claim task, and see progress synced

### Backend - Platform Integration

- [x] T091 [US5] Create HarvestTask model and BadgerDB caching layer
- [x] T092 [US5] Implement POST /auth/login proxy to data platform
- [x] T093 [US5] Create data platform API client with JWT authentication
- [x] T094 [US5] Implement GET /tasks endpoint (cached from platform)
- [x] T095 [US5] Create POST /tasks/{task_id}/claim endpoint
- [x] T096 [US5] Implement PUT /tasks/{task_id}/progress endpoint for syncing completion
- [x] T097 [US5] Add task status synchronization within 10 seconds of recording completion

### Frontend - Task Dashboard

- [x] T098 [US5] [P] Create LoginPage component with username/password form
- [x] T099 [US5] [P] Create TaskDashboardPage component
- [x] T100 [US5] [P] Create TaskList component with filtering by status/priority
- [x] T101 [US5] [P] Create TaskCard component showing task details
- [x] T102 [US5] Create TaskDetailView component with robot configuration
- [x] T103 [US5] Implement StartTaskButton with recording parameter pre-population
- [x] T104 [US5] Add task progress display (episodes completed, duration, quality metrics)

---

## Phase 8: Polish & Cross-Cutting Concerns

> **Goal**: Final polish, error handling, monitoring, and performance optimization
> **Prerequisites**: All user stories (US1-US5) complete
> **Independent Test**: Full end-to-end workflow executes successfully with graceful error handling

### Error Handling & Edge Cases

- [x] T105 Create global error boundary component for React
- [x] T106 Implement graceful handling for data platform unreachable during upload
- [x] T107 Add conversion script crash recovery with automatic restart
- [x] T108 Implement dataset validation failure handling with detailed error messages
- [x] T109 Add prerequisite verification failure diagnostics

### Performance Optimization

- [x] T110 Optimize WebSocket message serialization for <100ms latency
- [x] T111 Implement connection pooling for 20-30 concurrent robot support
- [x] T112 Add telemetry data buffering and batching strategies
- [x] T113 Optimize 3D rendering performance for smooth 60 FPS
- [x] T114 Implement upload chunk size optimization

### Monitoring & Observability

- [x] T115 Add structured logging for all operations
- [x] T116 Create system metrics collection (recording success rate, upload success rate)
- [x] T117 Implement health check endpoints for all services
- [x] T118 Add performance telemetry collection

### Documentation & Deployment

- [x] T119 [P] Create comprehensive README with setup instructions
- [x] T120 [P] Write API documentation for external consumers
- [x] T121 Create deployment configuration (Docker, docker-compose)
- [x] T122 Add environment-specific configuration examples

---

## Dependency Graph

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational)
    ↓
    ├──→ Phase 3 (US1: Prerequisites) ──→ Phase 4 (US2: Recording)
    │                                          ↓
    │                                     Phase 5 (US3: Conversion)
    │                                          ↓
    │                                     Phase 6 (US4: Upload)
    │
    └──→ Phase 7 (US5: Task Management) [parallel with US1-US4]
                                              ↓
                                         Phase 8 (Polish)
```

### Parallel Execution Opportunities

**Within Phase 1:**
- T001, T002, T003, T004 (project initialization) - all independent
- T005, T006 (protobuf generation) - independent once protos defined
- T007, T008, T009 (dependency installation) - independent
- T010, T011, T012 (tooling setup) - independent

**Within Phase 2:**
- T013-T017 (Go infrastructure) - sequential dependencies
- T018-T020 (WebSocket) - sequential
- T021-T023 (executor) - sequential
- T024-T027 (React infrastructure) - parallel with Go backend
- T028-T030 (UI foundation) - parallel with backend

**User Story Development:**
- US5 (Task Management) can be developed in parallel with US1-US4 after Phase 2
- US3 and US4 have strong sequential dependency (conversion → upload)
- US1 and US2 have strong sequential dependency (prerequisites → recording)

**MVP Scope Recommendation:**
Focus on US1 + US2 first - this delivers the core value of recording with prerequisites. Users can manually manage task tracking and upload separately initially.

---

## Summary

| Phase | Tasks | Story | Independent Test Criteria |
|-------|-------|-------|---------------------------|
| Phase 1 | 12 | Setup | All projects build, protobuf compiles |
| Phase 2 | 18 | Foundational | Server starts, WebSocket connects, health OK |
| Phase 3 | 12 | US1 | Checklist blocks recording until all verified |
| Phase 4 | 18 | US2 | Recording launches, telemetry streams, stops cleanly |
| Phase 5 | 15 | US3 | Raw recording converts to valid LeRobot V3.0 dataset |
| Phase 6 | 14 | US4 | Dataset uploads with retry on transient failures |
| Phase 7 | 13 | US5 | Login works, tasks display, claiming updates status |
| Phase 8 | 18 | Polish | End-to-end workflow executes, errors handled gracefully |

**Total Tasks**: 120
**MVP Tasks** (US1 + US2): 62 tasks
**Parallelizable Tasks**: 44 tasks (marked with [P])

---

## Next Steps

1. **Start with MVP**: Implement Phase 1-4 (US1 + US2) for core recording functionality
2. **Parallel Development**: US5 (Task Management) can be developed by separate team after Phase 2
3. **Incremental Delivery**: Each user story delivers independently testable functionality
4. **Run `/speckit.execute`**: To begin executing tasks from Phase 1
