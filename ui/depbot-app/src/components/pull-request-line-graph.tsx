import { LineChart } from '@mui/x-charts/LineChart';
import { observer } from 'mobx-react-lite';

import { PullRequest } from '../models/pull-request';
import { configStore } from '../store/configuration-store';
import { palette } from '../utils/theme';

const PullRequestLineGraph = observer(() => {
  const grouped = configStore.pullRequests.reduce((acc: Record<string, number>, pr: PullRequest) => {
    const createdDate = pr.createdAt.split('T')[0];
    acc[createdDate] = (acc[createdDate] ?? 0) + 1;
    return acc;
  }, {});

  // Sort by date so the x-axis reads left to right in time order.
  const labels = Object.keys(grouped).sort();
  const values = labels.map((date) => grouped[date]);

  return (
    <LineChart
      xAxis={[{ scaleType: 'point', data: labels }]}
      series={[{
        data: values,
        area: true,
        baseline: 'min',
        label: 'Pull requests opened',
        color: palette.accent,
        showMark: labels.length <= 30,
      }]}
      height={280}
      margin={{ left: 44, right: 20, top: 30, bottom: 34 }}
      sx={{
        '& .MuiAreaElement-root': { fillOpacity: 0.16 },
        '& .MuiChartsAxis-line, & .MuiChartsAxis-tick': { stroke: palette.borderSoft },
        '& .MuiChartsAxis-tickLabel': { fill: palette.textFaint, fontSize: 11 },
      }}
    />
  );
});

export default PullRequestLineGraph;
