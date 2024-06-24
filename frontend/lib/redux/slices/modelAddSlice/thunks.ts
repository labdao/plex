import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { createModel } from './asyncActions'
import { setAddModelError, setAddModelSuccess } from './slice'

interface ModelPayload {
  modelJson: { [key: string]: any }
}

export const createModelThunk = createAppAsyncThunk(
  'model/addModel',
  async ({ modelJson }: ModelPayload, { dispatch }) => {
    try {
      const response = await createModel({ modelJson })
      if (response && response.cid) {
        dispatch(setAddModelSuccess(true))
      } else {
        console.log('Failed to add model.', response)
        dispatch(setAddModelError('Failed to add model.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add model.', error)
      if (error instanceof Error) {
        dispatch(setAddModelError(error.message))
      } else {
        dispatch(setAddModelError('Failed to add model.'))
      }
      return false
    }
  }
)
