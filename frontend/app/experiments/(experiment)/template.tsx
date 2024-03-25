"use client";

import { notFound, useParams } from "next/navigation";
import React, { ReactNode, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { tasks } from "@/app/tasks/taskList";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { AppDispatch, flowDetailThunk, resetFlowDetail, selectFlowDetail, selectToolDetail } from "@/lib/redux";

import ExperimentDetail from "./ExperimentDetail";
import ExperimentForm from "./ExperimentForm";
import ExperimentResults from "./ExperimentResults";
import ModelInfo from "./ModelInfo";
type LayoutProps = {
  children: ReactNode;
  list: any;
  add: any;
  params: { slug: string };
};

export default function Layout({ children }: LayoutProps) {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const tool = useSelector(selectToolDetail);
  const { taskSlug, flowID } = useParams();

  let activeTool = tool;
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

  if (flowID) {
    // We are on the results page
    task = tasks.find((task) => task.slug === activeTool?.ToolJson?.taskCategory);
    isNew = false;
    activeTool = flow?.Jobs?.[0]?.Tool;
  } else {
    // We are on the new experiment page
    task = tasks.find((task) => task.slug === taskSlug);
    isNew = true;
    dispatch(resetFlowDetail());
  }

  if (task?.name) {
    breadcrumbItems.push({ name: task.name, href: `/experiments/${task.slug}` });
  }

  if (isNew) {
    breadcrumbItems.push({
      name: "New",
      href: "",
    });
  }

  return activeTool ? (
    <div className="relative flex flex-col h-screen max-w-full grow">
      <Breadcrumbs items={breadcrumbItems} />
      <TransactionSummaryInfo className="px-4 rounded-b-none" />
      <div className="flex flex-col-reverse min-h-screen lg:flex-row">
        <div className="p-2 space-y-3 shrink-0 grow basis-2/3">
          {!isNew && <ExperimentDetail />}
          <ExperimentForm task={task} showNameField={isNew} />
          <ExperimentResults />
        </div>
        <ModelInfo tool={activeTool} task={task} defaultOpen={isNew} showSelect={isNew} />
      </div>
    </div>
  ) : null;
}
