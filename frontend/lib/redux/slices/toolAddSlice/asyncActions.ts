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

    if (!response.ok) {
        let errorText = "Failed to create Tool"
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
