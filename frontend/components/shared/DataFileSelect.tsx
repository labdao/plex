"use client";

import { Check, ChevronsUpDown } from "lucide-react";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { AppDispatch, dataFileListThunk, selectDataFileList, selectDataFileListError } from "@/lib/redux";
import { DataFile } from "@/lib/redux/slices/dataFileListSlice/slice";
import { cn } from "@/lib/utils";

interface DataFileSelectProps {
  onSelect: (value: string) => void;
  value: string;
  dataFiles: { CID: string; Filename: string }[];
  label: string;
  globPatterns?: [];
}

export function DataFileSelect({ onSelect, value, label, globPatterns }: DataFileSelectProps) {
  const dispatch = useDispatch<AppDispatch>();

  const dataFiles = useSelector(selectDataFileList);
  const dataFileListError = useSelector(selectDataFileListError);

  const [open, setOpen] = useState(false);

  useEffect(() => {
    dispatch(dataFileListThunk(globPatterns));
  }, [dispatch, globPatterns]);

  const getDataFileValue = (dataFile: { CID: string; Filename: string }) => dataFile.CID + "/" + dataFile.Filename;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          className="justify-between w-full font-sans font-normal tracking-normal normal-case"
          variant="outline"
          role="combobox"
          aria-expanded={open}
        >
          {value ? dataFiles.find((dataFile) => getDataFileValue(dataFile) === value)?.Filename : `Select ${label} file...`}
          <ChevronsUpDown className="w-4 h-4 ml-2 opacity-50 shrink-0" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0 " style={{ width: `var(--radix-popover-trigger-width)` }}>
        <Command>
          <CommandInput placeholder={`Search files...`} />
          <CommandEmpty>No file found.</CommandEmpty>
          {dataFileListError && <CommandItem>Error fetching data files.</CommandItem>}
          <CommandGroup className="overflow-y-auto max-h-[300px]">
            {dataFiles.map((dataFile) => (
              <CommandItem
                key={dataFile.CID}
                value={getDataFileValue(dataFile)}
                onSelect={() => {
                  //currentValue doesn't work here
                  onSelect(getDataFileValue(dataFile) === value ? "" : getDataFileValue(dataFile));
                  setOpen(false);
                }}
              >
                <Check className={cn("mr-2 h-4 w-4", value === getDataFileValue(dataFile) ? "opacity-100" : "opacity-0")} />
                {dataFile.Filename}
              </CommandItem>
            ))}
          </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
