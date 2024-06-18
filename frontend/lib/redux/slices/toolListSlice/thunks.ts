import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { AppDispatch } from '../..'
import { listTools } from './asyncActions'
import { setToolList,setToolListError, setToolListSuccess } from './slice'


export const toolListThunk = createAppAsyncThunk(
  'tool/listTools',
  async (taskSlug: string | undefined, { dispatch }) => {
    console.log('toolListThunk')
    try {
      const response = await listTools(taskSlug)
      if (response) {
        dispatch(setToolListSuccess(true))
        dispatch(setToolList(response))
      } else {
        console.log('Failed to list Tools.', response)
        dispatch(setToolListError('Failed to list Tools.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list Experiments.', error)
      if (error instanceof Error) {
        dispatch(setToolListError(error.message))
      } else {
        dispatch(setToolListError('Failed to add tool.'))
      }
      return false
    }
  }
)
