"use client";

import { ColumnDef } from "@tanstack/react-table";
import Link from "next/link";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { DataTable } from "@/components/ui/data-table";
import { AppDispatch, flowListThunk, selectFlowList, selectWalletAddress } from "@/lib/redux";

export default function ListFlowFiles() {
  interface Flow {
    CID: string;
    Name: string;
    WalletAddress: string;
  }

  const shortenAddressOrCid = (addressOrCid: string) => {
    if (addressOrCid.length) {
      return `${addressOrCid.substring(0, 6)}...${addressOrCid.substring(addressOrCid.length - 4)}`;
    } else {
      return "";
    }
  };

  const columns: ColumnDef<Flow>[] = [
    {
      accessorKey: "Name",
      header: "Experiment",
      enableSorting: true,
      cell: ({ row }) => {
        return <Link href={`/experiments/${row.getValue("CID")}`}>{row.getValue("Name")}</Link>;
      },
    },
    {
      accessorKey: "CID",
      header: "CID",
      cell: ({ row }) => {
        return shortenAddressOrCid(row.getValue("CID"));
      },
    },
    {
      accessorKey: "WalletAddress",
      header: "User",
      cell: ({ row }) => {
        return shortenAddressOrCid(row.getValue("WalletAddress"));
      },
    },
  ];

  const dispatch = useDispatch<AppDispatch>();
  const flows = useSelector(selectFlowList);
  const walletAddress = useSelector(selectWalletAddress);

  const [sorting, setSorting] = useState([{ id: "Name", desc: false }]);

  useEffect(() => {
    if (walletAddress) {
      dispatch(flowListThunk(walletAddress));
    }
  }, [dispatch, walletAddress]);

  return (
    <ProtectedComponent method="hide" message="Log in to view your experiments">
      <div className="overflow-hidden border rounded-lg">
        <DataTable columns={columns} data={flows} sorting={sorting} />
      </div>
    </ProtectedComponent>
  );
}
