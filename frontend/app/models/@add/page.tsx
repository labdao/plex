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
  createModelThunk,
  dataFileListThunk,
  selectAddModelError,
  selectAddModelJson,
  selectAddModelLoading,
  selectAddModelSuccess,
  setAddModelError,
  setAddModelJson,
  setAddModelLoading,
  setAddModelSuccess,
  modelListThunk,
} from "@/lib/redux";

export default function AddModel() {
  const [open, setOpen] = React.useState(false);
  const { user } = usePrivy();
  const dispatch = useDispatch<AppDispatch>();

  const walletAddress = user?.wallet?.address;
  const loading = useSelector(selectAddModelLoading);
  const error = useSelector(selectAddModelError);
  const modelJson = useSelector(selectAddModelJson);
  const modelSuccess = useSelector(selectAddModelSuccess);

  useEffect(() => {
    if (modelSuccess) {
      dispatch(setAddModelSuccess(false));
      dispatch(setAddModelJson(""));
      setOpen(false);
      //TODO: Update the list
      return;
    }
    dispatch(modelListThunk());
    dispatch(dataFileListThunk({}));
  }, [modelSuccess, dispatch]);

  const handleModelJsonChange = (modelJsonInput: string) => {
    dispatch(setAddModelJson(modelJsonInput));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log("Submitting model.json: ", modelJson);
    dispatch(setAddModelLoading(true));
    dispatch(setAddModelError(""));
    try {
      const modelJsonParsed = JSON.parse(modelJson);
      if (!walletAddress) {
        dispatch(setAddModelError("Wallet address missing"));
        return;
      }
      await dispatch(createModelThunk({ modelJson: modelJsonParsed }));
    } catch (error) {
      console.error("Error creating model", error);
      dispatch(setAddModelError("Error creating model"));
    }
    dispatch(setAddModelLoading(false));
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
                  label="Model Definition"
                  placeholder="Paste your model's JSON definition here."
                  validationError="Invalid JSON"
                  autosize
                  minRows={10}
                  value={modelJson}
                  onChange={handleModelJsonChange}
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
