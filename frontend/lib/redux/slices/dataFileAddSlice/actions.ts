import backendUrl from "lib/backendUrl";

export const saveDataFileToServer = async (
  files: File[] | FileList,
  metadata: { [key: string]: any }
): Promise<{ filename: string, cid: string }[]> => {
  const formData = new FormData();

  const filesArray = Array.isArray(files) ? files : Array.from(files);

  filesArray.forEach((file) => {
    formData.append('files', file, file.name);
  });

  Object.keys(metadata).forEach(key => {
    formData.append(key, metadata[key]);
  });

  const response = await fetch(`${backendUrl()}/datafiles`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    let errorMsg = 'An error occurred while uploading the data files';
    try {
      const errorResult = await response.json();
      errorMsg = errorResult.message || errorMsg;
    } catch (e) {
      // Parsing JSON failed, retain the default error message.
    }
    console.error('errorMsg', errorMsg);
    throw new Error(errorMsg);
  }

  const result = await response.json();
  console.log('result', result);
  return result;
};