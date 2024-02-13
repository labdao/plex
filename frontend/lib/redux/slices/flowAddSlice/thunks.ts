import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { createFlow } from './asyncActions'
import { Kwargs,setFlowAddCid, setFlowAddError, setFlowAddID,setFlowAddSuccess } from './slice'

interface FlowPayload {
  name: string,
  toolCid: string,
  scatteringMethod: string,
  kwargs: Kwargs
}

export const addFlowThunk = createAppAsyncThunk(
  'flow/addFlow',
  async ({ name, toolCid, scatteringMethod, kwargs }: FlowPayload, { dispatch }) => {
    try {
      const response = await createFlow({ name, toolCid, scatteringMethod, kwargs })
      if (response && response.cid) {
        dispatch(setFlowAddSuccess(true))
        dispatch(setFlowAddCid(response.CID))
        dispatch(setFlowAddID(response.ID))
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
