import { useEffect, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Alert, Box, Button, Chip, Paper, Stack, Typography } from '@mui/material'
import { api } from '../api/client'
import { useAuthStore } from '../store/auth'

type Floor = {
  id: string
  code: string
  name: string
  widthM: number
  heightM: number
}

type Zone = {
  id: string
  code: string
  name: string
  minX: number
  minY: number
  maxX: number
  maxY: number
}

type LivePosition = {
  tagId: string
  tagCode: string
  floorId: string
  x: number
  y: number
  zoneCode?: string
  updatedAt: string
}

const MAP_SCALE = 8

function wsURL(token: string) {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/api/v1/ws/locations?token=${encodeURIComponent(token)}`
}

export default function LiveMapPage() {
  const token = useAuthStore((s) => s.token)
  const queryClient = useQueryClient()
  const [positions, setPositions] = useState<Record<string, LivePosition>>({})
  const [wsState, setWsState] = useState<'connecting' | 'open' | 'closed'>('connecting')

  const floors = useQuery({
    queryKey: ['floors'],
    queryFn: () => api<Floor[]>('/api/v1/floors', {}, token),
  })

  const floor = floors.data?.[0]
  const zones = useQuery({
    queryKey: ['zones', floor?.id],
    enabled: !!floor?.id,
    queryFn: () => api<Zone[]>(`/api/v1/floors/${floor!.id}/zones`, {}, token),
  })

  const simStatus = useQuery({
    queryKey: ['simulator-status'],
    queryFn: () => api<{ running: boolean }>('/api/v1/simulator/status', {}, token),
    refetchInterval: 5000,
  })

  const startSim = useMutation({
    mutationFn: () => api('/api/v1/simulator/start', { method: 'POST' }, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['simulator-status'] }),
  })

  const stopSim = useMutation({
    mutationFn: () => api('/api/v1/simulator/stop', { method: 'POST' }, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['simulator-status'] }),
  })

  useEffect(() => {
    if (!token) return
    let alive = true
    api<LivePosition[]>('/api/v1/locations/latest', {}, token)
      .then((items) => {
        if (!alive) return
        const next: Record<string, LivePosition> = {}
        for (const item of items) next[item.tagId] = item
        setPositions(next)
      })
      .catch(() => {})
    return () => {
      alive = false
    }
  }, [token])

  useEffect(() => {
    if (!token) return
    const socket = new WebSocket(wsURL(token))
    setWsState('connecting')

    socket.onopen = () => setWsState('open')
    socket.onclose = () => setWsState('closed')
    socket.onerror = () => setWsState('closed')
    socket.onmessage = (msg) => {
      try {
        const pos = JSON.parse(msg.data) as LivePosition
        setPositions((prev) => ({ ...prev, [pos.tagId]: pos }))
      } catch {
      }
    }

    return () => socket.close()
  }, [token])

  const view = useMemo(() => {
    if (!floor) return null
    return {
      width: floor.widthM * MAP_SCALE,
      height: floor.heightM * MAP_SCALE,
      scale: MAP_SCALE,
    }
  }, [floor])

  const list = Object.values(positions)

  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Live map</Typography>
        <Typography color="text.secondary">
          Watch demo tags move. Data comes from MQTT simulator over WebSocket.
        </Typography>
      </div>

      <Stack direction="row" gap={1} flexWrap="wrap" alignItems="center">
        <Chip
          size="small"
          label={`WebSocket: ${wsState}`}
          color={wsState === 'open' ? 'success' : 'default'}
          variant="outlined"
        />
        <Chip
          size="small"
          label={`Simulator: ${simStatus.data?.running ? 'running' : 'stopped'}`}
          color={simStatus.data?.running ? 'success' : 'default'}
          variant="outlined"
        />
        <Button size="small" variant="contained" onClick={() => startSim.mutate()}>
          Start simulator
        </Button>
        <Button size="small" variant="outlined" onClick={() => stopSim.mutate()}>
          Stop
        </Button>
      </Stack>

      {wsState === 'closed' && (
        <Alert severity="warning">WebSocket is closed. Is the API running on port 8090?</Alert>
      )}

      {floor && view && (
        <Paper sx={{ p: 2, overflow: 'auto' }}>
          <Typography variant="subtitle1" sx={{ mb: 1 }}>
            {floor.name} ({floor.code}) · {list.length} tags
          </Typography>
          <Box
            sx={{
              width: view.width,
              height: view.height,
              position: 'relative',
              background:
                'repeating-linear-gradient(0deg, transparent, transparent 39px, #D8CFE6 40px), repeating-linear-gradient(90deg, transparent, transparent 39px, #D8CFE6 40px), #FCFAFD',
              border: '1px solid #1A1028',
            }}
          >
            {(zones.data ?? []).map((z) => (
              <Box
                key={z.id}
                title={`${z.code} — ${z.name}`}
                sx={{
                  position: 'absolute',
                  left: z.minX * view.scale,
                  top: z.minY * view.scale,
                  width: (z.maxX - z.minX) * view.scale,
                  height: (z.maxY - z.minY) * view.scale,
                  bgcolor: 'rgba(200, 75, 255, 0.16)',
                  border: '2px solid #1A1028',
                  display: 'flex',
                  alignItems: 'flex-start',
                  p: 0.5,
                }}
              >
                <Typography
                  variant="caption"
                  sx={{
                    fontWeight: 700,
                    color: '#1A1028',
                    fontFamily: '"JetBrains Mono", monospace',
                  }}
                >
                  {z.code}
                </Typography>
              </Box>
            ))}

            {list.map((pos) => (
              <Box
                key={pos.tagId}
                title={`${pos.tagCode} ${pos.zoneCode ? '(' + pos.zoneCode + ')' : ''}`}
                sx={{
                  position: 'absolute',
                  left: pos.x * view.scale - 8,
                  top: pos.y * view.scale - 8,
                  width: 16,
                  height: 16,
                  borderRadius: '50%',
                  bgcolor: '#C84BFF',
                  border: '2px solid #1A1028',
                  transition: 'left 0.8s linear, top 0.8s linear',
                  zIndex: 2,
                }}
              >
                <Typography
                  sx={{
                    position: 'absolute',
                    top: 16,
                    left: '50%',
                    transform: 'translateX(-50%)',
                    fontSize: 10,
                    fontFamily: '"JetBrains Mono", monospace',
                    whiteSpace: 'nowrap',
                    color: '#1A1028',
                    fontWeight: 700,
                  }}
                >
                  {pos.tagCode}
                </Typography>
              </Box>
            ))}
          </Box>
        </Paper>
      )}
    </Stack>
  )
}
