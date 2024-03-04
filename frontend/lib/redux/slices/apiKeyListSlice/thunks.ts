import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { listApiKeys } from './asyncActions'
import { setApiKeyList, setApiKeyListError, setApiKeyListSuccess } from './slice'

export const apiKeyListThunk = createAppAsyncThunk(
  'apiKey/apiKeyList',
  async (_, { dispatch }) => {
    try {
      const response = await listApiKeys()
      if (response) {
        dispatch(setApiKeyListSuccess(true))
        dispatch(setApiKeyList(response))
      } else {
        console.log('Failed to list API Keys.', response)
        dispatch(setApiKeyListError('Failed to list API Keys.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list API Keys.', error)
      if (error instanceof Error) {
        dispatch(setApiKeyListError(error.message))
      } else {
        dispatch(setApiKeyListError('Failed to list API Keys.'))
      }
      return false
    }
  }
)