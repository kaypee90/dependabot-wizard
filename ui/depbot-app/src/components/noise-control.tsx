import React from 'react';
import { observer } from 'mobx-react-lite';
import Accordion from '@mui/material/Accordion';
import AccordionSummary from '@mui/material/AccordionSummary';
import AccordionDetails from '@mui/material/AccordionDetails';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import FilterListIcon from '@mui/icons-material/FilterList';
import Typography from '@mui/material/Typography';
import Stack from '@mui/material/Stack';
import Chip from '@mui/material/Chip';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';

import { configStore } from '../store/configuration-store';
import { generateYaml, GROUPING_LABELS, GroupingStrategy } from './config-generator';
import ConfigField from './config-field';
import { palette } from '../utils/theme';

const STRATEGIES: GroupingStrategy[] = ['none', 'all', 'by-dependency-type'];

/**
 * Grouping and ignore rules. These are the two settings that decide whether a
 * repository ends up with a reviewable stream of updates or a queue nobody
 * touches, and they are exactly what the health report recommends — so the
 * diagnosis and the fix live in the same tool.
 */
const NoiseControl: React.FC = observer(() => {
  const { grouping, ignoreMajor } = configStore;

  const changeGrouping = (event: SelectChangeEvent) => {
    configStore.setGrouping(event.target.value as GroupingStrategy);
    generateYaml();
  };

  const toggleIgnoreMajor = (_event: React.ChangeEvent<HTMLInputElement>, checked: boolean) => {
    configStore.setIgnoreMajor(checked);
    generateYaml();
  };

  const active = (grouping !== 'none' ? 1 : 0) + (ignoreMajor ? 1 : 0);

  return (
    <Accordion>
      <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="noise-content" id="noise-header">
        <Stack direction="row" spacing={1.5} alignItems="center" sx={{ width: '100%' }}>
          <FilterListIcon sx={{ fontSize: 20, color: palette.accent }} />
          <Stack sx={{ flexGrow: 1 }}>
            <Typography variant="h6">Pull Request Volume</Typography>
            <Typography variant="caption" sx={{ color: palette.textFaint }}>
              Batch updates and skip the ones you would close anyway
            </Typography>
          </Stack>
          {active > 0 && (
            <Chip
              size="small"
              label={`${active} active`}
              sx={{ mr: 1, background: palette.accentSoft, color: palette.accent, border: `1px solid ${palette.accent}` }}
            />
          )}
        </Stack>
      </AccordionSummary>
      <AccordionDetails>
        <Stack spacing={2.5}>
          <ConfigField
            label="Grouping"
            description="Without grouping every dependency gets its own pull request, which is what buries review queues"
          >
            <FormControl fullWidth>
              <InputLabel id="grouping-label">Grouping</InputLabel>
              <Select
                labelId="grouping-label"
                id="grouping"
                value={grouping}
                label="Grouping"
                onChange={changeGrouping}
              >
                {STRATEGIES.map((strategy) => (
                  <MenuItem key={strategy} value={strategy}>{GROUPING_LABELS[strategy]}</MenuItem>
                ))}
              </Select>
            </FormControl>
          </ConfigField>

          <ConfigField
            label="Major versions"
            description="Major bumps are the ones most often closed unmerged; skipping them keeps the queue to updates you will actually take"
          >
            <FormControlLabel
              control={<Switch checked={ignoreMajor} onChange={toggleIgnoreMajor} />}
              label={
                <Typography variant="body2">
                  {ignoreMajor
                    ? 'Skipping major version updates'
                    : 'Opening pull requests for major versions'}
                </Typography>
              }
            />
          </ConfigField>
        </Stack>
      </AccordionDetails>
    </Accordion>
  );
});

export default NoiseControl;
