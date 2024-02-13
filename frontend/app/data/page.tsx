"use client";

import { getAccessToken } from "@privy-io/react-auth";
import { ColumnDef } from "@tanstack/react-table";
import dayjs from "dayjs";
import backendUrl from "lib/backendUrl";
import { UploadIcon } from "lucide-react";
import React, { useEffect, useState } from "react";

import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { CopyToClipboard } from "@/components/shared/CopyToClipboard";
import { TruncatedString } from "@/components/shared/TruncatedString";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { DataPagination } from "@/components/ui/data-pagination";
import { DataTable } from "@/components/ui/data-table";
import { DataTableColumnHeader } from "@/components/ui/data-table-column-header";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";

import AddDataFileForm from "./AddDataFileForm";

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

  const columns: ColumnDef<DataFile>[] = [
    {
      accessorKey: "Filename",
      header: ({ column }) => <DataTableColumnHeader column={column} title="File" />,
      sortingFn: "alphanumeric",
      cell: ({ row }) => {
        let cid: string = row.getValue("CID");
        if (!cid) {
          cid = "null";
        }

        const handleDownloadClick = async (event: React.MouseEvent<HTMLAnchorElement>) => {
          event.preventDefault();
          const authToken = await getAccessToken();
          const response = await fetch(`${backendUrl()}/datafiles/${cid}/download`, {
            headers: {
              'Authorization': `Bearer ${authToken}`,
            },
          });
          if (!response.ok) {
            console.error('Failed to download file');
            return;
          } else {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = row.getValue("Filename");
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
          }
        };

        return (
          <div>
            <a target="#" onClick={handleDownloadClick} style={{ cursor: 'pointer' }}>
              <TruncatedString value={row.getValue("Filename")} trimLength={20} />
            </a>
            <div className="text-xs truncate max-w-[10rem] text-muted-foreground/50">
              {/* {row.getValue("CID")} */}
              <CopyToClipboard string={cid}>
                <span className="cursor-pointer">
                  cid: <TruncatedString value={cid} />
                </span>
              </CopyToClipboard>
            </div>
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
          <div className="flex flex-wrap gap-1">
            {tags.map((tag, index) => (
              <Badge variant="outline" key={index}>
                {tag.Name}
              </Badge>
            ))}
          </div>
        );
      },
    },
    {
      accessorKey: "CID",
      header: "CID",
      cell: ({ row }) => {
        const cid: string = row.getValue("CID");
        return <TruncatedString value={row.getValue("CID")} />;
      },
    },
    {
      accessorKey: "Timestamp",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Created" />,
      // @TODO: Need sorting added to API endpoint, this just sorts the current page.
      sortingFn: "datetime",
      cell: ({ row }) => {
        return dayjs(row.getValue("Timestamp")).format("YYYY-MM-DD HH:mm:ss");
      },
    },
  ];

  const [dataFiles, setDataFiles] = useState<DataFile[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const pageSize = 50;

  useEffect(() => {
    const fetchDataFiles = async () => {
      setLoading(true);
      try {
        const authToken = await getAccessToken();
        const response = await fetch(`${backendUrl()}/datafiles?page=${currentPage}&pageSize=${pageSize}`, {
          headers: {
            'Authorization': `Bearer ${authToken}`,
          },
        });
        if (!response.ok) {
          throw new Error('Failed to fetch data files');
        }
        const responseJson = await response.json();
        setDataFiles(responseJson.data);
        setTotalPages(Math.ceil(responseJson.pagination.totalCount / pageSize));
      } catch (error) {
        console.error("Error fetching data files:", error);
      } finally {
        setLoading(false);
      }
    };
  
    fetchDataFiles();
  }, [currentPage]);

  return (
    <div className="relative flex flex-col h-screen max-w-full grow">
      <Breadcrumbs
        items={[{ name: "My Files", href: "/data" }]}
        actions={
          <AddDataFileForm
            trigger={
              <Button size="sm">
                <UploadIcon />
                Upload Files
              </Button>
            }
          />
        }
      />
      <ScrollArea className="bg-white grow w-[calc(100vw-12rem)]">
        <DataTable columns={columns} data={dataFiles} sorting={[{ id: "Timestamp", desc: true }]} loading={loading} />
        <ScrollBar orientation="horizontal" />
        <ScrollBar orientation="vertical" />
      </ScrollArea>

      <DataPagination
        className="absolute bottom-0 z-10 w-full px-4 overflow-hidden border-t h-14 bg-background"
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={(page) => setCurrentPage(page)}
      />
    </div>
  );
}
