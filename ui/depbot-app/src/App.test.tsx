import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders depbot header', () => {
  render(<App />);
  const header = screen.getByText(/Depbot/i);
  expect(header).toBeInTheDocument();
});
