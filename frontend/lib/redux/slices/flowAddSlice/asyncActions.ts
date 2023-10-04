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

    if (!response) {
        throw new Error("Failed to create Flow")
    }

    const result = await response.json()
    return result;
}
