const baseUrl = 'http://localhost:3001/api/';

export const postAsync = async (data: any, url: string) => {
  const response = await fetch(baseUrl + url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data)
  });
  return response;
}

export const fetchPullRequests = async (payload: any) => {
  try {
    const response = await postAsync(payload, "pull-requests");

    if (!response.ok) {
      throw new Error(`Error! status: ${response.status}`);
    }
    const result = await response.json();
    return result.data.repository.pullRequests.edges.map((pr: any) => pr.node);
  } catch (error) {
    console.error('Error:', error);
  }
};


