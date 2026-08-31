import { createTheme } from '@mui/material/styles';

/**
 * One dark theme for the whole application. The configuration screen and the
 * pull request dashboard previously rendered in different modes, which read
 * as two products bolted together.
 *
 * The palette is built around #0E2535, the navy the app has always used, with
 * a cyan accent that stays legible on it and semantic colours kept separate
 * from that accent so state never competes with emphasis.
 */

export const palette = {
  ground: '#0B1A26',
  surface: '#12293A',
  raised: '#173347',
  border: '#22475F',
  borderSoft: '#1B3A50',
  accent: '#45B6D8',
  accentSoft: '#123243',
  text: '#E7EFF5',
  textMuted: '#93A9B9',
  textFaint: '#6D8496',
  success: '#4CAF8B',
  warning: '#D7A45A',
};

const SANS = "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif";
export const MONO = "'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, Consolas, monospace";

export const appTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: { main: '#45B6D8', contrastText: '#06131C' },
    secondary: { main: '#4CAF8B' },
    success: { main: '#4CAF8B' },
    warning: { main: '#D7A45A' },
    error: { main: '#E27A6C' },
    info: { main: '#45B6D8' },
    background: { default: '#0B1A26', paper: '#12293A' },
    text: { primary: '#E7EFF5', secondary: '#93A9B9', disabled: '#6D8496' },
    divider: '#22475F',
  },

  shape: { borderRadius: 8 },

  typography: {
    fontFamily: SANS,
    h1: { fontSize: '1.6rem', fontWeight: 700, letterSpacing: '-0.02em' },
    h6: { fontSize: '1.05rem', fontWeight: 600, letterSpacing: '-0.01em' },
    subtitle1: { fontSize: '0.95rem', fontWeight: 600, letterSpacing: '-0.005em' },
    subtitle2: { fontSize: '0.78rem', fontWeight: 600, letterSpacing: '0.06em', textTransform: 'uppercase' },
    body1: { fontSize: '0.92rem', lineHeight: 1.6 },
    body2: { fontSize: '0.85rem', lineHeight: 1.6 },
    button: { textTransform: 'none', fontWeight: 600, letterSpacing: 0 },
    caption: { fontSize: '0.76rem', lineHeight: 1.5 },
  },

  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: { backgroundColor: '#0B1A26' },
      },
    },

    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#0B1A26',
          backgroundImage: 'none',
          borderBottom: '1px solid #22475F',
          boxShadow: 'none',
        },
      },
    },

    MuiAccordion: {
      styleOverrides: {
        root: {
          backgroundColor: '#12293A',
          backgroundImage: 'none',
          border: '1px solid #22475F',
          borderRadius: 8,
          boxShadow: 'none',
          '&:before': { display: 'none' },
          '&.Mui-expanded': { margin: 0 },
        },
      },
    },
    MuiAccordionSummary: {
      styleOverrides: {
        root: { minHeight: 56, '&.Mui-expanded': { minHeight: 56 } },
        content: { margin: '12px 0', '&.Mui-expanded': { margin: '12px 0' } },
      },
    },
    MuiAccordionDetails: {
      styleOverrides: { root: { padding: '4px 20px 20px' } },
    },

    MuiButton: {
      styleOverrides: {
        root: { borderRadius: 8, paddingInline: 18, paddingBlock: 9 },
        containedPrimary: {
          '&:hover': { backgroundColor: '#63C6E3' },
        },
        outlined: { borderColor: '#22475F' },
      },
    },

    MuiTextField: {
      defaultProps: { size: 'small' },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: '#0E2231',
          '& fieldset': { borderColor: '#22475F' },
          '&:hover fieldset': { borderColor: '#2E5C79' },
        },
        input: { fontSize: '0.9rem' },
      },
    },
    MuiSelect: { defaultProps: { size: 'small' } },
    MuiFormHelperText: {
      styleOverrides: { root: { marginLeft: 2, fontSize: '0.74rem' } },
    },

    MuiChip: {
      styleOverrides: {
        root: { borderRadius: 6, fontWeight: 500 },
        sizeSmall: { fontSize: '0.75rem', height: 24 },
      },
    },

    MuiAlert: {
      styleOverrides: {
        root: { borderRadius: 8, fontSize: '0.85rem', border: '1px solid transparent' },
        standardInfo: { backgroundColor: '#123243', color: '#CFE7F1', borderColor: '#1E4B63' },
        standardError: { backgroundColor: '#3A1F1C', color: '#F0CFC9', borderColor: '#5A2F29' },
        standardSuccess: { backgroundColor: '#11302A', color: '#C9E9DD', borderColor: '#1E5142' },
      },
    },

    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#12293A',
          backgroundImage: 'none',
          border: '1px solid #22475F',
          boxShadow: 'none',
        },
      },
    },

    MuiDialog: {
      styleOverrides: {
        paper: { backgroundImage: 'none', border: '1px solid #22475F' },
      },
    },

    MuiTooltip: {
      styleOverrides: {
        tooltip: { backgroundColor: '#173347', fontSize: '0.76rem', border: '1px solid #22475F' },
      },
    },
  },
});

/** Kept for imports that still reference it; both surfaces are dark now. */
export const lightTheme = appTheme;
