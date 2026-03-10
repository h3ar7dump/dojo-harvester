# Quick Start: Task Kanban Board Development

**Feature**: Task Kanban Board for Data Platform
**Date**: 2026-03-09

## Prerequisites

### Required Tools
- **Node.js 20+** - Frontend runtime
- **npm 10+** or **yarn 4+** or **pnpm 8+** - Package manager
- **Git** - Version control

### Optional Tools
- **VS Code** - Recommended IDE
- **ESLint/Prettier extensions** - Code formatting
- **React Developer Tools** - Browser extension for debugging

---

## Repository Structure

```
dojo-harvester/
├── frontend/             # React Kanban Application
│   ├── src/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/
│   │   ├── stores/
│   │   └── types/
│   ├── public/
│   └── tests/
├── specs/002-task-kanban/
│   ├── spec.md
│   ├── plan.md
│   ├── research.md
│   ├── data-model.md
│   ├── contracts/
│   └── quickstart.md      # This file
└── .specify/
    └── memory/
        └── constitution.md
```

---

## Frontend Setup (React + Vite)

### 1. Initialize Project

```bash
# Using Vite with React and TypeScript
cd frontend
npm create vite@latest . -- --template react-ts
```

### 2. Install Dependencies

```bash
# Core
npm install react@19 react-dom@19

# HTTP Client
npm install axios axios-retry

# State Management
npm install zustand

# Drag and Drop
npm install @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities

# UI Components (Radix UI + shadcn)
npm install @radix-ui/react-dialog @radix-ui/react-select @radix-ui/react-toast

# Styling
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# Utilities
npm install date-fns clsx zod

# Icons (optional)
npm install lucide-react
```

### 3. Configure TailwindCSS

```javascript
// tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-in': 'slideIn 0.3s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideIn: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
```

```css
/* src/index.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer components {
  .skeleton {
    @apply animate-pulse bg-gray-200 rounded;
  }
}
```

### 4. Configure TypeScript

```json
// tsconfig.json (update compiler options)
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["src/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

### 5. Configure Vite

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
});
```

---

## Development Workflow

### 1. Start Development Server

```bash
cd frontend
npm run dev
```

**Default**: `http://localhost:5173`

### 2. Environment Configuration

```bash
# .env.development
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Task Kanban
```

```bash
# .env.production
VITE_API_URL=https://platform.example.com/api/v1
VITE_APP_NAME=Task Kanban
```

### 3. Test API Integration

```bash
# Test login endpoint
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "test"}'

# Test tasks endpoint (with token)
curl http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer {token}"
```

---

## Project Structure Setup

### Create Directory Structure

```bash
cd frontend/src

# Create directories
mkdir -p components/auth
mkdir -p components/kanban
mkdir -p components/task
mkdir -p components/filters
mkdir -p components/ui
mkdir -p hooks
mkdir -p services
mkdir -p stores
mkdir -p types
mkdir -p utils
```

### Key Files Template

```typescript
// src/types/index.ts
export interface HarvestTask {
  task_id: string;
  title: string;
  description: string;
  status: 'pending' | 'in_progress' | 'completed';
  priority: 'low' | 'medium' | 'high' | 'critical';
  assignee_id: string;
  due_date: string;
  required_episodes: number;
  completed_episodes: number;
  robot_configuration: Record<string, unknown>;
  objectives: string[];
  created_at: string;
  updated_at: string;
}

export interface UserSession {
  user_id: string;
  username: string;
  token: string;
  expires_at: number;
}
```

```typescript
// src/stores/authStore.ts
import { create } from 'zustand';
import { UserSession } from '@/types';

interface AuthState {
  user: UserSession | null;
  isLoading: boolean;
  error: string | null;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: () => boolean;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  isLoading: false,
  error: null,

  login: async (username, password) => {
    set({ isLoading: true, error: null });
    try {
      // API call implementation
      const response = await apiClient.login({ username, password });
      localStorage.setItem('kanban_auth_token', response.access_token);
      set({ user: response.user, isLoading: false });
    } catch (error) {
      set({ error: 'Login failed', isLoading: false });
    }
  },

  logout: () => {
    localStorage.removeItem('kanban_auth_token');
    set({ user: null, error: null });
  },

  isAuthenticated: () => {
    const token = localStorage.getItem('kanban_auth_token');
    return !!token && !!get().user;
  },
}));
```

```typescript
// src/services/api.ts
import axios from 'axios';
import axiosRetry from 'axios-retry';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

// Add JWT token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('kanban_auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle 401 errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('kanban_auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Configure retry
axiosRetry(apiClient, {
  retries: 3,
  retryDelay: (retryCount) => {
    const delays = [1000, 2000, 4000];
    return delays[retryCount - 1] || 4000;
  },
  retryCondition: (error) => {
    return !error.response || error.response.status >= 500;
  },
});

export default apiClient;
```

---

## Testing

### Install Test Dependencies

```bash
npm install -D vitest @testing-library/react @testing-library/jest-dom jsdom
```

### Configure Vitest

```typescript
// vite.config.ts (update)
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
});
```

### Run Tests

```bash
# Run all tests
npm test

# Run with watch mode
npm test -- --watch

# Run with coverage
npm test -- --coverage
```

---

## Build and Deploy

### Production Build

```bash
# Create production build
npm run build

# Preview production build locally
npm run preview
```

### Build Output

The build output will be in `frontend/dist/` directory, ready for deployment to:
- Static hosting (Netlify, Vercel, Cloudflare Pages)
- Container deployment (Docker + Nginx)
- CDN (AWS CloudFront, CloudFlare)

### Docker Deployment

```dockerfile
# Dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

```nginx
# nginx.conf
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api {
        proxy_pass http://platform-api:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }
}
```

---

## Architecture Overview

```
┌─────────────────────────────────────┐
│         React Kanban App            │
├─────────────────────────────────────┤
│  ┌─────────┐ ┌───────────────────┐ │
│  │  Login  │ │   Kanban Board    │ │
│  │  Page   │ │ ┌─────┬─────┬────┐│ │
│  └────┬────┘ │ │ To  │ In  │Done││ │
│       │      │ │ Do  │Prog │    ││ │
│       ▼      │ └─────┴─────┴────┘│ │
│  localStorage│        │          │ │
│       │      │        ▼          │ │
│       ▼      │  Drag & Drop      │ │
│  ┌─────────┐ │  (@dnd-kit)       │ │
│  │  Auth   │ │        │          │ │
│  │ Context │ └────────┼──────────┘ │
│  └────┬────┘          │            │
│       │               ▼            │
│       │         ┌─────────┐        │
│       │         │ Zustand │        │
│       │         │ Stores  │        │
│       │         └────┬────┘        │
│       │              │             │
│       └──────────────┘             │
│                      │             │
└──────────────────────┼─────────────┘
                       │
                       ▼
              ┌─────────────────┐
              │   Axios Client  │
              │ (axios-retry)   │
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  Data Platform  │
              │     API         │
              └─────────────────┘
```

---

## Next Steps

1. **Review the spec**: `specs/002-task-kanban/spec.md`
2. **Check the plan**: `specs/002-task-kanban/plan.md`
3. **Generate tasks**: Run `/speckit.tasks` to create implementation tasks
4. **Start coding**: Follow the task breakdown in `tasks.md`
