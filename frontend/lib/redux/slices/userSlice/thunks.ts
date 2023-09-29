import { createAppAsyncThunk } from '@/lib/redux/createAppAsyncThunk'
import { Web3Auth } from '@web3auth/modal';
import { setError } from './actions'

import { setIsLoggedIn } from '@/lib/redux'
import { saveUserDataToServer } from './actions'

interface UserPayload {
  username: string
  walletAddress: string
}

// export const initWeb3Auth = createAppAsyncThunk(
//   'user/initWeb3Auth',
//   async (_, { dispatch }) => {
//     try {
//       const web3AuthInstance = new Web3Auth({
//         clientId: "BKURHIghKRSWvu0c2IM8hrFtKRny4zLjBqO8Mr4fiIoGc2cSB8_il38d2T5fIxDvpIFQqyEZ5lbNswl_GITUZd0",
//         web3AuthNetwork: "sapphire_devnet",
//         chainConfig: {
//           chainId: "0x1a4",
//           chainNamespace: "other",
//           rpcTarget: "https://opt-goerli.g.alchemy.com/v2/tvdEoAYqtbNgXzL-Dma7if3i3NzUS3-N",
//         },
//         uiConfig: {
//           loginMethodsOrder: ["email_passwordless"],
//         },
//       });

//       await web3AuthInstance.initModal();
//       dispatch(setWeb3Auth(web3AuthInstance));
//     } catch (error: unknown) {
//       const errorMessage = typeof error === 'object' && error !== null && 'message' in error
//         ? (error as { message?: string }).message
//         : undefined;

//       dispatch(setError(errorMessage || 'An error occurred.'));
//     }
//   }
// )

export const saveUserAsync = createAppAsyncThunk(
  'user/saveUserData',
  async ({username, walletAddress}: UserPayload, { dispatch }) => {
    try {
      const response = await saveUserDataToServer(username, walletAddress)

      if (response.username && response.walletAddress) {
        localStorage.setItem('username', username)
        localStorage.setItem('walletAddress', walletAddress)
        dispatch(setIsLoggedIn(true))
      } else {
        dispatch(setError('Failed to save user data.'))
      }
      return response
    } catch (error: unknown) {
      const errorMessage = typeof error === 'object' && error !== null && 'message' in error 
        ? (error as { message?: string }).message 
        : undefined;

      dispatch(setError(errorMessage || 'An error occurred.'));
      return false;
    }
  }
)
