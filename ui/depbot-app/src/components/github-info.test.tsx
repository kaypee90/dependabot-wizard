import { render, screen } from '@testing-library/react';
import { ThemeProvider } from '@mui/material/styles';

import { GithubInfo } from './github-info';
import { configStore } from '../store/configuration-store';
import { appTheme } from '../utils/theme';
import { PullRequest } from '../models/pull-request';

const pr = (number: number, state: string, createdAt: string): PullRequest => ({
  number,
  state,
  title: `Bump dependency ${number}`,
  url: `https://example.test/${number}`,
  createdAt,
  author: { login: 'dependabot' },
});

const show = () => render(<ThemeProvider theme={appTheme}><GithubInfo /></ThemeProvider>);

describe('the pull request dashboard', () => {
  afterEach(() => {
    configStore.setPullRequests([]);
    configStore.setShowGithubInfo(false);
  });

  it('renders counts, chart and grid without throwing', () => {
    configStore.setPullRequests([
      pr(1, 'OPEN', '2026-08-01T10:00:00Z'),
      pr(2, 'MERGED', '2026-08-02T10:00:00Z'),
      pr(3, 'MERGED', '2026-08-03T10:00:00Z'),
      pr(4, 'CLOSED', '2026-08-04T10:00:00Z'),
    ], 'opened by app/dependabot');
    configStore.setShowGithubInfo(true);

    show();

    expect(screen.getByText('Dependabot pull requests')).toBeInTheDocument();
    expect(screen.getByText('Total')).toBeInTheDocument();
    expect(screen.getByText('Merge rate %')).toBeInTheDocument();
    // 2 merged of 3 resolved
    expect(screen.getByText('67')).toBeInTheDocument();
    expect(screen.getByText('Opened over time')).toBeInTheDocument();
    expect(screen.getByText('All pull requests')).toBeInTheDocument();
  });

  // The case that prompted the author-filter fix: nothing matched, and the
  // dashboard used to show four zeroes with no explanation.
  it('explains an empty result instead of showing zeroes', () => {
    configStore.setPullRequests([], 'opened by app/dependabot');
    configStore.setShowGithubInfo(true);

    show();

    expect(screen.getByText(/No pull requests opened by app\/dependabot were found/)).toBeInTheDocument();
    expect(screen.queryByText('Overview')).not.toBeInTheDocument();
  });
});
