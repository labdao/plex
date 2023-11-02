"use client";

import { Check, ChevronsUpDown } from "lucide-react";
import * as React from "react";

import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

interface DataFileSelectProps {
  onSelect: (value: string) => void;
  value: string;
  dataFiles: { CID: string; Filename: string }[];
  label: string;
}

export function DataFileSelect({ onSelect, value, dataFiles, label }: DataFileSelectProps) {
  const [open, setOpen] = React.useState(false);
  const dataFileValue = (dataFile: { CID: string; Filename: string }) => dataFile.CID + "/" + dataFile.Filename;

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          className="justify-between w-full font-sans font-normal tracking-normal normal-case"
          variant="outline"
          role="combobox"
          aria-expanded={open}
        >
          {value ? dataFiles.find((dataFile) => dataFileValue(dataFile) === value)?.Filename : `Select ${label} file...`}
          <ChevronsUpDown className="w-4 h-4 ml-2 opacity-50 shrink-0" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0 " style={{ width: `var(--radix-popover-trigger-width)` }}>
        <Command>
          <CommandInput placeholder={`Search files...`} />
          <CommandEmpty>No file found.</CommandEmpty>
          <CommandGroup className="overflow-y-auto  max-h-[300px]">
            {dataFiles.map((dataFile) => (
              <CommandItem
                key={dataFile.CID}
                value={dataFileValue(dataFile)}
                onSelect={() => {
                  //currentValue doesn't work here
                  onSelect(dataFileValue(dataFile) === value ? "" : dataFileValue(dataFile));
                  setOpen(false);
                }}
              >
                <Check className={cn("mr-2 h-4 w-4", value === dataFileValue(dataFile) ? "opacity-100" : "opacity-0")} />
                {dataFile.Filename}
              </CommandItem>
            ))}
          </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
