import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { getExperiment, patchExperiment } from "./asyncActions";
import { setExperimentDetail, setExperimentDetailError, setExperimentDetailLoading, setExperimentDetailSuccess } from "./slice";

export const experimentDetailThunk = createAppAsyncThunk("experiment/experimentDetail", async (experimentID: string, { dispatch }) => {
  try {
    dispatch(setExperimentDetailLoading(true));
    const responseData = await getExperiment(experimentID);
    dispatch(setExperimentDetailSuccess(true));
    dispatch(setExperimentDetail(responseData));
    dispatch(setExperimentDetailLoading(false));
    return responseData;
  } catch (error: unknown) {
    console.log("Failed to get Experiment.", error);
    if (error instanceof Error) {
      dispatch(setExperimentDetailError(error.message));
    } else {
      dispatch(setExperimentDetailError("Failed to get Experiment."));
    }
    dispatch(setExperimentDetailLoading(false));
    return false;
  }
});
