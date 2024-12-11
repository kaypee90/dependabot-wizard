import React, { useState } from 'react';
import Dialog from '@mui/material/Dialog';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid2';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import CloseIcon from '@mui/icons-material/Close';
import Slide from '@mui/material/Slide';
import { TransitionProps } from '@mui/material/transitions';
import { observer } from 'mobx-react-lite';
import { configStore } from '../store/configuration-store';
import PullRequestCount from './pull-request-count';
import PullRequestBarGraph from './pull-request-bar-graph';
import PullRequestLineGraph from './pull-request-line-graph';
import PullRequestDataGrid from './pull-request-data-grid';
import { ThemeProvider, createTheme } from '@mui/material/styles';

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & {
    children: React.ReactElement<unknown>;
  },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="left" ref={ref} {...props} />;
});

export const GithubInfo: React.FC = observer(() => {
  const handleClose = () => {
    configStore.setShowGithubInfo(false);
  };

  const lightTheme = createTheme({
    palette: {
      mode: 'light',
    },
  });

  const [all, setAll] = useState(0);
  const [opened, setOpened] = useState(0);
  const [merged, setMerged] = useState(0);
  const [closed, setClosed] = useState(0);

  return (
    <ThemeProvider theme={lightTheme}>
      <React.Fragment>
        <Dialog
          fullScreen
          open={configStore.showGithubInfo}
          onClose={handleClose}
          TransitionComponent={Transition}
        >
          <AppBar sx={{ position: 'relative', background: '#0E2535' }}>
            <Toolbar>

              <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
                GITHUB INFO
              </Typography>
              <IconButton
                edge="start"
                color="inherit"
                onClick={handleClose}
                aria-label="close"
              >
                <CloseIcon />
              </IconButton>
            </Toolbar>
          </AppBar>
          <Box sx={{ flexGrow: 1, background: '#fff', color: '#000' }}>
            <Grid container spacing={1} sx={{ marginTop: '2.5vh', marginLeft: '4vh', marginRight: '1vh' }}>
              <Grid size={{ xs: 6, md: 3 }}>
                <PullRequestCount title="All" count={all} />
              </Grid>
              <Grid size={{ xs: 6, md: 3 }}>
                <PullRequestCount title="Opened" count={opened} />            </Grid>
              <Grid size={{ xs: 6, md: 3 }}>
                <PullRequestCount title="Closed" count={closed} />             </Grid>
              <Grid size={{ xs: 6, md: 3 }}>
                <PullRequestCount title="Merged" count={merged} />
              </Grid>
            </Grid>

            <Grid container spacing={1} sx={{ marginTop: '2vh', marginLeft: '4vh', marginRight: '1vh' }}>
              <Grid size={{ xs: 6, md: 6 }}>
                <PullRequestBarGraph />
              </Grid>
              <Grid size={{ xs: 6, md: 6 }}>
                <PullRequestLineGraph />
              </Grid>
            </Grid>

            <Grid container spacing={1} sx={{ marginTop: '2vh', marginLeft: '4vh', marginRight: '1vh' }}>
              <PullRequestDataGrid />
            </Grid>

          </Box>
        </Dialog>
      </React.Fragment>
    </ThemeProvider>
  );
});

