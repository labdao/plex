// redux/flowNames/thunks.ts
import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import isToday from "dayjs/plugin/isToday";

import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";
import { listFlowNames } from './asyncActions';
import { FlowName, setFlowNames, setFlowNamesLoading, setFlowNamesError, setFlowNamesSuccess, setCategorizedFlowNames } from './slice';

dayjs.extend(isToday);
dayjs.extend(isBetween);

export const flowNamesThunk = createAppAsyncThunk("flow/flowNames", async (walletAddress: string, { dispatch }) => {
  dispatch(setFlowNamesLoading(true));
  try {
    const response = await listFlowNames(walletAddress);
    if (response) {
      console.log("test response", response);
      const today = dayjs();
      const categories = {
        today: [] as FlowName[],
        last7Days: [] as FlowName[],
        last30Days: [] as FlowName[],
        older: [] as FlowName[],
      };
      response.forEach((flowName: FlowName) => {
        const start = dayjs(flowName.StartTime);
        if (start.isToday()) {
          categories.today.push(flowName);
        } else if (start.isBetween(today.subtract(7, "day"), today)) {
          categories.last7Days.push(flowName);
        } else if (start.isBetween(today.subtract(30, "day"), today)) {
          categories.last30Days.push(flowName);
        } else {
          categories.older.push(flowName);
        }
      });

      dispatch(setCategorizedFlowNames(categories));
      dispatch(setFlowNamesSuccess(true));
      dispatch(setFlowNames(response));
    } else {
      console.log("Failed to list Flows.", response);
      dispatch(setFlowNamesError("Failed to list Flows."));
    }
    dispatch(setFlowNamesLoading(false));

    return response;
  } catch (error: unknown) {
    dispatch(setFlowNamesLoading(false));
    console.log("Failed to list Flows.", error);
    if (error instanceof Error) {
      dispatch(setFlowNamesError(error.message));
    } else {
      dispatch(setFlowNamesError("Failed to list Flows."));
    }
    return false;
  }
});
