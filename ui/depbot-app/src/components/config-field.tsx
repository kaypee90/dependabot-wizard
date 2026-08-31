import React from 'react';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';

import { palette } from '../utils/theme';

type ConfigFieldProps = {
  label: string;
  description: string;
  required?: boolean;
  children: React.ReactNode;
};

/**
 * One labelled question. Replaces the raw <h5> tags in a 4/8 grid, which
 * inherited browser heading margins and gave every row a different rhythm.
 */
const ConfigField: React.FC<ConfigFieldProps> = ({ label, description, required, children }) => (
  <Box>
    <Stack direction="row" alignItems="baseline" spacing={0.75} sx={{ mb: 0.25 }}>
      <Typography variant="subtitle1" sx={{ color: palette.text }}>{label}</Typography>
      {required && (
        <Typography component="span" sx={{ color: palette.accent, fontSize: '0.8rem' }} aria-hidden>
          *
        </Typography>
      )}
    </Stack>
    <Typography variant="caption" sx={{ color: palette.textFaint, display: 'block', mb: 1 }}>
      {description}
    </Typography>
    {children}
  </Box>
);

export default ConfigField;
