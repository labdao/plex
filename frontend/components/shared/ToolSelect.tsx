"use client";

import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AppDispatch, selectToolList, selectToolListError, toolListThunk } from "@/lib/redux";

interface ToolSelectProps {
  onChange: (value: string) => void;
  taskSlug?: string;
}

export function ToolSelect({ onChange, taskSlug }: ToolSelectProps) {
  const dispatch = useDispatch<AppDispatch>();
  const tools = useSelector(selectToolList);
  const [selectedToolCID, setSelectedToolCID] = useState("");
  const toolListError = useSelector(selectToolListError);

  useEffect(() => {
    if (taskSlug) {
      dispatch(toolListThunk(taskSlug));
    } else {
      dispatch(toolListThunk());
    }
  }, [dispatch, taskSlug]);

  useEffect(() => {
    const defaultTool = tools.find(tool => tool.DefaultTool);
    if (defaultTool) {
      setSelectedToolCID(defaultTool.CID);
      onChange(defaultTool.CID);
    }
  }, [tools, onChange]);

  useEffect(() => {
    console.log("Selected Tool CID:", selectedToolCID);
  }, [selectedToolCID]);  

  const handleSelectionChange = (value: string) => {
    setSelectedToolCID(value);
    onChange(value);
  };

  return (
    <Select onValueChange={handleSelectionChange} value={selectedToolCID}>
      <SelectTrigger>
        <SelectValue placeholder="Select a model" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {tools.map((tool) => (
            <SelectItem key={tool.CID} value={tool.CID}>
              {tool.ToolJson?.author || "unknown"}/{tool.Name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
