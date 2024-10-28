import React from "react";
import FormControl from '@mui/material/FormControl';
import { observer } from "mobx-react-lite";
import { configStore } from '../store/configuration-store';
import { Button } from "@mui/material";
import Snackbar from '@mui/material/Snackbar';

const CreateYamlFile: React.FC = observer(() => {
  const [open, setOpen] = React.useState(false);

  const isValidData = () => {
    if (configStore.ecosystem && configStore.directory && configStore.schedule && configStore.openPrLimit) {
      return true;
    }

    return false;
  }

  const onSubmit = () => {
    setOpen(false);
    if (isValidData()) {
      alert(configStore.yamlData);

    } else {
      setOpen(true);
      setTimeout(() => setOpen(false), 3000)
    }

  }

  return (
    <FormControl fullWidth>
      <Button sx={{ mt: 2, mb: 6, mr: 10, ml: 10 }} variant="contained" color="primary" onClick={onSubmit}>Create Yaml File</Button>
      <Snackbar
        open={open}
        autoHideDuration={3000}
        message="Failed! provide all required fields"
      />
    </FormControl>
  )
});

export default CreateYamlFile;
