import { useQuery } from '@tanstack/react-query'
import { Alert, Paper, Stack, Typography } from '@mui/material'
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

  return (
    <Stack spacing={3}>
      <div>
        <Typography variant="h4">Platform overview</Typography>
        <Typography color="text.secondary">
          Setup and phase status for tenant <strong>{tenantCode}</strong>.
        </Typography>
      </div>

      {(sites.isError || floors.isError) && (
        <Alert severity="warning">
          If the API is not ready yet, wait until Oracle is up (`make up`).
        </Alert>
      )}

      <Grid container spacing={2}>
        {[
          { label: 'Sites', value: sites.data?.length ?? '—' },
          { label: 'Floors', value: floors.data?.length ?? '—' },
          { label: 'Phase', value: '1 / 5' },
          { label: 'Engine', value: 'Oracle + Go' },
        ].map((card) => (
          <Grid key={card.label} size={{ xs: 12, sm: 6, md: 3 }}>
            <Paper
              sx={{
                p: 2.5,
                background: 'linear-gradient(145deg, #FFFFFF 0%, #F5F8FB 100%)',
                border: '1px solid #D5DEE8',
              }}
            >
              <Typography variant="overline" color="text.secondary">
                {card.label}
              </Typography>
              <Typography variant="h5">{card.value}</Typography>
            </Paper>
          </Grid>
        ))}
      </Grid>

      <Paper sx={{ p: 3, border: '1px solid #D5DEE8' }}>
        <Typography variant="h6" gutterBottom>
          Roadmap
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Phase 1: Auth + tenant + site hierarchy ✓ · Phase 2: Metadata engine · Phase 3: MQTT
          simulator + live map · Phase 4: Requirement / impact analysis · Phase 5: Production
          quality
        </Typography>
      </Paper>
    </Stack>
  )
}
