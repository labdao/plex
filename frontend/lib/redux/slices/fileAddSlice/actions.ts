import { getAccessToken } from "@privy-io/react-auth"
import backendUrl from "lib/backendUrl"

export const saveFileToServer = async (
    file: File,
    metadata: { [key: string]: any },
    isPublic: boolean
  ): Promise<{ filename: string, id: string }> => {
    const formData = new FormData()
    formData.append('file', file, file.name)
    formData.append('filename', file.name)
    formData.append('public', (isPublic ?? false).toString())

    for (const key in metadata) {
      formData.append(key, metadata[key])
    }

    let authToken
    try {
      authToken = await getAccessToken();
    } catch (error) {
      console.log("Failed to get access token: ", error)
      throw new Error("Authentication failed")
    }

    const response = await fetch(`${backendUrl()}/files`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`,
      },
      body: formData,
    })

    if (!response.ok) {
      let errorMsg = 'An error occurred while uploading the file'
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
    return result
  }
  