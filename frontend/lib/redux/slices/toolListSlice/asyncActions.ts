import backendUrl from "lib/backendUrl"

export const listTools = async (): Promise<any> => {
  const response = await fetch(`${backendUrl()}/tools`, {
    method: 'Get',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response) {
    throw new Error("Failed to list Tools")
  }

  const result = await response.json()
  return result;
}
