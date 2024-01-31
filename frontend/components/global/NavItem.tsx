import { Slot } from "@radix-ui/react-slot";
import { ChevronRightIcon } from "lucide-react";
import dynamicIconImports from "lucide-react/dynamicIconImports";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ReactNode } from "react";
import * as React from "react";

import { cn } from "@/lib/utils";

interface NavLinkProps {
  href: string;
  exact?: boolean;
  children?: ReactNode;
  title?: string;
  className?: string;
  icon?: ReactNode;
  target?: "_blank";
}

const NavLink = ({ href, exact, children, title, className = "", icon, ...props }: NavLinkProps) => {
  const pathname = usePathname();
  const isActive = exact ? pathname === href : pathname.startsWith(href);

  return (
    <Link
      href={href}
      {...props}
      className={cn(
        "w-full flex items-center text-sm px-3 py-2 hover:bg-muted rounded-full text-muted-foreground hover:text-foreground",
        isActive && "text-foreground",
        className
      )}
    >
      {icon && <span className="mr-2">{icon}</span>}
      {title || children}
    </Link>
  );
};

interface NavButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  asChild?: boolean;
  hasDropdown?: boolean;
  title?: string;
  className?: string;
  icon?: ReactNode;
}

const NavButton = React.forwardRef<HTMLButtonElement, NavButtonProps>(
  ({ className, icon, children, title, asChild = false, hasDropdown = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        className={cn(
          "w-full flex items-center px-3 py-2 text-sm rounded-full text-muted-foreground hover:text-foreground hover:bg-muted",
          "aria-expanded:text-foreground",
          className
        )}
        ref={ref}
        {...props}
      >
        <>
          {icon && <span className="mr-2">{icon}</span>}
          {<span className="truncate">{title}</span> || children}
          {hasDropdown && <ChevronRightIcon size={12} className="ml-auto opacity-50 shrink-0" />}
        </>
      </Comp>
    );
  }
);
NavButton.displayName = "Button";

export { NavButton, NavLink };
