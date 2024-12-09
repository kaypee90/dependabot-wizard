
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

import { MainConfiguration, GithubInfoButton, GithubInfo, GithubRepoDetails } from './components';

export default function ButtonAppBar() {
  const darkTheme = createTheme({
    palette: {
      mode: 'dark',
    },
  });


  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1 }}>
        <AppBar sx={{ background: '#0E2535' }} position="static">
          <Toolbar>
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              Depbot
            </Typography>
            <GithubInfoButton />
          </Toolbar>
        </AppBar>
        <GithubInfo />
        <GithubRepoDetails />
        <MainConfiguration />
      </Box>
    </ThemeProvider>
  );
}

