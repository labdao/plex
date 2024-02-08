import { getAccessToken } from "@privy-io/react-auth"
import backendUrl from "lib/backendUrl"

import { Kwargs } from "./slice"

export const createFlow = async (
    payload: { name: string, toolCid: string, scatteringMethod: string, kwargs: Kwargs }
): Promise<any> => {
    let authToken
    try {
        authToken = await getAccessToken();
    } catch (error) {
        console.log("Failed to get access token: ", error)
        throw new Error("Authentication failed")
    }

    const response = await fetch(`${backendUrl()}/flows`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response) {
        throw new Error("Failed to create Flow")
    }

    const result = await response.json()
    return result;
}
