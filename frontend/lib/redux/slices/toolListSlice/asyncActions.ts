import backendUrl from "lib/backendUrl"

export const listTools = async (taskSlug?: string): Promise<any> => {
  const url = taskSlug ? `${backendUrl()}/tools?taskCategory=${encodeURIComponent(taskSlug)}` : `${backendUrl()}/tools`;
  const response = await fetch(url, {
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
