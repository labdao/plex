"use client";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import { BadgeCheck, Dna } from "lucide-react";
import React, { useContext, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { TruncatedString } from "@/components/shared/TruncatedString";
import { Alert } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  AppDispatch,
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

import { ExperimentRenameForm } from "../(forms)/ExperimentRenameForm";
import { aggregateJobStatus, ExperimentStatus } from "../ExperimentStatus";
import { ExperimentUIContext } from "../ExperimentUIContext";
import ExperimentShare from "./ExperimentShare";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@radix-ui/react-tooltip";

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

  const { modelPanelOpen, setModelPanelOpen } = useContext(ExperimentUIContext);

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

  const tooltipMessage = () => {
    if (flow.Public) {
      return (<span>Experiment has already been published<br />and is now read-only</span>)
    } else if (status.status === "completed" || status.status === "partial-failure" || status.status === "failed") {
      return (<span>Click to publish this experiment<br />Warning: Once published, you can't make changes</span>)
    } else {
      return (<span>You can't publish until your experiments finish<br />Please try later</span>)
    }
  };

  return flow.Name ? (
    <div>
      <>
        <Card>
          <CardContent>
            {error && <Alert variant="destructive">{error}</Alert>}
            <div className="flex items-center justify-between">
              <div className="flex grow">
                <ExperimentStatus jobs={flow.Jobs} className="mr-1 mt-3.5" />
                <ExperimentRenameForm
                  initialName={flow.Name}
                  key={flow.Name}
                  inputProps={{ variant: "subtle", className: "text-xl shrink-0 font-heading w-full" }}
                />
              </div>
              <div className="flex justify-end space-x-2 ">
                {userWalletAddress === flow.WalletAddress && (
                  <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <span>
                        <Button 
                          variant="outline" 
                          className="text-sm" 
                          onClick={handlePublish} 
                          disabled={updateLoading || flow.Public || !(status.status === "completed" || status.status === "partial-failure" || status.status === "failed")}
                        >
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
                      </span>
                    </TooltipTrigger>
                    <TooltipContent>
                      <div className="text-xs">
                        <p className="mb-1 font-mono uppercase">
                        {tooltipMessage()}
                        </p>
                      </div>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
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
                <CopyToClipboard string={flow.RecordCID || flow.CID}>
                  {flow.RecordCID ? (
                    <>
                      Record ID: <TruncatedString value={flow.RecordCID} />
                    </>
                  ) : (
                    <>
                      Experiment ID: <TruncatedString value={flow.CID} />
                    </>
                  )}
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
                <Button variant="outline" size="xs" onClick={() => setModelPanelOpen(!modelPanelOpen)}>
                  {flow.Jobs?.[0]?.Tool?.Name}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </>
    </div>
  ) : null;
}
