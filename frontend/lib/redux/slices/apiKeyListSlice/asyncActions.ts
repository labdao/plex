import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";

export const listApiKeys = async (): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken();
    console.log('authToken: ', authToken);
  } catch (error) {
    console.log('Failed to get access token: ', error);
    throw new Error("Authentication failed");
  }

  const requestUrl = `${backendUrl()}/api-keys`;
  const requestOptions = {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  };
  const response = await fetch(requestUrl, requestOptions);

  if (!response.ok) {
    let errorText = "Failed to list API Keys";
    try {
      const errorResult = await response.json();
      errorText = errorResult.message || errorText;
      console.log(errorText);
    } catch (e) {
      console.log('Failed to parse error response: ', e);
    }
    throw new Error(errorText);
  }

  const result = await response.json();
  return result;
}