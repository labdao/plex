import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";

export const updateExperiment = async (experimentId: string, data: { name?: string; public?: boolean; }): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken()
  } catch (error) {
    console.log('Failed to get access token: ', error)
    throw new Error("Authentication failed");
  }

  const requestUrl = `${backendUrl()}/experiments/${experimentId}`;
  const requestOptions = {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data)
  };

  try {
    const response = await fetch(requestUrl, requestOptions);
    if (!response.ok) {
      throw new Error(`Failed to update Experiment: ${response.statusText}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to update Experiment:', error);
    throw new Error('Failed to update Experiment');
  }
};