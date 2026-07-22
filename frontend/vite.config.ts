import { readFileSync } from 'node:fs'
import { defineConfig, type Connect, type Plugin, type PreviewServer, type ViteDevServer } from 'vite'
import react from '@vitejs/plugin-react'

const pkg = JSON.parse(readFileSync(new URL('./package.json', import.meta.url), 'utf-8')) as {
  version: string
}

const FRONTEND_VERSION = {
  service: 'metartls-web',
  version: pkg.version,
}

const FRONTEND_CONFIG = {
  service: 'metartls-web',
  version: pkg.version,
  uiPort: 5173,
  apiBase: '',
  apiProxyTarget: 'http://localhost:8090',
}

function sendJSON(res: Connect.ServerResponse, body: unknown) {
  res.statusCode = 200
  res.setHeader('Content-Type', 'application/json; charset=utf-8')
  res.end(JSON.stringify({ data: body }))
}

function funcQueryPlugin(): Plugin {
  const handle: Connect.NextHandleFunction = (req, res, next) => {
    const host = req.headers.host ?? 'localhost'
    const url = new URL(req.url ?? '/', `http://${host}`)
    const func = (url.searchParams.get('func') ?? '').toLowerCase()

    if (func === 'getversion') {
      sendJSON(res, FRONTEND_VERSION)
      return
    }
    if (func === 'getconfig') {
      sendJSON(res, FRONTEND_CONFIG)
      return
    }
    if (func) {
      res.statusCode = 400
      res.setHeader('Content-Type', 'application/json; charset=utf-8')
      res.end(JSON.stringify({ error: 'unknown func' }))
      return
    }
    next()
  }

  return {
    name: 'func-query-api',
    configureServer(server: ViteDevServer) {
      server.middlewares.use(handle)
    },
    configurePreviewServer(server: PreviewServer) {
      server.middlewares.use(handle)
    },
  }
}

export default defineConfig({
  plugins: [react(), funcQueryPlugin()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8090',
        changeOrigin: true,
        ws: true,
      },
      '/health': 'http://localhost:8090',
      '/ready': 'http://localhost:8090',
    },
  },
})
