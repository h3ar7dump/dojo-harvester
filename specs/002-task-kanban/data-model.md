# Data Model: Task Kanban Board

**Feature**: Task Kanban Board for Data Platform
**Date**: 2026-03-09
**Phase**: Phase 1 - Design & Contracts

## Entity Definitions

### UserSession

Represents an authenticated user session with JWT token.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| user_id | string | required | User identifier from data platform |
| username | string | required | User's login username |
| token | string | JWT format | Authentication token |
| expires_at | number | Unix timestamp | Token expiration time |
| login_timestamp | number | Unix timestamp | When user logged in |
| is_authenticated | boolean | computed | True if token exists and not expired |

**Storage**: localStorage (browser)

**Lifecycle**:
- Created: On successful login
- Updated: Token refresh (if implemented)
- Deleted: On logout or token expiration

---

### HarvestTask

A data collection assignment from the platform.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| task_id | string | UUID, PK | Unique task identifier |
| title | string | required | Task name/title |
| description | string | required | Detailed task description |
| status | enum | pending, in_progress, completed | Current task status |
| priority | enum | low, medium, high, critical | Task priority level |
| assignee_id | string | required | Assigned operator ID |
| due_date | string | ISO 8601 date | Target completion date |
| required_episodes | number | integer > 0 | Episodes needed for task |
| completed_episodes | number | integer >= 0 | Episodes completed so far |
| robot_configuration | object | JSON | Robot config (model, joints, cameras) |
| objectives | array | string[] | List of task objectives |
| created_at | string | ISO 8601 | Task creation timestamp |
| updated_at | string | ISO 8601 | Last update timestamp |

**Status Mapping to Kanban Columns**:
- `pending` → "To Do" column
- `in_progress` → "In Progress" column
- `completed` → "Done" column

---

### KanbanColumn

Represents a column in the kanban board.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | string | enum | Column identifier (todo, in_progress, done) |
| title | string | required | Display title (To Do, In Progress, Done) |
| status | enum | pending, in_progress, completed | Maps to task status |
| order | number | integer | Column display order |
| task_ids | array | string[] | Ordered list of task IDs in column |

**Default Columns**:
| id | title | status | order |
|----|-------|--------|-------|
| todo | To Do | pending | 1 |
| in_progress | In Progress | in_progress | 2 |
| done | Done | completed | 3 |

---

### TaskFilter

Represents active filter and sort criteria.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| priority | enum | low, medium, high, critical, null | Filter by priority (null = all) |
| due_date_from | string | ISO 8601 date, nullable | Filter start date |
| due_date_to | string | ISO 8601 date, nullable | Filter end date |
| robot_type | string | nullable | Filter by robot model/type |
| sort_by | enum | due_date, priority, created_at | Sort field |
| sort_direction | enum | asc, desc | Sort direction |

**Default State**:
- priority: null (all priorities)
- due_date_from: null
- due_date_to: null
- robot_type: null
- sort_by: due_date
- sort_direction: asc

**Storage**: Zustand store (with optional localStorage persistence)

---

### ApiError

Represents an API error response.

| Field | Type | Description |
|-------|------|-------------|
| code | string | Error code (e.g., UNAUTHORIZED, NETWORK_ERROR) |
| message | string | Human-readable error message |
| retryable | boolean | Whether the request can be retried |
| timestamp | string | ISO 8601 error timestamp |

---

## Entity Relationships

```
UserSession (1) ───< Authenticates >─── (N) HarvestTask (via API)
    │
    └──< Contains >─── TaskFilter (settings)

KanbanColumn (3) ───< Groups >─── HarvestTask[]
    - todo: status = pending
    - in_progress: status = in_progress
    - done: status = completed

TaskFilter ───< Filters >─── HarvestTask[]
```

## Data Flow

1. **Authentication**:
   ```
   Login Form → POST /auth/login → UserSession (token in localStorage)
   ```

2. **Initial Load**:
   ```
   Kanban Mount → GET /tasks (with JWT) → Filter by assignee_id → KanbanStore.tasks
   ```

3. **Filtering**:
   ```
   User Changes Filter → TaskFilter Update → Filter tasks locally → Re-render board
   ```

4. **Drag-and-Drop Status Update**:
   ```
   Drag Task to Column → Optimistic UI Update → PUT /tasks/{id}/status → Confirm/Revert
   ```

## State Management Strategy

### Zustand Stores

**authStore**:
```typescript
interface AuthStore {
  user: UserSession | null;
  isLoading: boolean;
  error: ApiError | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  checkAuth: () => boolean;
}
```

**kanbanStore**:
```typescript
interface KanbanStore {
  tasks: HarvestTask[];
  columns: KanbanColumn[];
  filters: TaskFilter;
  isLoading: boolean;
  error: ApiError | null;

  // Actions
  fetchTasks: () => Promise<void>;
  moveTask: (taskId: string, newStatus: TaskStatus) => Promise<void>;
  updateFilters: (filters: Partial<TaskFilter>) => void;
  resetFilters: () => void;

  // Computed
  filteredTasks: () => HarvestTask[];
  tasksByColumn: () => Record<string, HarvestTask[]>;
}
```

### localStorage Keys

| Key | Purpose | Expires |
|-----|---------|---------|
| `kanban_auth_token` | JWT token | With token expiration |
| `kanban_auth_user` | User info (user_id, username) | With token expiration |
| `kanban_filters` | Last used filter settings | Never (user preference) |

## API Integration Points

| Endpoint | Method | Request | Response |
|----------|--------|---------|----------|
| /auth/login | POST | {username, password} | {token, expires_in, user} |
| /tasks | GET | Headers: Authorization | HarvestTask[] |
| /tasks/{id}/status | PUT | {status} | HarvestTask |
| /tasks/{id} | GET | - | HarvestTask (detail) |

## Validation Rules

**HarvestTask**:
- task_id must be valid UUID
- required_episodes >= 1
- completed_episodes >= 0 and <= required_episodes
- due_date must be valid ISO 8601 date
- status transitions: pending → in_progress → completed

**TaskFilter**:
- due_date_from <= due_date_to (if both provided)
- priority must be valid enum value or null
- sort_by must be valid enum value
