import { Box, Button, Container, Typography } from '@mui/material'
import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/auth'

const NAV_ITEMS = [
  { to: '/', label: 'Overview' },
  { to: '/sites', label: 'Sites & Zones' },
  { to: '/metadata', label: 'Metadata' },
  { to: '/live-map', label: 'Live Map' },
]

const SIDEBAR_WIDTH = 232

export default function AppShell() {
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const tenantCode = useAuthStore((s) => s.tenantCode)
  const clear = useAuthStore((s) => s.clear)

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <Box
        component="aside"
        sx={{
          width: SIDEBAR_WIDTH,
          flexShrink: 0,
          position: 'fixed',
          inset: 0,
          right: 'auto',
          display: 'flex',
          flexDirection: 'column',
          borderRight: '1px solid var(--line)',
          background: 'linear-gradient(180deg, #241338 0%, #1A1028 100%)',
          color: '#FCFAFD',
          zIndex: 20,
        }}
      >
        <Box sx={{ px: 2.5, pt: 3, pb: 2.5, display: 'flex', alignItems: 'center', gap: 1.1 }}>
          <Box className="signal-dot" />
          <Typography
            sx={{
              fontFamily: '"Syne", sans-serif',
              fontWeight: 800,
              letterSpacing: '-0.04em',
              fontSize: '1.25rem',
            }}
          >
            MetaRTLS
          </Typography>
        </Box>

        <Box sx={{ px: 1.5, display: 'flex', flexDirection: 'column', gap: 0.4, flex: 1 }}>
          {NAV_ITEMS.map((item) => (
            <Box
              key={item.to}
              component={NavLink}
              to={item.to}
              end={item.to === '/'}
              sx={{
                px: 1.5,
                py: 1.1,
                color: 'rgba(252,250,253,0.62)',
                textDecoration: 'none',
                fontFamily: '"Figtree", sans-serif',
                fontWeight: 600,
                fontSize: '0.95rem',
                borderLeft: '3px solid transparent',
                transition: 'color 0.15s ease, background 0.15s ease',
                '&.active': {
                  color: '#E9B7FF',
                  background: 'rgba(200,75,255,0.12)',
                  borderLeftColor: '#C84BFF',
                },
                '&:hover': { color: '#FCFAFD' },
              }}
            >
              {item.label}
            </Box>
          ))}
        </Box>

        <Box sx={{ px: 2.5, pb: 2.5 }}>
          <Typography
            sx={{
              fontFamily: '"JetBrains Mono", monospace',
              fontSize: '0.68rem',
              letterSpacing: '0.08em',
              color: 'rgba(200,75,255,0.85)',
              mb: 0.5,
              textTransform: 'uppercase',
            }}
          >
            {tenantCode}
          </Typography>
          <Typography variant="body2" sx={{ color: 'rgba(252,250,253,0.75)', mb: 1.5 }}>
            {user?.displayName}
          </Typography>
          <Button
            fullWidth
            variant="outlined"
            onClick={() => {
              clear()
              navigate('/login')
            }}
            sx={{
              color: '#FCFAFD',
              borderColor: 'rgba(252,250,253,0.28)',
              '&:hover': {
                borderColor: '#C84BFF',
                color: '#E9B7FF',
                background: 'transparent',
              },
            }}
          >
            Sign out
          </Button>
        </Box>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          ml: `${SIDEBAR_WIDTH}px`,
          minHeight: '100vh',
        }}
      >
        <Box
          sx={{
            borderBottom: '1px solid var(--line)',
            background: 'rgba(250,251,248,0.86)',
            backdropFilter: 'blur(8px)',
            px: 3,
            py: 1.75,
            position: 'sticky',
            top: 0,
            zIndex: 10,
          }}
        >
          <Typography
            sx={{
              fontFamily: '"JetBrains Mono", monospace',
              fontSize: '0.72rem',
              letterSpacing: '0.14em',
              color: 'var(--mute)',
              textTransform: 'uppercase',
            }}
          >
            Indoor location control
          </Typography>
        </Box>
        <Container maxWidth="lg" sx={{ py: 4 }}>
          <Outlet />
        </Container>
      </Box>
    </Box>
  )
}
