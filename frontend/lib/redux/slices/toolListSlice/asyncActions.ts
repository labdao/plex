import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl"

export const listTools = async (taskSlug?: string): Promise<any> => {
  const url = taskSlug ? `${backendUrl()}/models?taskCategory=${encodeURIComponent(taskSlug)}` : `${backendUrl()}/models`;
  const authToken = await getAccessToken()
  const response = await fetch(url, {
    method: 'Get',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response) {
    throw new Error("Failed to list Models")
  }

  const result = await response.json()
  return result;
}
