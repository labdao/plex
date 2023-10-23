import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { saveUserDataToServer } from './actions'

interface UserPayload {
  walletAddress: string,
  isMember: boolean,
}

export const saveUserAsync = createAppAsyncThunk(
  'user/saveUserDataToServer',
  async ({walletAddress, isMember}: UserPayload) => {
    const result = await saveUserDataToServer(walletAddress, isMember);
    return result;
  }
)