import * as React from 'react';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import Typography from '@mui/material/Typography';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Stack from '@mui/material/Stack';
import Alert from '@mui/material/Alert';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import TextField from '@mui/material/TextField';
import { observer } from 'mobx-react-lite';

import { generateYaml } from './config-generator';
import { configStore } from '../store/configuration-store';
import ConfigField from './config-field';
import { palette } from '../utils/theme';

const ECOSYSTEMS: Array<[string, string]> = [
  ['bundler', 'Bundler'], ['cargo', 'Cargo'], ['composer', 'Composer'],
  ['devcontainers', 'Dev Containers'], ['docker', 'Docker'], ['elm', 'Elm'],
  ['github-actions', 'GitHub Actions'], ['gitsubmodule', 'Git Submodule'],
  ['gomod', 'Go Mod'], ['gradle', 'Gradle'], ['maven', 'Maven'], ['mix', 'Mix'],
  ['npm', 'NPM'], ['nuget', 'Nuget'], ['pip', 'Pip'], ['poetry', 'Poetry'],
  ['pub', 'Pub'], ['swift', 'Swift'], ['terraform', 'Terraform'],
];

const PackageEcosystem: React.FC = observer(() => {
  const [ecosystem, setEcosystem] = React.useState('');
  const [directory, setDirectory] = React.useState('/');
  const [interval, setInterval] = React.useState('weekly');
  const [reviewer, setReviewer] = React.useState('');
  const [openPrLimit, setOpenPrLimit] = React.useState('5');
  const [configLabel, setConfigLabel] = React.useState('');

  // A single effect owns regeneration. The handlers below only set state:
  // calling the generator from a handler read the value being replaced and
  // wrote a stale preview that this effect then overwrote.
  React.useEffect(() => {
    configStore.setPackageConfiguration(ecosystem, directory, interval, reviewer, openPrLimit, configLabel);
    generateYaml();
  }, [ecosystem, directory, interval, reviewer, openPrLimit, configLabel]);

  const usingDetected = configStore.useDetected && configStore.detected.length > 0;
  const directoryInvalid = directory.trim() !== '' && !directory.trim().startsWith('/');

  return (
    <Accordion defaultExpanded>
      <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="config-content" id="config-header">
        <Stack>
          <Typography variant="h6">Package Configuration</Typography>
          <Typography variant="caption" sx={{ color: palette.textFaint }}>
            {usingDetected
              ? `Applied to all ${configStore.detected.length} scanned ecosystems`
              : 'Describe one package manager by hand'}
          </Typography>
        </Stack>
      </AccordionSummary>
      <AccordionDetails>
        <Stack spacing={2.5}>
          {usingDetected && (
            <Alert severity="info">
              Ecosystem and directory come from the scan. The schedule, reviewer, limit and label
              below apply to every one of them.
            </Alert>
          )}

          <ConfigField
            label="Package ecosystem"
            description="The package manager Dependabot should watch"
            required={!usingDetected}
          >
            <FormControl fullWidth disabled={usingDetected}>
              <InputLabel id="ecosystem-label">Ecosystem</InputLabel>
              <Select
                labelId="ecosystem-label"
                id="ecosystem"
                value={ecosystem}
                label="Ecosystem"
                onChange={(event: SelectChangeEvent) => setEcosystem(event.target.value)}
              >
                {ECOSYSTEMS.map(([value, display]) => (
                  <MenuItem key={value} value={value}>{display}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </ConfigField>

          <ConfigField
            label="Directory"
            description="Where the manifest lives, relative to the repository root"
            required={!usingDetected}
          >
            <TextField
              fullWidth
              label="Directory"
              value={directory}
              disabled={usingDetected}
              error={directoryInvalid}
              helperText={directoryInvalid ? 'Must start with /' : ' '}
              onChange={(event) => setDirectory(event.target.value)}
            />
          </ConfigField>

          <ConfigField label="Schedule" description="How often Dependabot checks for updates" required>
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
          </ConfigField>

          <ConfigField label="Open pull request limit" description="Caps how many updates are in flight at once">
            <TextField
              fullWidth
              type="number"
              label="Limit"
              value={openPrLimit}
              onChange={(event) => setOpenPrLimit(event.target.value)}
            />
          </ConfigField>

          <ConfigField label="Reviewer" description="GitHub username added to every update (optional)">
            <TextField fullWidth label="Reviewer" value={reviewer}
              onChange={(event) => setReviewer(event.target.value)} />
          </ConfigField>

          <ConfigField label="Label" description="Applied to every pull request Dependabot opens (optional)">
            <TextField fullWidth label="Label" value={configLabel}
              onChange={(event) => setConfigLabel(event.target.value)} />
          </ConfigField>
        </Stack>
      </AccordionDetails>
    </Accordion>
  );
});

export default PackageEcosystem;
