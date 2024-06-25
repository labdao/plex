// redux/experimentNames/thunks.ts
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import isToday from "dayjs/plugin/isToday";

import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";
import { listExperimentNames } from './asyncActions';
import { ExperimentName, setExperimentNames, setExperimentNamesLoading, setExperimentNamesError, setExperimentNamesSuccess, setCategorizedExperimentNames } from './slice';

dayjs.extend(isToday);
dayjs.extend(isBetween);

export const experimentNamesThunk = createAppAsyncThunk("experiment/experimentNames", async (walletAddress: string, { dispatch }) => {
  dispatch(setExperimentNamesLoading(true));
  try {
    const response = await listExperimentNames(walletAddress);
    if (response) {
      console.log("test response", response);
      const today = dayjs();
      const categories = {
        today: [] as ExperimentName[],
        last7Days: [] as ExperimentName[],
        last30Days: [] as ExperimentName[],
        older: [] as ExperimentName[],
      };
      response.forEach((experimentName: ExperimentName) => {
        const start = dayjs(experimentName.StartTime);
        if (start.isToday()) {
          categories.today.push(experimentName);
        } else if (start.isBetween(today.subtract(7, "day"), today)) {
          categories.last7Days.push(experimentName);
        } else if (start.isBetween(today.subtract(30, "day"), today)) {
          categories.last30Days.push(experimentName);
        } else {
          categories.older.push(experimentName);
        }
      });

      dispatch(setCategorizedExperimentNames(categories));
      dispatch(setExperimentNamesSuccess(true));
      dispatch(setExperimentNames(response));
    } else {
      console.log("Failed to list Experiments.", response);
      dispatch(setExperimentNamesError("Failed to list Experiments."));
    }
    dispatch(setExperimentNamesLoading(false));

    return response;
  } catch (error: unknown) {
    dispatch(setExperimentNamesLoading(false));
    console.log("Failed to list Experiments.", error);
    if (error instanceof Error) {
      dispatch(setExperimentNamesError(error.message));
    } else {
      dispatch(setExperimentNamesError("Failed to list Experiments."));
    }
    return false;
  }
});
