import React, { ReactNode } from "react";

type LayoutProps = {
  children: ReactNode;
};

export default function Layout({ children }: LayoutProps) {
  return <div className="relative flex flex-col max-w-full min-h-screen grow">{children}</div>;
}
