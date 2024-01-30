"use client";

import { ColumnDef } from "@tanstack/react-table";
import backendUrl from "lib/backendUrl";
import { RefreshCcw } from "lucide-react";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
import { AppDispatch, jobDetailThunk, selectJobDetail, selectJobDetailError, selectJobDetailLoading } from "@/lib/redux";

import LogViewer from "./LogViewer";

export default function JobDetail({ jobID }: { jobID: string }) {
  const dispatch = useDispatch<AppDispatch>();

  const job = useSelector(selectJobDetail);
  const loading = useSelector(selectJobDetailLoading);
  const error = useSelector(selectJobDetailError);

  interface File {
    CID: string;
    Filename: string;
    Tags: Tag[];
  }

  interface Tag {
    Name: string;
    Type: string;
  }

  const shortenAddressOrCid = (addressOrCid: string) => {
    if (addressOrCid) {
      if (addressOrCid.length) {
        return `${addressOrCid.substring(0, 6)}...${addressOrCid.substring(addressOrCid.length - 4)}`;
      } else {
        return "";
      }
    }
  };

  const columns: ColumnDef<File>[] = [
    {
      accessorKey: "Filename",
      header: "Filename",
      enableSorting: true,
      sortingFn: "alphanumeric",
      cell: ({ row }) => (
        <div>
          <a target="_blank" href={`${backendUrl()}/datafiles/${row.getValue("CID")}/download`}>
            {row.getValue("Filename")}
          </a>
          <div style={{ fontSize: "smaller", marginTop: "4px", color: "gray" }}>{row.getValue("CID")}</div>
        </div>
      ),
    },
    {
      accessorKey: "Tags",
      header: "Tags",
      cell: ({ row }) => {
        const tags: Tag[] = row.getValue("Tags") as Tag[];
        if (tags && tags.length > 0) {
          return (
            <div>
              {tags.map((tag, index) => (
                <div key={index}>{tag.Name}</div>
              ))}
            </div>
          );
        }
      },
    },
    {
      accessorKey: "CID",
      header: "CID",
      cell: ({ row }) => {
        return shortenAddressOrCid(row.getValue("CID"));
      },
    },
  ];

  useEffect(() => {
    console.log(`jobId is ${jobID}`);
    if (jobID) {
      dispatch(jobDetailThunk(`${jobID}`));
    }
  }, [dispatch]);

  return (
    <>
      <div className="mt-8">
        <Card className="pt-4">
          <CardTitle className="flex items-center justify-between px-4 pb-4 mb-4 border-b">
            Job {job.BacalhauJobID}{" "}
            <div className="flex gap-2">
              <Button variant="ghost" onClick={() => dispatch(jobDetailThunk(`${job.ID}`))} disabled={loading}>
                <RefreshCcw size={20} className="mr-2" /> {loading ? "Updating..." : "Update"}
              </Button>
            </div>
          </CardTitle>
          <CardContent>
            {error && <Alert variant="destructive">{error}</Alert>}
            <div className="py-4 border-b">
              Bacalhau ID: <strong>{job.BacalhauJobID}</strong>
            </div>
            <div className="py-4 border-b">
              Status: <strong className="capitalize">{job.State}</strong>
            </div>
            <div className="py-4 border-b">
              Error: <strong>{job.Error || "None"}</strong>
            </div>
            <div className="py-4 border-b">
              <strong>
                <a target="_blank" href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${job.ToolID}/`}>
                  🔬 Tool
                </a>
              </strong>
            </div>
            <div className="py-4">
              <strong>
                <a target="_blank" href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${job.FlowID}/`}>
                  🔍 Experimental Parameters
                </a>
              </strong>
            </div>
          </CardContent>
        </Card>
        <Card className="mt-4">
          <div className="p-4 font-bold uppercase">Logs</div>
          <div className="px-4 pb-6 bg-gray-50">
            <LogViewer />
          </div>
        </Card>
        <Card className="mt-4">
          <div className="p-4 font-bold uppercase">Inputs</div>
          <div className="bg-gray-50">
            <DataTable columns={columns} data={job.InputFiles} />
          </div>
        </Card>
        <Card className="mt-4">
          <div className="p-4 font-bold uppercase">Outputs</div>
          <div className="bg-gray-50">
            <DataTable columns={columns} data={job.OutputFiles} />
          </div>
        </Card>
      </div>
    </>
  );
}
