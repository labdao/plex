import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { createApiKey } from './asyncActions'
import { setApiKey, setApiKeyError, setApiKeySuccess } from './slice'

interface ApiKeyPayload {
  name: string;
  // Add any other properties that are needed for creating an API key
}

export const addApiKeyThunk = createAppAsyncThunk(
  'apiKey/addApiKey',
  async (payload: ApiKeyPayload, { dispatch }) => {
    try {
      const response = await createApiKey(payload)
      if (response && response.id) { // Assuming the response will have an 'id' field on successful creation
        dispatch(setApiKeySuccess(true))
        dispatch(setApiKey(response.key)) // Assuming you want to store the key in the state
      } else {
        console.log('Failed to add API key.', response)
        dispatch(setApiKeyError('Failed to add API key.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add API key.', error)
      if (error instanceof Error) {
        dispatch(setApiKeyError(error.message))
      } else {
        dispatch(setApiKeyError('Failed to add API key.'))
      }
      return false
    }
  }
)