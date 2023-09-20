import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { setErrorToolAddSlice, startFileUpload, endFileUpload, setIsUploaded } from './toolSlice' // Importing actions from toolSlice.ts
import { addToolToServer } from './actions'

interface ToolPayload {
  toolData: { [key: string]: any }
  walletAddress: string
}

export const addToolAsync = createAppAsyncThunk(
  'tool/addTool',
  async ({ toolData, walletAddress }: ToolPayload, { dispatch }) => {
    try {
      dispatch(startFileUpload());
      const response = await addToolToServer({ toolData, walletAddress });

      if (response && response.filename && response.cid) {
        // Optionally, you could store something in localStorage or perform other operations.
        dispatch(setIsUploaded(true));
      } else {
        dispatch(setErrorToolAddSlice('Failed to add tool.'));
      }
      dispatch(endFileUpload());
      return response;
    } catch (error: unknown) {
      dispatch(endFileUpload());
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error 
        ? (error as { message?: string }).message 
        : undefined;

      dispatch(setErrorToolAddSlice(errorMessage || 'An error occurred while adding the tool.'));
      return false;
    }
  }
);
