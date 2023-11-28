import backendUrl from "lib/backendUrl";

export const getTool = async (CID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/tools/${CID}`, {
    method: "Get",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to get tool: ${response.status} ${response.statusText}`);
  }

  const result = await response.json();
  return result;
};

export const patchTool = async (CID: string): Promise<any> => {
  const response = await fetch(`${backendUrl()}/tools/${CID}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to patch tool: ${response.status} ${response.statusText}`);
  }

  const result = await response.json();
  return result;
};
