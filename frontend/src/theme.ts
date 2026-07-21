import { createTheme } from '@mui/material/styles'

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
    fontFamily: '"Figtree", "Helvetica Neue", sans-serif',
    h3: {
      fontFamily: '"Syne", sans-serif',
      fontWeight: 770,
      letterSpacing: '-0.04em',
    },
    h4: {
      fontFamily: '"Syne", sans-serif',
      fontWeight: 740,
      letterSpacing: '-0.03em',
    },
    h5: {
      fontFamily: '"Syne", sans-serif',
      fontWeight: 700,
      letterSpacing: '-0.02em',
    },
    h6: {
      fontFamily: '"Syne", sans-serif',
      fontWeight: 700,
    },
    overline: {
      fontFamily: '"JetBrains Mono", monospace',
      letterSpacing: '0.14em',
      fontWeight: 500,
    },
    button: {
      textTransform: 'none',
      fontWeight: 700,
      letterSpacing: '-0.01em',
    },
  },
  shape: { borderRadius: 4 },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundImage:
            'radial-gradient(circle at 1px 1px, rgba(26,16,40,0.06) 1px, transparent 0)',
          backgroundSize: '18px 18px',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        containedPrimary: {
          backgroundColor: '#1A1028',
          color: '#F5E9FF',
          borderRadius: 999,
          paddingLeft: 22,
          paddingRight: 22,
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
          borderRadius: 4,
          '& fieldset': { borderColor: '#D8CFE6' },
          '&:hover fieldset': { borderColor: '#1A1028' },
          '&.Mui-focused fieldset': { borderColor: '#8B2FD6', borderWidth: 1.5 },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 4,
          fontFamily: '"JetBrains Mono", monospace',
          fontSize: '0.75rem',
        },
      },
    },
  },
})
