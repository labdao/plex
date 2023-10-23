import backendUrl from "lib/backendUrl"
import { createAction } from "@reduxjs/toolkit"

export const setError = createAction<string | null>('user/setError')

export const saveUserDataToServer = async (
  walletAddress: string,
  isMember: boolean,
): Promise<{ walletAddress: string, isMember: boolean }> => {
  console.log('Entering saveUserDataToServer', walletAddress, isMember)

  try {
    const response = await fetch(`${backendUrl()}/user`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ walletAddress, isMember: false }),
    })

    if (!response.ok) {
      let errorMsg = 'An error occurred'
      try {
        const errorResult = await response.json()
        errorMsg = errorResult.message || errorMsg;
      } catch (e) {
        console.log('Error parsing JSON:', e)
      }
      console.log('Error message:', errorMsg)
      throw new Error(errorMsg)
    }

    const result = await response.json()
    console.log('Result:', result)
    return result;
  } catch (e) {
    console.log('Error in saveUserDataToServer:', e)
    throw e
  }
}