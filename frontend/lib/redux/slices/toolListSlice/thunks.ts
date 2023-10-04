import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setToolListError, setToolListSuccess, setToolList } from './slice'
import { listTools } from './asyncActions'
import { AppDispatch } from '../..'


export const toolListThunk = createAppAsyncThunk(
  'tool/listTools',
  async (_, { dispatch }) => {
    console.log('toolListThunk')
    try {
      const response = await listTools()
      if (response) {
        dispatch(setToolListSuccess(true))
        dispatch(setToolList(response))
      } else {
        console.log('Failed to list Tools.', response)
        dispatch(setToolListError('Failed to list Tools.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list Flows.', error)
      if (error instanceof Error) {
        dispatch(setToolListError(error.message))
      } else {
        dispatch(setToolListError('Failed to add tool.'))
      }
      return false
    }
  }
)
