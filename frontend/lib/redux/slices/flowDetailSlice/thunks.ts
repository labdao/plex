import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { getFlow, patchFlow } from "./asyncActions";
import { setFlowDetail, setFlowDetailError, setFlowDetailLoading, setFlowDetailSuccess } from "./slice";

export const flowDetailThunk = createAppAsyncThunk("flow/flowDetail", async (flowID: string, { dispatch }) => {
  try {
    dispatch(setFlowDetailLoading(true));
    const responseData = await getFlow(flowID);
    dispatch(setFlowDetailSuccess(true));
    dispatch(setFlowDetail(responseData));
    dispatch(setFlowDetailLoading(false));
    return responseData;
  } catch (error: unknown) {
    console.log("Failed to get Flow.", error);
    if (error instanceof Error) {
      dispatch(setFlowDetailError(error.message));
    } else {
      dispatch(setFlowDetailError("Failed to get Flow."));
    }
    dispatch(setFlowDetailLoading(false));
    return false;
  }
});
