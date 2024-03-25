"use client";

import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AppDispatch, selectToolDetail, selectToolList, selectToolListError, toolDetailThunk, toolListThunk } from "@/lib/redux";

interface ToolSelectProps {
  onChange: (value: string) => void;
  taskSlug?: string;
}

export function ToolSelect({ onChange, taskSlug }: ToolSelectProps) {
  const dispatch = useDispatch<AppDispatch>();

  const tool = useSelector(selectToolDetail);
  const tools = useSelector(selectToolList);
  const toolListError = useSelector(selectToolListError);

  useEffect(() => {
    if (taskSlug) {
      dispatch(toolListThunk(taskSlug));
    } else {
      dispatch(toolListThunk());
    }
  }, [dispatch, taskSlug]);

  useEffect(() => {
    const defaultTool = tools.find((tool) => tool.DefaultTool === true);
    const defaultToolCID = defaultTool?.CID;

    if (defaultToolCID) {
      dispatch(toolDetailThunk(defaultToolCID));
    }
  }, [dispatch, tools]);

  const handleSelectionChange = (value: string) => {
    onChange(value);
  };

  return (
    <Select onValueChange={handleSelectionChange} value={tool?.CID}>
      <SelectTrigger className="font-mono text-sm font-bold uppercase rounded-full border-primary hover:bg-primary/10">
        <SelectValue placeholder="Select a model" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {tools.map((tool, index) => {
            return (
              <SelectItem key={index} value={tool?.CID}>
                {tool.Name}
              </SelectItem>
            );
          })}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
