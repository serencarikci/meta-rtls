import {
  AppBar,
  Box,
  Button,
  Container,
  Drawer,
  List,
  ListItemButton,
  ListItemText,
  Toolbar,
  Typography,
} from '@mui/material'
import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../store/auth'

const NAV_ITEMS = [
  { to: '/', label: 'Overview' },
  { to: '/sites', label: 'Sites & Zones' },
  { to: '/live-map', label: 'Live Map' },
]

export default function AppShell() {
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const tenantCode = useAuthStore((s) => s.tenantCode)
  const clear = useAuthStore((s) => s.clear)

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <AppBar position="fixed" elevation={0}>
        <Toolbar>
          <Typography variant="h6" sx={{ flexGrow: 1, fontFamily: '"IBM Plex Mono", monospace' }}>
            MetaRTLS
          </Typography>
          <Typography variant="body2" sx={{ mr: 2, opacity: 0.9 }}>
            {tenantCode} · {user?.displayName}
          </Typography>
          <Button
            color="inherit"
            onClick={() => {
              clear()
              navigate('/login')
            }}
          >
            Sign out
          </Button>
        </Toolbar>
      </AppBar>

      <Drawer
        variant="permanent"
        sx={{
          width: 240,
          [`& .MuiDrawer-paper`]: {
            width: 240,
            mt: 8,
            borderRight: '1px solid #D5DEE8',
            background: 'linear-gradient(180deg, #FFFFFF 0%, #EEF2F6 100%)',
          },
        }}
      >
        <List>
          {NAV_ITEMS.map((item) => (
            <ListItemButton
              key={item.to}
              component={NavLink}
              to={item.to}
              end={item.to === '/'}
              sx={{
                '&.active': { bgcolor: 'rgba(22,50,79,0.08)', borderRight: '3px solid #16324F' },
              }}
            >
              <ListItemText primary={item.label} />
            </ListItemButton>
          ))}
        </List>
      </Drawer>

      <Box component="main" sx={{ flexGrow: 1, mt: 10, mb: 4 }}>
        <Container maxWidth="lg">
          <Outlet />
        </Container>
      </Box>
    </Box>
  )
}
