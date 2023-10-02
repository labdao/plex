import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { saveUserDataToServer } from './actions'

interface UserPayload {
  walletAddress: string,
  emailAddress: string
}

export const saveUserAsync = createAppAsyncThunk(
  'user/saveUserDataToServer',
  async ({walletAddress, emailAddress}: {walletAddress: string, emailAddress: string}) => {
    const result = await saveUserDataToServer(walletAddress, emailAddress);
    return result;
  }
)
