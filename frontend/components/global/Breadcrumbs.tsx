import Link from "next/link";

interface BreadcrumbsProps {
  items: {
    name: string;
    href?: string;
  }[];
  actions?: React.ReactNode;
}

export function Breadcrumbs({ items, actions }: BreadcrumbsProps) {
  return (
    <div className="z-10 flex items-center justify-between w-full px-4 py-2 font-mono font-bold uppercase border-b shadow-sm whitespace-nowrap h-14 bg-background">
      <div className="flex overflow-hidden">
        {items.map((item, idx) => {
          if (idx === items.length - 1 || !item.href)
            return (
              <div className="truncate " key={idx}>
                {item.name}/
              </div>
            );
          return (
            <Link className="text-muted-foreground hover:text-foreground" href={item.href} key={idx}>
              {item.name}/
            </Link>
          );
        })}
      </div>
      <div>{actions}</div>
    </div>
  );
}
