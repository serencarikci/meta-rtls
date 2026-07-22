# MetaRTLS Frontend

This is the React UI for MetaRTLS.
It talks to the Go API with REST and WebSocket.

## What you need

- Node.js 22+
- Backend API running on http://localhost:8090
- Oracle ready (`/ready` must work)

## How to start the frontend

First start Docker + backend (from the **repo root**):

```bash
cp config/config-temp.env config/config.env
make up
make backend-run
```

Wait until http://localhost:8090/ready is OK.

Then start the UI:

```bash
cd frontend
npm install
npm run dev
```

Or from the repo root:

```bash
make frontend-run
```

Open: http://localhost:5173

## Version and config URLs

- http://localhost:5173?func=getversion
- http://localhost:5173?func=getconfig

Backend (same style):
- http://localhost:8090?func=getversion
- http://localhost:8090?func=getconfig

In code use `Services` (`src/services/Services.ts`): `getVersion()` and `getConfig()`.

Current frontend version: `0.1.0` (`package.json`).

## Important

Vite sends these paths to the API:
- `/api`
- `/health`
- `/ready`

So the UI needs the backend on port `8090`.

If login fails:
1. Check that the API is running
2. Open http://localhost:8090/ready
3. Wait for Oracle if it is still starting

## Demo login

| Tenant | Email | Password |
|--------|-------|----------|
| warehouse-s | admin@warehouse-s.demo | MetaRTLS!2026 |
| hospital-m | admin@hospital-m.demo | MetaRTLS!2026 |
| factory-l | admin@factory-l.demo | MetaRTLS!2026 |

## Main screens

- Login
- Overview
- Sites & Zones
- Metadata
- Live Map (moving tags)
- Analysis

## Scripts

```bash
npm run dev       # start local UI
npm run build     # production build
npm run format    # Prettier
npm run preview   # preview build
```

## Stack

- React + TypeScript
- Vite
- Material UI
- TanStack Query
- Zustand
- React Router
- Prettier

## Folder layout

```text
src/
  api/           API helper
  layout/        top menu and page shell
  pages/         screens
  services/      version and config helpers
  store/         auth state
  styles/        global CSS
  types/         shared types
  theme.ts       MUI theme
```
