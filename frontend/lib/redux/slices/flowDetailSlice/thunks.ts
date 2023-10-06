import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setFlowDetailError, setFlowDetailSuccess, setFlowDetail } from './slice'
import { getFlow } from './asyncActions'


export const flowDetailThunk = createAppAsyncThunk(
  'flow/flowDetail',
  async (flowCid: string, { dispatch }) => {
    try {
      const responseData = await getFlow(flowCid)
      dispatch(setFlowDetailSuccess(true))
      dispatch(setFlowDetail(responseData))
      return responseData
    } catch (error: unknown) {
      console.log('Failed to get Flow.', error)
      if (error instanceof Error) {
        dispatch(setFlowDetailError(error.message))
      } else {
        dispatch(setFlowDetailError('Failed to get Flow.'))
      }
      return false
    }
  }
)
