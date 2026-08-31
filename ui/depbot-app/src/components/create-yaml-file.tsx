import React from "react";
import { observer } from "mobx-react-lite";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';

import { configStore } from '../store/configuration-store';
import { saveConfiguration } from "../api/api-client";
import { isStaticBuild } from "../utils/runtime";
import { palette } from '../utils/theme';

type Feedback = { message: string; severity: 'success' | 'error' | 'info' };

/**
 * Writes the configuration through the local API, or — in the static build,
 * where there is no depbot process to write anything — copies it instead.
 */
const CreateYamlFile: React.FC = observer(() => {
  const [feedback, setFeedback] = React.useState<Feedback | null>(null);
  const staticBuild = isStaticBuild();
  const ready = configStore.isValid;

  const copyToClipboard = async () => {
    try {
      if (!navigator.clipboard) {
        throw new Error('unavailable');
      }
      await navigator.clipboard.writeText(configStore.yamlData);
      setFeedback({ message: 'Copied. Save it as .github/dependabot.yml', severity: 'success' });
    } catch {
      setFeedback({ message: 'Could not copy automatically. Select the config and copy it manually.', severity: 'error' });
    }
  };

  const writeToDisk = async () => {
    try {
      const path = await saveConfiguration(configStore.yamlData);
      setFeedback({ message: `Wrote ${path}`, severity: 'success' });
    } catch (error) {
      setFeedback({
        message: error instanceof Error ? error.message : 'Could not create the file.',
        severity: 'error',
      });
    }
  };

  const onSubmit = async () => {
    setFeedback(null);

    if (!ready) {
      setFeedback({
        message: staticBuild
          ? 'Fill in the required fields first.'
          : 'Fill in the required fields, or scan the project first.',
        severity: 'info',
      });
      return;
    }

    await (staticBuild ? copyToClipboard() : writeToDisk());
  };

  return (
    <Box
      sx={{
        p: 2,
        border: `1px solid ${palette.border}`,
        borderRadius: 2,
        background: palette.surface,
      }}
    >
      <Stack direction="row" alignItems="center" spacing={2}>
        <Stack sx={{ flexGrow: 1 }}>
          <Typography variant="subtitle1">
            {staticBuild ? 'Copy the configuration' : 'Write the configuration'}
          </Typography>
          <Typography variant="caption" sx={{ color: palette.textFaint }}>
            {ready
              ? staticBuild
                ? 'Validated and ready to paste'
                : 'Validated, then written to .github/dependabot.yml'
              : 'Complete the required fields to continue'}
          </Typography>
        </Stack>
        <Button variant="contained" onClick={onSubmit} disabled={!ready}>
          {staticBuild ? 'Copy YAML' : 'Create Yaml File'}
        </Button>
      </Stack>

      <Snackbar
        open={feedback !== null}
        autoHideDuration={4000}
        onClose={() => setFeedback(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert
          severity={feedback?.severity ?? 'info'}
          variant="filled"
          onClose={() => setFeedback(null)}
          sx={{ width: '100%' }}
        >
          {feedback?.message}
        </Alert>
      </Snackbar>
    </Box>
  );
});

export default CreateYamlFile;
