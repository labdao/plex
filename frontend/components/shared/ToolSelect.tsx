"use client";

import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AppDispatch, selectToolList, selectToolListError, toolListThunk } from "@/lib/redux";

interface ToolSelectProps {
  onChange: (value: string) => void;
  taskSlug?: string;
  defaultValue?: string; // New prop for the default value
}

export function ToolSelect({ onChange, taskSlug, defaultValue }: ToolSelectProps) {
  const dispatch = useDispatch<AppDispatch>();
  const tools = useSelector(selectToolList);
  const [selectedToolCID, setSelectedToolCID] = useState("");
  const [userHasSelected, setUserHasSelected] = useState(false);
  const toolListError = useSelector(selectToolListError);

  useEffect(() => {
    if (taskSlug) {
      dispatch(toolListThunk(taskSlug));
    } else {
      dispatch(toolListThunk());
    }
  }, [dispatch, taskSlug]);

  // useEffect(() => {
  //   // Only set default tool if user hasn't made a selection
  //   console.log("problematic useeffect is called. User has selected:", userHasSelected);
  //   if (!userHasSelected) {
  //     const defaultTool = tools.find(tool => tool.DefaultTool);
  //     if (defaultTool) {
  //       setSelectedToolCID(defaultTool.CID);
  //       onChange(defaultTool.CID);
  //     }
  //   }
  // }, [tools, onChange, userHasSelected]); // Add userHasSelected to dependency array

  useEffect(() => {
    console.log("Selected Tool CID:", defaultValue);
  }, [defaultValue]);  

  const handleSelectionChange = (value: string) => {
    console.log("handling selection change");
    setSelectedToolCID(value);
    onChange(value);
    setUserHasSelected(true); // Set to true when user makes a selection
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
