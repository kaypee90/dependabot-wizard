import { useEffect, useState } from 'react';
import { LineChart } from '@mui/x-charts/LineChart';
import { PullRequest } from '../models/pull-request';
import { configStore } from '../store/configuration-store';


const PullRequestLineGraph = () => {
  const [data, setData] = useState<number[]>([]);
  const [xLabels, setXLabels] = useState<string[]>([]);

  useEffect(() => {
    if (configStore.pullRequests.length === 0) {
      return;
    }

    const prGroupedByDate = configStore.pullRequests.reduce((acc: Record<string, number>, pr: PullRequest) => {
      const createdDate = pr.createdAt.split('T')[0];
      if (!acc[createdDate]) {
        acc[createdDate] = 1;
      } else {
        acc[createdDate]++;
      }

      return acc;
    }, {} as Record<string, number>);

    const labels = Object.keys(prGroupedByDate);
    const values = Object.values(prGroupedByDate);
    setData(values);
    setXLabels(labels);
  }, []);

  return (
    <LineChart
      xAxis={[{ scaleType: 'point', data: xLabels }]}
      series={[
        {
          data: data,
          area: true,
          baseline: 'min',
          label: 'Pull Requests Created by Dependabot Over Time',
        },
      ]}
      height={300}
    />
  );
}

export default PullRequestLineGraph;
