import * as React from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import { observer } from 'mobx-react-lite';
import { configStore } from '../store/configuration-store';
import { fetchPullRequests } from '../api/api-client';
import { ThemeProvider, createTheme } from '@mui/material/styles';

export const GithubRepoDetails: React.FC = observer(() => {
  const handleClose = () => {
    configStore.setShowGithubRepoDetails(false);
  };

  const lightTheme = createTheme({
    palette: {
      mode: 'light',
    },
  });

  return (
    <ThemeProvider theme={lightTheme}>
      <React.Fragment>
        <Dialog
          open={configStore.showGithubRepoDetails}
          onClose={handleClose}
          PaperProps={{
            component: 'form',
            onSubmit: (event: React.FormEvent<HTMLFormElement>) => {
              event.preventDefault();
              const formData = new FormData(event.currentTarget);
              const formJson = Object.fromEntries((formData as any).entries());

              const data = {
                repo: formJson.repoName,
                owner: formJson.repoOwner,
                token: formJson.repoToken
              }

              fetchPullRequests(data)
                .then((data) => {
                  configStore.setPullRequests(data);
                  configStore.setHasSavedRepoDetails(true);
                  handleClose();

                  configStore.setShowGithubInfo(true);
                })
                .catch((error) => {
                  console.error("Error fetching data:", error);
                });

            },
          }}
        >
          <DialogTitle>Repository Details</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Provide the repository details.
            </DialogContentText>
            <TextField
              autoFocus
              required
              margin="dense"
              id="name"
              name="repoName"
              label="Name"
              type="text"
              fullWidth
              variant="standard"
            />
            <TextField
              autoFocus
              required
              margin="dense"
              id="owner"
              name="repoOwner"
              label="Owner"
              type="text"
              fullWidth
              variant="standard"
            />
            <TextField
              autoFocus
              required
              margin="dense"
              id="token"
              name="repoToken"
              label="Token"
              type="text"
              fullWidth
              variant="standard"
            />

          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Cancel</Button>
            <Button type="submit">Submit</Button>
          </DialogActions>
        </Dialog>
      </React.Fragment>
    </ThemeProvider>
  );
});

