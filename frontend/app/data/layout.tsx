import React, { ReactNode } from "react";

import { Breadcrumbs } from "@/components/global/Breadcrumbs";

type LayoutProps = {
  children: ReactNode;
  list: any;
  add: any;
};

export default async function Layout({ children, list, add }: LayoutProps) {
  return (
    <>
      <Breadcrumbs items={[{ name: "My Files", href: "/data" }]} actions={<div className="flex justify-end my-8"> {add}</div>} />
      {list}
      {children}
    </>
  );
}
