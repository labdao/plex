import backendUrl from "lib/backendUrl"

export const getJob = async (jobID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/jobs/${jobID}`, {
    method: 'Get',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to get job: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}

export const patchJob = async (jobID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/jobs/${jobID}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to patch job: ${response.status} ${response.statusText}`);
  }

  const result = await response.json()
  return result
}
