import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";

export const updateFlow = async (flowId: string): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken()
  } catch (error) {
    console.log('Failed to get access token: ', error)
    throw new Error("Authentication failed");
  }

  const requestUrl = `${backendUrl()}/flows/${flowId}`;
  const requestOptions = {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  };

  try {
    const response = await fetch(requestUrl, requestOptions);
    if (!response.ok) {
      throw new Error(`Failed to update Flow: ${response.statusText}`);
    }
    const result = await response.json();
    return result;
  } catch (error) {
    console.error('Failed to update Flow:', error);
    throw new Error('Failed to update Flow');
  }
};