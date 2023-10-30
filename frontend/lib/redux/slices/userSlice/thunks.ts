import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { saveUserDataToServer } from './actions'
import backendUrl from 'lib/backendUrl'

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

export const fetchUserMemberStatusAsync = createAppAsyncThunk(
  'user/fetchUserMemberStatus',
  async (walletAddress: string) => {
    const response = await fetch(`${backendUrl()}/user/${walletAddress}`)
    if (!response.ok) {
      throw new Error('Failed to fetch member status');
    }
    const { isMember } = await response.json();
    return isMember;
  }
)
