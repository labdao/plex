import { getAccessToken } from "@privy-io/react-auth"
import backendUrl from "lib/backendUrl"

export interface ApiKeyPayload {
  name: string;
  // Add any other properties that are needed for creating an API key
}

export const createApiKey = async (
    payload: ApiKeyPayload
): Promise<any> => {
    let authToken
    try {
        authToken = await getAccessToken();
    } catch (error) {
        console.log("Failed to get access token: ", error)
        throw new Error("Authentication failed")
    }

    const response = await fetch(`${backendUrl()}/api-keys`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response.ok) {
        const errorResult = await response.json();
        throw new Error(errorResult.message || "Failed to create API Key")
    }

    const result = await response.json()
    return result;
}