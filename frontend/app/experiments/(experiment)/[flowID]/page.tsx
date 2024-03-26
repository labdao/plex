"use client";

import { notFound, useParams } from "next/navigation";
import React, { ReactNode, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { tasks } from "@/app/tasks/taskList";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { AppDispatch, flowDetailThunk, resetFlowDetail, selectFlowDetail, selectToolDetail, toolDetailThunk } from "@/lib/redux";

import ExperimentForm from "../(forms)/NewExperimentForm";
import ExperimentDetail from "../ExperimentDetail";
import ExperimentResults from "../ExperimentResults";
import ModelInfo from "../ModelInfo";

type ExperimentDetailProps = {
  params: { flowID: string };
};

export default function Layout({ params }: ExperimentDetailProps) {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const tool = useSelector(selectToolDetail);
  const { flowID } = params;

  let task;
  let isNew = false;
  let breadcrumbItems = [{ name: "Experiments", href: "/experiments" }];

  useEffect(() => {
    if (flowID) {
      if (typeof flowID === "string") {
        dispatch(flowDetailThunk(flowID));
      }
    }
  }, [dispatch, flowID]);

  useEffect(() => {
    //if (!!flow.Jobs?.length) {
    dispatch(toolDetailThunk(flow.Jobs?.[0]?.Tool?.CID));
    //}
  }, [dispatch, flow.Jobs]);

  task = tasks.find((task) => task.slug === tool?.ToolJson?.taskCategory);

  if (task?.name) {
    breadcrumbItems.push({ name: task.name, href: `/experiments/new/${task.slug}` });
  }

  if (flow?.Name) {
    breadcrumbItems.push({ name: flow.Name, href: "" });
  }

  return (
    <>
      <Breadcrumbs items={breadcrumbItems} />
      <TransactionSummaryInfo className="px-4 rounded-b-none" />
      <div className="flex flex-col-reverse min-h-screen lg:flex-row">
        <div className="p-2 space-y-3 shrink-0 grow basis-2/3">
          <ExperimentDetail />
          <ExperimentForm task={task} />
          <ExperimentResults />
        </div>
        <ModelInfo task={task} defaultOpen={false} />
      </div>
    </>
  );
}
