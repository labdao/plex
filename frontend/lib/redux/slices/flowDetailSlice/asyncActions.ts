import backendUrl from "lib/backendUrl"

export const getFlow = async (flowCid: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/flows/${flowCid}`, {
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

export const patchFlow = async (flowCid: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/flows/${flowCid}`, {
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
