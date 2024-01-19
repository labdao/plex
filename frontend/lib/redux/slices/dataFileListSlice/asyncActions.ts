import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl"

export const listDataFiles = async ({ page = 1, pageSize = 50, filters = {} }: { page?: number, pageSize?: number, filters?: Record<string, string | undefined> }): Promise<any> => {
  console.log("------- listDataFiles BEING CALLED -------")
  const queryParams = new URLSearchParams({ ...filters, page: page.toString(), pageSize: pageSize.toString() });
  let authToken;
  try {
    authToken = await getAccessToken()
  } catch (error) {
    console.log('Failed to get access token: ', error)
    throw new Error("Authentication failed");
  }

  const response = await fetch(`${backendUrl()}/datafiles?${queryParams}`, {
    method: 'Get',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error("Failed to list DataFiles");
  }

  return await response.json();
}
