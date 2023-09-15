import backendUrl from "lib/backendUrl";

export const saveDataFileToServer = async (
    file: File,
    metadata: { [key: string]: any }
  ): Promise<{ filename: string, cid: string }> => {
    const formData = new FormData();
    formData.append('file', file, file.name);
    formData.append('filename', file.name)
    formData.append('isPublic', String(metadata.isPublic))
    formData.append('isVisible', String(metadata.isVisible))

    for (const key in metadata) {
      formData.append(key, metadata[key]);
    }
  
    const response = await fetch(`${backendUrl()}/add-datafile`, {
      method: 'POST',
      body: formData,
    })
  
    if (!response.ok) {
      let errorMsg = 'An error occurred while uploading the data file'
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
  