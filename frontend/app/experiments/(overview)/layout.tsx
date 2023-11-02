import React, { ReactNode } from "react";

type LayoutProps = {
  children: ReactNode;
  list: any;
  add: any;
};

export default async function Layout({ children, list, add }: LayoutProps) {
  return (
    <>
      <div className="container mt-8">
        <div className="flex justify-end my-8"> {add}</div>
        {list}
        {children}
      </div>
    </>
  );
}
