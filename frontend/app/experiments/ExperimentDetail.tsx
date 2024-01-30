"use client";

import { ColumnDef } from "@tanstack/react-table";
import { ExternalLink, RefreshCcw } from "lucide-react";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
import { AppDispatch, flowDetailThunk, selectFlowDetail, selectFlowDetailError, selectFlowDetailLoading } from "@/lib/redux";

export default function ExperimentDetail() {
  const dispatch = useDispatch<AppDispatch>();
  const flow = useSelector(selectFlowDetail);
  const loading = useSelector(selectFlowDetailLoading);
  const error = useSelector(selectFlowDetailError);

  interface Job {
    BacalhauJobID: string;
    Tool: {
      CID: string;
      Name: string;
      WalletAddress: string;
    };
    State: string;
  }

  const columns: ColumnDef<Job>[] = [
    {
      accessorKey: "ID",
      header: "Job ID",
      cell: ({ row }) => {
        return <Link href={`/jobs/${row.getValue("ID")}/`}>{row.getValue("ID")}</Link>;
      },
    },
    {
      id: "tool",
      header: "Model",
      cell: ({ row }) => {
        const toolName = row.original.Tool.Name;
        const toolCID = row.original.Tool.CID;
        const toolCIDUrl = `${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}/${toolCID}`;

        return (
          <a href={toolCIDUrl} target="_blank">
            {toolName}
          </a>
        );
      },
    },
    {
      accessorKey: "State",
      header: "Status",
    },
  ];

  useEffect(() => {
    const flowID = window.location.href.split("/").pop();
    if (flowID) {
      dispatch(flowDetailThunk(flowID));
    }
  }, [dispatch]);

  return (
    <div>
      <Card>
        <CardTitle className="flex items-center justify-between px-4 pb-4 mb-4 border-b">
          <span className="font-heading">{flow.Name}</span>
        </CardTitle>
        <CardContent>
          {error && <Alert variant="destructive">{error}</Alert>}
          <div className="py-4 border-b">
            CID: <strong>{flow.CID}</strong>
          </div>
          <div className="py-4">
            Wallet Address: <strong>{flow.WalletAddress}</strong>
          </div>
        </CardContent>
      </Card>
      <Card className="mt-4">
        <div className="p-4 font-bold font-heading">Jobs</div>
        <div className="bg-gray-50">
          <DataTable columns={columns} data={flow.Jobs} />
        </div>
      </Card>
    </div>
  );
}
