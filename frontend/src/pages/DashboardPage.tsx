import { useQuery } from '@tanstack/react-query'
import { Alert, Box, Paper, Stack, Typography } from '@mui/material'
import Grid from '@mui/material/Grid2'
import { api } from '../api/client'
import { useAuthStore } from '../store/auth'

type Site = { id: string; code: string; name: string; timezone: string }
type Floor = { id: string; code: string; name: string; widthM: number; heightM: number }

export default function DashboardPage() {
  const token = useAuthStore((s) => s.token)
  const tenantCode = useAuthStore((s) => s.tenantCode)

  const sites = useQuery({
    queryKey: ['sites'],
    queryFn: () => api<Site[]>('/api/v1/sites', {}, token),
  })
  const floors = useQuery({
    queryKey: ['floors'],
    queryFn: () => api<Floor[]>('/api/v1/floors', {}, token),
  })

  const stats = [
    { label: 'Sites', value: sites.data?.length ?? '—' },
    { label: 'Floors', value: floors.data?.length ?? '—' },
    { label: 'Stack', value: 'Go + React' },
    { label: 'Database', value: 'Oracle' },
  ]

  return (
    <Stack spacing={3} className="rise-in">
      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: { xs: '1fr', md: '1.2fr 0.8fr' },
          gap: 2,
          alignItems: 'stretch',
        }}
      >
        <Paper
          sx={{
            p: { xs: 3, md: 4 },
            border: '1px solid var(--line)',
            background:
              'linear-gradient(135deg, #FCFAFD 0%, #F3ECF9 55%, rgba(200,75,255,0.16) 100%)',
          }}
        >
          <Typography variant="overline">Workspace</Typography>
          <Typography variant="h3" sx={{ mt: 0.5, mb: 1.5, maxWidth: 520 }}>
            {tenantCode}
            <Box component="span" sx={{ color: 'var(--signal-deep)' }}>
              .
            </Box>
          </Typography>
          <Typography color="text.secondary" sx={{ maxWidth: 460 }}>
            Metadata fields, sites and live zones for this tenant — ready for the RTLS signal path.
          </Typography>
        </Paper>

        <Paper
          sx={{
            p: 3,
            background: '#1A1028',
            color: '#FCFAFD',
            border: 'none',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'space-between',
            minHeight: 180,
          }}
        >
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box className="signal-dot" />
            <Typography sx={{ fontFamily: '"JetBrains Mono", monospace', fontSize: '0.75rem' }}>
              SIGNAL READY
            </Typography>
          </Box>
          <Typography
            sx={{ fontFamily: '"Syne", sans-serif', fontWeight: 740, fontSize: '1.6rem' }}
          >
            Live tags online.
            <br />
            Metadata aware.
          </Typography>
        </Paper>
      </Box>

      {(sites.isError || floors.isError) && (
        <Alert severity="warning">
          If the API is not ready yet, wait until Oracle is up (`make up`).
        </Alert>
      )}

      <Grid container spacing={2}>
        {stats.map((stat) => (
          <Grid key={stat.label} size={{ xs: 12, sm: 6, md: 3 }}>
            <Paper sx={{ p: 2.5, borderTop: '3px solid var(--signal)' }}>
              <Typography variant="overline" color="text.secondary">
                {stat.label}
              </Typography>
              <Typography variant="h5" sx={{ mt: 0.5 }}>
                {stat.value}
              </Typography>
            </Paper>
          </Grid>
        ))}
      </Grid>
    </Stack>
  )
}
