"use client";

import { UploadIcon } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import AddDataFileForm from "@/app/data/AddDataFileForm";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandSeparator } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { listDataFiles } from "@/lib/redux/slices/dataFileListSlice/asyncActions";
import { selectDataFileList, selectDataFileListPagination } from "@/lib/redux/slices/dataFileListSlice/selectors";
import { dataFileListThunk } from "@/lib/redux/slices/dataFileListSlice/thunks";
import { cn } from "@/lib/utils";

interface DataFile {
  CID: string;
  Filename: string;
  S3URI: string;
}

interface DataFileFilters {
  filename?: string;
  cid?: string;
}

interface DataFileSelectProps {
  onValueChange: (value: string) => void;
  value: string;
  label: string;
  globPatterns?: string[];
}

export function DataFileSelect({ onValueChange, value, label, globPatterns }: DataFileSelectProps) {
  const dispatch = useDispatch();
  const dataFiles = useSelector(selectDataFileList);
  const pagination = useSelector(selectDataFileListPagination);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [open, setOpen] = useState<boolean>(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchTerm, setSearchTerm] = useState("");

  const cidPattern = /^[a-zA-Z0-9]{46}$/;

  let filenameFilter = globPatterns?.map((pattern) => pattern.replace("*", "%")).join("|") || "%";
  let fileAcceptExtensions = globPatterns?.map((pattern) => pattern.replace("*", "")).join(",") || "*";
  let filters: DataFileFilters = useMemo(() => ({}), []);

  if (searchTerm && cidPattern.test(searchTerm)) {
    filters.cid = searchTerm;
    filters.filename = filenameFilter;
  } else {
    filters.filename = searchTerm ? `%${searchTerm}%` : filenameFilter;
  }

  useEffect(() => {
    setLoading(true);
    setError(null);
    // @ts-ignore
    dispatch(dataFileListThunk({ page: currentPage, pageSize: 10, filters }))
      .unwrap()
      .then(() => setLoading(false))
      .catch((fetchError: Error) => {
        setError(fetchError.message);
        setLoading(false);
      });
  }, [dispatch, currentPage, filters, searchTerm, open]);

  const getDataFileValue = (dataFile: DataFile): string => `${dataFile?.S3URI}`;

  const handleUpload = async (cid: string) => {
    // Since we only know the cid and not the filename after upload,
    // search for the file and set the input value from what we find
    const files = await listDataFiles({ page: 1, pageSize: 10, filters: { cid: cid, filename: filenameFilter } });
    if (files?.data?.length && files?.data?.[0]?.CID === cid) {
      onValueChange(getDataFileValue(files?.data[0]));
    }
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <div className="relative flex items-center gap-2">
        <PopoverTrigger asChild>
          <Button className="justify-between p-3 text-base grow" variant="input" role="combobox" aria-expanded={open}>
            <span className={cn(!value && "text-muted-foreground")}>{value ? value.split("/")?.[1] : `Select ${label} file...`}</span>
          </Button>
        </PopoverTrigger>
        <AddDataFileForm
          trigger={
            <Button variant="secondary" size="xs" className="absolute z-30 right-5 group hover:bg-secondary">
              <span className="w-0 overflow-hidden group-hover:w-12 transition-[width]">Upload</span>
              <UploadIcon className="w-4 h-4 ml-2" />
            </Button>
          }
          accept={fileAcceptExtensions}
          onUpload={handleUpload}
        />
      </div>
      <PopoverContent className="p-0 " style={{ width: `var(--radix-popover-trigger-width)` }}>
        <Command>
          <CommandInput placeholder="Search files..." onValueChange={setSearchTerm} value={searchTerm} />
          <CommandEmpty>No file found.</CommandEmpty>
          {error && <CommandItem>Error fetching data files.</CommandItem>}
          <CommandGroup>
            {dataFiles.map((dataFile) => (
              <CommandItem
                key={dataFile.CID}
                value={getDataFileValue(dataFile)}
                onSelect={() => {
                  onValueChange(getDataFileValue(dataFile) === value ? "" : getDataFileValue(dataFile));
                  setOpen(false);
                }}
              >
                {dataFile.Filename}
              </CommandItem>
            ))}
          </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
