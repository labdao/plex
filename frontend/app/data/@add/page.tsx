"use client";

import { useRouter } from "next/navigation";
import React, { useState } from "react";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  endLoading,
  saveDataFilesAsync,
  selectDataFileError,
  selectDataFileIsLoading,
  selectWalletAddress,
  setError,
  startLoading,
  useDispatch,
  useSelector,
} from "@/lib/redux";

export default function DataFileForm() {
  const [open, setOpen] = React.useState(false);

  const dispatch = useDispatch();

  const router = useRouter();
  const errorMessage = useSelector(selectDataFileError);
  const isLoading = useSelector(selectDataFileIsLoading);
  const walletAddress = useSelector(selectWalletAddress);

  const [files, setFiles] = useState<File[]>([]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      setFiles(Array.from(e.target.files));
    }
  };

  const handleSuccess = () => {
    setOpen(false);
    //TODO: Update the list
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (files.length === 0) {
      dispatch(setError("Please select at least one file"));
      return;
    }

    dispatch(startLoading());
    dispatch(setError(null));
    const metadata = { walletAddress };

    try {
      await dispatch(saveDataFilesAsync({ files, metadata, handleSuccess }));
    } catch (error) {
      dispatch(setError("Error uploading files"));
    } finally {
      dispatch(endLoading());
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger>
        <Button size="lg">Add Data</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Data</DialogTitle>
          <DialogDescription>
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <Input type="file" onChange={handleFileChange} multiple/>
              <Button type="submit">{isLoading ? "Submitting..." : "Submit"}</Button>
              {errorMessage && <Alert variant="destructive">{errorMessage}</Alert>}
            </form>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
