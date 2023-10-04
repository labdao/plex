import backendUrl from "lib/backendUrl"

export const listFlows = async (): Promise<any> => {
  const response = await fetch(`${backendUrl()}/flows`, {
    method: 'Get',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    let errorText = "Failed to list Flows"
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
