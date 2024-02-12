"use client";

import { useEffect, useState } from "react";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export function ViewportWarning() {
  const [open, setOpen] = useState(undefined as boolean | undefined);

  useEffect(() => {
    const handleResize = (width: number) => {
      width < 900 ? setOpen(true) : setOpen(false);
    };

    handleResize(window.innerWidth);
    window.addEventListener("resize", () => {
      handleResize(window.innerWidth);
    });
  }, []);

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Optimized for larger screens</AlertDialogTitle>
          <AlertDialogDescription>
            We&lsquo;re still working on optimizing the UI for smaller screens. For the best experience, please use a device with a larger screen.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogAction>Continue Anyway</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
