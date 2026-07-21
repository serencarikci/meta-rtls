import { createTheme } from '@mui/material/styles'

export const APP_THEME = createTheme({
  palette: {
    mode: 'light',
    primary: { main: '#16324F' },
    secondary: { main: '#E09F3E' },
    background: { default: '#EEF2F6', paper: '#FFFFFF' },
    text: { primary: '#0F1720', secondary: '#5B6B7C' },
  },
  typography: {
    fontFamily: '"IBM Plex Sans", "Segoe UI", sans-serif',
    h4: { fontWeight: 700, letterSpacing: '-0.02em' },
    h5: { fontWeight: 650 },
    button: { textTransform: 'none', fontWeight: 600 },
  },
  shape: { borderRadius: 8 },
  components: {
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundImage: 'linear-gradient(115deg, #0F1C2E 0%, #16324F 50%, #1F4B73 100%)',
        },
      },
    },
  },
})
