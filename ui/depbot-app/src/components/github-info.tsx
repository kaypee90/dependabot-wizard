import React from 'react';
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
import PullRequestLineGraph from './pull-request-line-graph';
import PullRequestDataGrid from './pull-request-data-grid';
import { ThemeProvider } from '@mui/material/styles';
import { lightTheme } from '../utils/theme';
import { PullRequestState } from '../models/pull-request';

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

  // Derived during render rather than copied into state by an effect: this
  // component observes the store, so a second fetch updates the counts on
  // its own. The effect version left them showing the previous repository.
  const countOf = (state: string) =>
    configStore.pullRequests.filter((pr) => pr.state.toLowerCase() === state).length;

  const all = configStore.pullRequests.length;
  const opened = countOf(PullRequestState.OPEN);
  const merged = countOf(PullRequestState.MERGED);
  const closed = countOf(PullRequestState.CLOSED);

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
                DEPENDABOT PRs
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
                <PullRequestCount title="Opened" count={opened} />
              </Grid>
              <Grid size={{ xs: 6, md: 3 }}>
                <PullRequestCount title="Closed" count={closed} />
              </Grid>
              <Grid size={{ xs: 6, md: 3 }}>
                <PullRequestCount title="Merged" count={merged} />
              </Grid>
            </Grid>

            <Grid container spacing={1} sx={{ marginTop: '2vh', marginLeft: '4vh', marginRight: '1vh' }}>
              <PullRequestLineGraph />
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

