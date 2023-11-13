"use client";

import { Check, ChevronsUpDown } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

import { listDataFiles } from "@/lib/redux/slices/dataFileListSlice/asyncActions";

interface DataFile {
  CID: string;
  Filename: string;
}

interface DataFileSelectProps {
  onChange: (value: string) => void;
  value: string;
  label: string;
  globPatterns?: string[]; // assuming globPatterns is an array of string patterns
}

export function DataFileSelect({ onChange, value, label, globPatterns }: DataFileSelectProps) {
  const [dataFiles, setDataFiles] = useState<DataFile[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [open, setOpen] = useState<boolean>(false);

  useEffect(() => {
    setLoading(true);
    setError(null);

    listDataFiles(globPatterns)
      .then((fetchedDataFiles: DataFile[]) => {
        setDataFiles(fetchedDataFiles);
        setLoading(false);
      })
      .catch((fetchError: Error) => {
        setError(fetchError.message);
        setLoading(false);
      });

    return () => {
      setDataFiles([]);
      setError(null);
      setLoading(false);
    };
  }, [globPatterns]);

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
          <ChevronsUpDown className="w-4 h-4 ml-2 shrink-0" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0 " style={{ width: `var(--radix-popover-trigger-width)` }}>
        <Command>
          <CommandInput placeholder={`Search files...`} />
          <CommandEmpty className="lowercase">No file found.</CommandEmpty>
          {error && <CommandItem>Error fetching data files.</CommandItem>}
          <CommandGroup className="overflow-y-auto max-h-[300px]">
            {dataFiles.map((dataFile) => (
              <CommandItem
                key={dataFile.CID}
                value={getDataFileValue(dataFile)}
                onSelect={() => {
                  //currentValue doesn't work here
                  onChange(getDataFileValue(dataFile) === value ? "" : getDataFileValue(dataFile));
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
