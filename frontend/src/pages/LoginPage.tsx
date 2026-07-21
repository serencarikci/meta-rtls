import { FormEvent, useState } from 'react'
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  MenuItem,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
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
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'grid',
        placeItems: 'center',
        background:
          'radial-gradient(1200px 600px at 10% 0%, #D7E6F5 0%, transparent 55%), linear-gradient(160deg, #EEF2F6 0%, #D9E3EE 100%)',
        px: 2,
      }}
    >
      <Card sx={{ width: '100%', maxWidth: 440, boxShadow: '0 18px 50px rgba(15,28,46,0.12)' }}>
        <CardContent sx={{ p: 4 }}>
          <Typography variant="overline" color="text.secondary">
            Metadata-driven RTLS
          </Typography>
          <Typography variant="h4" sx={{ mb: 1, fontFamily: '"IBM Plex Mono", monospace' }}>
            MetaRTLS
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
            Multi-tenant location platform. Sign in with a demo tenant.
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
        </CardContent>
      </Card>
    </Box>
  )
}
