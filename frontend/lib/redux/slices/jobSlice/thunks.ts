// gateway/frontend/lib/redux/slices/jobSlice/thunks.ts

import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setError, startJobInitialization, endJobInitialization, setIsInitialized } from './jobSlice' // Importing actions from jobSlice.ts
import { initJobOnServer } from './actions'

interface JobPayload {
  jobData: { [key: string]: any }
}

export const initJobAsync = createAppAsyncThunk(
  'job/initJob',
  async ({ jobData }: JobPayload, { dispatch }) => {
    try {
      dispatch(startJobInitialization());
      const response = await initJobOnServer(jobData);

      if (response && response.cid) {
        // Optionally, you could store something in localStorage or perform other operations.
        dispatch(setIsInitialized(true));
      } else {
        dispatch(setError('Failed to initialize job.'));
      }
      dispatch(endJobInitialization());
      return response;
    } catch (error: unknown) {
      dispatch(endJobInitialization());
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error 
        ? (error as { message?: string }).message 
        : undefined;

      dispatch(setError(errorMessage || 'An error occurred while initializing the job.'));
      return false;
    }
  }
);