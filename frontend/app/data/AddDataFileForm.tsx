"use client";

import { usePrivy } from "@privy-io/react-auth";
import { useRouter } from "next/navigation";
import React, { useState } from "react";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  endFileUploadDataSlice,
  saveDataFileAsync,
  selectDataFileError,
  selectDataFileIsLoading,
  setError,
  startFileUploadDataSlice,
  useDispatch,
  useSelector,
} from "@/lib/redux";

interface AddDataFileFormProps {
  trigger: React.ReactNode;
}

export default function AddDataFileForm({ trigger }: AddDataFileFormProps) {
  const [open, setOpen] = React.useState(false);
  const { user } = usePrivy();
  const dispatch = useDispatch();

  const router = useRouter();
  const errorMessage = useSelector(selectDataFileError);
  const isLoading = useSelector(selectDataFileIsLoading);
  const walletAddress = user?.wallet?.address;

  const [file, setFile] = useState<File | null>(null);
  const [isPublic, setIsPublic] = useState(false);

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

  const handlePublicChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setIsPublic(e.target.checked);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (file === null) {
      dispatch(setError("Please select a file"));
      return;
    }

    dispatch(startFileUploadDataSlice());
    dispatch(setError(null));
    const metadata = { walletAddress };

    try {
      await dispatch(saveDataFileAsync({ file, metadata, isPublic, handleSuccess }));
      dispatch(endFileUploadDataSlice());
    } catch (error) {
      dispatch(setError("Error uploading file"));
      dispatch(endFileUploadDataSlice());
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upload Files</DialogTitle>
          <DialogDescription>
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <Input type="file" onChange={handleFileChange} />
              {file && (
                <label>
                  <input 
                    type="checkbox" 
                    checked={isPublic}
                    onChange={handlePublicChange} 
                  />
                  <span className="ml-2 font-bold">
                    Mark Public.
                  </span>
                  <span className="ml-1 text-sm text-gray-600">
                    Once a datafile is public, this cannot be undone.
                  </span>
                </label>
              )}
              <Button type="submit">{isLoading ? "Submitting..." : "Submit"}</Button>
              {errorMessage && <Alert variant="destructive">{errorMessage}</Alert>}
            </form>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}
