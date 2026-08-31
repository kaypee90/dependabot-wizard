import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

import { MainConfiguration, GithubInfoButton, GithubInfo, GithubRepoDetails } from './components';
import { appTheme, palette } from './utils/theme';
import { isStaticBuild } from './utils/runtime';

const Wordmark = () => (
  <Stack direction="row" alignItems="baseline" spacing={1.25}>
    <Typography
      component="span"
      sx={{ fontSize: '1.15rem', fontWeight: 700, letterSpacing: '-0.02em', color: palette.text }}
    >
      Depbot
    </Typography>
    <Typography
      component="span"
      sx={{
        display: { xs: 'none', sm: 'block' },
        fontSize: '0.82rem',
        color: palette.textFaint,
        borderLeft: `1px solid ${palette.border}`,
        pl: 1.5,
      }}
    >
      dependabot.yml, for every ecosystem you actually have
    </Typography>
  </Stack>
);

export default function App() {
  // The pull request dashboard proxies through the local API, so it only
  // exists when a depbot binary is serving the page.
  const showDashboard = !isStaticBuild();

  return (
    <ThemeProvider theme={appTheme}>
      <CssBaseline />
      <Box sx={{ minHeight: '100vh', background: palette.ground }}>
        <AppBar position="sticky">
          <Toolbar sx={{ minHeight: 60, px: { xs: 2, md: 3 } }}>
            <Box sx={{ flexGrow: 1 }}>
              <Wordmark />
            </Box>
            <Stack direction="row" spacing={1} alignItems="center">
              {showDashboard && <GithubInfoButton />}
              <Button
                size="small"
                href="https://github.com/kaypee90/dependabot-wizard"
                target="_blank"
                rel="noreferrer"
                sx={{ color: palette.textMuted, '&:hover': { color: palette.text } }}
              >
                GitHub
              </Button>
            </Stack>
          </Toolbar>
        </AppBar>

        {showDashboard && <GithubInfo />}
        {showDashboard && <GithubRepoDetails />}
        <MainConfiguration />
      </Box>
    </ThemeProvider>
  );
}
