import React from 'react';
import { observer } from 'mobx-react-lite';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Alert from '@mui/material/Alert';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';

import { configStore } from '../store/configuration-store';
import { fetchDetectedEcosystems } from '../api/api-client';
import { generateYaml } from './config-generator';

/**
 * Scans the project depbot was started in and offers the ecosystems it found.
 * This is the path that produces a complete multi-ecosystem config; the manual
 * form below it still handles a single block by hand.
 */
const ProjectScan: React.FC = observer(() => {
  const scan = React.useCallback(async () => {
    configStore.setScanState('scanning');
    try {
      const detected = await fetchDetectedEcosystems();
      configStore.setDetected(detected);
      configStore.setScanState('done');
      generateYaml();
    } catch (error) {
      configStore.setScanState('failed', error instanceof Error ? error.message : 'The scan failed.');
    }
  }, []);

  const toggleUseDetected = (_event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => {
    configStore.setUseDetected(checked);
    generateYaml();
  };

  const { scanState, scanError, detected, useDetected } = configStore;

  return (
    <Accordion defaultExpanded sx={{ background: '#0E2535', color: '#fff' }}>
      <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="scan-content" id="scan-header">
        <Typography>Scan This Project</Typography>
      </AccordionSummary>
      <AccordionDetails>
        <Typography variant="body2" sx={{ mb: 2, color: '#B7C4CD' }}>
          Look for package manifests in the directory depbot was started in and build a
          configuration covering every ecosystem found.
        </Typography>

        <Button variant="outlined" color="inherit" onClick={scan} disabled={scanState === 'scanning'}>
          {scanState === 'scanning' ? 'Scanning…' : 'Scan project'}
        </Button>

        {scanState === 'failed' && (
          <Alert severity="error" sx={{ mt: 2 }}>{scanError}</Alert>
        )}

        {scanState === 'done' && detected.length === 0 && (
          <Alert severity="info" sx={{ mt: 2 }}>
            No package manifests found here. Fill in the form below instead.
          </Alert>
        )}

        {detected.length > 0 && (
          <>
            <Stack direction="row" spacing={1} useFlexGap flexWrap="wrap" sx={{ mt: 2 }}>
              {detected.map((item) => (
                <Chip
                  key={`${item.ecosystem}:${item.directory}`}
                  label={`${item.ecosystem} · ${item.directory}`}
                  title={`Found ${item.evidence}`}
                  size="small"
                  sx={{ background: '#1B3A50', color: '#fff' }}
                />
              ))}
            </Stack>
            <FormControlLabel
              sx={{ mt: 1 }}
              control={<Switch checked={useDetected} onChange={toggleUseDetected} />}
              label={`Use these ${detected.length} ecosystems`}
            />
          </>
        )}
      </AccordionDetails>
    </Accordion>
  );
});

export default ProjectScan;
