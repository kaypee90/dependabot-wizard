import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';

import { palette } from '../utils/theme';

type PullRequestCountProps = {
  title: string;
  count: number;
  /** Draws the eye to the number that needs action. */
  emphasis?: boolean;
};

const PullRequestCount: React.FC<PullRequestCountProps> = ({ title, count, emphasis }) => (
  <Box
    sx={{
      p: 2,
      height: '100%',
      borderRadius: 2,
      background: palette.surface,
      border: `1px solid ${emphasis && count > 0 ? palette.accent : palette.border}`,
    }}
  >
    <Typography
      variant="subtitle2"
      sx={{ color: palette.textFaint, mb: 0.5 }}
    >
      {title}
    </Typography>
    <Typography
      sx={{
        fontSize: '1.9rem',
        fontWeight: 700,
        lineHeight: 1.1,
        fontVariantNumeric: 'tabular-nums',
        color: emphasis && count > 0 ? palette.accent : palette.text,
      }}
    >
      {count}
    </Typography>
  </Box>
);

export default PullRequestCount;
