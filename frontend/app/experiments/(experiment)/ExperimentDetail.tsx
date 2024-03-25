"use client";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { BadgeCheck, Dna, Share2 } from "lucide-react";
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
  setFlowDetailPublic,
} from "@/lib/redux";
import { flowUpdateThunk } from "@/lib/redux/slices/flowUpdateSlice/thunks";

import ExperimentShare from "./ExperimentShare";
import { aggregateJobStatus, ExperimentStatus } from "./ExperimentStatus";
import JobDetail from "./JobDetail";
import MetricsVisualizer from "./MetricsVisualizer";

dayjs.extend(relativeTime);

export default function ExperimentDetail() {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const loading = useSelector(selectFlowDetailLoading);
  const error = useSelector(selectFlowDetailError);

  const status = aggregateJobStatus(flow.Jobs);

  const [isDelaying, setIsDelaying] = useState(false);
  const updateLoading = useSelector(selectFlowUpdateLoading);
  const updateError = useSelector(selectFlowUpdateError);
  const updateSuccess = useSelector(selectFlowUpdateSuccess);

  const userWalletAddress = useSelector(selectUserWalletAddress);

  const experimentID = flow.ID?.toString();

  const handlePublish = () => {
    setIsDelaying(true);
    if (experimentID) {
      dispatch(flowUpdateThunk({ flowId: experimentID })).then(() => {
        setTimeout(() => {
          dispatch(setFlowDetailPublic(true));
          setIsDelaying(false);
        }, 2000);
      });
      setTimeout(() => {
        setIsDelaying(false);
      }, 2000);
    }
  };

  const isButtonDisabled = updateLoading || isDelaying;

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
        <Card>
          <CardContent>
            {error && <Alert variant="destructive">{error}</Alert>}
            <div className="flex items-center justify-between">
              <div className="flex text-xl">
                <ExperimentStatus jobs={flow.Jobs} className="mr-4 mt-2.5" />
                <span className="font-heading">{flow.Name}</span>
              </div>
              <div className="flex justify-end space-x-2 ">
                {userWalletAddress === flow.WalletAddress && (
                  <Button variant="outline" className="text-sm" onClick={handlePublish} disabled={updateLoading || flow.Public}>
                    {updateLoading || isDelaying ? (
                      <>
                        <Dna className="w-4 h-4 ml-2 animate-spin" />
                        <span>Publishing...</span>
                      </>
                    ) : flow.Public ? (
                      <>
                        <BadgeCheck className="w-4 h-4 mr-2" /> Published
                      </>
                    ) : (
                      <>
                        <Dna className="w-4 h-4 mr-2" /> Publish
                      </>
                    )}
                  </Button>
                )}
                {flow.Public && experimentID && (
                  // <Button variant="outline" className="text-sm">
                  //   <Share2 className="w-4 h-4 mr-2" /> Share
                  // </Button>
                  <ExperimentShare experimentID={experimentID} />
                )}
              </div>
            </div>
            <div className="py-4 space-y-1 text-xs pl-7">
              <div className="opacity-70">
                Started by <TruncatedString value={flow.WalletAddress} trimLength={4} />{" "}
                <span className="text-muted-foreground" suppressHydrationWarning>
                  {dayjs().to(dayjs(flow.StartTime))}
                </span>
              </div>
              <div className="opacity-50">
                <CopyToClipboard string={flow.CID}>
                  Experiment ID: <TruncatedString value={flow.CID} />
                </CopyToClipboard>
              </div>
            </div>
            <div className="space-y-2 font-mono text-sm uppercase pl-7">
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
                <Button size="xs" variant="outline">
                  {flow.Jobs?.[0]?.Tool?.Name}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </>
    </div>
  );
}
