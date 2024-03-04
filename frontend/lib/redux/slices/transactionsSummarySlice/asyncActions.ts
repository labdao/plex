import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";

export const getTransactionsSummary = async (): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken();
  } catch (error) {
    console.log("Failed to get access token: ", error);
    throw new Error("Authentication failed");
  }

  const response = await fetch(`${backendUrl()}/transactions-summary`, {
    method: "Get",
    headers: {
      Authorization: `Bearer ${authToken}`,
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Problem fetching transactions: ${response.status} ${response.statusText}`);
  }

  const result = await response.json();
  return result;
};
