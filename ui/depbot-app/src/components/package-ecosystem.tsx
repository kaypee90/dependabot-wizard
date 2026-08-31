import * as React from 'react';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import Typography from '@mui/material/Typography';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Grid from '@mui/material/Grid';
import Alert from '@mui/material/Alert';

import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import TextField from '@mui/material/TextField';
import { generateYaml } from './config-generator';
import { observer } from 'mobx-react-lite';
import { configStore } from '../store/configuration-store';

const ECOSYSTEMS: Array<[string, string]> = [
  ['bundler', 'Bundler'],
  ['cargo', 'Cargo'],
  ['composer', 'Composer'],
  ['devcontainers', 'Dev Containers'],
  ['docker', 'Docker'],
  ['elm', 'Elm'],
  ['github-actions', 'GitHub Actions'],
  ['gitsubmodule', 'Git Submodule'],
  ['gomod', 'Go Mod'],
  ['gradle', 'Gradle'],
  ['maven', 'Maven'],
  ['mix', 'Mix'],
  ['npm', 'NPM'],
  ['nuget', 'Nuget'],
  ['pip', 'Pip'],
  ['poetry', 'Poetry'],
  ['pub', 'Pub'],
  ['swift', 'Swift'],
  ['terraform', 'Terraform'],
];

const PackageEcosystem: React.FC = observer(() => {
  const [ecosystem, setEcosystem] = React.useState('');
  const [directory, setDirectory] = React.useState('/');
  const [interval, setInterval] = React.useState('');
  const [reviewer, setReviewer] = React.useState('');
  const [openPrLimit, setOpenPrLimit] = React.useState('');
  const [configLabel, setConfigLabel] = React.useState('');

  // A single effect owns regeneration. The handlers below only set state:
  // calling the generator from a handler read the value being replaced and
  // wrote a stale preview that this effect then overwrote.
  React.useEffect(() => {
    configStore.setPackageConfiguration(ecosystem, directory, interval, reviewer, openPrLimit, configLabel);
    generateYaml();
  }, [ecosystem, directory, interval, reviewer, openPrLimit, configLabel]);

  const usingDetected = configStore.useDetected && configStore.detected.length > 0;

  return (
    <Accordion sx={{ background: '#0E2535', color: '#fff' }}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel1-content" id="panel1-header">
        <Typography>Package Configuration</Typography>
      </AccordionSummary>
      <AccordionDetails>

        {usingDetected && (
          <Alert severity="info" sx={{ mb: 2 }}>
            Using the {configStore.detected.length} scanned ecosystems. Interval, reviewer, limit and
            label below still apply to all of them; ecosystem and directory are taken from the scan.
          </Alert>
        )}

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
                  disabled={usingDetected}
                  onChange={(event: SelectChangeEvent) => setEcosystem(event.target.value)}
                >
                  {ECOSYSTEMS.map(([value, display]) => (
                    <MenuItem key={value} value={value}>{display}</MenuItem>
                  ))}
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
                <TextField
                  label="Directory"
                  variant="outlined"
                  value={directory}
                  disabled={usingDetected}
                  error={directory.trim() !== '' && !directory.trim().startsWith('/')}
                  helperText={directory.trim() !== '' && !directory.trim().startsWith('/') ? 'Must start with /' : ' '}
                  onChange={(event) => setDirectory(event.target.value)}
                />
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
                  onChange={(event: SelectChangeEvent) => setInterval(event.target.value)}
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
                <TextField label="Reviewer" variant="outlined" value={reviewer}
                  onChange={(event) => setReviewer(event.target.value)} />
              </FormControl>
            </div>
          </Grid>
        </Grid>

        <Grid container spacing={3}>
          <Grid item xs={4} sx={{ textAlign: 'left' }}>
            <h5>Open Pull Requests Limit</h5>
          </Grid>
          <Grid item xs={8}>
            <div>
              <h5>Set the maximum number of pull requests to be opened</h5>
              <FormControl fullWidth>
                <TextField type="number" label="Opened PR Limits" variant="outlined" value={openPrLimit}
                  onChange={(event) => setOpenPrLimit(event.target.value)} />
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
                <TextField label="Label" variant="outlined" value={configLabel}
                  onChange={(event) => setConfigLabel(event.target.value)} />
              </FormControl>
            </div>
          </Grid>
        </Grid>

      </AccordionDetails>
    </Accordion>
  );
});

export default PackageEcosystem;
