import { createTheme } from '@mui/material/styles'

const FONT = '"Source Sans 3", "Segoe UI", sans-serif'

export const APP_THEME = createTheme({
  palette: {
    mode: 'light',
    primary: { main: '#1A1028', contrastText: '#F5E9FF' },
    secondary: { main: '#C84BFF', contrastText: '#1A1028' },
    error: { main: '#E4572E' },
    background: { default: '#EFEAF6', paper: '#FCFAFD' },
    text: { primary: '#1A1028', secondary: '#746887' },
    divider: '#D8CFE6',
  },
  typography: {
    fontFamily: FONT,
    fontSize: 15,
    htmlFontSize: 16,
    body1: { fontFamily: FONT, fontSize: '1rem', lineHeight: 1.65 },
    body2: { fontFamily: FONT, fontSize: '0.95rem', lineHeight: 1.6 },
    h3: { fontFamily: FONT, fontWeight: 700, letterSpacing: '-0.01em', lineHeight: 1.25 },
    h4: { fontFamily: FONT, fontWeight: 700, letterSpacing: '-0.01em', lineHeight: 1.3 },
    h5: { fontFamily: FONT, fontWeight: 650, letterSpacing: 0, lineHeight: 1.35 },
    h6: { fontFamily: FONT, fontWeight: 650, letterSpacing: 0 },
    overline: {
      fontFamily: FONT,
      letterSpacing: '0.06em',
      fontWeight: 600,
      fontSize: '0.72rem',
    },
    button: {
      fontFamily: FONT,
      textTransform: 'none',
      fontWeight: 600,
      letterSpacing: 0,
    },
  },
  shape: { borderRadius: 12 },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          fontFamily: FONT,
          backgroundImage:
            'radial-gradient(circle at 1px 1px, rgba(26,16,40,0.05) 1px, transparent 0)',
          backgroundSize: '20px 20px',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: { fontFamily: FONT },
        containedPrimary: {
          backgroundColor: '#1A1028',
          color: '#F5E9FF',
          borderRadius: 999,
          paddingLeft: 20,
          paddingRight: 20,
          boxShadow: 'none',
          '&:hover': {
            backgroundColor: '#2B1848',
            boxShadow: 'none',
          },
        },
        outlined: {
          borderColor: '#1A1028',
          borderRadius: 999,
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          backgroundImage: 'none',
          border: '1px solid #D8CFE6',
          boxShadow: 'none',
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: '#FFFFFF',
          borderRadius: 12,
          '& fieldset': { borderColor: '#D8CFE6' },
          '&:hover fieldset': { borderColor: '#1A1028' },
          '&.Mui-focused fieldset': { borderColor: '#8B2FD6', borderWidth: 1.5 },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 999,
          fontFamily: FONT,
          fontSize: '0.8rem',
          fontWeight: 600,
        },
      },
    },
  },
})
