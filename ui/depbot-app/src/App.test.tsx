import { render, screen } from '@testing-library/react';
import App from './App';

test('renders the app bar title', () => {
  render(<App />);

  // Exact match: several panels mention depbot in their body copy.
  expect(screen.getByText('Depbot')).toBeInTheDocument();
});

test('offers both the project scan and the manual form', () => {
  render(<App />);

  expect(screen.getByText('Scan This Project')).toBeInTheDocument();
  expect(screen.getByText('Package Configuration')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: /create yaml file/i })).toBeInTheDocument();
});
