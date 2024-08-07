"use client";

import { getAccessToken } from "@privy-io/react-auth";
import { ColumnDef } from "@tanstack/react-table";
import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { DataTable } from "@/components/ui/data-table";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

export default function ListModelFiles() {
  interface Model {
    ID: string;
    Name: string;
  }

  const columns: ColumnDef<Model>[] = [
    {
      accessorKey: "Name",
      header: "Name",
    },
    {
      accessorKey: "ID",
      header: "Model ID",
      cell: ({ row }) => {
        return (
          <a target="_blank" href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${row.getValue("ID")}/`}>
            {row.getValue("ID")}
          </a>
        );
      },
    },
    {
      accessorKey: "WalletAddress",
      header: "User",
      cell: ({ row }) => {
        return row.getValue("WalletAddress");
      },
    },
  ];

  const [models, setModels] = useState<Model[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const authToken = await getAccessToken();
        const response = await fetch(`${backendUrl()}/models`, {
          headers: {
            Authorization: `Bearer ${authToken}`,
          },
        });

        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }

        const data = await response.json();
        console.log("Fetched models:", data);
        setModels(data);
      } catch (error) {
        console.error("Error fetching models:", error);
      }
    };

    fetchData();
  }, []);

  return (
    <ScrollArea className="w-full bg-white grow">
      <DataTable columns={columns} data={models} /> <ScrollBar orientation="horizontal" />
      <ScrollBar orientation="vertical" />
    </ScrollArea>
  );
}
