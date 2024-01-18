"use client";

import { usePrivy } from "@privy-io/react-auth";
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
  selectFlowAddID,
  selectFlowAddName,
  selectFlowAddTool,
  selectToolList,
  selectToolListError,
  setFlowAddCid,
  setFlowAddError,
  setFlowAddKwargs,
  setFlowAddLoading,
  setFlowAddName,
  setFlowAddSuccess,
  setFlowAddTool,
  setFlowAddID,
  toolListThunk,
} from "@/lib/redux";
import { DataFile } from "@/lib/redux/slices/dataFileListSlice/slice";

export default function AddGraph() {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const { user } = usePrivy();

  const walletAddress = user?.wallet?.address;
  const name = useSelector(selectFlowAddName);
  const loading = useSelector(selectFlowAddLoading);
  const error = useSelector(selectFlowAddError);
  const kwargs = useSelector(selectFlowAddKwargs);
  const cid = useSelector(selectFlowAddCid);
  const flowID = useSelector(selectFlowAddID);

  interface ToolInput {
    position?: string;
    glob?: string[];
    type: string;
    default?: string;
  }

  interface ToolJson {
    inputs: Record<string, ToolInput>;
  }

  const selectedTool = useSelector(selectFlowAddTool);
  const toolListError = useSelector(selectToolListError);
  const dataFileListError = useSelector(selectDataFileListError);
  const dataFiles = useSelector(selectDataFileList);
  const tools = useSelector(selectToolList);

  const [selectedToolIndex, setSelectedToolIndex] = useState("");
  const [inputDataFiles, setInputDataFiles] = useState<Record<string, DataFile[]>>({});

  useEffect(() => {
    if (flowID != null) {
      dispatch(setFlowAddSuccess(false))
      dispatch(setFlowAddKwargs({}))
      dispatch(
        setFlowAddTool({
          CID: "",
          WalletAddress: "",
          Name: "",
          ToolJson: { inputs: {}, outputs: {}, name: "", author: "", description: "", github: "", paper: "" },
        })
      );      dispatch(setFlowAddError(null))
      dispatch(setFlowAddName(""))
      dispatch(setFlowAddCid(""))
      dispatch(setFlowAddID(null))
      router.push(`/experiments/${flowID}`)
      return
    }
    dispatch(toolListThunk())
    dispatch(dataFileListThunk({}))
  }, [flowID, dispatch, router])

  const handleToolChange = async (value: string) => {
    const selectedTool = tools[parseInt(value)];
    dispatch(setFlowAddTool(selectedTool));
    setSelectedToolIndex(value);

    const newInputDataFiles: Record<string, DataFile[]> = {};

    for (const inputKey in selectedTool.ToolJson.inputs) {
      const input = (selectedTool.ToolJson.inputs as Record<string, { glob: string[] }>)[inputKey];
      if (typeof input === "object" && input !== null && "glob" in input) {
        const globPatterns = input.glob;
        // @ts-ignore
        const action = (await dispatch(dataFileListThunk(globPatterns))) as PayloadAction<DataFile[]>;
        newInputDataFiles[inputKey] = action.payload;
      }
    }

    setInputDataFiles(newInputDataFiles);
  };

  const handleKwargsChange = (value: string, key: string) => {
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
    if (!walletAddress) {
      dispatch(setFlowAddError("Wallet address missing"));
      return;
    }

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

              {Object.keys(selectedTool.ToolJson.inputs)
                .sort((a, b) => {
                  // @ts-ignore
                  const positionA = parseInt(a[1].position, 10);
                  // @ts-ignore
                  const positionB = parseInt(b[1].position, 10);
                  return positionA - positionB;
                })
                .map((key) => {
                  // @ts-ignore
                  const inputDetail = selectedTool.ToolJson.inputs[key];
                  return (
                    <div key={key}>
                      <Label>
                        {key}
                        {inputDetail.glob && ` (Glob: ${inputDetail.glob.join(", ")})`}
                      </Label>
                      <DataFileSelect onValueChange={(value) => handleKwargsChange(value, key)} value={kwargs[key]?.[0]} label={key} />
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
