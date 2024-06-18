import { getAccessToken } from "@privy-io/react-auth"
import backendUrl from "lib/backendUrl"

import { Kwargs } from "./slice"

export const createExperiment = async (
    payload: { name: string, toolCid: string, scatteringMethod: string, kwargs: Kwargs }
): Promise<any> => {
    let authToken
    try {
        authToken = await getAccessToken();
    } catch (error) {
        console.log("Failed to get access token: ", error)
        throw new Error("Authentication failed")
    }

    const response = await fetch(`${backendUrl()}/experiments`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${authToken}`,
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
    })

    if (!response) {
        throw new Error("Failed to create Experiment")
    }

    const result = await response.json()
    return result;
}

export const addJobToExperiment = async (experimentId: number, payload: { name: string, toolCid: string, scatteringMethod: string, kwargs: Kwargs }): Promise<any> => {
    let authToken;
    try {
      authToken = await getAccessToken()
    } catch (error) {
      console.log('Failed to get access token: ', error)
      throw new Error("Authentication failed");
    }
  
    const requestUrl = `${backendUrl()}/experiments/${experimentId}/add-job`;
    const requestOptions = {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    };
  
    try {
      const response = await fetch(requestUrl, requestOptions);
      if (!response.ok) {
        throw new Error(`Failed to add job to Experiment: ${response.statusText}`);
      }
      const result = await response.json();
      return result;
    } catch (error) {
      console.error('Failed to add job to Experiment:', error);
      throw new Error('Failed to add job to Experiment');
    }
  };