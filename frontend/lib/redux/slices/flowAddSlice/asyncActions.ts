import backendUrl from "lib/backendUrl"

import { Kwargs } from "./slice"

export const createFlow = async (
    payload: { name: string, walletAddress: string, toolCid: string, scatteringMethod: string, kwargs: Kwargs }
): Promise<any> => {
    const response = await fetch(`${backendUrl()}/flows`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response.ok) {
        let errorText = "Failed to create Flow"
        try {
            errorText = await response.text()
            console.log(errorText)
        } catch (e) {
            // Parsing JSON failed, retain the default error message.
        }
        throw new Error(errorText)
    }

    const result = await response.json()
    return result;
}
