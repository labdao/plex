"use client";

//Workaround for intercepted route not being dismissed
// https://github.com/vercel/next.js/issues/49662

import { usePathname } from "next/navigation";

export default function Layout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  return pathname.startsWith("/experiments/") ? children : null;
}
