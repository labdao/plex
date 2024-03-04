"use client";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { ChevronsUpDownIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { PageLoader } from "@/components/shared/PageLoader";
import { TruncatedString } from "@/components/shared/TruncatedString";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import { Alert } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { AppDispatch, flowDetailThunk, selectFlowDetail, selectFlowDetailError, selectFlowDetailLoading } from "@/lib/redux";
import { cn } from "@/lib/utils";

import { aggregateJobStatus, ExperimentStatus } from "./ExperimentStatus";
import JobDetail from "./JobDetail";

dayjs.extend(relativeTime);

export default function ExperimentDetail({ experimentID }: { experimentID: string }) {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const loading = useSelector(selectFlowDetailLoading);
  const error = useSelector(selectFlowDetailError);

  const status = aggregateJobStatus(flow.Jobs);

  useEffect(() => {
    if (experimentID) {
      dispatch(flowDetailThunk(experimentID));
    }
  }, [dispatch, experimentID]);

  useEffect(() => {
    if (["running", "queued"].includes(status.status)) {
      const interval = setInterval(() => {
        console.log("Checking for new results");
        dispatch(flowDetailThunk(experimentID));
      }, 15000);

      return () => clearInterval(interval);
    }
  }, [dispatch, experimentID, status]);

  return (
    <div>
      {loading && (
        <div className="absolute inset-0 z-50 bg-background/70">
          <PageLoader />
        </div>
      )}
      {flow.Name && (
        <>
          <Card>
            <CardContent className="pb-0">
              {error && <Alert variant="destructive">{error}</Alert>}
              <div className="flex text-xl">
                <ExperimentStatus jobs={flow.Jobs} className="mr-2 mt-2.5" />
                <span className="font-heading">{flow.Name}</span>
              </div>
              <div className="py-4 pl-5 space-y-1 text-xs">
                <div className="opacity-70">
                  started by <TruncatedString value={flow.WalletAddress} trimLength={4} />{" "}
                  <span className="text-muted-foreground" suppressHydrationWarning>
                    {dayjs().to(dayjs(flow.StartTime))}
                  </span>
                </div>
                <div className="opacity-50">
                  <CopyToClipboard string={flow.CID}>
                    experiment id: <TruncatedString value={flow.CID} />
                  </CopyToClipboard>
                </div>
              </div>
            </CardContent>
            <CardContent className="border-t bg-muted/30 border-border/50">
              <div className="pl-5 space-y-2 font-mono text-sm uppercase">
                <div>
                  <strong>Queued: </strong>
                  {dayjs(flow.StartTime).format("YYYY-MM-DD HH:mm:ss")}
                </div>
                {/*@TODO: Endtime currently doesn't show a correct datetime and Runtime is missing
                <div>
                  <strong>Completed: </strong>
                  {dayjs(flow.EndTime).format("YYYY-MM-DD HH:mm:ss")}
                </div>             
                */}
                <div>
                  <strong>Model: </strong>
                  <Button size="xs" variant="outline" asChild>
                    <Link href={`/tasks/protein-binder-design/${flow.Jobs?.[0]?.Tool?.CID}`}>{flow.Jobs?.[0]?.Tool?.Name}</Link>
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>

          <Accordion {...(flow?.Jobs?.length > 1 ? { type: "multiple" } : { type: "single", defaultValue: "0", collapsible: true })}>
            {flow.Jobs?.map((job, index) => {
              const validStates = ["queued", "running", "failed", "completed"];
              const status = (validStates.includes(job.State) ? job.State : "unknown") as "queued" | "running" | "failed" | "completed" | "unknown";

              return (
                <AccordionItem value={index.toString()} className="border-0 [&[data-state=open]>div]:shadow-lg" key={job.ID}>
                  <Card className="my-2 shadow-sm">
                    <AccordionTrigger className="flex items-center justify-between w-full px-6 py-3 text-left hover:no-underline [&[data-state=open]]:bg-muted">
                      <div className="flex items-center gap-2">
                        <div className="w-28">
                          <div>condition {index + 1}</div>
                          <div className="flex gap-1 text-xs text-muted-foreground/70">
                            id: {job.BacalhauJobID ? <TruncatedString value={job.BacalhauJobID} /> : "n/a"}
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
      )}
    </div>
  );
}
