"use client";

import { DataTable } from "@/components/ui/data-table";
import { ColumnDef } from "@tanstack/react-table";
import Link from "next/link";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { AppDispatch, flowListThunk, selectFlowList } from "@/lib/redux";

export default function ListFlowFiles() {
  interface Flow {
    CID: string;
    Name: string;
    WalletAddress: string;
  }

  const columns: ColumnDef<Flow>[] = [
    {
      accessorKey: "Name",
      header: "Name",
      cell: ({ row }) => {
        return <Link href={`/experiments/${row.getValue("CID")}`}>{row.getValue("Name")}</Link>;
      },
    },
    {
      accessorKey: "CID",
      header: "CID",
      cell: ({ row }) => {
        return (
          <a target="_blank" href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${row.getValue("CID")}/`}>
            {row.getValue("CID")}
          </a>
        );
      },
    },
    {
      accessorKey: "WalletAddress",
      header: "Uploader Wallet Address",
    },
  ];

  const dispatch = useDispatch<AppDispatch>();

  const flows = useSelector(selectFlowList);

  useEffect(() => {
    dispatch(flowListThunk());
  }, [dispatch]);

  return (
    <div className="border rounded-lg overflow-hidden">
      <DataTable columns={columns} data={flows} />
    </div>
  );
}
