import backendUrl from "lib/backendUrl"

export const listDataFiles = async (): Promise<any> => {
  const response = await fetch(`${backendUrl()}/datafiles`, {
    method: 'Get',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response) {
    throw new Error("Failed to list DataFiles")
  }

  const result = await response.json()
  return result;
}
