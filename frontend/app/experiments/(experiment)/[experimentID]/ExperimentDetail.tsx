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
  selectExperimentDetail,
  selectExperimentDetailError,
  selectExperimentDetailLoading,
  selectExperimentUpdateError,
  selectExperimentUpdateLoading,
  selectExperimentUpdateSuccess,
  selectUserWalletAddress,
  setExperimentDetailPublic,
} from "@/lib/redux";
import { experimentUpdateThunk } from "@/lib/redux/slices/experimentUpdateSlice/thunks";

import { ExperimentRenameForm } from "../(forms)/ExperimentRenameForm";
import { aggregateJobStatus, ExperimentStatus } from "../ExperimentStatus";
import { ExperimentUIContext } from "../ExperimentUIContext";
import ExperimentShare from "./ExperimentShare";

dayjs.extend(relativeTime);

export default function ExperimentDetail() {
  const dispatch = useDispatch<AppDispatch>();
  const experiment = useSelector(selectExperimentDetail);
  const loading = useSelector(selectExperimentDetailLoading);
  const error = useSelector(selectExperimentDetailError);

  const status = aggregateJobStatus(experiment.Jobs);

  const [isDelaying, setIsDelaying] = useState(false);
  const updateLoading = useSelector(selectExperimentUpdateLoading);
  const updateError = useSelector(selectExperimentUpdateError);
  const updateSuccess = useSelector(selectExperimentUpdateSuccess);

  const userWalletAddress = useSelector(selectUserWalletAddress);

  const { modelPanelOpen, setModelPanelOpen } = useContext(ExperimentUIContext);

  const experimentID = experiment.ID?.toString();

  const handlePublish = () => {
    setIsDelaying(true);
    if (experimentID) {
      dispatch(experimentUpdateThunk({ experimentId: experimentID, updates: { public: true } }))
      .unwrap()
      .then(() => {
        setTimeout(() => {
          dispatch(setExperimentDetailPublic(true));
          setIsDelaying(false);
        }, 2000);
      });
      setTimeout(() => {
        setIsDelaying(false);
      }, 2000);
    }
  };

  const isButtonDisabled = updateLoading || isDelaying;

  return experiment.Name && experimentID ? (
    <div>
      <>
        <Card>
          <CardContent>
            {error && <Alert variant="destructive">{error}</Alert>}
            <div className="flex items-center justify-between">
              <div className="flex grow">
                <ExperimentStatus jobs={experiment.Jobs} className="mr-1 mt-3.5" />
                <ExperimentRenameForm
                  initialName={experiment.Name}
                  experimentId={experimentID}
                  inputProps={{ variant: "subtle", className: "text-xl shrink-0 font-heading w-full" }}
                />
              </div>
              <div className="flex justify-end space-x-2 ">
                {userWalletAddress === experiment.WalletAddress && (
                  <Button variant="outline" className="text-sm" onClick={handlePublish} disabled={updateLoading || experiment.Public}>
                    {updateLoading || isDelaying ? (
                      <>
                        <Dna className="w-4 h-4 ml-2 animate-spin" />
                        <span>Publishing...</span>
                      </>
                    ) : experiment.Public ? (
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
                {experiment.Public && experimentID && (
                  // <Button variant="outline" className="text-sm">
                  //   <Share2 className="w-4 h-4 mr-2" /> Share
                  // </Button>
                  <ExperimentShare experimentID={experimentID} />
                )}
              </div>
            </div>
            <div className="py-4 space-y-1 text-xs pl-7">
              <div className="opacity-70">
                Started by <TruncatedString value={experiment.WalletAddress} trimLength={4} />{" "}
                <span className="text-muted-foreground" suppressHydrationWarning>
                  {dayjs().to(dayjs(experiment.StartTime))}
                </span>
              </div>
              <div className="opacity-50">
                <CopyToClipboard string={experiment.RecordCID || String(experiment.ID)}>
                  {experiment.RecordCID ? (
                    <>
                      Record ID: <TruncatedString value={experiment.RecordCID} />
                    </>
                  ) : (
                    <>
                      Experiment ID: <TruncatedString value={String(experiment.ID)} />
                    </>
                  )}
                </CopyToClipboard>
              </div>
            </div>
            <div className="space-y-2 font-mono text-sm uppercase pl-7">
              <div>
                <strong>Queued: </strong>
                {dayjs(experiment.StartTime).format("YYYY-MM-DD HH:mm:ss")}
              </div>
              {/*@TODO: Endtime currently doesn't show a correct datetime and Runtime is missing
                <div>
                  <strong>Completed: </strong>
                  {dayjs(experiment.EndTime).format("YYYY-MM-DD HH:mm:ss")}
                </div>             
                */}
              <div>
                <strong>Model: </strong>
                <Button variant="outline" size="xs" onClick={() => setModelPanelOpen(!modelPanelOpen)}>
                  {experiment.Jobs?.[0]?.Model?.Name}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </>
    </div>
  ) : null;
}
