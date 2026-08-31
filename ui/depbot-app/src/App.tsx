import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

import { MainConfiguration, GithubInfoButton, GithubInfo, GithubRepoDetails } from './components';
import { isStaticBuild } from './utils/runtime';

export default function ButtonAppBar() {
  const darkTheme = createTheme({
    palette: {
      mode: 'dark',
    },
  });

  // The pull request dashboard proxies through the local API, so it only
  // exists when a depbot binary is serving the page.
  const showDashboard = !isStaticBuild();

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1 }}>
        <AppBar sx={{ background: '#0E2535' }} position="static">
          <Toolbar>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              Depbot
            </Typography>
            {showDashboard && <GithubInfoButton />}
          </Toolbar>
        </AppBar>
        {showDashboard && <GithubInfo />}
        {showDashboard && <GithubRepoDetails />}
        <MainConfiguration />
      </Box>
    </ThemeProvider>
  );
}
