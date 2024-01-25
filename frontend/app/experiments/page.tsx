"use client";

import { usePrivy } from "@privy-io/react-auth";
import { ColumnDef } from "@tanstack/react-table";
import dayjs from "dayjs";
import { PlusIcon } from "lucide-react";
import Link from "next/link";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import TasksMenu from "@/app/tasks/TasksMenu";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { TruncatedString } from "@/components/shared/TruncatedString";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/ui/data-table";
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header";
import { AppDispatch, Flow, flowListThunk, selectFlowList, selectFlowListLoading } from "@/lib/redux";

import { ExperimentStatus } from "./ExperimentStatus";

export default function ListExperiments() {
  const { user } = usePrivy();

  const columns: ColumnDef<Flow>[] = [
    {
      accessorKey: "ID",
    },
    {
      accessorKey: "Jobs",
      header: "Status",
      maxSize: 42,
      cell: ({ row }) => {
        return (
          <>
            <ExperimentStatus jobs={row.getValue("Jobs")} />
          </>
        );
      },
    },
    {
      accessorKey: "Name",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Experiment" />,
      cell: ({ row }) => {
        return <Link href={`/experiments/${row.getValue("ID")}`}>{row.getValue("Name")}</Link>;
      },
    },
    {
      accessorKey: "CID",
      header: "CID",
      cell: ({ row }) => {
        return <TruncatedString value={row.getValue("CID")} />;
      },
    },
    {
      accessorKey: "StartTime",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Start Time" />,
      sortingFn: "datetime",
      cell: ({ row }) => {
        return row.getValue("StartTime") ? dayjs(row.getValue("StartTime")).format("YYYY-MM-DD HH:mm:ss") : "-";
      },
    },
  ];

  const dispatch = useDispatch<AppDispatch>();
  const flows = useSelector(selectFlowList);
  const loading = useSelector(selectFlowListLoading);
  const walletAddress = user?.wallet?.address;

  useEffect(() => {
    if (walletAddress) {
      dispatch(flowListThunk(walletAddress));
    }
  }, [dispatch, walletAddress]);

  return (
    <>
      <Breadcrumbs
        items={[{ name: "My Experiments", href: "/experiments" }]}
        actions={
          <div>
            <TasksMenu
              trigger={
                <Button size="sm">
                  <PlusIcon /> Run Experiment
                </Button>
              }
            />
          </div>
        }
      />
      <ProtectedComponent method="hide" message="Log in to view your experiments">
        <DataTable columns={columns} data={flows} sorting={[{ id: "StartTime", desc: true }]} loading={loading} />
      </ProtectedComponent>
    </>
  );
}
