import React from "react";
import FormControl from '@mui/material/FormControl';
import { observer } from "mobx-react-lite";
import { configStore } from '../store/configuration-store';
import { Button } from "@mui/material";
import Snackbar from '@mui/material/Snackbar';

const CreateYamlFile: React.FC = observer(() => {
  const [open, setOpen] = React.useState(false);
  const [message, setMessage] = React.useState('');

  const isValidData = () => {
    if (configStore.ecosystem && configStore.directory && configStore.schedule && configStore.openPrLimit) {
      return true;
    }

    return false;
  }

  const showMessage = (message: string) => {
    setMessage(message);
    setOpen(true);
    setTimeout(() => setOpen(false), 3000);

  }

  const handlePostRequest = async (data: string) => {
    try {
      const response = await fetch('http://localhost:3001/api/configurations', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ configuration: data })
      });

      if (!response.ok) {
        throw new Error(`Error! status: ${response.status}`);
      }

      showMessage('Yaml file created successfully');
    } catch (error) {
      console.error('Error:', error);
      showMessage('An error occurred while creating the yaml file, please try again');
    }
  };


  const onSubmit = () => {
    setOpen(false);
    if (isValidData()) {
      handlePostRequest(configStore.yamlData);

    } else {
      showMessage('Please fill out all required fields');
    }

  }

  return (
    <FormControl fullWidth>
      <Button sx={{ mt: 2, mb: 6, mr: 10, ml: 10 }} variant="contained" color="primary" onClick={onSubmit}>Create Yaml File</Button>
      <Snackbar
        open={open}
        autoHideDuration={3000}
        message={message} />
    </FormControl>
  )
});

export default CreateYamlFile;
