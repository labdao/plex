// frontend/app/api/@page.tsx
"use client";

import { usePrivy } from "@privy-io/react-auth";
import { ColumnDef } from "@tanstack/react-table";
import { format } from "date-fns";
import backendUrl from "lib/backendUrl";
import React, { useEffect, useState } from "react";

import { DataTable } from "@/components/ui/data-table";
import { Pagination } from "@/components/ui/pagination";

export default function ApiPage() {
  // Define your API data interface here
  interface ApiData {
    // Your fields here
  }

  // Define your columns based on the ApiData interface
  const columns: ColumnDef<ApiData>[] = [
    // Your column definitions here
  ];

  const [apiData, setApiData] = useState<ApiData[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const pageSize = 50; // Adjust if needed
  const [sorting, setSorting] = useState([{ id: "someField", desc: true }]); // Adjust the sorting field

  const { getAccessToken } = usePrivy();

  useEffect(() => {
    const fetchApiData = async () => {
      // Fetch your API data here
    };

    fetchApiData();
  }, [currentPage, getAccessToken]);

  return (
    <div>
      <div className="border rounded-lg overflow-hidden">
        <DataTable columns={columns} data={apiData} sorting={sorting} />
      </div>
      <Pagination
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={(page) => setCurrentPage(page)}
      />
    </div>
  );
}