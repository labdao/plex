// redux/experimentNames/asyncActions.ts
import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";

export const listExperimentNames = async (walletAddress: string): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken();
  } catch (error) {
    console.log('Failed to get access token: ', error);
    throw new Error("Authentication failed");
  }

  const requestUrl = `${backendUrl()}/experiments?fields=name&walletAddress=${encodeURIComponent(walletAddress)}`;
  const requestOptions = {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  };
  const response = await fetch(requestUrl, requestOptions);

  if (!response) {
    let errorText = "Failed to list Experiments";
    try {
      console.log(errorText);
    } catch (e) {
      // Parsing JSON failed, retain the default error message.
    }
    throw new Error(errorText);
  }

  const result = await response.json();
  return result;
};
