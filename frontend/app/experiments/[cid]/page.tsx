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
import { AppDispatch, flowDetailThunk, flowPatchDetailThunk, selectFlowDetail, selectFlowDetailError, selectFlowDetailLoading } from "@/lib/redux";

export default function FlowDetail() {
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
      accessorKey: "BacalhauJobID",
      header: "Bacalhau ID",
      cell: ({ row }) => {
        return <Link href={`/jobs/${row.getValue("BacalhauJobID")}/`}>{row.getValue("BacalhauJobID")}</Link>;
      },
    },
    {
      accessorKey: "Tool.Name",
      header: "Tool Name",
    },
    {
      accessorKey: "Tool.CID",
      accessorFn: (row) => row.Tool.CID,
      header: "Tool CID",
      cell: ({ cell }) => {
        return (
          <a target="_blank" href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}/${cell.getValue()}`}>
            {cell.getValue() as string}
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
    const flowCid = window.location.href.split("/").pop();
    if (flowCid) {
      dispatch(flowDetailThunk(flowCid));
    }
  }, [dispatch]);

  return (
    <div className="container mt-8">
      <Card className="pt-4">
        <CardTitle className="px-4 flex justify-between items-center border-b pb-4 mb-4">
          {flow.Name}{" "}
          <div className="flex gap-2">
            <Button variant="ghost" onClick={() => dispatch(flowPatchDetailThunk(flow.CID))} disabled={loading}>
              <RefreshCcw size={20} className="mr-2" /> {loading ? "Updating..." : "Update"}
            </Button>
            <Button variant="ghost" asChild>
              <a href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${flow.CID}/`}>
                <ExternalLink size={20} className="mr-2" />
                Open
              </a>
            </Button>
          </div>
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
        <div className="p-4 font-medium uppercase">Jobs</div>
        <div className="bg-gray-50">
          <DataTable columns={columns} data={flow.Jobs} />
        </div>
      </Card>
    </div>
  );
}
