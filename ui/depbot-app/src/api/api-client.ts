const baseUrl = 'http://localhost:3001/api/';

export const postAsync = async (data: string, url: string) => {
  const response = await fetch(baseUrl + url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ configuration: data })
  });
  return response;
}

export const getAsync = async (url: string) => {
  const response = await fetch(baseUrl + url);
  return response;
}
