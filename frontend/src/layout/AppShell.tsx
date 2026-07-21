import { Box, Button, Container, Typography } from '@mui/material'
import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/auth'

const NAV_ITEMS = [
  { to: '/', label: 'Overview' },
  { to: '/sites', label: 'Sites & Zones' },
  { to: '/metadata', label: 'Metadata' },
  { to: '/live-map', label: 'Live Map' },
  { to: '/analysis', label: 'Analysis' },
]

export default function AppShell() {
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const tenantCode = useAuthStore((s) => s.tenantCode)
  const clear = useAuthStore((s) => s.clear)

  return (
    <Box sx={{ minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <Box
        component="header"
        sx={{
          position: 'sticky',
          top: 0,
          zIndex: 30,
          borderBottom: '1px solid rgba(26,16,40,0.12)',
          background: '#1A1028',
          color: '#FCFAFD',
        }}
      >
        <Box
          sx={{
            maxWidth: 1180,
            mx: 'auto',
            px: { xs: 2, md: 3 },
            py: 1.35,
            display: 'flex',
            alignItems: 'center',
            gap: 2,
            flexWrap: 'wrap',
          }}
        >
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mr: { md: 1 } }}>
            <Box className="signal-dot" />
            <Typography sx={{ fontWeight: 700, fontSize: '1.15rem', letterSpacing: '-0.01em' }}>
              MetaRTLS
            </Typography>
          </Box>

          <Box
            component="nav"
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 0.5,
              flex: 1,
              minWidth: 0,
              overflowX: 'auto',
              py: 0.25,
            }}
          >
            {NAV_ITEMS.map((item) => (
              <Box
                key={item.to}
                component={NavLink}
                to={item.to}
                end={item.to === '/'}
                sx={{
                  px: 1.5,
                  py: 0.75,
                  borderRadius: 999,
                  color: 'rgba(252,250,253,0.72)',
                  textDecoration: 'none',
                  fontWeight: 650,
                  fontSize: '0.94rem',
                  whiteSpace: 'nowrap',
                  transition: 'background 0.15s ease, color 0.15s ease',
                  '&:hover': {
                    color: '#FCFAFD',
                    background: 'rgba(252,250,253,0.1)',
                  },
                  '&.active': {
                    color: '#1A1028',
                    background: '#C84BFF',
                    fontWeight: 700,
                  },
                }}
              >
                {item.label}
              </Box>
            ))}
          </Box>

          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 1.25,
              ml: 'auto',
            }}
          >
            <Box sx={{ textAlign: 'right', display: { xs: 'none', sm: 'block' } }}>
              <Typography sx={{ fontSize: '0.78rem', color: '#E9B7FF', fontWeight: 700 }}>
                {tenantCode}
              </Typography>
              <Typography sx={{ fontSize: '0.84rem', color: 'rgba(252,250,253,0.7)', lineHeight: 1.2 }}>
                {user?.displayName}
              </Typography>
            </Box>
            <Button
              size="small"
              onClick={() => {
                clear()
                navigate('/login')
              }}
              sx={{
                borderRadius: 999,
                px: 1.75,
                color: '#FCFAFD',
                border: '1px solid rgba(252,250,253,0.35)',
                '&:hover': {
                  borderColor: '#C84BFF',
                  color: '#E9B7FF',
                  background: 'rgba(200,75,255,0.12)',
                },
              }}
            >
              Sign out
            </Button>
          </Box>
        </Box>
      </Box>

      <Box component="main" sx={{ flex: 1 }}>
        <Container maxWidth="lg" sx={{ py: { xs: 3, md: 4 } }}>
          <Outlet />
        </Container>
      </Box>
    </Box>
  )
}
