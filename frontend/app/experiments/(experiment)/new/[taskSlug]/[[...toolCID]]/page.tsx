"use client";

import React, { ReactNode, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import NewExperimentForm from "@/app/experiments/(experiment)/(forms)/NewExperimentForm";
import ExperimentResults from "@/app/experiments/(experiment)/(results)/ExperimentResults";
import ModelInfo from "@/app/experiments/(experiment)/ModelPanel";
import { tasks } from "@/app/tasks/taskList";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import PoweredByLogo from "@/components/global/PoweredByLogo";
import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { AppDispatch, resetExperimentDetail, selectExperimentDetail } from "@/lib/redux";

type NewExperimentProps = {
  params: { taskSlug: string };
};

export default function NewExperiment({ params }: NewExperimentProps) {
  const dispatch = useDispatch<AppDispatch>();
  const experiment = useSelector(selectExperimentDetail);
  const { taskSlug } = params;
  const task = tasks.find((task) => task.slug === taskSlug);

  useEffect(() => {
    dispatch(resetExperimentDetail());
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

  if (experiment?.Name) {
    breadcrumbItems.push({ name: experiment.Name, href: "" });
  }

  return (
    <>
      <ProtectedComponent method="hide" message="Log in to run an experiment">
        <Breadcrumbs items={breadcrumbItems} />
        <div>
          {/* <TransactionSummaryInfo className="px-4 rounded-b-none" /> */}
          <div className="flex flex-col-reverse min-h-screen lg:flex-row">
            <div className="max-w-4xl p-2 mx-auto space-y-3 shrink-0 grow basis-2/3">
              <NewExperimentForm task={task} />
              <ExperimentResults />
              <PoweredByLogo />
            </div>
            <ModelInfo task={task} defaultOpen showSelect />
          </div>
        </div>
      </ProtectedComponent>
    </>
  );
}
