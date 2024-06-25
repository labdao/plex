"use client";

import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { tasks } from "@/app/tasks/taskList";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import PoweredByLogo from "@/components/global/PoweredByLogo";
import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { AppDispatch, flowDetailThunk, selectFlowDetail, selectToolDetail, setToolDetail, toolDetailThunk } from "@/lib/redux";

import RerunExperimentForm from "../(forms)/RerunExperimentForm";
import ExperimentResults from "../(results)/ExperimentResults";
import ModelPanel from "../ModelPanel";
import ExperimentDetail from "./ExperimentDetail";

type ExperimentDetailProps = {
  params: { flowID: string };
};

export default function Layout({ params }: ExperimentDetailProps) {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const tool = useSelector(selectToolDetail);
  const { flowID } = params;

  const task = tasks.find((task) => task.slug === tool?.ToolJson?.taskCategory);
  let breadcrumbItems = [{ name: "Experiments", href: "/experiments" }];

  useEffect(() => {
    if (flowID) {
      dispatch(flowDetailThunk(flowID));
    }
  }, [dispatch, flowID]);

  useEffect(() => {
    if (!!flow.Jobs?.length) {
      //Update redux with the tool stored in the flow rather than making a separate request
      dispatch(setToolDetail(flow.Jobs?.[0]?.Tool));
    }
  }, [dispatch, flow.Jobs]);

  if (task?.name) {
    breadcrumbItems.push({ name: task.name, href: `/experiments/new/${task.slug}` });
  }

  if (flow?.Name) {
    breadcrumbItems.push({ name: flow.Name, href: "" });
  }

  return flowID ? (
    <>
      <ProtectedComponent method="hide" message="Log in to continue">
        <Breadcrumbs items={breadcrumbItems} />
        <TransactionSummaryInfo className="px-4 rounded-b-none" />
        <div className="flex flex-col-reverse min-h-screen lg:flex-row">
          <div className="max-w-4xl p-2 mx-auto space-y-3 shrink-0 grow basis-2/3">
            <ExperimentDetail />
            <RerunExperimentForm key={flow.ID} />
            <ExperimentResults />
            <PoweredByLogo />
          </div>
          <ModelPanel task={task} defaultOpen={false} />
        </div>
      </ProtectedComponent>
    </>
  ) : null;
}