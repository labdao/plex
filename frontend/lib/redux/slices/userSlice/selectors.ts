import type { ReduxState } from "@/lib/redux";

export const selectUserError = (state: ReduxState) => state.user.error;
export const selectUserWalletAddress = (state: ReduxState) => state.user.walletAddress;
export const selectUserDID = (state: ReduxState) => state.user.did;
export const selectUserTier = (state: ReduxState) => state.user.tier;
export const selectUserIsAdmin = (state: ReduxState) => state.user.isAdmin;