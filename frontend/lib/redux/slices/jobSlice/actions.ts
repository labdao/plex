export const initJobOnServer = async (
    jobData: { [key: string]: any }
): Promise<any> => {
    const response = await fetch('http://localhost:8080/init-job', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            },
            body: JSON.stringify(jobData),
        })

    if (!response.ok) {
        let errorMsg = 'An error occurred while initializing the job'
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