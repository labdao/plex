"use client";

import { ColumnDef } from "@tanstack/react-table";
import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { DataTable } from "@/components/ui/data-table";

export default function ListToolFiles() {
  interface Tool {
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
  }

  const columns: ColumnDef<Tool>[] = [
    {
      accessorKey: "Name",
      header: "Name",
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
      accessorKey: "WalletAddress",
      header: "User",
      cell: ({ row }) => {
        return shortenAddressOrCid(row.getValue("WalletAddress"));
      }
    },
  ];

  const [tools, setTools] = useState<Tool[]>([]);

  const [sorting, setSorting] = useState([{ id: "Name", desc: false }])

  useEffect(() => {
    fetch(`${backendUrl()}/tools`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        console.log("Fetched tools:", data);
        setTools(data);
      })
      .catch((error) => {
        console.error("Error fetching tools:", error);
      });
  }, []);

  return (
    <div className="border rounded-lg overflow-hidden">
      <DataTable columns={columns} data={tools} sorting={sorting}/>
    </div>
  );
}
