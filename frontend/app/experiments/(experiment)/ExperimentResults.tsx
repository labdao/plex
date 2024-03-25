"use client";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { TruncatedString } from "@/components/shared/TruncatedString";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Card } from "@/components/ui/card";
import {
  AppDispatch,
  flowDetailThunk,
  selectFlowDetail,
  selectFlowDetailError,
  selectFlowDetailLoading,
  selectFlowUpdateError,
  selectFlowUpdateLoading,
  selectFlowUpdateSuccess,
  selectUserWalletAddress,
} from "@/lib/redux";

import { aggregateJobStatus } from "./ExperimentStatus";
import JobDetail from "./JobDetail";
import MetricsVisualizer from "./MetricsVisualizer";

dayjs.extend(relativeTime);

export default function ExperimentDetail() {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const loading = useSelector(selectFlowDetailLoading);
  const error = useSelector(selectFlowDetailError);

  const status = aggregateJobStatus(flow.Jobs);
  const userWalletAddress = useSelector(selectUserWalletAddress);

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
      <>
        {/* add if checkpoint compatible later!! */}
        <MetricsVisualizer flow={flow} />
        <Accordion {...(flow?.Jobs?.length > 1 ? { type: "multiple" } : { type: "single", defaultValue: "0", collapsible: true })}>
          {flow.Jobs?.map((job, index) => {
            const validStates = ["queued", "running", "failed", "completed"];
            const status = (validStates.includes(job.State) ? job.State : "unknown") as "queued" | "running" | "failed" | "completed" | "unknown";

            return (
              <AccordionItem value={index.toString()} className="border-0 [&[data-state=open]>div]:shadow-lg" key={job.ID}>
                <Card className="my-2 shadow-sm">
                  <AccordionTrigger className="flex items-center justify-between w-full px-6 py-3 text-left hover:no-underline [&[data-state=open]]:bg-muted">
                    <div className="flex items-center gap-2">
                      <div className="w-30">
                        <div>condition {index + 1}</div>
                        <div className="flex gap-1 text-xs text-muted-foreground/70">
                          Job ID: {job.BacalhauJobID ? <TruncatedString value={job.BacalhauJobID} /> : "n/a"}
                        </div>
                      </div>
                      <Badge status={status} variant="outline">
                        {job.State}
                      </Badge>
                    </div>
                  </AccordionTrigger>
                  <AccordionContent className="pb-0">
                    <JobDetail jobID={job.ID} />
                  </AccordionContent>
                </Card>
              </AccordionItem>
            );
          })}
        </Accordion>
      </>
    </div>
  );
}
