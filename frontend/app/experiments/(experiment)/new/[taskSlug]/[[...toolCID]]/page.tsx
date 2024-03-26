"use client";

import { notFound, useParams } from "next/navigation";
import React, { ReactNode, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { tasks } from "@/app/tasks/taskList";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { AppDispatch, flowDetailThunk, resetFlowDetail, selectFlowDetail, selectToolDetail } from "@/lib/redux";

import NewExperimentForm from "../../../(forms)/NewExperimentForm";
import ExperimentResults from "../../../ExperimentResults";
import ModelInfo from "../../../ModelInfo";

type NewExperimentProps = {
  params: { taskSlug: string };
};

export default function NewExperiment({ params }: NewExperimentProps) {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const { taskSlug } = params;
  const task = tasks.find((task) => task.slug === taskSlug);

  useEffect(() => {
    dispatch(resetFlowDetail());
  }, [dispatch]);

  let breadcrumbItems = [
    { name: "Experiments", href: "/experiments" },
    {
      name: "New",
      href: "",
    },
  ];

  if (task?.name) {
    breadcrumbItems.push({ name: task.name, href: `/experiments/new/${task.slug}` });
  }

  if (flow?.Name) {
    breadcrumbItems.push({ name: flow.Name, href: "" });
  }

  return (
    <>
      <Breadcrumbs items={breadcrumbItems} />
      <ProtectedComponent method="hide" message="Log in to run an experiment">
        <div>
          <TransactionSummaryInfo className="px-4 rounded-b-none" />
          <div className="flex flex-col-reverse min-h-screen lg:flex-row">
            <div className="p-2 space-y-3 shrink-0 grow basis-2/3">
              <NewExperimentForm task={task} />
              <ExperimentResults />
            </div>
            <ModelInfo task={task} defaultOpen showSelect />
          </div>
        </div>
      </ProtectedComponent>
    </>
  );
}
