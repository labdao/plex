"use client";

import { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { DataPagination } from "@/components/ui/data-pagination";
import { DataTable } from "@/components/ui/data-table";

export default function ListDataFiles() {
  interface Tag {
    Name: string;
    Type: string;
  }

  interface DataFile {
    CID: string;
    WalletAddress: string;
    Filename: string;
    Timestamp: string;
    Tags: Tag[];
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

  const columns: ColumnDef<DataFile>[] = [
    {
      accessorKey: "Filename",
      header: "File",
      enableSorting: true,
      sortingFn: "alphanumeric",
      cell: ({ row }) => {
        let cid = row.getValue("CID");
        if (!cid) {
          cid = "null";
        }
        return (
          <div>
            <a target="_blank" href={`${backendUrl()}/datafiles/${row.getValue("CID")}/download`}>
              {row.getValue("Filename")}
            </a>
            <div style={{ fontSize: "smaller", marginTop: "4px", color: "gray" }}>{row.getValue("CID")}</div>
          </div>
        );
      },
    },
    {
      accessorKey: "Tags",
      header: "Tags",
      cell: ({ row }) => {
        const tags: Tag[] = row.getValue("Tags") as Tag[];
        return (
          <div>
            {tags.map((tag, index) => (
              <div key={index}>{tag.Name}</div>
            ))}
          </div>
        );
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
      accessorKey: "Timestamp",
      header: "Date",
      enableSorting: true,
      sortingFn: "datetime",
      cell: ({ row }) => {
        return format(new Date(row.getValue("Timestamp")), "yyyy-MM-dd HH:mm:ss");
      },
    },
  ];

  const [dataFiles, setDataFiles] = useState<DataFile[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const pageSize = 50;
  const [sorting, setSorting] = useState([{ id: "Timestamp", desc: true }]);
  console.log("dataFiles", dataFiles);
  useEffect(() => {
    fetch(`${backendUrl()}/datafiles?page=${currentPage}&pageSize=${pageSize}`)
      .then((response) => response.json())
      .then((responseJson) => {
        setDataFiles(responseJson.data);
        setTotalPages(Math.ceil(responseJson.pagination.totalCount / pageSize));
      })
      .catch((error) => console.error("Error fetching data files:", error));
  }, [currentPage]);

  return (
    <div className="pb-14">
      <DataTable columns={columns} data={dataFiles} sorting={sorting} />
      <DataPagination
        className="absolute bottom-0 z-10 w-full px-4 border-t h-14 bg-background"
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={(page) => setCurrentPage(page)}
      />
    </div>
  );
}
