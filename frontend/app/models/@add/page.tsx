"use client";

import { JsonInput } from "@mantine/core";
import { MantineProvider } from "@mantine/core";
import { usePrivy } from "@privy-io/react-auth";
import { PlusIcon } from "lucide-react";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import {
  AppDispatch,
  createToolThunk,
  dataFileListThunk,
  selectAddToolError,
  selectAddToolJson,
  selectAddToolLoading,
  selectAddToolSuccess,
  setAddToolError,
  setAddToolJson,
  setAddToolLoading,
  setAddToolSuccess,
  toolListThunk,
} from "@/lib/redux";

export default function AddTool() {
  const [open, setOpen] = React.useState(false);
  const { user } = usePrivy();
  const dispatch = useDispatch<AppDispatch>();

  const walletAddress = user?.wallet?.address;
  const loading = useSelector(selectAddToolLoading);
  const error = useSelector(selectAddToolError);
  const toolJson = useSelector(selectAddToolJson);
  const toolSuccess = useSelector(selectAddToolSuccess);

  useEffect(() => {
    if (toolSuccess) {
      dispatch(setAddToolSuccess(false));
      dispatch(setAddToolJson(""));
      setOpen(false);
      //TODO: Update the list
      return;
    }
    dispatch(toolListThunk());
    dispatch(dataFileListThunk({}));
  }, [toolSuccess, dispatch]);

  const handleToolJsonChange = (toolJsonInput: string) => {
    dispatch(setAddToolJson(toolJsonInput));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Submitting tool.json: ", toolJson);
    dispatch(setAddToolLoading(true));
    dispatch(setAddToolError(""));
    try {
      const toolJsonParsed = JSON.parse(toolJson);
      if (!walletAddress) {
        dispatch(setAddToolError("Wallet address missing"));
        return;
      }
      await dispatch(createToolThunk({ toolJson: toolJsonParsed }));
    } catch (error) {
      console.error("Error creating tool", error);
      dispatch(setAddToolError("Error creating tool"));
    }
    dispatch(setAddToolLoading(false));
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <PlusIcon />
          Add Model
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Model</DialogTitle>
          <DialogDescription>
            <form onSubmit={handleSubmit}>
              <MantineProvider>
                <JsonInput
                  label="Tool Definition"
                  placeholder="Paste your tool's JSON definition here."
                  validationError="Invalid JSON"
                  autosize
                  minRows={10}
                  value={toolJson}
                  onChange={handleToolJsonChange}
                  styles={{
                    input: { width: "100%" },
                  }}
                />
              </MantineProvider>
              <Button type="submit">{loading ? "Submitting..." : "Submit"}</Button>
              {error && (
                <Alert variant="destructive" className="my-4">
                  {error}
                </Alert>
              )}
            </form>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
