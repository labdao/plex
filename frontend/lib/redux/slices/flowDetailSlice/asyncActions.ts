import backendUrl from "lib/backendUrl"

export const getFlow = async (flowID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/flows/${flowID}`, {
    method: 'Get',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to get flow: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}

export const patchFlow = async (flowID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/flows/${flowID}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to patch flow: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}
