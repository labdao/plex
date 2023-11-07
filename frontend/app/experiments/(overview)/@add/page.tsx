"use client";

import { PayloadAction } from "@reduxjs/toolkit";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";

import { DataFileSelect } from "@/components/shared/DataFileSelect";
import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
  addFlowThunk,
  AppDispatch,
  dataFileListThunk,
  selectDataFileList,
  selectDataFileListError,
  selectFlowAddCid,
  selectFlowAddError,
  selectFlowAddKwargs,
  selectFlowAddLoading,
  selectFlowAddName,
  selectFlowAddTool,
  selectToolList,
  selectToolListError,
  selectWalletAddress,
  setFlowAddCid,
  setFlowAddError,
  setFlowAddKwargs,
  setFlowAddLoading,
  setFlowAddName,
  setFlowAddSuccess,
  setFlowAddTool,
  toolListThunk,
} from "@/lib/redux";
import { DataFile } from '@/lib/redux/slices/dataFileListSlice/slice';

export default function AddGraph() {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();

  const walletAddress = useSelector(selectWalletAddress);
  const name = useSelector(selectFlowAddName);
  const loading = useSelector(selectFlowAddLoading);
  const error = useSelector(selectFlowAddError);
  const kwargs = useSelector(selectFlowAddKwargs);
  const cid = useSelector(selectFlowAddCid);
  const selectedTool = useSelector(selectFlowAddTool);
  const toolListError = useSelector(selectToolListError);
  const dataFileListError = useSelector(selectDataFileListError);
  const dataFiles = useSelector(selectDataFileList);
  const tools = useSelector(selectToolList);

  const [selectedToolIndex, setSelectedToolIndex] = useState("");
  const [inputDataFiles, setInputDataFiles] = useState<Record<string, DataFile[]>>({});

  useEffect(() => {
    if (cid !== "") {
      dispatch(setFlowAddSuccess(false));
      dispatch(setFlowAddKwargs({}));
      dispatch(setFlowAddTool({ CID: "", WalletAddress: "", Name: "", ToolJson: { inputs: {} } }));
      dispatch(setFlowAddError(null));
      dispatch(setFlowAddName(""));
      dispatch(setFlowAddCid(""));
      router.push(`/experiments/${cid}`);
      return;
    }
    dispatch(toolListThunk());
    dispatch(dataFileListThunk());
  }, [cid, dispatch, router]);

  const handleToolChange = async (value: string) => {
    const selectedTool = tools[parseInt(value)];
    dispatch(setFlowAddTool(selectedTool));
    setSelectedToolIndex(value);
  
    const newInputDataFiles: Record<string, DataFile[]> = {};
  
    for (const inputKey in selectedTool.ToolJson.inputs) {
      const input = (selectedTool.ToolJson.inputs as Record<string, { glob: string[] }>)[inputKey];
      if (typeof input === 'object' && input !== null && 'glob' in input) {
        const globPatterns = input.glob;
        const action = await dispatch(dataFileListThunk(globPatterns)) as PayloadAction<DataFile[]>;
        newInputDataFiles[inputKey] = action.payload;
      }
    }
  
    setInputDataFiles(newInputDataFiles);
  };

  const handleKwargsChange = (value: string, key: string) => {
    console.log(value);
    console.log(key);
    const updatedKwargs = { ...kwargs, [key]: [value] };
    dispatch(setFlowAddKwargs(updatedKwargs));
  };

  const isValidForm = (): boolean => {
    if (selectedTool.CID === "") return false;
    for (const key in kwargs) {
      if (kwargs[key].length === 0) return false;
    }
    return true;
  };

  const handleSubmit = async (event: any) => {
    event.preventDefault();
    console.log("Submitting flow");
    dispatch(setFlowAddLoading(true));
    dispatch(setFlowAddError(null));
    await dispatch(
      addFlowThunk({
        name,
        walletAddress,
        toolCid: selectedTool.CID,
        scatteringMethod: "dotProduct",
        kwargs,
      })
    );
    dispatch(setFlowAddLoading(false));
  };

  return (
    <Dialog>
      <DialogTrigger>
        <Button size="lg">Add Experiment</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Experiment</DialogTitle>
          <DialogDescription>
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <div>
                <Label>Name</Label>
                <Input value={name} onChange={(e) => dispatch(setFlowAddName(e.target.value))} />
              </div>
              <div>
                <Label>Tool</Label>
                <Select onValueChange={handleToolChange} defaultValue={selectedToolIndex}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a tool" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {tools.map((tool, index) => {
                        return (
                          <SelectItem key={index} value={index.toString()}>
                            {tool.Name}
                          </SelectItem>
                        );
                      })}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>

              {Object.keys(selectedTool.ToolJson.inputs).map((key) => {
                // @ts-ignore
                const inputDetail = selectedTool.ToolJson.inputs[key];
                return (
                  <div key={key}>
                    <Label>
                      {key}
                      {inputDetail.glob && ` (Glob: ${inputDetail.glob.join(', ')})`}
                    </Label>
                    {/* <DataFileSelect onSelect={(value) => handleKwargsChange(value, key)} value={kwargs[key]?.[0]} dataFiles={inputDataFiles[key]} label={key} /> */}
                    <DataFileSelect 
                      onSelect={(value) => handleKwargsChange(value, key)} 
                      value={kwargs[key]?.[0]} 
                      dataFiles={inputDataFiles[key] ? inputDataFiles[key] : []} 
                      label={key} 
                    />
                  </div>
                );
              })}
              <Button type="submit" disabled={loading || !isValidForm()}>
                {loading ? "Submitting..." : "Submit"}
              </Button>
              {error && <Alert variant="destructive">{error}</Alert>}
              {toolListError && <Alert variant="destructive">{toolListError}</Alert>}
              {dataFileListError && <Alert variant="destructive">{dataFileListError}</Alert>}
            </form>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
