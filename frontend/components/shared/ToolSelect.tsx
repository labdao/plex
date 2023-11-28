"use client";

import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AppDispatch, selectToolList, selectToolListError, toolListThunk } from "@/lib/redux";

interface ToolSelectProps {
  onChange: (value: string) => void;
  defaultValue: string;
}

export function ToolSelect({ onChange, defaultValue }: ToolSelectProps) {
  const dispatch = useDispatch<AppDispatch>();

  const tools = useSelector(selectToolList);
  const toolListError = useSelector(selectToolListError);

  useEffect(() => {
    dispatch(toolListThunk());
  }, [dispatch]);

  return (
    <Select onValueChange={onChange} defaultValue={defaultValue}>
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
