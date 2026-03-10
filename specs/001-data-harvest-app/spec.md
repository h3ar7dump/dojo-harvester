# Feature Specification: Data Harvest App for Humanoid Robot

**Feature Branch**: `001-data-harvest-app`
**Created**: 2026-03-09
**Status**: Draft
**Input**: User description: "Build an user-friendly data harvest app for humanoid robot in dojo to help user easy and guided to prepare recording prerequisties and launch recording and conversion by external script, then upload converted dataset to customized data platform. and user can login to data platfrom data harvestion task and proceed the task."

## Clarifications

### Session 2026-03-09

- **Q**: What authentication method should the app use for data platform access? → **A**: Username/password with JWT token returned in JSON response
- **Q**: What is the scale and concurrency model for recordings? → **A**: Single operator per robot, single recording per robot at a time, supporting 20-30 concurrent robots per dojo
- **Q**: What should happen when the robot disconnects mid-recording? → **A**: Immediately abort and discard partial recording
- **Q**: How should the system handle insufficient storage space during recording? → **A**: Alert operator and gracefully stop while discarding current recording
- **Q**: What is the policy when the data platform is completely unreachable for an extended period during upload? → **A**: Retry with backoff for a fixed period (e.g., 24 hours), then mark as failed and alert the operator for manual intervention.
- **Q**: How should the system handle and recover from a crashed conversion script? → **A**: Mark the conversion as failed immediately and require the operator to manually review logs and restart it.
- **Q**: How does the system coordinate 20-30 concurrent recording sessions across multiple robots in the same dojo? → **A**: Each robot runs its own completely independent, isolated instance of the app.
- **Q**: What happens when a dataset fails LeRobot V3.0 validation before upload? → **A**: Mark the dataset as "Invalid", quarantine it locally, and require manual operator inspection/repair.
- **Q**: What happens when an auto-verifiable prerequisite item fails or cannot be automatically verified? → **A**: The system allows the operator to manually override and force-verify the prerequisite, with a mandatory override log/reason.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Guided Recording Preparation (Priority: P1)

As a robot operator in the dojo, I want to be guided through a pre-recording checklist so that I can ensure all prerequisites are met before starting data collection.

**Why this priority**: Without proper preparation, recordings may fail or produce invalid data, wasting time and resources. This is the foundation of the entire data harvesting workflow.

**Independent Test**: Can be fully tested by walking through the prerequisite checklist without actually recording and verifying all items must be confirmed before the record button becomes active.

**Acceptance Scenarios**:

1. **Given** the operator has opened the data harvest app, **When** they select "Start New Recording", **Then** they are presented with a guided checklist of prerequisites (robot powered on, cameras calibrated, storage space verified, etc.)
2. **Given** the operator is viewing the prerequisite checklist, **When** they check off all required items, **Then** the "Launch Recording" button becomes enabled
3. **Given** the operator has not completed all prerequisites, **When** they attempt to proceed, **Then** the app prevents progression and highlights the missing items
4. **Given** a prerequisite item requires verification (e.g., storage space), **When** the operator requests automatic verification, **Then** the app checks system status and marks the item accordingly

---

### User Story 2 - Recording Execution and Monitoring (Priority: P1)

As a robot operator, I want to launch recording via external scripts and monitor the progress in real-time so that I can verify data is being captured correctly.

**Why this priority**: The core value of the app is enabling data collection. Real-time monitoring ensures data quality and allows operators to intervene if issues arise.

**Independent Test**: Can be fully tested by launching a recording session and verifying the UI displays real-time status, progress indicators, and telemetry data from the robot.

**Acceptance Scenarios**:

1. **Given** all prerequisites are met, **When** the operator clicks "Launch Recording", **Then** the app invokes the external recording script and displays a real-time progress dashboard
2. **Given** a recording is in progress, **When** the operator views the dashboard, **Then** they see live telemetry (joint positions, camera feeds, sensor data) and recording duration
3. **Given** a recording is in progress, **When** the operator clicks "Stop Recording", **Then** the external script receives the stop signal and the app transitions to the conversion phase
4. **Given** the recording script reports an error, **When** the error occurs, **Then** the app displays a clear error message with troubleshooting guidance

---

### User Story 3 - Data Conversion and Formatting (Priority: P2)

As a robot operator, I want recorded data to be automatically converted to LeRobot V3.0 format so that it is compatible with standardized training pipelines.

**Why this priority**: Raw recordings need transformation before they are useful for machine learning. Automated conversion reduces manual effort and ensures format consistency.

**Independent Test**: Can be fully tested by completing a recording and verifying the output is structured as LeRobot V3.0 format with proper metadata, statistics, and video files.

**Acceptance Scenarios**:

1. **Given** a recording has been completed, **When** the operator initiates conversion, **Then** the app invokes the external conversion script with the recorded data path
2. **Given** conversion is in progress, **When** the operator views the status, **Then** they see progress percentage, estimated time remaining, and any conversion warnings
3. **Given** conversion has completed successfully, **When** the operator reviews the output, **Then** they see a valid LeRobot V3.0 dataset structure (meta/, data/, videos/ directories)
4. **Given** conversion fails, **When** the error is detected, **Then** the app displays detailed logs and allows the operator to retry or abort

---

### User Story 4 - Dataset Upload to Data Platform (Priority: P2)

As a robot operator, I want to upload converted datasets to the remote data platform so that they are available for training and collaboration.

**Why this priority**: Uploading makes the data accessible to the broader team and enables downstream machine learning workflows. Resilient upload ensures no data is lost.

**Independent Test**: Can be fully tested by initiating an upload and verifying the dataset appears on the remote platform with all metadata intact.

**Acceptance Scenarios**:

1. **Given** a dataset has been successfully converted, **When** the operator clicks "Upload to Platform", **Then** the app validates the dataset structure and begins upload with progress tracking
2. **Given** an upload is in progress, **When** the operator views the upload status, **Then** they see bytes transferred, upload speed, and estimated completion time
3. **Given** an upload fails due to network issues, **When** the failure occurs, **Then** the app automatically retries with exponential backoff and resumes from the last successful chunk
4. **Given** the upload completes successfully, **When** the operator checks the data platform, **Then** the dataset is visible with all metadata (session info, episode count, duration)

---

### User Story 5 - Data Platform Task Management (Priority: P3)

As a robot operator, I want to log in to the data platform to view and manage my data harvest tasks so that I can track progress and prioritize work.

**Why this priority**: Task management provides visibility into the overall data collection workflow and enables coordination across the team.

**Independent Test**: Can be fully tested by logging into the platform and verifying task lists, statuses, and the ability to claim or update tasks.

**Acceptance Scenarios**:

1. **Given** the operator has valid credentials, **When** they log in to the data platform, **Then** they see a dashboard of their assigned and available harvest tasks
2. **Given** the operator is viewing the task list, **When** they select a specific task, **Then** they see task details (robot configuration, objectives, required episodes, current progress)
3. **Given** a task is in "Pending" status, **When** the operator clicks "Start Task", **Then** the task status updates to "In Progress" and the app pre-populates recording parameters
4. **Given** the operator has completed episodes for a task, **When** they view the task on the platform, **Then** they see updated progress (episodes completed, total duration, data quality metrics)

---

### Edge Cases

- When the robot disconnects mid-recoring, the system MUST immediately abort the recording and discard partial data to maintain data integrity
- When storage space becomes insufficient during recording, the system MUST alert the operator and gracefully stop while discarding the current recording
- When the data platform is completely unreachable during upload, the system MUST retry with exponential backoff for a fixed period (e.g., 24 hours). If it still fails, mark as failed and alert the operator for manual intervention.
- Coordination of 20-30 concurrent recording sessions across multiple robots is handled by running independent, isolated instances of the app on each robot.
- When an auto-verifiable prerequisite item fails or cannot be automatically verified, the system MUST allow the operator to manually override and force-verify the prerequisite, with a mandatory override log/reason.
- When a conversion script crashes, the system MUST mark the conversion as failed immediately and require the operator to manually review logs and restart it.
- When a dataset fails LeRobot V3.0 validation before upload, the system MUST mark the dataset as "Invalid", quarantine it locally, and require manual operator inspection/repair.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The app MUST present a guided checklist interface for recording prerequisites before allowing recording to start
- **FR-002**: The app MUST disable the recording launch button until all required prerequisites are confirmed
- **FR-003**: The app MUST support automatic verification of system-level prerequisites (storage space, robot connectivity, camera status)
- **FR-004**: The app MUST invoke external recording scripts via shell execution with explicit arguments (session ID, output path, duration limits)
- **FR-005**: The app MUST display real-time recording telemetry including joint positions, camera feeds, and sensor data
- **FR-006**: The app MUST allow operators to gracefully stop ongoing recordings
- **FR-007**: The app MUST invoke external conversion scripts to transform raw recordings into LeRobot V3.0 format
- **FR-008**: The app MUST validate converted datasets contain required LeRobot V3.0 structure (meta/, data/, videos/ directories with proper files)
- **FR-009**: The app MUST upload validated datasets to the remote data platform with progress tracking
- **FR-010**: The app MUST implement resilient upload with exponential backoff retry for failed transfers up to a fixed period (e.g., 24 hours) before marking as failed
- **FR-011**: The app MUST authenticate operators with the data platform using username/password credentials and receive a JWT token in JSON response for subsequent authorized requests
- **FR-012**: The app MUST display a task dashboard showing assigned tasks, their status, and progress metrics
- **FR-013**: The app MUST allow operators to claim and start data harvest tasks from the platform
- **FR-014**: The app MUST synchronize task progress back to the data platform after recording completion
- **FR-015**: The app MUST persist upload state in local storage to enable recovery from app restarts
- **FR-016**: The system MUST support 20-30 concurrent robot recording sessions per dojo, with each robot managed by a single operator and limited to one active recording at a time
- **FR-017**: The system MUST immediately abort and discard partial recordings when robot disconnection is detected to maintain data integrity
- **FR-018**: The system MUST continuously monitor available storage space during recording and gracefully abort with operator notification when space falls below a safe threshold

### Key Entities *(include if feature involves data)*

- **RecordingSession**: Represents a single data collection episode. Attributes: session ID, start time, duration, status (preparing/recording/converting/uploaded), associated task ID, storage location.
- **PrerequisiteItem**: A checklist item for recording preparation. Attributes: name, description, verification type (manual/auto), verification status, required/optional flag.
- **Dataset**: Converted LeRobot V3.0 dataset. Attributes: dataset ID, source session ID, storage paths, validation status, metadata (info.json, stats.json, tasks.json), video files list.
- **UploadJob**: Tracks dataset upload to remote platform. Attributes: job ID, dataset ID, upload status (pending/in-progress/completed/failed), progress percentage, retry count, error logs.
- **HarvestTask**: A data collection assignment from the platform. Attributes: task ID, title, description, robot configuration, required episodes, completed episodes, task status, assignee, priority, due date.
- **RobotTelemetry**: Real-time robot state during recording. Attributes: timestamp, joint positions, camera frames, sensor readings, battery level, error flags.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Operators can complete prerequisite checklist and start recording within 5 minutes of app launch
- **SC-002**: 95% of recordings complete successfully without manual intervention
- **SC-003**: Converted datasets achieve 100% LeRobot V3.0 format compliance on validation
- **SC-004**: Dataset uploads achieve 99.5% success rate with automatic retry handling transient failures
- **SC-005**: Operators can view real-time telemetry with latency under 100ms from robot state change to UI update
- **SC-006**: Task status synchronization between app and data platform completes within 10 seconds of recording completion
- **SC-007**: 90% of operators successfully complete their first data harvest task without requiring support
- **SC-008**: The app gracefully handles and recovers from 95% of network interruptions during upload
