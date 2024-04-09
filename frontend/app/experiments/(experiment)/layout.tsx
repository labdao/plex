"use client";

import React, { ReactNode } from "react";

import { ExperimentUIContextProvider } from "@/app/experiments/(experiment)/ExperimentUIContext";

type LayoutProps = {
  children: ReactNode;
};

export default function Layout({ children }: LayoutProps) {
  return (
    <ExperimentUIContextProvider>
      <div className="relative flex flex-col max-w-full min-h-screen grow">{children}</div>
    </ExperimentUIContextProvider>
  );
}
