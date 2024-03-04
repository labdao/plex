import React, { ReactNode } from "react";

import { Breadcrumbs } from "@/components/global/Breadcrumbs";

type LayoutProps = {
  children: ReactNode;
  list: any;
  add: any;
};

export default async function Layout({ children, list, add }: LayoutProps) {
  return (
    <div className="relative flex flex-col h-screen max-w-full grow">
      <Breadcrumbs items={[{ name: "Models", href: "/models" }]} actions={<div>{add}</div>} />
      {list}
      {children}
    </div>
  );
}
