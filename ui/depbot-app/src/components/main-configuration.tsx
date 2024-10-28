import React, { useState } from 'react';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';

import PackageEcosystem from './package-ecosystem';
import ConfigDivider from './config-divider';
import ConfigYaml from './config-yaml';

export function MainConfiguration() {
  const yamlData = `
  version: 2
  updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    reviewers:
    - kaypee90
    open-pull-requests-limit: 5
    labels:
    - go mod`;

  const [yamlCode, setYamlCode] = useState<string>(yamlData);


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
          <ConfigYaml code={yamlCode} />
        </Grid>
      </Grid>
    </Box>
  );
}

