import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid2';

import PackageEcosystem from './package-ecosystem';
import ConfigDivider from './config-divider';
import ConfigYaml from './config-yaml';
import CreateYamlFile from './create-yaml-file';

export function MainConfiguration() {

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Grid container spacing={1} sx={{ height: '98vh' }}>
        <Grid sx={{ background: '#0E2535' }} size={{ xs: 6, md: 6 }}>
          <ConfigDivider />
          <PackageEcosystem />
          <ConfigDivider />
          <CreateYamlFile />
        </Grid>

        <Grid size={{ xs: 6, md: 6 }} sx={{ background: '#fff', color: '#000' }}>
          <ConfigYaml />
        </Grid>
      </Grid>
    </Box>
  );
}

