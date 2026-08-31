import { LineChart } from '@mui/x-charts/LineChart';
import { observer } from 'mobx-react-lite';
import { PullRequest } from '../models/pull-request';
import { configStore } from '../store/configuration-store';

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
        label: 'Pull Requests Created by Dependabot Over Time',
      }]}
      height={300}
    />
  );
});

export default PullRequestLineGraph;
