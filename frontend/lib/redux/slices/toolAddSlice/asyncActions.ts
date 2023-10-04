import backendUrl from "lib/backendUrl"

export const createTool = async (
    payload: { toolJson: { [key: string]: any }, walletAddress: string }
): Promise<any> => {
    const response = await fetch(`${backendUrl()}/tool`, {
        method: 'POST',
        headers: {
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
