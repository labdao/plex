import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { AppDispatch } from '../..'
import { listModels } from './asyncActions'
import { setModelList,setModelListError, setModelListSuccess } from './slice'


export const modelListThunk = createAppAsyncThunk(
  'model/listModels',
  async (taskSlug: string | undefined, { dispatch }) => {
    console.log('modelListThunk')
    try {
      const response = await listModels(taskSlug)
      if (response) {
        dispatch(setModelListSuccess(true))
        dispatch(setModelList(response))
      } else {
        console.log('Failed to list Models.', response)
        dispatch(setModelListError('Failed to list Models.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list Experiments.', error)
      if (error instanceof Error) {
        dispatch(setModelListError(error.message))
      } else {
        dispatch(setModelListError('Failed to add model.'))
      }
      return false
    }
  }
)
