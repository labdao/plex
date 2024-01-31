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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AppDispatch, jobDetailThunk, selectJobDetail, selectJobDetailError, selectJobDetailLoading } from "@/lib/redux";

import LogViewer from "./LogViewer";

interface JobDetailProps {
  jobID: number;
}

export default function JobDetail({ jobID }: JobDetailProps) {
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
  }, [dispatch, jobID]);

  console.log(job);

  return (
    <Tabs defaultValue="parameters" className="w-full @container">
      <TabsList className="justify-start w-full px-6 pt-0 rounded-t-none">
        <TabsTrigger value="parameters">Parameters</TabsTrigger>
        <TabsTrigger value="outputs">Outputs</TabsTrigger>
        <TabsTrigger value="inputs">Inputs</TabsTrigger>
        <TabsTrigger value="logs">Logs</TabsTrigger>
      </TabsList>
      <TabsContent value="parameters">parameters</TabsContent>
      <TabsContent value="outputs">outputs</TabsContent>
      <TabsContent value="inputs">inputs</TabsContent>
      <TabsContent value="logs">
        <div className="w-full">
          <LogViewer />
        </div>
      </TabsContent>
    </Tabs>
  );
}
