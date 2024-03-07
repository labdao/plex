import dayjs from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import isToday from "dayjs/plugin/isToday";

import { createAppAsyncThunk } from "@/lib/redux/createAppAsyncThunk";

import { listFlows } from "./asyncActions";
import { Flow, setCategorizedFlows, setFlowList, setFlowListError, setFlowListLoading, setFlowListSuccess } from "./slice";

dayjs.extend(isToday);
dayjs.extend(isBetween);

export const flowListThunk = createAppAsyncThunk("flow/flowList", async (walletAddress: string, { dispatch }) => {
  dispatch(setFlowListLoading(true));
  try {
    const response = await listFlows(walletAddress);
    if (response) {
      const today = dayjs();
      const categories = {
        today: [] as Flow[],
        last7Days: [] as Flow[],
        last30Days: [] as Flow[],
        older: [] as Flow[],
      };
      response.forEach((flow: Flow) => {
        const start = dayjs(flow.StartTime);
        if (start.isToday()) {
          categories.today.push(flow);
        } else if (start.isBetween(today.subtract(7, "day"), today)) {
          categories.last7Days.push(flow);
        } else if (start.isBetween(today.subtract(30, "day"), today)) {
          categories.last30Days.push(flow);
        } else {
          categories.older.push(flow);
        }
      });

      dispatch(setCategorizedFlows(categories));
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
