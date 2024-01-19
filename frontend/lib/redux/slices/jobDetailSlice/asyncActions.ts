import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl"

export const getJob = async (bacalhauJobID: string): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken()
  } catch (error) {
    console.log('Failed to get access token: ', error)
    throw new Error("Authentication failed");
  }

  const response = await fetch(`${backendUrl()}/jobs/${bacalhauJobID}`, {
    method: 'Get',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to get job: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}

export const patchJob = async (bacalhauJobID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/jobs/${bacalhauJobID}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to patch job: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}
