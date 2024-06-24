import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { createTool } from './asyncActions'
import { setAddToolError, setAddToolSuccess } from './toolSlice'

interface ToolPayload {
  toolJson: { [key: string]: any }
}

export const createToolThunk = createAppAsyncThunk(
  'model/addTool',
  async ({ toolJson }: ToolPayload, { dispatch }) => {
    try {
      const response = await createTool({ toolJson })
      if (response && response.cid) {
        dispatch(setAddToolSuccess(true))
      } else {
        console.log('Failed to add model.', response)
        dispatch(setAddToolError('Failed to add model.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add model.', error)
      if (error instanceof Error) {
        dispatch(setAddToolError(error.message))
      } else {
        dispatch(setAddToolError('Failed to add model.'))
      }
      return false
    }
  }
)
