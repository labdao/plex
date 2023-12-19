import backendUrl from "lib/backendUrl"

export const listDataFiles = async ({ page = 1, pageSize = 50, filters = {} }: { page?: number, pageSize?: number, filters?: Record<string, string | undefined> }): Promise<any> => {
  const queryParams = new URLSearchParams({ ...filters, page: page.toString(), pageSize: pageSize.toString() });

  const response = await fetch(`${backendUrl()}/datafiles?${queryParams}`, {
    method: 'Get',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error("Failed to list DataFiles");
  }

  return await response.json();
}

// export const listDataFiles = async (globPatterns?: string[]): Promise<any> => {
//   let dataFiles: any[] = [];

//   if (globPatterns) {
//     for (const pattern of globPatterns) {
//       const filename = pattern.replace(/\*/g, '');

//       const response = await fetch(`${backendUrl()}/datafiles?filename=${encodeURIComponent(filename)}`, {
//         method: 'Get',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//       });

//       if (!response.ok) {
//         throw new Error("Failed to list DataFiles");
//       }

//       const result = await response.json();
//       dataFiles = [...dataFiles, ...result];
//     }
//   } else {
//     const response = await fetch(`${backendUrl()}/datafiles`, {
//       method: 'Get',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//     });

//     if (!response.ok) {
//       throw new Error("Failed to list DataFiles");
//     }

//     dataFiles = await response.json();
//   }

//   return dataFiles;
// }
