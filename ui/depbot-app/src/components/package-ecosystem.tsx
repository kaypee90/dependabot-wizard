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
import { generateYaml } from './config-generator';
import { observer } from 'mobx-react-lite';
import { configStore } from '../store/configuration-store';

const PackageEcosystem: React.FC = observer(() => {
  const [ecosystem, setEcosystem] = React.useState('');
  const [directory, setDirectory] = React.useState('');
  const [interval, setInterval] = React.useState('');
  const [reviewer, setReviewer] = React.useState('');
  const [openPrLimit, setOpenPrLimit] = React.useState('');
  const [configLabel, setConfigLabel] = React.useState('');

  const generateConfigYaml = React.useCallback(() => {
    configStore.setPackageConfiguration(ecosystem, directory, interval, reviewer, openPrLimit, configLabel);
    generateYaml();
  }, [ecosystem, directory, interval, reviewer, openPrLimit, configLabel]);

  const ecosystemHandleChange = (event: SelectChangeEvent) => {
    setEcosystem(event.target.value as string);
    generateConfigYaml();
  };

  const directoryHandleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setDirectory(event.target.value as string);
    generateConfigYaml();
  };

  const intervalHandleChange = (event: SelectChangeEvent) => {
    setInterval(event.target.value as string);
    generateConfigYaml();
  };

  const reviewerHandleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setReviewer(event.target.value as string);
    generateConfigYaml();
  };

  const openPrLimitHandleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setOpenPrLimit(event.target.value as string);
    generateConfigYaml();
  };

  const configLabelHandleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setConfigLabel(event.target.value as string);
    generateConfigYaml();
  };

  React.useEffect(() => {
    generateConfigYaml();
  }, [generateConfigYaml]);

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
            <h5>Package Ecosystem *</h5>
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
            <h5>Directory *</h5>
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


        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Interval (Schedule) *</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>How often should Dependabot check for updates</h5>

              <FormControl fullWidth>
                <InputLabel id="interval-label">Interval</InputLabel>
                <Select
                  labelId="interval-label"
                  id="interval"
                  value={interval}
                  label="Interval"
                  onChange={intervalHandleChange}
                >
                  <MenuItem value={'daily'}>Daily</MenuItem>
                  <MenuItem value={'weekly'}>Weekly</MenuItem>
                  <MenuItem value={'monthly'}>Monthly</MenuItem>
                </Select>
              </FormControl>
            </div>
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Reviewer</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>Enter GitHub username to add as reviewer</h5>

              <FormControl fullWidth>
                <TextField label="Reviewer" variant="outlined" value={reviewer} onChange={reviewerHandleChange} />
              </FormControl>
            </div>
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Open Pull Requests Limit *</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>Set the maximum number of pull requests to be opened</h5>

              <FormControl fullWidth>
                <TextField type="number" label="Opened PR Limits" variant="outlined" value={openPrLimit} onChange={openPrLimitHandleChange} />
              </FormControl>
            </div>
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Label</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>Label to set on pull requests</h5>

              <FormControl fullWidth>
                <TextField label="Label" variant="outlined" value={configLabel} onChange={configLabelHandleChange} />
              </FormControl>
            </div>
          </Grid>
        </Grid>

      </AccordionDetails>
    </Accordion>
  );
});

export default PackageEcosystem;

