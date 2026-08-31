import React from 'react';
import Dialog from '@mui/material/Dialog';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import CloseIcon from '@mui/icons-material/Close';
import Slide from '@mui/material/Slide';
import Alert from '@mui/material/Alert';
import { TransitionProps } from '@mui/material/transitions';
import { observer } from 'mobx-react-lite';

import { configStore } from '../store/configuration-store';
import PullRequestCount from './pull-request-count';
import PullRequestLineGraph from './pull-request-line-graph';
import PullRequestDataGrid from './pull-request-data-grid';
import { PullRequestState } from '../models/pull-request';
import { palette } from '../utils/theme';

const Transition = React.forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement<unknown> },
  ref: React.Ref<unknown>,
) {
  return <Slide direction="left" ref={ref} {...props} />;
});

const SectionHeading = ({ children }: { children: React.ReactNode }) => (
  <Typography variant="subtitle2" sx={{ color: palette.textFaint, mb: 1.5 }}>
    {children}
  </Typography>
);

export const GithubInfo: React.FC = observer(() => {
  const handleClose = () => configStore.setShowGithubInfo(false);

  // Derived during render rather than copied into state by an effect: this
  // component observes the store, so a second fetch updates the counts on
  // its own. The effect version left them showing the previous repository.
  const countOf = (state: string) =>
    configStore.pullRequests.filter((pr) => pr.state.toLowerCase() === state).length;

  const all = configStore.pullRequests.length;
  const opened = countOf(PullRequestState.OPEN);
  const merged = countOf(PullRequestState.MERGED);
  const closed = countOf(PullRequestState.CLOSED);
  const resolved = merged + closed;
  const mergeRate = resolved > 0 ? Math.round((merged / resolved) * 100) : 0;

  return (
    <Dialog fullScreen open={configStore.showGithubInfo} onClose={handleClose} TransitionComponent={Transition}>
      <AppBar sx={{ position: 'relative' }}>
        <Toolbar sx={{ minHeight: 60 }}>
          <Stack sx={{ flex: 1 }}>
            <Typography variant="h6">Dependabot pull requests</Typography>
            <Typography variant="caption" sx={{ color: palette.textFaint }}>
              {configStore.pullRequestFilter || 'Matched by author'}
            </Typography>
          </Stack>
          <IconButton edge="end" color="inherit" onClick={handleClose} aria-label="close">
            <CloseIcon />
          </IconButton>
        </Toolbar>
      </AppBar>

      <Box sx={{ flexGrow: 1, background: palette.ground, p: { xs: 2, md: 3 }, overflowY: 'auto' }}>
        <Box sx={{ maxWidth: 1280, mx: 'auto' }}>
          {all === 0 ? (
            <Alert severity="info">
              No pull requests {configStore.pullRequestFilter || 'matched'} were found. If Dependabot
              is configured here, check the owner and repository, and that your token can see this
              repository.
            </Alert>
          ) : (
            <Stack spacing={4}>
              <Box>
                <SectionHeading>Overview</SectionHeading>
                <Box sx={{
                  display: 'grid',
                  gridTemplateColumns: { xs: 'repeat(2, 1fr)', md: 'repeat(5, 1fr)' },
                  gap: 2,
                }}>
                  <PullRequestCount title="Total" count={all} />
                  <PullRequestCount title="Open" count={opened} emphasis />
                  <PullRequestCount title="Merged" count={merged} />
                  <PullRequestCount title="Closed" count={closed} />
                  <PullRequestCount title="Merge rate %" count={mergeRate} />
                </Box>
              </Box>

              <Box>
                <SectionHeading>Opened over time</SectionHeading>
                <Box sx={{
                  border: `1px solid ${palette.border}`,
                  borderRadius: 2,
                  background: palette.surface,
                  p: 1,
                }}>
                  <PullRequestLineGraph />
                </Box>
              </Box>

              <Box>
                <SectionHeading>All pull requests</SectionHeading>
                <PullRequestDataGrid />
              </Box>
            </Stack>
          )}
        </Box>
      </Box>
    </Dialog>
  );
});
