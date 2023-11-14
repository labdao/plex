"use client";

import { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { DataTable } from "@/components/ui/data-table";

export default function ListDataFiles() {
  interface DataFile {
    CID: string;
    WalletAddress: string;
    Filename: string;
    Timestamp: string;
  }

  const shortenAddressOrCid = (addressOrCid: string) => {
    if (addressOrCid.length) {
      return `${addressOrCid.substring(0, 6)}...${addressOrCid.substring(addressOrCid.length - 4)}`;
    } else {
      return "";
    }
  }

  const columns: ColumnDef<DataFile>[] = [
    {
      accessorKey: "Filename",
      header: "Filename",
      enableSorting: true,
      sortingFn: "alphanumeric",
      cell: ({ row }) => {
        return (
          <a target="_blank" href={`${backendUrl()}/datafiles/${row.getValue("CID")}/download`}>
            {row.getValue("Filename")}
          </a>
        );
      },
    },
    {
      accessorKey: "CID",
      header: "CID",
      cell: ({ row }) => {
        return (
          <a target="_blank" href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${row.getValue("CID")}/`}>
            {shortenAddressOrCid(row.getValue("CID"))}
          </a>
        );
      },
    },
    {
      accessorKey: "Timestamp",
      header: "Date",
      enableSorting: true,
      sortingFn: "datetime",
      cell: ({ row }) => {
        return format(new Date(row.getValue("Timestamp")), "yyyy-MM-dd HH:mm:ss")
      }
    },
  ];

  const [datafiles, setDataFiles] = useState<DataFile[]>([]);
  const [sorting, setSorting] = useState([{ id: "Timestamp", desc: true }]);

  useEffect(() => {
    fetch(`${backendUrl()}/datafiles`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        console.log("Fetched datafiles:", data);
        setDataFiles(data);
      });
  }, []);

  return (
    <div className="border rounded-lg overflow-hidden">
      <DataTable columns={columns} data={datafiles} sorting={sorting} />
    </div>
  );
}
