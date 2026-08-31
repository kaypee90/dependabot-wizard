import * as React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import Stack from '@mui/material/Stack';
import Alert from '@mui/material/Alert';
import { observer } from 'mobx-react-lite';

import { configStore } from '../store/configuration-store';
import { fetchPullRequests } from '../api/api-client';
import { palette } from '../utils/theme';

export const GithubRepoDetails: React.FC = observer(() => {
  const [error, setError] = React.useState('');
  const [loading, setLoading] = React.useState(false);

  const handleClose = () => {
    setError('');
    configStore.setShowGithubRepoDetails(false);
  };

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    setLoading(true);

    const formData = new FormData(event.currentTarget);
    const formJson = Object.fromEntries((formData as any).entries());

    try {
      const { pullRequests, filter } = await fetchPullRequests({
        repo: formJson.repoName,
        owner: formJson.repoOwner,
        token: formJson.repoToken,
        label: formJson.repoLabel,
      });

      configStore.setPullRequests(pullRequests, filter);
      configStore.setHasSavedRepoDetails(true);
      handleClose();
      configStore.setShowGithubInfo(true);
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : 'Could not load pull requests.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog
      open={configStore.showGithubRepoDetails}
      onClose={handleClose}
      maxWidth="xs"
      fullWidth
      PaperProps={{ component: 'form', onSubmit }}
    >
      <DialogTitle sx={{ pb: 1 }}>Repository details</DialogTitle>
      <DialogContent>
        <DialogContentText sx={{ fontSize: '0.85rem', color: palette.textMuted, mb: 2 }}>
          Pull requests opened by Dependabot are matched by author, so they are found whether or not
          they carry a label. Add a label only to narrow the results further.
        </DialogContentText>

        <Stack spacing={2}>
          <TextField autoFocus required name="repoOwner" label="Owner" fullWidth
            placeholder="kaypee90" helperText="The user or organisation" />
          <TextField required name="repoName" label="Repository" fullWidth placeholder="zebrok" />
          <TextField required name="repoToken" label="Token" type="password" fullWidth
            helperText="Needs Pull requests: Read-only. Held in memory only." />
          <TextField name="repoLabel" label="Label (optional)" fullWidth
            helperText="Leave blank to include every Dependabot pull request" />
          {error && <Alert severity="error">{error}</Alert>}
        </Stack>
      </DialogContent>
      <DialogActions sx={{ px: 3, pb: 2.5 }}>
        <Button onClick={handleClose} sx={{ color: palette.textMuted }}>Cancel</Button>
        <Button type="submit" variant="contained" disabled={loading}>
          {loading ? 'Loading…' : 'View pull requests'}
        </Button>
      </DialogActions>
    </Dialog>
  );
});
