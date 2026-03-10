# Implementation Plan: Task Kanban Board for Data Platform

**Branch**: `002-task-kanban` | **Date**: 2026-03-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-task-kanban/spec.md`

## Summary

Build a React-based kanban board interface for robot operators to view and manage their assigned data harvest tasks. The application authenticates with a remote data platform using JWT tokens, displays tasks in a kanban layout (To Do, In Progress, Done), supports drag-and-drop status updates, and provides filtering/sorting capabilities. This is a frontend-only feature per constitutional mandate - task management communicates directly with the data platform via HTTP REST, not through the Agent Backend.

## Technical Context

**Language/Version**: React 19 + TypeScript 5.x
**Primary Dependencies**:
- UI Framework: React 19, TypeScript 5.x, Vite
- Drag-and-Drop: @dnd-kit/core + @dnd-kit/sortable
- HTTP Client: axios with retry logic
- State Management: Zustand
- Styling: TailwindCSS, shadcn/ui, Radix UI
- Testing: Vitest, React Testing Library
**Storage**: Browser localStorage (for JWT token), in-memory state (Zustand)
**Testing**: Vitest for unit tests, React Testing Library for component tests
**Target Platform**: Web browsers (Chrome, Firefox, Safari, Edge)
**Project Type**: Frontend web application
**Performance Goals**: Login <10s, board load <3s, filter response <1s, drag-and-drop sync <2s
**Constraints**: JWT in localStorage (XSS vulnerable), exponential backoff retry (3 retries)
**Scale/Scope**: Single user per session, handles up to 100 tasks per operator

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. Serialization Protocol** | ✅ PASS | JSON for HTTP/REST API communication with data platform |
| **II. Task Management Isolation** | ✅ PASS | Kanban UI communicates directly via HTTP REST to Data Platform (NOT through Agent Backend) |
| **III. Execution Abstraction** | N/A | Frontend-only feature, no shell script execution |
| **IV. 3D Asset Delivery** | N/A | No 3D assets in this feature |
| **React Stack Mandates** | ✅ PASS | Will use React 19, TypeScript 5.x, Zustand, shadcn/ui, Radix UI, TailwindCSS |
| **Upload Resilience** | ✅ PASS | Exponential backoff retry per FR-015 |
| **Real-time Feedback** | ✅ PASS | Skeleton screens per FR-016, drag-and-drop feedback |

**Gate Result**: ✅ ALL CHECKS PASSED - Proceed to Phase 0

## Project Structure

### Documentation (this feature)

```text
specs/002-task-kanban/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
frontend/                # React Kanban Application
├── src/
│   ├── components/     # React components
│   │   ├── auth/       # Login form, auth context
│   │   ├── kanban/     # Kanban board, columns, cards
│   │   ├── task/       # Task detail view
│   │   ├── filters/    # Filter and sort controls
│   │   └── ui/         # Reusable UI components
│   ├── hooks/          # Custom React hooks
│   ├── services/       # API client services
│   ├── stores/         # Zustand state stores
│   ├── types/          # TypeScript type definitions
│   └── utils/          # Utility functions
├── public/
└── tests/
    ├── unit/
    └── integration/
```

**Structure Decision**: Single frontend project structure. This feature is purely a React web application that communicates directly with the data platform API. No Go backend or Python toolkit needed as this follows the Task Management Isolation principle from the constitution.

## Complexity Tracking

> No constitutional violations requiring justification.

## Phase 0: Research & Technical Decisions

*See [research.md](./research.md) for detailed findings*

### Research Areas

1. **JWT Authentication Patterns**: Best practices for storing JWT in localStorage, handling token expiration, and automatic redirection
2. **Drag-and-Drop Libraries**: Comparison of @dnd-kit vs react-beautiful-dnd for kanban boards
3. **Skeleton Screen Implementation**: React skeleton component patterns for loading states
4. **Axios Retry Logic**: Implementing exponential backoff with axios-retry or custom interceptor
5. **Zustand Best Practices**: State structure for kanban data with filtering and sorting

### Key Decisions

| Decision | Rationale | Alternatives Rejected |
|----------|-----------|----------------------|
| @dnd-kit for drag-and-drop | Modern, accessible, TypeScript-native, composable API | react-beautiful-dnd - less active maintenance |
| axios-retry for exponential backoff | Built-in retry mechanism, configurable delays | Custom implementation - more code to maintain |
| Zustand for state management | Lightweight, minimal boilerplate, devtools support | Redux - too complex for single-feature scope |
| localStorage for JWT | Spec requirement per clarification | httpOnly cookies - spec chose localStorage |

## Phase 1: Design & Contracts

### Data Model

*See [data-model.md](./data-model.md) for entity definitions*

### API Contracts

*See [contracts/](./contracts/) for interface definitions*

### Quick Start

*See [quickstart.md](./quickstart.md) for development setup*

## Post-Design Constitution Re-check

*Phase 1 design complete - all principles verified*

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. Serialization Protocol** | ✅ PASS | JSON for HTTP/REST API communication per contracts/api-contracts.md |
| **II. Task Management Isolation** | ✅ PASS | Kanban UI communicates directly via HTTP REST to Data Platform, NOT through Agent Backend (per spec FR-002, FR-006) |
| **React Stack Mandates** | ✅ PASS | Will use React 19, TypeScript 5.x, Zustand, shadcn/ui, Radix UI, TailwindCSS per quickstart.md |

**Final Gate Result**: ✅ ALL CHECKS PASSED - Ready for task generation
