// The UI is served by the depbot binary on the same origin, so requests are
// relative. This keeps the port configurable via --port with no rebuild.
const baseUrl = '/api/';

const readError = async (response: Response, fallback: string): Promise<string> => {
  try {
    const body = await response.json();
    if (body && typeof body.error === 'string') {
      return body.error;
    }
  } catch {
    // Fall through to the generic message below.
  }

  return fallback;
};

export const postAsync = async (data: any, url: string) => {
  return fetch(baseUrl + url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data)
  });
};

export const saveConfiguration = async (configuration: string): Promise<string> => {
  const response = await postAsync({ configuration }, 'configurations');

  if (!response.ok) {
    throw new Error(await readError(response, `Could not save the file (status ${response.status}).`));
  }

  const body = await response.json();
  return body.path as string;
};

export const fetchPullRequests = async (payload: any) => {
  const response = await postAsync(payload, 'pull-requests');

  if (!response.ok) {
    throw new Error(await readError(response, `Could not load pull requests (status ${response.status}).`));
  }

  const result = await response.json();
  return result?.data?.repository?.pullRequests?.edges?.map((pr: any) => pr.node) ?? [];
};

export const fetchDetectedEcosystems = async () => {
  const response = await fetch(baseUrl + 'detect');

  if (!response.ok) {
    throw new Error(await readError(response, `Could not scan the project (status ${response.status}).`));
  }

  const result = await response.json();
  return result?.ecosystems ?? [];
};
