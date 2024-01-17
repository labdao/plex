import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { getJob, patchJob } from './asyncActions'
import { setJobDetail, setJobDetailError, setJobDetailLoading,setJobDetailSuccess } from './slice'


export const jobDetailThunk = createAppAsyncThunk(
  'job/jobDetail',
  async (jobID: string, { dispatch }) => {
    dispatch(setJobDetailLoading(true))
    try {
      const responseData = await getJob(jobID)
      dispatch(setJobDetailSuccess(true))
      dispatch(setJobDetail(responseData))
      dispatch(setJobDetailLoading(false))
      return responseData
    } catch (error: unknown) {
      console.log('Failed to get Job.', error)
      if (error instanceof Error) {
        dispatch(setJobDetailError(error.message))
      } else {
        dispatch(setJobDetailError('Failed to get Job.'))
      }
      dispatch(setJobDetailLoading(false))
      return false
    }
  }
)

export const jobPatchDetailThunk = createAppAsyncThunk(
  'job/jobPatchDetail',
  async (jobID: string, { dispatch }) => {
    dispatch(setJobDetailLoading(true))
    try {
      const responseData = await patchJob(jobID)
      dispatch(setJobDetailLoading(false))
      dispatch(setJobDetailSuccess(true))
      dispatch(setJobDetail(responseData))
      return responseData
    } catch (error: unknown) {
      console.log('Failed to get Job.', error)
      if (error instanceof Error) {
        dispatch(setJobDetailError(error.message))
      } else {
        dispatch(setJobDetailError('Failed to patch Job.'))
      }
      dispatch(setJobDetailLoading(false))
      return false
    }
  }
)
