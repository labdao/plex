"use client";

import { ColumnDef } from "@tanstack/react-table";
import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { DataTable } from "@/components/ui/data-table";

export default function ListDataFiles() {
  interface DataFile {
    CID: string;
    WalletAddress: string;
    Filename: string;
  }

  const columns: ColumnDef<DataFile>[] = [
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
    {
      accessorKey: "Filename",
      header: "Filename",
    },
  ];

  const [datafiles, setDataFiles] = useState<DataFile[]>([]);

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
      <DataTable columns={columns} data={datafiles} />
    </div>
  );
}
