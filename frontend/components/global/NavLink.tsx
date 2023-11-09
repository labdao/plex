import Link from "next/link";
import { usePathname } from "next/navigation";
import { ReactNode } from "react";

import { cn } from "@/lib/utils";

interface NavLinkProps {
  href: string;
  exact?: boolean;
  children: ReactNode;
  className?: string;
}

export function NavLink({ href, exact, children, className = "", ...props }: NavLinkProps): JSX.Element {
  const pathname = usePathname();
  const isActive = exact ? pathname === href : pathname.startsWith(href);

  return (
    <Link href={href} {...props} className={cn(!isActive && "text-muted", className)}>
      {children}
    </Link>
  );
}
