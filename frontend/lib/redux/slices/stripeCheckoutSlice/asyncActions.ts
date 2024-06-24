import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl";

interface CheckoutPayload {
  modelCid: string;
  scatteringMethod: string;
  kwargs: string;
}

export const getCheckoutURL = async (payload: CheckoutPayload): Promise<any> => {
  let authToken;
  try {
    authToken = await getAccessToken();
  } catch (error) {
    console.log("Failed to get access token: ", error);
    throw new Error("Authentication failed");
  }

  const response = await fetch(`${backendUrl()}/stripe/checkout`, {
    method: "POST",
    headers: {
      Authorization: `Bearer ${authToken}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error(`Problem getting checkout URL: ${response.status} ${response.statusText}`);
  }

  const result = await response.json();
  return result;
};