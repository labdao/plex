import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { createExperiment } from './asyncActions'
import { Kwargs, setExperimentAddError, setExperimentAddID,setExperimentAddSuccess } from './slice'

interface ExperimentPayload {
  name: string,
  toolCid: string,
  scatteringMethod: string,
  kwargs: Kwargs
}

export const addExperimentThunk = createAppAsyncThunk(
  'experiment/addExperiment',
  async ({ name, toolCid, scatteringMethod, kwargs }: ExperimentPayload, { dispatch }) => {
    try {
      const response = await createExperiment({ name, toolCid, scatteringMethod, kwargs })
      if (response && response.cid) {
        dispatch(setExperimentAddSuccess(true))
        dispatch(setExperimentAddID(response.ID))
      } else {
        console.log('Failed to add tool.', response)
        dispatch(setExperimentAddError('Failed to add tool.'))
      }
      return response
    } catch (error: unknown) {
      console.log('Failed to add Experiment.', error)
      if (error instanceof Error) {
        dispatch(setExperimentAddError(error.message))
      } else {
        dispatch(setExperimentAddError('Failed to add tool.'))
      }
      return false
    }
  }
)
