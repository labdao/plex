import { cn } from "@lib/utils";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ReactNode } from "react";

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
    <Link href={href} {...props} className={cn(!isActive && "text-muted-foreground", className)}>
      {children}
    </Link>
  );
}
