import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setFlowListError, setFlowListSuccess } from './slice'
import { listFlows } from './asyncActions'


export const FlowListThunk = createAppAsyncThunk(
  'flow/flowList',
  async ({}, { dispatch }) => {
    try {
      const response = await listFlows()
      if (response && response.ok) {
        dispatch(setFlowListSuccess(true))
      } else {
        console.log('Failed to list Flows.', response)
        dispatch(setFlowListError('Failed to list DataFiles.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to list Flows.', error)
      if (error instanceof Error) {
        dispatch(setFlowListError(error.message))
      } else {
        dispatch(setFlowListError('Failed to add tool.'))
      }
      return false
    }
  }
)
