import backendUrl from "lib/backendUrl";

export const listFlows = async (walletAddress: string): Promise<any> => {
  //Temporarily fetch all experiments for dev, change back before merging!
  //const response = await fetch(`${backendUrl()}/flows?walletAddress=${encodeURIComponent(walletAddress)}`, {
  const response = await fetch(`${backendUrl()}/flows`, {
    method: "Get",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response) {
    let errorText = "Failed to list Flows";
    try {
      console.log(errorText);
    } catch (e) {
      // Parsing JSON failed, retain the default error message.
    }
    throw new Error(errorText);
  }

  const result = await response.json();
  return result;
};
