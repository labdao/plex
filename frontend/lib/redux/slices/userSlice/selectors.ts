import type { ReduxState } from '@/lib/redux'

export const selectWalletAddress = (state: ReduxState) => state.user.walletAddress
export const selectIsLoggedIn = (state: ReduxState) => state.user.isLoggedIn
export const selectIsMember = (state: ReduxState) => state.user.isMember
