import React, { useState } from 'react';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';

import PackageEcosystem from './package-ecosystem';
import ConfigDivider from './config-divider';
import ConfigYaml from './config-yaml';

export function MainConfiguration() {

  return (
    <Box sx={{ flexGrow: 1 }}>
      <Grid sx={{ height: '98vh' }} container spacing={1}>
        <Grid sx={{ background: '#0E2535' }} item xs={6}>
          <div>
            <ConfigDivider />
            <PackageEcosystem />
            <ConfigDivider />
          </div>
        </Grid>

        <Grid item xs={6} sx={{ background: '#fff', color: '#000' }}>
          <ConfigYaml />
        </Grid>
      </Grid>
    </Box>
  );
}

