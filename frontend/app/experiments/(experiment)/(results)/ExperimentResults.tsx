"use client";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { AppDispatch, flowDetailThunk, selectFlowDetail, selectFlowDetailError, selectFlowDetailLoading, selectToolDetail } from "@/lib/redux";

import { aggregateJobStatus } from "../ExperimentStatus";
import { ActiveResultContextProvider } from "./ActiveResultContext";
import MetricsVisualizer from "./MetricsVisualizer";
import JobsAccordion from "./JobsAccordion";

dayjs.extend(relativeTime);

export default function ExperimentDetail() {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const tool = useSelector(selectToolDetail);
  const loading = useSelector(selectFlowDetailLoading);
  const error = useSelector(selectFlowDetailError);

  const status = aggregateJobStatus(flow.Jobs);

  const experimentID = flow.ID?.toString();

  useEffect(() => {
    if (["running", "queued"].includes(status.status) && experimentID) {
      const interval = setInterval(() => {
        console.log("Checking for new results");
        dispatch(flowDetailThunk(experimentID));
      }, 15000);

      return () => clearInterval(interval);
    }
  }, [dispatch, experimentID, status]);

  return (
    <div>
      <ActiveResultContextProvider>
        {tool?.ToolJson?.checkpointCompatible && <MetricsVisualizer flow={flow} />}
        <JobsAccordion flow={flow} />
      </ActiveResultContextProvider>
    </div>
  );
}