import * as React from 'react';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import Typography from '@mui/material/Typography';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Grid from '@mui/material/Grid';

import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import TextField from '@mui/material/TextField';

export default function PackageEcosystem() {
  const [ecosystem, setEcosystem] = React.useState('');
  const [directory, setDirectory] = React.useState('');

  const ecosystemHandleChange = (event: SelectChangeEvent) => {
    setEcosystem(event.target.value as string);
  };

  const directoryHandleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setDirectory(event.target.value);
  };

  return (

    <Accordion sx={{ background: '#0E2535', color: '#fff' }}>
      <AccordionSummary
        expandIcon={<ExpandMoreIcon />}
        aria-controls="panel1-content"
        id="panel1-header"
      >
        <Typography>Package Configuration</Typography>
      </AccordionSummary>
      <AccordionDetails>

        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Package Ecosystem</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>Select package manager to use</h5>

              <FormControl fullWidth>
                <InputLabel id="ecosystem-label">Ecosystem</InputLabel>
                <Select
                  labelId="ecosystem-label"
                  id="ecosystem"
                  value={ecosystem}
                  label="Ecosystem"
                  onChange={ecosystemHandleChange}
                >
                  <MenuItem value={'cargo'}>Cargo</MenuItem>
                  <MenuItem value={'docker'}>Docker</MenuItem>
                  <MenuItem value={'gomod'}>Go Mod</MenuItem>
                </Select>
              </FormControl>
            </div>
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Directory</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>Enter location of package manifest</h5>

              <FormControl fullWidth>
                <TextField label="Directory" variant="outlined" value={directory} onChange={directoryHandleChange} />
              </FormControl>
            </div>
          </Grid>
        </Grid>


      </AccordionDetails>
    </Accordion>
  );
}

