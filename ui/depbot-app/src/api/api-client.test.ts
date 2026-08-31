import { fetchPullRequests } from './api-client';

const mockFetch = (body: unknown, ok = true, status = 200) => {
  global.fetch = jest.fn().mockResolvedValue({
    ok,
    status,
    json: async () => body,
  }) as unknown as typeof fetch;
};

describe('fetchPullRequests', () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it('reads the flat list the server returns', async () => {
    mockFetch({
      pullRequests: [{ number: 12, state: 'OPEN' }, { number: 14, state: 'MERGED' }],
      filter: 'opened by app/dependabot',
    });

    const { pullRequests, filter } = await fetchPullRequests({ owner: 'o', repo: 'r', token: 't' });

    expect(pullRequests).toHaveLength(2);
    expect(pullRequests[0].number).toBe(12);
    expect(filter).toBe('opened by app/dependabot');
  });

  it('sends the label only when one was given', async () => {
    mockFetch({ pullRequests: [], filter: '' });

    await fetchPullRequests({ owner: 'o', repo: 'r', token: 't', label: 'go mod' });

    const [, init] = (global.fetch as jest.Mock).mock.calls[0];
    expect(JSON.parse(init.body)).toEqual({ owner: 'o', repo: 'r', token: 't', label: 'go mod' });
  });

  it('survives an empty result without throwing', async () => {
    mockFetch({ pullRequests: [], filter: 'opened by app/dependabot' });

    const { pullRequests } = await fetchPullRequests({ owner: 'o', repo: 'r', token: 't' });

    expect(pullRequests).toEqual([]);
  });

  it('raises the error message the server sent', async () => {
    mockFetch({ error: 'GitHub reported: Could not resolve to a Repository' }, false, 502);

    await expect(fetchPullRequests({ owner: 'o', repo: 'r', token: 't' }))
      .rejects.toThrow(/Could not resolve to a Repository/);
  });
});
