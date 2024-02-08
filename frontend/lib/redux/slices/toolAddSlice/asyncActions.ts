import { getAccessToken } from "@privy-io/react-auth"
import backendUrl from "lib/backendUrl"

export const createTool = async (
    payload: { toolJson: { [key: string]: any } }
): Promise<any> => {
    let authToken
    try {
        authToken = await getAccessToken();
    } catch (error) {
        console.log("Failed to get access token: ", error)
        throw new Error("Authentication failed")
    }

    const response = await fetch(`${backendUrl()}/tools`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response) {
        throw new Error("Failed to create Tool")
    }

    const result = await response.json()
    return result;
}
