import { Column } from "@tanstack/react-table";
import { ChevronDownIcon, ChevronsUpDownIcon, ChevronUpIcon, EyeOffIcon } from "lucide-react";

import { cn } from "@/lib/utils";

import { Button } from "./button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from "./dropdown-menu";

interface DataTableColumnHeaderProps<TData, TValue> extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>;
  title: string;
}

export function DataTableColumnHeader<TData, TValue>({ column, title, className }: DataTableColumnHeaderProps<TData, TValue>) {
  if (!column.getCanSort()) {
    return <div className={cn(className)}>{title}</div>;
  }

  return (
    <div className={cn("flex items-center space-x-2 whitespace-nowrap", className)}>
      <span>{title}</span>
      {column.getIsSorted() === "desc" ? (
        <Button variant="ghost" size="icon" onClick={() => column.toggleSorting(false)}>
          <ChevronDownIcon className="w-4 h-4 ml-2" />
        </Button>
      ) : column.getIsSorted() === "asc" ? (
        <Button variant="ghost" size="icon" onClick={() => column.toggleSorting(true)}>
          <ChevronUpIcon className="w-4 h-4 ml-2" />
        </Button>
      ) : (
        <Button variant="ghost" size="icon" onClick={() => column.toggleSorting(false)}>
          <ChevronsUpDownIcon className="w-4 h-4 ml-2" />
        </Button>
      )}
    </div>
  );
}
