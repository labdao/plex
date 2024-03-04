import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { listFlows } from "./asyncActions";
import { setFlowList, setFlowListError, setFlowListLoading, setFlowListSuccess } from "./slice";

export const flowListThunk = createAppAsyncThunk("flow/flowList", async (walletAddress: string, { dispatch }) => {
  dispatch(setFlowListLoading(true));
  try {
    const response = await listFlows(walletAddress);
    if (response) {
      dispatch(setFlowListSuccess(true));
      dispatch(setFlowList(response));
    } else {
      console.log("Failed to list Flows.", response);
      dispatch(setFlowListError("Failed to list Flows."));
    }
    dispatch(setFlowListLoading(false));

    return response;
  } catch (error: unknown) {
    dispatch(setFlowListLoading(false));
    console.log("Failed to list Flows.", error);
    if (error instanceof Error) {
      dispatch(setFlowListError(error.message));
    } else {
      dispatch(setFlowListError("Failed to list Flows."));
    }
    return false;
  }
});
