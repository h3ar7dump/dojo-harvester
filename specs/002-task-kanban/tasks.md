# Tasks: Task Kanban Board for Data Platform

**Feature**: Task Kanban Board for Data Platform
**Branch**: `002-task-kanban`
**Date**: 2026-03-09
**Plan**: [plan.md](./plan.md) | **Spec**: [spec.md](./spec.md)

---

## Phase 1: Setup & Project Initialization

> **Goal**: Initialize React project with all dependencies and tooling
> **Prerequisites**: None (first phase)
> **Independent Test**: Project builds successfully, dev server starts

### Project Initialization

- [ ] T001 Create frontend directory structure per implementation plan
- [ ] T002 Initialize Vite + React + TypeScript project
- [ ] T003 Configure TypeScript compiler options (strict mode, path aliases)
- [ ] T004 Configure Vite with proxy for API requests

### Dependency Installation

- [ ] T005 [P] Install React 19 and React DOM
- [ ] T006 [P] Install TypeScript and type definitions
- [ ] T007 [P] Install axios and axios-retry for HTTP client with retry logic
- [ ] T008 [P] Install Zustand for state management
- [ ] T009 [P] Install @dnd-kit/core, @dnd-kit/sortable, @dnd-kit/utilities for drag-and-drop
- [ ] T010 [P] Install TailwindCSS, PostCSS, Autoprefixer
- [ ] T011 [P] Install Radix UI primitives (@radix-ui/react-dialog, @radix-ui/react-select, @radix-ui/react-toast)
- [ ] T012 [P] Install utility libraries (date-fns, clsx, zod)
- [ ] T013 [P] Install icon library (lucide-react)

### Configuration Setup

- [ ] T014 Configure TailwindCSS with custom theme and animations
- [ ] T015 [P] Configure ESLint with TypeScript and React rules
- [ ] T016 [P] Configure Prettier for code formatting
- [ ] T017 [P] Configure Vitest for unit testing
- [ ] T018 Create environment variable templates (.env.development, .env.production)

---

## Phase 2: Foundational Components (Blocking Prerequisites)

> **Goal**: Build foundational components required by all user stories
> **Prerequisites**: Phase 1 complete
> **Independent Test**: API client connects, stores initialize, routing works

### Type Definitions & Types

- [ ] T019 Create TypeScript types for UserSession in src/types/index.ts
- [ ] T020 Create TypeScript types for HarvestTask in src/types/index.ts
- [ ] T021 Create TypeScript types for KanbanColumn in src/types/index.ts
- [ ] T022 Create TypeScript types for TaskFilter in src/types/index.ts
- [ ] T023 Create TypeScript types for ApiError in src/types/index.ts

### API Client Setup

- [ ] T024 Create axios instance with base URL configuration in src/services/api.ts
- [ ] T025 Implement JWT token interceptor for Authorization header
- [ ] T026 Implement 401 error interceptor with redirect to login
- [ ] T027 Configure axios-retry with exponential backoff (1s, 2s, 4s delays)
- [ ] T028 Implement API client methods (login, getTasks, updateTaskStatus, getTask)

### State Management (Zustand)

- [ ] T029 Create authStore with login/logout actions in src/stores/authStore.ts
- [ ] T030 Create kanbanStore with task state and actions in src/stores/kanbanStore.ts
- [ ] T031 Create filterStore with filter state in src/stores/filterStore.ts
- [ ] T032 Implement localStorage persistence for JWT token
- [ ] T033 Implement localStorage persistence for filter preferences

### UI Component Foundation

- [ ] T034 [P] Create Button component in src/components/ui/Button.tsx
- [ ] T035 [P] Create Input component in src/components/ui/Input.tsx
- [ ] T036 [P] Create Card component in src/components/ui/Card.tsx
- [ ] T037 [P] Create Modal/Dialog component in src/components/ui/Modal.tsx
- [ ] T038 [P] Create Toast/Notification component in src/components/ui/Toast.tsx
- [ ] T039 [P] Create Skeleton component in src/components/ui/Skeleton.tsx
- [ ] T040 [P] Create LoadingSpinner component in src/components/ui/LoadingSpinner.tsx

### Routing & Navigation

- [ ] T041 Implement React Router setup in src/App.tsx
- [ ] T042 Create ProtectedRoute component for authenticated routes
- [ ] T043 Create PublicRoute component for login page
- [ ] T044 Implement automatic redirect based on auth status

---

## Phase 3: User Story 1 - Login and JWT Authentication (P1)

> **Goal**: Implement login form and JWT authentication flow
> **Prerequisites**: Phase 2 complete
> **Independent Test**: User can log in with valid/invalid credentials, JWT stored in localStorage

### Login UI Components

- [ ] T045 [US1] [P] Create LoginPage component in src/components/auth/LoginPage.tsx
- [ ] T046 [US1] [P] Create LoginForm component with username/password fields
- [ ] T047 [US1] [P] Create form validation using zod schema
- [ ] T048 [US1] Implement error message display for invalid credentials
- [ ] T049 [US1] Implement loading state during login submission
- [ ] T050 [US1] [P] Style login page with TailwindCSS

### Authentication Logic

- [ ] T051 [US1] Implement login action in authStore with API integration
- [ ] T052 [US1] Store JWT token in localStorage on successful login
- [ ] T053 [US1] Store user info in localStorage on successful login
- [ ] T054 [US1] Implement redirect to kanban board after login
- [ ] T055 [US1] Implement JWT expiration checking logic
- [ ] T056 [US1] Implement automatic redirect to login when JWT expires
- [ ] T057 [US1] Implement logout functionality
- [ ] T058 [US1] Implement session persistence on page refresh

---

## Phase 4: User Story 2 - Kanban Board Task Visualization (P1)

> **Goal**: Display tasks in kanban board layout with columns and cards
> **Prerequisites**: Phase 3 (US1) complete - needs authentication
> **Independent Test**: Authenticated user sees tasks in correct columns with task details

### Kanban Board Components

- [ ] T059 [US2] Create KanbanBoardPage component in src/components/kanban/KanbanBoardPage.tsx
- [ ] T060 [US2] [P] Create KanbanColumn component in src/components/kanban/KanbanColumn.tsx
- [ ] T061 [US2] [P] Create TaskCard component in src/components/kanban/TaskCard.tsx
- [ ] T062 [US2] [P] Create TaskCardContent component showing title, priority, due date
- [ ] T063 [US2] [P] Create ProgressIndicator component for episodes (completed/total)
- [ ] T064 [US2] [P] Create PriorityBadge component for priority levels
- [ ] T065 [US2] Implement column layout (To Do, In Progress, Done)
- [ ] T066 [US2] Implement task distribution logic (group by status)

### Data Fetching & Loading States

- [ ] T067 [US2] Implement fetchTasks action in kanbanStore with API call
- [ ] T068 [US2] Implement skeleton screen loading for columns (FR-016)
- [ ] T069 [US2] Implement skeleton screen loading for task cards (FR-016)
- [ ] T070 [US2] Implement error handling for failed task fetch
- [ ] T071 [US2] Implement empty state display when no tasks assigned (FR-017)
- [ ] T072 [US2] Implement automatic task fetching on kanban mount

### Task Detail View

- [ ] T073 [US2] [P] Create TaskDetailModal component in src/components/task/TaskDetailModal.tsx
- [ ] T074 [US2] [P] Create TaskDetailContent component showing robot configuration
- [ ] T075 [US2] [P] Create TaskObjectivesList component
- [ ] T076 [US2] Implement click handler on task cards to open detail view
- [ ] T077 [US2] Implement task detail API fetch
- [ ] T078 [US2] [P] Style task detail modal

---

## Phase 5: User Story 3 - Task Status Management via Drag-and-Drop (P2)

> **Goal**: Enable drag-and-drop to update task status
> **Prerequisites**: Phase 4 (US2) complete - needs kanban board with tasks
> **Independent Test**: User can drag task between columns, status updates persist

### Drag-and-Drop Infrastructure

- [ ] T079 [US3] Install and configure @dnd-kit/core DndContext
- [ ] T080 [US3] Configure DndContext sensors (pointer, keyboard)
- [ ] T081 [US3] Implement useDraggable hook for TaskCard
- [ ] T082 [US3] Implement useDroppable hook for KanbanColumn
- [ ] T083 [US3] Implement useSortable for tasks within columns

### Drag-and-Drop UI Components

- [ ] T084 [US3] [P] Create DraggableTaskCard wrapper component
- [ ] T085 [US3] [P] Create DroppableColumn wrapper component
- [ ] T086 [US3] Implement drag overlay for visual feedback during drag
- [ ] T087 [US3] Implement drop indicator (highlight column on hover)
- [ ] T088 [US3] Implement drag handle on task cards

### Status Update Logic

- [ ] T089 [US3] Implement onDragEnd handler in KanbanBoardPage
- [ ] T090 [US3] Implement optimistic UI update (update local state immediately)
- [ ] T091 [US3] Implement moveTask action in kanbanStore with API call
- [ ] T092 [US3] Implement status validation (pending → in_progress → completed)
- [ ] T093 [US3] Implement revert on API failure with error toast
- [ ] T094 [US3] Implement success toast on status update completion
- [ ] T095 [US3] Handle network interruption during drag with retry logic

---

## Phase 6: User Story 4 - Task Filtering and Sorting (P3)

> **Goal**: Add filter and sort controls for task management
> **Prerequisites**: Phase 4 (US2) complete - needs tasks to filter
> **Independent Test**: User can filter by priority, sort by due date, see filtered results

### Filter UI Components

- [ ] T096 [US4] [P] Create FilterPanel component in src/components/filters/FilterPanel.tsx
- [ ] T097 [US4] [P] Create PriorityFilter dropdown component
- [ ] T098 [US4] [P] Create DueDateRangeFilter component with date pickers
- [ ] T099 [US4] [P] Create RobotTypeFilter dropdown component
- [ ] T100 [US4] [P] Create SortControl component (field + direction)
- [ ] T101 [US4] [P] Create ClearFilters button component
- [ ] T102 [US4] [P] Style filter panel with TailwindCSS

### Filter Logic

- [ ] T103 [US4] Implement filterTasks function with priority filter
- [ ] T104 [US4] Implement filterTasks function with due date range filter
- [ ] T105 [US4] Implement filterTasks function with robot type filter
- [ ] T106 [US4] Implement sortTasks function (due date, priority, created_at)
- [ ] T107 [US4] Implement combined filter and sort in kanbanStore
- [ ] T108 [US4] Implement filter persistence to localStorage
- [ ] T109 [US4] Implement filter reset functionality
- [ ] T110 [US4] Ensure filters apply across all columns
- [ ] T111 [US4] Ensure filtered tasks maintain drag-and-drop functionality

---

## Phase 7: Polish & Cross-Cutting Concerns

> **Goal**: Final polish, error handling, accessibility, and performance optimization
> **Prerequisites**: All user stories (US1-US4) complete
> **Independent Test**: Full end-to-end workflow executes with graceful error handling

### Error Handling & Edge Cases

- [ ] T112 Implement global error boundary component
- [ ] T113 Handle API unreachable during login with user-friendly error
- [ ] T114 Handle JWT expiration during active session with auto-redirect
- [ ] T115 Handle task deleted from platform while being viewed
- [ ] T116 Handle malformed task data from API with validation
- [ ] T117 Implement network connectivity detection
- [ ] T118 Implement offline indicator component

### Accessibility (a11y)

- [ ] T119 Add ARIA labels to all interactive elements
- [ ] T120 Implement keyboard navigation for kanban board
- [ ] T121 Implement keyboard shortcuts for common actions
- [ ] T122 Ensure focus management in modals
- [ ] T123 Implement screen reader announcements for status updates
- [ ] T124 Add focus visible styles for keyboard navigation

### Performance Optimization

- [ ] T125 Implement React.memo for TaskCard to prevent unnecessary re-renders
- [ ] T126 Implement useMemo for filtered/sorted task lists
- [ ] T127 Implement useCallback for event handlers
- [ ] T128 Add virtualization for task lists >50 tasks (react-window)
- [ ] T129 Implement debounce for filter changes
- [ ] T130 Optimize skeleton screen rendering

### Documentation & Build

- [ ] T131 [P] Create comprehensive README with setup instructions
- [ ] T132 [P] Write API integration documentation
- [ ] T133 Configure production build with environment variables
- [ ] T134 Configure Docker containerization
- [ ] T135 Add nginx configuration for SPA routing

---

## Dependency Graph

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational)
    ├── Type Definitions (T019-T023)
    ├── API Client (T024-T028)
    ├── State Management (T029-T033)
    ├── UI Foundation (T034-T040)
    └── Routing (T041-T044)
    ↓
Phase 3 (US1: Authentication)
    ├── Login UI (T045-T050)
    └── Auth Logic (T051-T058)
    ↓
Phase 4 (US2: Kanban Visualization)
    ├── Kanban Components (T059-T066)
    ├── Data Fetching (T067-T072)
    └── Task Detail (T073-T078)
    ↓
Phase 5 (US3: Drag-and-Drop)
    ├── DnD Infrastructure (T079-T083)
    ├── DnD UI (T084-T088)
    └── Status Update Logic (T089-T095)
    ↓
Phase 6 (US4: Filtering)
    ├── Filter UI (T096-T102)
    └── Filter Logic (T103-T111)
    ↓
Phase 7 (Polish)
    ├── Error Handling (T112-T118)
    ├── Accessibility (T119-T124)
    ├── Performance (T125-T130)
    └── Documentation (T131-T135)
```

### Parallel Execution Opportunities

**Within Phase 1:**
- T001-T004 (project setup) - sequential
- T005-T013 (dependency installation) - all parallel
- T014-T018 (configuration) - parallel after dependencies

**Within Phase 2:**
- T019-T023 (types) - sequential
- T024-T028 (API client) - sequential
- T029-T033 (stores) - sequential, parallel with API client
- T034-T040 (UI components) - all parallel
- T041-T044 (routing) - sequential

**User Story Development:**
- US3 (Drag-and-Drop) can start after US2 kanban components are complete
- US4 (Filtering) can be developed in parallel with US3 after US2

**MVP Scope Recommendation:**
Focus on Phase 1-4 (US1 + US2) first - this delivers core authentication and kanban visualization. Drag-and-drop and filtering can be added incrementally.

---

## Summary

| Phase | Tasks | Story | Independent Test Criteria |
|-------|-------|-------|---------------------------|
| Phase 1 | 18 | Setup | Project builds, dev server starts |
| Phase 2 | 26 | Foundational | API client connects, stores initialize |
| Phase 3 | 14 | US1 (P1) | Login works, JWT stored, redirect functional |
| Phase 4 | 20 | US2 (P1) | Kanban displays tasks in correct columns |
| Phase 5 | 17 | US3 (P2) | Drag-and-drop updates status and persists |
| Phase 6 | 16 | US4 (P3) | Filters and sorting work correctly |
| Phase 7 | 24 | Polish | Errors handled gracefully, a11y compliance |

**Total Tasks**: 135
**MVP Tasks** (US1 + US2): 78 tasks
**Parallelizable Tasks**: 54 tasks (marked with [P])

---

## Next Steps

1. **Start with MVP**: Implement Phase 1-4 (US1 + US2) for core authentication and kanban functionality
2. **Parallel Development**: US3 and US4 can be developed incrementally after US2
3. **Incremental Delivery**: Each user story delivers independently testable functionality
4. **Run `/speckit.execute`**: To begin executing tasks from Phase 1
