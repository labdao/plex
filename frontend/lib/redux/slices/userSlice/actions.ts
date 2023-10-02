import backendUrl from "lib/backendUrl"
import { createAction } from "@reduxjs/toolkit"

export const setError = createAction<string | null>('user/setError')

export const saveUserDataToServer = async (
  walletAddress: string,
  emailAddress: string,
): Promise<{ walletAddress: string, emailAddress: string }> => {
  const response = await fetch(`${backendUrl()}/user`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ walletAddress, emailAddress }),
  })

  if (!response.ok) {
    let errorMsg = 'An error occurred'
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
