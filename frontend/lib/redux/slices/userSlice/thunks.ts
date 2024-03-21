import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'

import { fetchUserData, saveUserDataToServer } from './actions'

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

export const fetchUserDataAsync = createAppAsyncThunk(
  'user/fetchUserData',
  async () => {
    const result = await fetchUserData();
    return result;
  }
)