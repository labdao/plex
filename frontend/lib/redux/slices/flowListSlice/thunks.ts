import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { listFlows } from './asyncActions'
import { setFlowList,setFlowListError, setFlowListSuccess } from './slice'


export const flowListThunk = createAppAsyncThunk(
  'flow/flowList',
  async (_, { dispatch }) => {
    try {
      const response = await listFlows()
      if (response) {
        dispatch(setFlowListSuccess(true))
        dispatch(setFlowList(response))
      } else {
        console.log('Failed to list Flows.', response)
        dispatch(setFlowListError('Failed to list Flows.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list Flows.', error)
      if (error instanceof Error) {
        dispatch(setFlowListError(error.message))
      } else {
        dispatch(setFlowListError('Failed to list Flows.'))
      }
      return false
    }
  }
)
