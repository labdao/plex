import { getAccessToken } from "@privy-io/react-auth"
import backendUrl from "lib/backendUrl"

export const createModel = async (
    payload: { modelJson: { [key: string]: any } }
): Promise<any> => {
    let authToken
    try {
        authToken = await getAccessToken();
    } catch (error) {
        console.log("Failed to get access token: ", error)
        throw new Error("Authentication failed")
    }

    const response = await fetch(`${backendUrl()}/models`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response) {
        throw new Error("Failed to create Model")
    }

    const result = await response.json()
    return result;
}
