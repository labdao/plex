export const addToolToServer = async (
    payload: { toolData: { [key: string]: any }, walletAddress: string }
): Promise<any> => {
    const response = await fetch('http://localhost:8080/add-tool', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response.ok) {
        let errorMsg = 'An error occurred while adding the tool'
        try {
            const errorResult = await response.json()
            errorMsg = errorResult.message || errorMsg;
        } catch (e) {
            // Parsing JSON failed, retain the default error message.
        }
        console.log('errorMsg', errorMsg)
        throw new Error(errorMsg)
    }

    const result = await response.json()
    console.log('result', result)
    return result;
}