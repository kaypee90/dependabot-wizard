import { render, screen } from '@testing-library/react';
import App from './App';

// REACT_APP_STATIC is injected at build time. These tests set it directly so
// both shapes of the app are covered from one source tree.
const withStaticBuild = (enabled: boolean) => {
  const previous = process.env.REACT_APP_STATIC;

  beforeEach(() => {
    process.env.REACT_APP_STATIC = enabled ? 'true' : '';
  });

  afterEach(() => {
    process.env.REACT_APP_STATIC = previous;
  });
};

describe('static build (published to Pages, no depbot process)', () => {
  withStaticBuild(true);

  it('offers to copy the config rather than write it', () => {
    render(<App />);

    expect(screen.getByRole('button', { name: /copy yaml/i })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /create yaml file/i })).not.toBeInTheDocument();
  });

  it('hides the pull request dashboard, which needs the local API', () => {
    render(<App />);

    expect(screen.queryByRole('button', { name: /pull requests/i })).not.toBeInTheDocument();
  });

  it('explains how to get scanning instead of offering a broken button', () => {
    render(<App />);

    expect(screen.queryByRole('button', { name: /scan project/i })).not.toBeInTheDocument();
    expect(screen.getByText(/go install github.com\/kaypee90\/dependabot-wizard/)).toBeInTheDocument();
  });
});

describe('served by the depbot binary', () => {
  withStaticBuild(false);

  it('writes the file and exposes scanning and the dashboard', () => {
    render(<App />);

    expect(screen.getByRole('button', { name: /create yaml file/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /scan project/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /pull requests/i })).toBeInTheDocument();
  });
});
