import { ChevronFirstIcon, ChevronLastIcon, ChevronLeftIcon, ChevronRightIcon, MoreHorizontal } from "lucide-react";
import React from "react";

import { cn } from "@/lib/utils";

import { Button } from "./button";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  className?: string;
}

export const DataPagination: React.FC<PaginationProps> = ({ currentPage, totalPages, onPageChange, className }) => {
  const renderPageNumbers = () => {
    const pageNumbers: (number | string)[] = [];

    pageNumbers.push(1);

    let rangeStart: number, rangeEnd: number;

    if (totalPages <= 10) {
      rangeStart = 2;
      rangeEnd = totalPages - 1;
    } else {
      if (currentPage <= 6) {
        rangeStart = 2;
        rangeEnd = 10;
      } else if (currentPage >= totalPages - 5) {
        rangeStart = totalPages - 9;
        rangeEnd = totalPages - 1;
      } else {
        rangeStart = currentPage - 4;
        rangeEnd = currentPage + 4;
      }
    }

    if (rangeStart > 2) {
      pageNumbers.push("...");
    }

    for (let i = rangeStart; i <= rangeEnd; i++) {
      if (i !== 1 && i !== totalPages) {
        pageNumbers.push(i);
      }
    }

    if (rangeEnd < totalPages - 1) {
      pageNumbers.push("...");
    }

    if (totalPages !== 1) {
      pageNumbers.push(totalPages);
    }

    return pageNumbers;
  };

  return (
    <div className={cn("flex items-center justify-between whitespace-nowrap", className)}>
      <div>
        <Button variant="ghost" size="icon" onClick={() => onPageChange(1)} disabled={currentPage === 1}>
          <ChevronFirstIcon />
        </Button>
        <Button variant="ghost" size="icon" onClick={() => onPageChange(currentPage - 1)} disabled={currentPage === 1}>
          <ChevronLeftIcon />
        </Button>
      </div>
      <div className="flex items-center">
        {renderPageNumbers().map((page, index) =>
          typeof page === "number" ? (
            <Button
              className="h-auto px-3 py-1 shrink min-w-[40px]"
              size="sm"
              variant={currentPage === page ? "outline" : "ghost"}
              key={index}
              onClick={() => onPageChange(page)}
            >
              {page}
            </Button>
          ) : (
            <span key={index} className="flex items-center justify-center w-[40px]">
              <MoreHorizontal />
            </span>
          )
        )}
      </div>
      <div>
        <Button variant="ghost" size="icon" onClick={() => onPageChange(currentPage + 1)} disabled={currentPage === totalPages}>
          <ChevronRightIcon />
        </Button>
        <Button variant="ghost" size="icon" onClick={() => onPageChange(totalPages)} disabled={currentPage === totalPages}>
          <ChevronLastIcon />
        </Button>
      </div>
    </div>
  );
};
