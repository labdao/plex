import type { ReduxState } from "@/lib/redux";

export const selectUserError = (state: ReduxState) => state.user.error;
