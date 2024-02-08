// // frontend/app/api/@page.tsx
// "use client";

// import { usePrivy } from "@privy-io/react-auth";
// import { ColumnDef } from "@tanstack/react-table";
// import { format } from "date-fns";
// import React, { useEffect, useState } from "react";
// import { useDispatch, useSelector } from "react-redux";

// import { DataTable } from "@/components/ui/data-table";
// import { Pagination } from "@/components/ui/pagination";
// import { Button } from "@/components/ui/button";
// import { AppDispatch, apiKeyListThunk, selectApiKeyList, ApiKey, addApiKeyThunk } from "@/lib/redux";

// export default function ApiPage() {
//   const dispatch = useDispatch<AppDispatch>();
//   const apiKeys = useSelector(selectApiKeyList);
//   const { getAccessToken } = usePrivy();

//   // Define your columns based on the ApiKey interface
//   const columns: ColumnDef<ApiKey>[] = [
//     {
//         accessorKey: 'key',
//         header: 'API Key',
//         cell: info => <span>{info.getValue()}</span>,
//       },
//       {
//         accessorKey: 'scope',
//         header: 'Scope',
//         cell: info => <span>{info.getValue()}</span>,
//       },
//       {
//         accessorKey: 'createdAt',
//         header: 'Created At',
//         cell: info => <span>{format(new Date(info.getValue()), "yyyy-MM-dd HH:mm:ss")}</span>,
//       },
//       {
//         accessorKey: 'expiresAt',
//         header: 'Expires At',
//         cell: info => <span>{format(new Date(info.getValue()), "yyyy-MM-dd HH:mm:ss")}</span>,
//       },
//   ];


//   useEffect(() => {
//     dispatch(apiKeyListThunk());
//   }, [dispatch]);

//   const handleGenerateApiKey = async () => {
//     // Dispatch the thunk to generate a new API key
//     dispatch(addApiKeyThunk());
//   };

//   return (
//     <div>
//       <div className="border rounded-lg overflow-hidden">
//         <DataTable columns={columns} data={apiKeys} sorting={sorting} />
//       </div>
//       <div className="my-4">
//         <Button onClick={handleGenerateApiKey}>Add API Key</Button>
//       </div>
//       <Pagination
//         currentPage={currentPage}
//         totalPages={totalPages}
//         onPageChange={(page) => setCurrentPage(page)}
//       />
//     </div>
//   );
// }