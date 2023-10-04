import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setFlowAddError, setFlowAddSuccess, Kwargs } from './slice'
import { createFlow } from './asyncActions'

interface FlowPayload {
  name: string,
  walletAddress: string,
  toolCid: string,
  scatteringMethod: string,
  kwargs: Kwargs
}

export const addFlowThunk = createAppAsyncThunk(
  'flow/addFlow',
  async ({ name, walletAddress, toolCid, scatteringMethod, kwargs }: FlowPayload, { dispatch }) => {
    try {
      const response = await createFlow({ name, walletAddress, toolCid, scatteringMethod, kwargs })
      if (response && response.cid) {
        dispatch(setFlowAddSuccess(true))
      } else {
        console.log('Failed to add tool.', response)
        dispatch(setFlowAddError('Failed to add tool.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add Flow.', error)
      if (error instanceof Error) {
        dispatch(setFlowAddError(error.message))
      } else {
        dispatch(setFlowAddError('Failed to add tool.'))
      }
      return false
    }
  }
)
