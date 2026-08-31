import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';

import PackageEcosystem from './package-ecosystem';
import ProjectScan from './project-scan';
import NoiseControl from './noise-control';
import ConfigYaml from './config-yaml';
import CreateYamlFile from './create-yaml-file';

/**
 * Two columns: the questions on the left, the file they produce on the right.
 * The preview is sticky so it stays in view while the form is filled in.
 */
export function MainConfiguration() {
  return (
    <Box
      sx={{
        display: 'grid',
        gridTemplateColumns: { xs: '1fr', md: 'minmax(0, 1fr) minmax(0, 1fr)' },
        gap: { xs: 2, md: 3 },
        p: { xs: 2, md: 3 },
        alignItems: 'start',
        maxWidth: 1440,
        mx: 'auto',
      }}
    >
      <Stack spacing={2}>
        <ProjectScan />
        <PackageEcosystem />
        <NoiseControl />
        <CreateYamlFile />
      </Stack>

      <Box sx={{ position: { md: 'sticky' }, top: { md: 84 } }}>
        <ConfigYaml />
      </Box>
    </Box>
  );
}
