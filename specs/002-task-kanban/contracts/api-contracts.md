# API Contracts: Task Kanban Board

**Feature**: Task Kanban Board for Data Platform
**Date**: 2026-03-09

## Data Platform API Integration

### Base URL
```
{DATA_PLATFORM_URL}/api/v1
```

### Authentication

All endpoints (except login) require JWT token in Authorization header:
```
Authorization: Bearer {jwt_token}
```

---

## Authentication Endpoints

### POST /auth/login

Authenticate user and receive JWT token.

**Request**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response (200 OK)**:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "user_id": "uuid",
    "username": "operator_001",
    "email": "operator@example.com"
  }
}
```

**Response (401 Unauthorized)**:
```json
{
  "error": "INVALID_CREDENTIALS",
  "message": "Username or password is incorrect"
}
```

**Error Codes**:
| Code | HTTP Status | Description |
|------|-------------|-------------|
| INVALID_CREDENTIALS | 401 | Username or password incorrect |
| ACCOUNT_LOCKED | 403 | Account temporarily locked |
| RATE_LIMITED | 429 | Too many login attempts |

---

## Task Management Endpoints

### GET /tasks

Retrieve all harvest tasks assigned to the authenticated user.

**Query Parameters**:
- `status` (optional): Filter by status (pending, in_progress, completed)
- `priority` (optional): Filter by priority (low, medium, high, critical)
- `due_date_from` (optional): Filter start date (ISO 8601)
- `due_date_to` (optional): Filter end date (ISO 8601)
- `sort_by` (optional): Sort field (due_date, priority, created_at)
- `sort_direction` (optional): asc or desc

**Response (200 OK)**:
```json
{
  "tasks": [
    {
      "task_id": "uuid",
      "title": "Grasp Training Dataset",
      "description": "Collect grasp motions for training",
      "status": "in_progress",
      "priority": "high",
      "assignee_id": "operator_001",
      "due_date": "2026-03-15",
      "required_episodes": 10,
      "completed_episodes": 5,
      "robot_configuration": {
        "model": "humanoid_v2",
        "joints": ["shoulder", "elbow", "wrist"],
        "cameras": ["cam_0", "cam_1"]
      },
      "objectives": ["Collect left hand grasps", "Collect right hand grasps"],
      "created_at": "2026-03-01T10:00:00Z",
      "updated_at": "2026-03-09T14:30:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "per_page": 50
}
```

**Response (401 Unauthorized)**:
```json
{
  "error": "UNAUTHORIZED",
  "message": "Invalid or expired token"
}
```

---

### GET /tasks/{task_id}

Get detailed information about a specific task.

**Response (200 OK)**:
```json
{
  "task_id": "uuid",
  "title": "Grasp Training Dataset",
  "description": "Collect grasp motions for training",
  "status": "in_progress",
  "priority": "high",
  "assignee_id": "operator_001",
  "due_date": "2026-03-15",
  "required_episodes": 10,
  "completed_episodes": 5,
  "robot_configuration": {
    "model": "humanoid_v2",
    "joint_count": 7,
    "cameras": [
      {"id": "cam_0", "resolution": "1920x1080"},
      {"id": "cam_1", "resolution": "1920x1080"}
    ]
  },
  "objectives": [
    "Collect left hand grasps",
    "Collect right hand grasps",
    "Record from multiple angles"
  ],
  "recording_parameters": {
    "fps": 30,
    "duration_per_episode": 60,
    "data_format": "lerobot_v3"
  },
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-03-09T14:30:00Z",
  "completed_at": null
}
```

---

### PUT /tasks/{task_id}/status

Update the status of a task (for drag-and-drop).

**Request**:
```json
{
  "status": "in_progress"
}
```

**Response (200 OK)**:
```json
{
  "task_id": "uuid",
  "status": "in_progress",
  "updated_at": "2026-03-09T15:00:00Z"
}
```

**Response (400 Bad Request)**:
```json
{
  "error": "INVALID_STATUS_TRANSITION",
  "message": "Cannot transition from completed to pending"
}
```

**Response (403 Forbidden)**:
```json
{
  "error": "NOT_ASSIGNED",
  "message": "Task is not assigned to you"
}
```

**Status Transition Rules**:
- pending → in_progress (allowed)
- in_progress → completed (allowed)
- pending → completed (not allowed, must go through in_progress)
- completed → any (not allowed)

---

### POST /tasks/{task_id}/claim

Claim a pending task.

**Response (200 OK)**:
```json
{
  "task_id": "uuid",
  "status": "in_progress",
  "assignee_id": "operator_001",
  "claimed_at": "2026-03-09T15:00:00Z"
}
```

---

## Error Response Format

All error responses follow this structure:

```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable description",
  "details": {},
  "timestamp": "2026-03-09T15:00:00Z",
  "request_id": "uuid"
}
```

**Standard Error Codes**:
| Code | HTTP Status | Description | Retryable |
|------|-------------|-------------|-----------|
| UNAUTHORIZED | 401 | Invalid or expired token | No (re-auth required) |
| FORBIDDEN | 403 | Insufficient permissions | No |
| NOT_FOUND | 404 | Resource not found | No |
| VALIDATION_ERROR | 400 | Request validation failed | No |
| INVALID_STATUS_TRANSITION | 400 | Invalid task status change | No |
| RATE_LIMITED | 429 | Too many requests | Yes (with backoff) |
| NETWORK_ERROR | 503 | Service temporarily unavailable | Yes |
| INTERNAL_ERROR | 500 | Internal server error | Yes |

---

## Retry Strategy

Per FR-015, implement exponential backoff retry:

```typescript
const retryConfig = {
  retries: 3,
  retryDelay: (retryCount: number) => {
    const delays = [1000, 2000, 4000]; // 1s, 2s, 4s
    return delays[retryCount - 1] || 4000;
  },
  retryCondition: (error: AxiosError) => {
    // Retry on network errors and 5xx responses
    return !error.response || error.response.status >= 500;
  }
};
```

**Non-retryable errors** (fail immediately):
- 401 Unauthorized (redirect to login)
- 403 Forbidden (show permission error)
- 400 Validation Error (show validation message)
- 404 Not Found (show not found message)

---

## Frontend API Client Structure

```typescript
// api/client.ts
class ApiClient {
  private axiosInstance: AxiosInstance;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({ baseURL });
    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor - add JWT token
    this.axiosInstance.interceptors.request.use((config) => {
      const token = localStorage.getItem('kanban_auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Response interceptor - handle 401
    this.axiosInstance.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('kanban_auth_token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );

    // Setup retry
    axiosRetry(this.axiosInstance, retryConfig);
  }

  // Auth
  async login(credentials: LoginRequest): Promise<LoginResponse>;

  // Tasks
  async getTasks(filters?: TaskFilters): Promise<TaskListResponse>;
  async getTask(taskId: string): Promise<HarvestTask>;
  async updateTaskStatus(taskId: string, status: TaskStatus): Promise<HarvestTask>;
  async claimTask(taskId: string): Promise<HarvestTask>;
}
```
