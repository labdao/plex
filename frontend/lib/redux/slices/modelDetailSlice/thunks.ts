import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { getModel, patchModel } from "./asyncActions";
import { setModelDetail, setModelDetailError, setModelDetailLoading, setModelDetailSuccess } from "./slice";

export const modelDetailThunk = createAppAsyncThunk("model/modelDetail", async (bacalhauModelID: string, { dispatch }) => {
  dispatch(setModelDetailLoading(true));
  try {
    const responseData = await getModel(bacalhauModelID);
    dispatch(setModelDetailSuccess(true));
    dispatch(setModelDetail(responseData));
    dispatch(setModelDetailLoading(false));
    return responseData;
  } catch (error: unknown) {
    console.log("Failed to get Model.", error);
    if (error instanceof Error) {
      dispatch(setModelDetailError(error.message));
    } else {
      dispatch(setModelDetailError("Failed to get Model."));
    }
    dispatch(setModelDetailLoading(false));
    return false;
  }
});

export const modelPatchDetailThunk = createAppAsyncThunk("model/modelPatchDetail", async (bacalhauModelID: string, { dispatch }) => {
  dispatch(setModelDetailLoading(true));
  try {
    const responseData = await patchModel(bacalhauModelID);
    dispatch(setModelDetailLoading(false));
    dispatch(setModelDetailSuccess(true));
    dispatch(setModelDetail(responseData));
    return responseData;
  } catch (error: unknown) {
    console.log("Failed to get Model.", error);
    if (error instanceof Error) {
      dispatch(setModelDetailError(error.message));
    } else {
      dispatch(setModelDetailError("Failed to patch Model."));
    }
    dispatch(setModelDetailLoading(false));
    return false;
  }
});
