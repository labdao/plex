"use client";

import { usePrivy } from "@privy-io/react-auth";
import { useRouter } from "next/navigation";
import React, { useState } from "react";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  endLoading,
  saveDataFileAsync,
  selectCID,
  selectDataFileError,
  selectDataFileIsLoading,
  setError,
  startLoading,
  useDispatch,
  useSelector,
} from "@/lib/redux";

export default function DataFileForm() {
  const [open, setOpen] = React.useState(false);
  const { user } = usePrivy();
  const dispatch = useDispatch();

  const router = useRouter();
  const errorMessage = useSelector(selectDataFileError);
  const isLoading = useSelector(selectDataFileIsLoading);
  const walletAddress = user?.wallet?.address;

  const [file, setFile] = useState<File | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const uploadedFile = e.target.files && e.target.files[0];
    if (uploadedFile) {
      setFile(uploadedFile);
    }
  };

  const handleSuccess = () => {
    setOpen(false);
    //TODO: Update the list
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (file === null) {
      dispatch(setError("Please select a file"));
      return;
    }

    dispatch(startLoading());
    dispatch(setError(null));
    const metadata = { walletAddress };

    try {
      await dispatch(saveDataFileAsync({ file, metadata, handleSuccess }));
      dispatch(endLoading());
    } catch (error) {
      dispatch(setError("Error uploading file"));
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
              <Input type="file" onChange={handleFileChange} />
              <Button type="submit">{isLoading ? "Submitting..." : "Submit"}</Button>
              {errorMessage && <Alert variant="destructive">{errorMessage}</Alert>}
            </form>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
