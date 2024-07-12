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

export const refreshUserDataThunk = createAppAsyncThunk(
  'user/refreshUserData',
  async (_, { dispatch }) => {
    try {
      await dispatch(fetchUserDataAsync()).unwrap();
    } catch (error) {
      console.error('Failed to refresh user data: ', error);
      throw error;
    }
  }
)