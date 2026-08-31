import React from "react";
import FormControl from '@mui/material/FormControl';
import { observer } from "mobx-react-lite";
import { Button } from "@mui/material";
import Snackbar from '@mui/material/Snackbar';

import { configStore } from '../store/configuration-store';
import { saveConfiguration } from "../api/api-client";
import { isStaticBuild } from "../utils/runtime";

/**
 * Writes the configuration through the local API, or — in the static build,
 * where there is no depbot process to write anything — copies it instead.
 */
const CreateYamlFile: React.FC = observer(() => {
  const [open, setOpen] = React.useState(false);
  const [message, setMessage] = React.useState('');
  const staticBuild = isStaticBuild();

  const showMessage = (text: string) => {
    setMessage(text);
    setOpen(true);
  };

  const copyToClipboard = async () => {
    try {
      if (!navigator.clipboard) {
        throw new Error('unavailable');
      }
      await navigator.clipboard.writeText(configStore.yamlData);
      showMessage('Copied. Save it as .github/dependabot.yml');
    } catch {
      showMessage('Could not copy automatically. Select the config and copy it manually.');
    }
  };

  const writeToDisk = async () => {
    try {
      const path = await saveConfiguration(configStore.yamlData);
      showMessage(`Wrote ${path}`);
    } catch (error) {
      showMessage(error instanceof Error ? error.message : 'Could not create the file.');
    }
  };

  const onSubmit = async () => {
    setOpen(false);

    if (!configStore.isValid) {
      showMessage(staticBuild
        ? 'Fill in the required fields first.'
        : 'Fill in the required fields, or scan the project first.');
      return;
    }

    await (staticBuild ? copyToClipboard() : writeToDisk());
  };

  return (
    <FormControl fullWidth>
      <Button sx={{ mt: 2, mb: 6, mr: 10, ml: 10 }} variant="contained" color="primary" onClick={onSubmit}>
        {staticBuild ? 'Copy YAML' : 'Create Yaml File'}
      </Button>
      <Snackbar
        open={open}
        autoHideDuration={4000}
        onClose={() => setOpen(false)}
        message={message} />
    </FormControl>
  );
});

export default CreateYamlFile;
