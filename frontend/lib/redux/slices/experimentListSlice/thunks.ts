import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import isToday from "dayjs/plugin/isToday";

import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { listExperiments } from "./asyncActions";
import { Experiment, setCategorizedExperiments, setExperimentList, setExperimentListError, setExperimentListLoading, setExperimentListSuccess } from "./slice";

dayjs.extend(isToday);
dayjs.extend(isBetween);

export const experimentListThunk = createAppAsyncThunk("experiment/experimentList", async (walletAddress: string, { dispatch }) => {
  dispatch(setExperimentListLoading(true));
  try {
    const response = await listExperiments(walletAddress);
    if (response) {
      const today = dayjs();
      const categories = {
        today: [] as Experiment[],
        last7Days: [] as Experiment[],
        last30Days: [] as Experiment[],
        older: [] as Experiment[],
      };
      response.forEach((experiment: Experiment) => {
        const start = dayjs(experiment.StartTime);
        if (start.isToday()) {
          categories.today.push(experiment);
        } else if (start.isBetween(today.subtract(7, "day"), today)) {
          categories.last7Days.push(experiment);
        } else if (start.isBetween(today.subtract(30, "day"), today)) {
          categories.last30Days.push(experiment);
        } else {
          categories.older.push(experiment);
        }
      });

      dispatch(setCategorizedExperiments(categories));
      dispatch(setExperimentListSuccess(true));
      dispatch(setExperimentList(response));
    } else {
      console.log("Failed to list Experiments.", response);
      dispatch(setExperimentListError("Failed to list Experiments."));
    }
    dispatch(setExperimentListLoading(false));

    return response;
  } catch (error: unknown) {
    dispatch(setExperimentListLoading(false));
    console.log("Failed to list Experiments.", error);
    if (error instanceof Error) {
      dispatch(setExperimentListError(error.message));
    } else {
      dispatch(setExperimentListError("Failed to list Experiments."));
    }
    return false;
  }
});
