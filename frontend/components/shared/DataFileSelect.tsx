"use client";

import { useEffect, useState } from "react";
import { useDispatch, useSelector } from 'react-redux';

import { Button } from "@/components/ui/button";
import { Command, CommandEmpty,CommandGroup, CommandInput, CommandItem } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { selectDataFileList, selectDataFileListPagination } from "@/lib/redux/slices/dataFileListSlice/selectors";
import { dataFileListThunk } from "@/lib/redux/slices/dataFileListSlice/thunks";
import { cn } from "@/lib/utils";

interface DataFile {
  CID: string;
  Filename: string;
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
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    setLoading(true);
    setError(null);
  
    const cidPattern = /^[a-zA-Z0-9]{46}$/;
  
    let filenameFilter = globPatterns?.map(pattern => pattern.replace('*', '%')).join('|') || '%';
    let filters: DataFileFilters = {};
  
    if (searchTerm && cidPattern.test(searchTerm)) {
      filters.cid = searchTerm;
      filters.filename = filenameFilter;
    } else {
      filters.filename = searchTerm ? `%${searchTerm}%` : filenameFilter;
    }
  
    // @ts-ignore
    dispatch(dataFileListThunk({ page: currentPage, pageSize: 10, filters }))
      .unwrap()
      .then(() => setLoading(false))
      .catch((fetchError: Error) => {
        setError(fetchError.message);
        setLoading(false);
      });
  }, [dispatch, currentPage, globPatterns, searchTerm]);

  const getDataFileValue = (dataFile: DataFile): string => `${dataFile.CID}/${dataFile.Filename}`;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          className="justify-between w-full p-6 text-base font-normal tracking-normal lowercase font-body"
          variant="outline"
          role="combobox"
          aria-expanded={open}
        >
          <span className={cn(!value && "text-muted-foreground")}>
            {value ? dataFiles.find((dataFile) => getDataFileValue(dataFile) === value)?.Filename : `Select ${label} file...`}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0 " style={{ width: `var(--radix-popover-trigger-width)` }}>
        <Command>
          <CommandInput 
            placeholder="Search files..."
            onValueChange={setSearchTerm}
          />
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
