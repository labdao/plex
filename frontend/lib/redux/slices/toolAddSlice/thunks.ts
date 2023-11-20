import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { createTool } from './asyncActions'
import { setAddToolError, setAddToolSuccess } from './toolSlice'

interface ToolPayload {
  toolJson: { [key: string]: any }
  walletAddress: string
}

export const createToolThunk = createAppAsyncThunk(
  'tool/addTool',
  async ({ toolJson, walletAddress }: ToolPayload, { dispatch }) => {
    try {
      const response = await createTool({ toolJson, walletAddress })
      if (response && response.cid) {
        dispatch(setAddToolSuccess(true))
      } else {
        console.log('Failed to add tool.', response)
        dispatch(setAddToolError('Failed to add tool.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add tool.', error)
      if (error instanceof Error) {
        dispatch(setAddToolError(error.message))
      } else {
        dispatch(setAddToolError('Failed to add tool.'))
      }
      return false
    }
  }
)
