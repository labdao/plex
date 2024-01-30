"use client";

import { ColumnDef } from "@tanstack/react-table";
import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { PageLoader } from "@/components/shared/PageLoader";
import { TruncatedString } from "@/components/shared/TruncatedString";
import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
import { AppDispatch, flowDetailThunk, selectFlowDetail, selectFlowDetailError, selectFlowDetailLoading } from "@/lib/redux";

import { ExperimentStatus } from "./ExperimentStatus";

dayjs.extend(relativeTime);

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

  console.log(flow);

  return (
    <div>
      {loading && (
        <div className="absolute inset-0 z-50 bg-background/70">
          <PageLoader />
        </div>
      )}
      <Card>
        <CardContent className="pb-0">
          {error && <Alert variant="destructive">{error}</Alert>}
          <div className="flex text-xl">
            <ExperimentStatus jobs={flow.Jobs} className="mr-2 mt-2.5" />
            <span className="font-heading">{flow.Name}</span>
          </div>
          <div className="py-4 pl-5 space-y-1 text-xs">
            <div className="opacity-70">
              started by <TruncatedString value={flow.WalletAddress} trimLength={4} />{" "}
              <span className="text-muted-foreground" suppressHydrationWarning>
                {dayjs().to(dayjs(flow.StartTime))}
              </span>
            </div>
            <div className="opacity-50">
              <CopyToClipboard string={flow.CID}>
                cid: <TruncatedString value={flow.CID} />
              </CopyToClipboard>
            </div>
          </div>
        </CardContent>
        <CardContent className="border-t bg-muted/50 border-border/50">
          <div className="pl-5 space-y-2 font-mono text-sm uppercase">
            <div>
              <strong>Queued: </strong>
              {dayjs(flow.StartTime).format("YYYY-MM-DD HH:mm:ss")}
            </div>
            {/*@TODO: Completed currently doesn't show a correct datetime and Runtime is missing
            <dt>Completed:</dt>
            <dd>{dayjs(flow.EndTime).format("YYYY-MM-DD HH:mm:ss")}</dd>
             */}
            <div>
              <strong>Model: </strong>
              <Button size="xs" variant="outline" asChild>
                <Link href={`/tasks/protein-binder-design/${flow.Jobs?.[0]?.Tool?.CID}`}>{flow.Jobs?.[0]?.Tool?.Name}</Link>
              </Button>
            </div>
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
