import { getAccessToken } from "@privy-io/react-auth";
import backendUrl from "lib/backendUrl"

export const listTools = async (taskSlug?: string): Promise<any> => {
  const url = taskSlug ? `${backendUrl()}/tools?taskCategory=${encodeURIComponent(taskSlug)}` : `${backendUrl()}/tools`;
  const authToken = await getAccessToken()
  const response = await fetch(url, {
    method: 'Get',
    headers: {
      'Authorization': `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response) {
    throw new Error("Failed to list Tools")
  }

  const result = await response.json()
  return result;
}
