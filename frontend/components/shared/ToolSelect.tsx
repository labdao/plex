"use client";

import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AppDispatch, selectToolList, selectToolListError, toolListThunk } from "@/lib/redux";

interface ToolSelectProps {
  onChange: (value: string) => void;
  taskSlug?: string;
  defaultValue?: string;
}

export function ToolSelect({ onChange, taskSlug, defaultValue }: ToolSelectProps) {
  const dispatch = useDispatch<AppDispatch>();

  const tools = useSelector(selectToolList);
  const toolListError = useSelector(selectToolListError);

  useEffect(() => {
    if (taskSlug) {
      dispatch(toolListThunk(taskSlug));
    } else {
      dispatch(toolListThunk());
    }
  }, [dispatch, taskSlug]);

  const handleSelectionChange = (value: string) => {
    onChange(value);
  };

  return (
    <Select onValueChange={handleSelectionChange} defaultValue={defaultValue}>
      <SelectTrigger>
        <SelectValue placeholder="Select a model" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {tools.map((tool, index) => {
            return (
              <SelectItem key={index} value={tool?.CID}>
                {tool?.ToolJson?.author || "unknown"}/{tool.Name}
              </SelectItem>
            );
          })}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
