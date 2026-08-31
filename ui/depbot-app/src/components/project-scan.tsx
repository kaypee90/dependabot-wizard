import React from 'react';
import { observer } from 'mobx-react-lite';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import RadarIcon from '@mui/icons-material/TrackChanges';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Alert from '@mui/material/Alert';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';
import Box from '@mui/material/Box';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';
import CircularProgress from '@mui/material/CircularProgress';

import { configStore } from '../store/configuration-store';
import { fetchDetectedEcosystems } from '../api/api-client';
import { generateYaml } from './config-generator';
import { isStaticBuild, CLI_INSTALL_COMMAND } from '../utils/runtime';
import { palette, MONO } from '../utils/theme';

const CommandBlock = ({ children }: { children: React.ReactNode }) => (
  <Box
    component="pre"
    sx={{
      m: 0, p: 1.5, borderRadius: 1.5, overflowX: 'auto',
      background: '#0A1B27', border: `1px solid ${palette.borderSoft}`,
      fontSize: '0.76rem', color: '#C6D6E0',
    }}
  >
    {children}
  </Box>
);

/**
 * Scans the project depbot was started in and offers the ecosystems it found.
 * This is the path that produces a complete multi-ecosystem config; the manual
 * form below it still handles a single block by hand.
 *
 * Scanning reads the filesystem, so it only exists when a depbot binary is
 * serving the page. The static build explains that instead of failing.
 */
const ProjectScan: React.FC = observer(() => {
  const staticBuild = isStaticBuild();

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
  const scanning = scanState === 'scanning';

  return (
    <Accordion defaultExpanded>
      <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="scan-content" id="scan-header">
        <Stack direction="row" spacing={1.5} alignItems="center" sx={{ width: '100%' }}>
          <RadarIcon sx={{ fontSize: 20, color: palette.accent }} />
          <Stack sx={{ flexGrow: 1 }}>
            <Typography variant="h6">Scan This Project</Typography>
            <Typography variant="caption" sx={{ color: palette.textFaint }}>
              {detected.length > 0
                ? `${detected.length} ecosystems found`
                : 'Find every package manifest automatically'}
            </Typography>
          </Stack>
          {detected.length > 0 && (
            <Chip
              size="small"
              label={useDetected ? 'In use' : 'Off'}
              sx={{
                mr: 1,
                background: useDetected ? palette.accentSoft : 'transparent',
                color: useDetected ? palette.accent : palette.textFaint,
                border: `1px solid ${useDetected ? palette.accent : palette.border}`,
              }}
            />
          )}
        </Stack>
      </AccordionSummary>
      <AccordionDetails>
        {staticBuild ? (
          <Stack spacing={2}>
            <Typography variant="body2" sx={{ color: palette.textMuted }}>
              Scanning reads the files in your project, which a web page cannot do. Run depbot
              locally and it finds every manifest and builds the whole configuration for you.
            </Typography>
            <CommandBlock>{CLI_INSTALL_COMMAND}</CommandBlock>
            <Typography variant="body2" sx={{ color: palette.textFaint }}>
              Meanwhile, fill in the form below and copy the result.
            </Typography>
          </Stack>
        ) : (
          <Stack spacing={2}>
            <Typography variant="body2" sx={{ color: palette.textMuted }}>
              Looks for package manifests in the directory depbot was started in, skipping
              <code> node_modules</code> and <code> vendor</code>, and builds one update block per
              ecosystem at the right directory.
            </Typography>

            <Box>
              <Button
                variant="contained"
                onClick={scan}
                disabled={scanning}
                startIcon={scanning ? <CircularProgress size={15} color="inherit" /> : undefined}
              >
                {scanning ? 'Scanning' : detected.length > 0 ? 'Scan again' : 'Scan project'}
              </Button>
            </Box>

            {scanState === 'failed' && <Alert severity="error">{scanError}</Alert>}

            {scanState === 'done' && detected.length === 0 && (
              <Alert severity="info">
                No package manifests found here. Fill in the form below instead.
              </Alert>
            )}

            {detected.length > 0 && (
              <Stack spacing={1.5}>
                <Stack direction="row" spacing={1} useFlexGap flexWrap="wrap">
                  {detected.map((item) => (
                    <Chip
                      key={`${item.ecosystem}:${item.directory}`}
                      label={`${item.ecosystem} · ${item.directory}`}
                      title={`Found ${item.evidence}`}
                      size="small"
                      sx={{
                        background: palette.raised,
                        color: palette.text,
                        border: `1px solid ${palette.border}`,
                        fontFamily: MONO,
                      }}
                    />
                  ))}
                </Stack>
                <FormControlLabel
                  control={<Switch checked={useDetected} onChange={toggleUseDetected} />}
                  label={
                    <Typography variant="body2">
                      Use these {detected.length} ecosystems
                    </Typography>
                  }
                />
              </Stack>
            )}
          </Stack>
        )}
      </AccordionDetails>
    </Accordion>
  );
});

export default ProjectScan;
