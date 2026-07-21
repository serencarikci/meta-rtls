import { FormEvent, useState } from 'react'
import { Alert, Box, Button, MenuItem, Stack, TextField, Typography } from '@mui/material'
import { useNavigate } from 'react-router-dom'
import { api } from '../api/client'
import { useAuthStore, AuthUser } from '../store/auth'

type LoginResponse = {
  accessToken: string
  expiresAt: string
  user: AuthUser
}

const DEMO_TENANTS = [
  { code: 'warehouse-s', label: 'Warehouse (Small)' },
  { code: 'hospital-m', label: 'Hospital (Medium)' },
  { code: 'factory-l', label: 'Factory (Large)' },
]

const DEMO_PASSWORD = 'MetaRTLS!2026'

export default function LoginPage() {
  const navigate = useNavigate()
  const setSession = useAuthStore((s) => s.setSession)
  const [tenantCode, setTenantCode] = useState('warehouse-s')
  const [email, setEmail] = useState('admin@warehouse-s.demo')
  const [password, setPassword] = useState(DEMO_PASSWORD)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      const data = await api<LoginResponse>('/api/v1/auth/login', {
        method: 'POST',
        body: JSON.stringify({ tenantCode, email, password }),
      })
      setSession(data.accessToken, data.user, tenantCode)
      navigate('/')
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Login failed'
      const lower = message.toLowerCase()
      if (
        message === 'Not Found' ||
        message === 'Internal Server Error' ||
        lower.includes('failed to fetch') ||
        lower.includes('econnrefused')
      ) {
        setError('API/Oracle is not ready. Run: make up, then make backend-run (port 8090)')
      } else {
        setError(message)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'grid',
        gridTemplateColumns: { xs: '1fr', md: '1.15fr 0.85fr' },
      }}
    >
      <Box
        className="surface-grid"
        sx={{
          position: 'relative',
          display: { xs: 'none', md: 'flex' },
          flexDirection: 'column',
          justifyContent: 'space-between',
          p: { md: 6, lg: 8 },
          color: '#FAFBF8',
          overflow: 'hidden',
        }}
      >
        <Box className="rise-in" sx={{ display: 'flex', alignItems: 'center', gap: 1.25 }}>
          <Box className="signal-dot" />
          <Typography sx={{ fontWeight: 700, fontSize: '1.35rem' }}>MetaRTLS</Typography>
        </Box>

        <Box className="rise-in-delay" sx={{ maxWidth: 520 }}>
          <Typography
            sx={{
              fontWeight: 700,
              fontSize: { md: '2.8rem', lg: '3.25rem' },
              lineHeight: 1.18,
              letterSpacing: '-0.02em',
              mb: 2.5,
            }}
          >
            Locate
            <br />
            every asset
            <Box component="span" sx={{ color: 'var(--signal)' }}>
              .
            </Box>
          </Typography>
          <Typography
            sx={{
              color: 'rgba(252,250,253,0.78)',
              maxWidth: 400,
              fontSize: '1.05rem',
              lineHeight: 1.65,
            }}
          >
            Metadata-driven indoor tracking for warehouses, hospitals and factories — one model,
            many tenants.
          </Typography>
        </Box>

        <Typography sx={{ fontSize: '0.9rem', color: 'rgba(233,183,255,0.9)', fontWeight: 500 }}>
          Live signal · schema aware · multi-tenant
        </Typography>
      </Box>

      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          px: { xs: 2.5, sm: 4 },
          py: 6,
          background: 'linear-gradient(165deg, #FCFAFD 0%, #EFEAF6 55%, #E4DCF0 100%)',
        }}
      >
        <Box className="rise-in" sx={{ width: '100%', maxWidth: 400 }}>
          <Box sx={{ display: { xs: 'flex', md: 'none' }, alignItems: 'center', gap: 1, mb: 3 }}>
            <Box className="signal-dot" />
            <Typography sx={{ fontWeight: 700 }}>MetaRTLS</Typography>
          </Box>

          <Typography variant="overline" sx={{ color: 'var(--mute)' }}>
            Sign in
          </Typography>
          <Typography variant="h4" sx={{ mb: 1 }}>
            Enter workspace
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3.5 }}>
            Pick a demo tenant and continue.
          </Typography>

          <Stack component="form" spacing={2} onSubmit={onSubmit}>
            <TextField
              select
              label="Tenant"
              value={tenantCode}
              onChange={(e) => {
                const code = e.target.value
                setTenantCode(code)
                setEmail(`admin@${code}.demo`)
              }}
            >
              {DEMO_TENANTS.map((demoTenant) => (
                <MenuItem key={demoTenant.code} value={demoTenant.code}>
                  {demoTenant.label}
                </MenuItem>
              ))}
            </TextField>
            <TextField label="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
            <TextField
              label="Password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            {error && <Alert severity="error">{error}</Alert>}
            <Button type="submit" variant="contained" size="large" disabled={loading}>
              {loading ? 'Signing in…' : 'Sign in'}
            </Button>
          </Stack>
        </Box>
      </Box>
    </Box>
  )
}
