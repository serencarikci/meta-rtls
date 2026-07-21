# MetaRTLS Frontend

React UI for MetaRTLS. Talks to the Go API with REST and WebSocket.

## Stack

- React + TypeScript
- Vite
- Material UI
- TanStack Query
- Zustand
- React Router
- Prettier

## Run

From the repo root:

```bash
make up
make backend-run
cd frontend
npm install
npm run dev
```

UI: http://localhost:5173

Vite proxies `/api` and `/health` to `http://localhost:8090`.

## Main screens

- Login (demo tenants)
- Overview
- Sites & Zones
- Metadata
- Live Map (moving tags)
- Analysis (compare tenants + impact score)

## Scripts

```bash
npm run dev       # local UI
npm run build     # production build
npm run format    # Prettier
npm run preview   # preview build
```

## Folder layout

```text
src/
  api/           fetch helper
  layout/        app shell and nav
  pages/         screens
  store/         auth state
  styles/        global CSS
  theme.ts       MUI theme
```

## Demo login

| Tenant | Email | Password |
|--------|-------|----------|
| warehouse-s | admin@warehouse-s.demo | MetaRTLS!2026 |
| hospital-m | admin@hospital-m.demo | MetaRTLS!2026 |
| factory-l | admin@factory-l.demo | MetaRTLS!2026 |

If login fails, check that the backend is running on port 8090.
