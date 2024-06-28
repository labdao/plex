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

import AddFileForm from "./AddFileForm";

export default function ListFiles() {
  interface Tag {
    Name: string;
    Type: string;
  }

  interface File {
    ID: number;
    WalletAddress: string;
    Filename: string;
    Timestamp: string;
    Tags: Tag[];
  }

  const columns: ColumnDef<File>[] = [
    {
      accessorKey: "Filename",
      header: ({ column }) => <DataTableColumnHeader column={column} title="File" />,
      sortingFn: "alphanumeric",
      cell: ({ row }) => {
        let ID: number = row.getValue("ID");
        if (!ID) {
          ID = 0;
        }

        const handleDownloadClick = async (event: React.MouseEvent<HTMLAnchorElement>) => {
          event.preventDefault();
          const authToken = await getAccessToken();
          const response = await fetch(`${backendUrl()}/files/${ID}/download`, {
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
              <CopyToClipboard string={ID.toString()}>
                <span className="cursor-pointer">
                  File ID: <TruncatedString value={ID.toString()} />
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
      accessorKey: "ID",
      header: "File ID",
      cell: ({ row }) => {
        const ID: number = row.getValue("ID");
        return <TruncatedString value={ID.toString()} />;
      },
    },
    {
      accessorKey: "CreatedAt",
      header: ({ column }) => <DataTableColumnHeader column={column} title="Created" />,
      // @TODO: Need sorting added to API endpoint, this just sorts the current page.
      sortingFn: "datetime",
      cell: ({ row }) => {
        return dayjs(row.getValue("CreatedAt")).format("YYYY-MM-DD HH:mm:ss");
      },
    },
  ];

  const [files, setFiles] = useState<File[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const pageSize = 50;

  useEffect(() => {
    const fetchFiles = async () => {
      setLoading(true);
      try {
        const authToken = await getAccessToken();
        const response = await fetch(`${backendUrl()}/files?page=${currentPage}&pageSize=${pageSize}`, {
          headers: {
            'Authorization': `Bearer ${authToken}`,
          },
        });
        if (!response.ok) {
          throw new Error('Failed to fetch files');
        }
        const responseJson = await response.json();
        setFiles(responseJson.data);
        setTotalPages(Math.ceil(responseJson.pagination.totalCount / pageSize));
      } catch (error) {
        console.error("Error fetching files:", error);
      } finally {
        setLoading(false);
      }
    };
  
    fetchFiles();
  }, [currentPage]);

  return (
    <div className="relative flex flex-col h-screen max-w-full grow">
      <Breadcrumbs
        items={[{ name: "Files", href: "/data" }]}
        actions={
          <AddFileForm
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
        <DataTable columns={columns} data={files} sorting={[{ id: "CreatedAt", desc: true }]} loading={loading} />
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
