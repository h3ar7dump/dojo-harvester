# Research: Task Kanban Board Technical Decisions

**Feature**: Task Kanban Board for Data Platform
**Date**: 2026-03-09
**Phase**: Phase 0 - Research & Technical Decisions

## Research Areas

### 1. JWT Authentication Patterns

**Context**: The application must authenticate with the data platform using username/password and store JWT in localStorage.

**Decision**: Implement JWT storage in localStorage with automatic expiration handling and redirect to login

**Rationale**:
- localStorage is accessible across browser tabs
- Simple to implement with React Context
- Must implement XSS protection measures (CSP, input sanitization)
- Automatic token expiration detection via interceptors

**Implementation Pattern**:
```typescript
// Auth context with automatic expiration handling
interface AuthState {
  token: string | null;
  expiresAt: number | null;
  isAuthenticated: boolean;
}

// Axios interceptor for automatic redirect on 401
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear token and redirect to login
      authStore.clearToken();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

**Security Considerations**:
- Content Security Policy (CSP) headers to mitigate XSS
- Input sanitization on all user inputs
- Token validation before each API request
- Secure password entry (no autocomplete if needed)

---

### 2. Drag-and-Drop Libraries

**Context**: The kanban board requires drag-and-drop for moving tasks between columns.

**Decision**: Use @dnd-kit/core with @dnd-kit/sortable

**Rationale**:
- Modern, actively maintained library
- First-class TypeScript support
- Composable API (useDraggable, useDroppable hooks)
- Built-in accessibility support
- Works with React 19
- Better performance than alternatives

**Alternatives Considered**:
- **react-beautiful-dnd**: Rejected - less active maintenance, some React 18+ compatibility issues
- **react-dnd**: Rejected - more complex API, requires backend setup
- **HTML5 native drag-and-drop**: Rejected - inconsistent browser support, limited customization

**Implementation Pattern**:
```typescript
import { DndContext, useDraggable, useDroppable } from '@dnd-kit/core';
import { SortableContext, useSortable } from '@dnd-kit/sortable';

// Wrap board in DndContext
// Use SortableContext per column
// Implement useSortable for task cards
```

---

### 3. Skeleton Screen Implementation

**Context**: Need skeleton loading states for kanban columns and task cards per FR-016.

**Decision**: Use CSS-based skeleton components with TailwindCSS

**Rationale**:
- Pure CSS solution, no additional dependencies
- Tailwind animate-pulse utility for shimmer effect
- Composable skeleton atoms (SkeletonText, SkeletonCard)
- Matches shadcn/ui design system

**Implementation Pattern**:
```typescript
// Skeleton components
const SkeletonCard = () => (
  <div className="p-4 bg-white rounded-lg shadow animate-pulse">
    <div className="h-4 bg-gray-200 rounded w-3/4 mb-2" />
    <div className="h-3 bg-gray-200 rounded w-1/2" />
  </div>
);

// Usage during loading state
{isLoading ? <SkeletonCard /> : <TaskCard task={task} />}
```

---

### 4. Axios Retry Logic

**Context**: Need exponential backoff retry (1s, 2s, 4s) for API requests per FR-015.

**Decision**: Use axios-retry library with custom retry condition

**Rationale**:
- Battle-tested library with exponential backoff support
- Configurable retry delays and conditions
- Works seamlessly with axios interceptors
- TypeScript definitions available

**Implementation Pattern**:
```typescript
import axios from 'axios';
import axiosRetry from 'axios-retry';

axiosRetry(axios, {
  retries: 3,
  retryDelay: (retryCount) => {
    const delays = [1000, 2000, 4000]; // 1s, 2s, 4s
    return delays[retryCount - 1] || 4000;
  },
  retryCondition: (error) => {
    return axiosRetry.isNetworkOrIdempotentRequestError(error);
  },
});
```

**Alternatives Considered**:
- **Custom implementation**: Rejected - more code to maintain, edge cases to handle
- **fetch with retry**: Rejected - axios provides better error handling and interceptors

---

### 5. Zustand State Management

**Context**: Need global state for auth, tasks, filters, and kanban board state.

**Decision**: Use Zustand with persistence middleware

**Rationale**:
- Minimal boilerplate compared to Redux
- Built-in TypeScript support
- DevTools middleware for debugging
- Persistence middleware for filter state
- Small bundle size (~1KB)

**State Structure**:
```typescript
// Stores
interface AuthStore {
  token: string | null;
  setToken: (token: string) => void;
  clearToken: () => void;
}

interface KanbanStore {
  tasks: HarvestTask[];
  columns: KanbanColumn[];
  filters: TaskFilter;
  setTasks: (tasks: HarvestTask[]) => void;
  moveTask: (taskId: string, newStatus: TaskStatus) => void;
  updateFilters: (filters: TaskFilter) => void;
}
```

---

## Integration Patterns

### Data Platform API Integration
- **Protocol**: HTTP REST
- **Data Format**: JSON
- **Auth**: JWT token in Authorization header
- **Base URL**: Configurable via environment variable

### Authentication Flow
```
Login Form → POST /auth/login → JWT Response → Store in localStorage → Redirect to Kanban
```

### Task Data Flow
```
Kanban Mount → GET /tasks (with JWT) → Store in Zustand → Render Board
Drag Task → PUT /tasks/{id}/status → Update Platform → Update Local State
```

---

## Performance Budgets

| Metric | Target | Justification |
|--------|--------|---------------|
| Login Complete | <10s | SC-001 requirement |
| Board Load | <3s | SC-002 requirement |
| Filter Response | <1s | SC-004 requirement |
| Drag-and-Drop Sync | <2s | SC-003 requirement |
| API First Attempt Success | 99% | SC-005 requirement |
| Task Detail View | <1s | SC-006 requirement |
| Concurrent Users | 1 per session | Single operator per browser |
| Max Tasks Displayed | 100 | Per operator limit |

---

## Risk Mitigation

| Risk | Mitigation Strategy |
|------|-------------------|
| XSS via localStorage | Implement CSP headers, input sanitization |
| API rate limiting | Implement request queuing and caching |
| Token expiration mid-session | Automatic redirect with toast notification |
| Large task lists | Virtualization for >50 tasks per column |
| Drag-and-drop performance | Debounce API calls, optimistic UI updates |
| Network failures | Exponential backoff, offline indicators |

---

## Dependencies Summary

**Core**:
- react@19
- react-dom@19
- typescript@5.x
- vite@5.x

**State Management**:
- zustand

**HTTP Client**:
- axios
- axios-retry

**Drag and Drop**:
- @dnd-kit/core
- @dnd-kit/sortable
- @dnd-kit/utilities

**UI Components**:
- @radix-ui/react-dialog
- @radix-ui/react-select
- @radix-ui/react-toast
- tailwindcss
- shadcn/ui (via CLI)

**Utilities**:
- date-fns (date formatting)
- clsx (conditional classes)
- zod (validation)

**Testing**:
- vitest
- @testing-library/react
- @testing-library/jest-dom
