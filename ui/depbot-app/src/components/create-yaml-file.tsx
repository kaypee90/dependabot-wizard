import React from "react";
import FormControl from '@mui/material/FormControl';
import { observer } from "mobx-react-lite";
import { Button } from "@mui/material";
import Snackbar from '@mui/material/Snackbar';

import { configStore } from '../store/configuration-store';
import { saveConfiguration } from "../api/api-client";

const CreateYamlFile: React.FC = observer(() => {
  const [open, setOpen] = React.useState(false);
  const [message, setMessage] = React.useState('');

  const showMessage = (text: string) => {
    setMessage(text);
    setOpen(true);
  };

  const onSubmit = async () => {
    setOpen(false);

    if (!configStore.isValid) {
      showMessage('Fill in the required fields, or scan the project first.');
      return;
    }

    try {
      const path = await saveConfiguration(configStore.yamlData);
      showMessage(`Wrote ${path}`);
    } catch (error) {
      showMessage(error instanceof Error ? error.message : 'Could not create the file.');
    }
  };

  return (
    <FormControl fullWidth>
      <Button sx={{ mt: 2, mb: 6, mr: 10, ml: 10 }} variant="contained" color="primary" onClick={onSubmit}>
        Create Yaml File
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
