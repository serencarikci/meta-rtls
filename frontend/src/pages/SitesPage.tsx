import { useQuery } from '@tanstack/react-query'
import {
  Chip,
  CircularProgress,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material'
import { api } from '../api/client'
import { useAuthStore } from '../store/auth'
import type { Building, Floor, Site } from '../types/rtls'

export default function SitesPage() {
  const token = useAuthStore((s) => s.token)
  const sites = useQuery({
    queryKey: ['sites'],
    queryFn: () => api<Site[]>('/api/v1/sites', {}, token),
  })
  const buildings = useQuery({
    queryKey: ['buildings'],
    queryFn: () => api<Building[]>('/api/v1/buildings', {}, token),
  })
  const floors = useQuery({
    queryKey: ['floors'],
    queryFn: () => api<Floor[]>('/api/v1/floors', {}, token),
  })

  if (sites.isLoading || buildings.isLoading || floors.isLoading) {
    return <CircularProgress />
  }

  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Sites & zones</Typography>
        <Typography color="text.secondary">
          Site → Building → Floor → Zone hierarchy (tenant isolated).
        </Typography>
      </div>

      <Paper sx={{ p: 2, border: '1px solid #D5DEE8' }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Sites
        </Typography>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Code</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Timezone</TableCell>
              <TableCell>Status</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(sites.data ?? []).map((s) => (
              <TableRow key={s.id}>
                <TableCell>{s.code}</TableCell>
                <TableCell>{s.name}</TableCell>
                <TableCell>{s.timezone}</TableCell>
                <TableCell>
                  <Chip size="small" label={s.status} color="success" variant="outlined" />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Paper>

      <Paper sx={{ p: 2, border: '1px solid #D5DEE8' }}>
        <Typography variant="h6" sx={{ mb: 1 }}>
          Buildings / Floors
        </Typography>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>Building</TableCell>
              <TableCell>Floor</TableCell>
              <TableCell>Size (m)</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(floors.data ?? []).map((f) => {
              const building = (buildings.data ?? []).find((b) => b.id === f.buildingId)
              return (
                <TableRow key={f.id}>
                  <TableCell>{building?.name ?? building?.code ?? '—'}</TableCell>
                  <TableCell>
                    {f.code} · {f.name} (L{f.levelIndex})
                  </TableCell>
                  <TableCell>
                    {f.widthM} × {f.heightM}
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </Paper>
    </Stack>
  )
}
