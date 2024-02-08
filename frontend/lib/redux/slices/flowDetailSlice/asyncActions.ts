import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl"

export const getFlow = async (flowCid: string): Promise<any> => {
  let authToken;

  try {
    authToken = await getAccessToken()
  } catch (error) {
    console.log('Failed to get access token: ', error)
    throw new Error("Authentication failed");
  }

  console.log(`Fetching flow details with CID: ${flowCid}`); // Log the flow CID
  console.log(`Using Authorization token: ${authToken}`); // Log the token for debugging
  console.log(`URL: ${backendUrl()}/flows/${flowCid}`); // Log the URL for debugging

  const response = await fetch(`${backendUrl()}/flows/${flowCid}`, {
    method: 'Get',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to get flow: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}

export const patchFlow = async (flowCid: string): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken();
  } catch (error) {
    console.log('Failed to get access token: ', error)
    throw new Error("Authentication failed");
  }

  const response = await fetch(`${backendUrl()}/flows/${flowCid}`, {
    method: 'PATCH',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to patch flow: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}
