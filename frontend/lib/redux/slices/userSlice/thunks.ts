import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { saveUserDataToServer } from './actions'

interface UserPayload {
  walletAddress: string,
}

export const saveUserAsync = createAppAsyncThunk(
  'user/saveUserDataToServer',
  async ({walletAddress}: {walletAddress: string}) => {
    const result = await saveUserDataToServer(walletAddress);
    return result;
  }
)
