"use client";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { AppDispatch, experimentDetailThunk, selectExperimentDetail, selectExperimentDetailError, selectExperimentDetailLoading, selectToolDetail } from "@/lib/redux";

import { aggregateJobStatus } from "../ExperimentStatus";
import JobsAccordion from "./JobsAccordion";
import MetricsVisualizer from "./MetricsVisualizer";

dayjs.extend(relativeTime);

export default function ExperimentDetail() {
  const dispatch = useDispatch<AppDispatch>();
  const experiment = useSelector(selectExperimentDetail);
  const tool = useSelector(selectToolDetail);
  const loading = useSelector(selectExperimentDetailLoading);
  const error = useSelector(selectExperimentDetailError);

  const status = aggregateJobStatus(experiment.Jobs);

  const experimentID = experiment.ID?.toString();

  useEffect(() => {
    if (["running", "queued"].includes(status.status) && experimentID) {
      const interval = setInterval(() => {
        console.log("Checking for new results");
        dispatch(experimentDetailThunk(experimentID));
      }, 15000);

      return () => clearInterval(interval);
    }
  }, [dispatch, experimentID, status]);

  return (
    <div>
      {tool?.ToolJson?.checkpointCompatible && <MetricsVisualizer experiment={experiment} key={experiment.ID} />}
      <JobsAccordion experiment={experiment} />
    </div>
  );
}
