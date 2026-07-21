import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Alert, Box, Paper, Stack, Typography } from '@mui/material'
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

const MAP_SCALE = 8

export default function LiveMapPage() {
  const token = useAuthStore((s) => s.token)
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

  const view = useMemo(() => {
    if (!floor) return null
    return {
      width: floor.widthM * MAP_SCALE,
      height: floor.heightM * MAP_SCALE,
      scale: MAP_SCALE,
    }
  }, [floor])

  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Live map</Typography>
        <Typography color="text.secondary">
          Phase 1: SVG floor plan + zone boxes. Phase 3: WebSocket location points.
        </Typography>
      </div>

      <Alert severity="info">MQTT simulator and live tag positions will connect in Phase 3.</Alert>

      {floor && view && (
        <Paper sx={{ p: 2, border: '1px solid #D5DEE8', overflow: 'auto' }}>
          <Typography variant="subtitle1" sx={{ mb: 1 }}>
            {floor.name} ({floor.code})
          </Typography>
          <Box
            sx={{
              width: view.width,
              height: view.height,
              position: 'relative',
              background:
                'repeating-linear-gradient(0deg, transparent, transparent 39px, #D5DEE8 40px), repeating-linear-gradient(90deg, transparent, transparent 39px, #D5DEE8 40px), #F8FAFC',
              border: '1px solid #9AAEBF',
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
                  bgcolor: 'rgba(224, 159, 62, 0.18)',
                  border: '2px solid #E09F3E',
                  display: 'flex',
                  alignItems: 'flex-start',
                  p: 0.5,
                }}
              >
                <Typography variant="caption" sx={{ fontWeight: 700, color: '#16324F' }}>
                  {z.code}
                </Typography>
              </Box>
            ))}
          </Box>
        </Paper>
      )}
    </Stack>
  )
}
