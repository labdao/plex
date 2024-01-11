import Image from "next/image";
import Link from "next/link";
import { ReactNode } from "react";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardTitle } from "@/components/ui/card";

const CardWithHref = ({ href, children }: { href: string | undefined; children: ReactNode }) =>
  href ? (
    <Link href={href}>
      <Card className="h-full border border-border hover:border-ring">{children}</Card>
    </Link>
  ) : (
    <Card className="relative h-full border cursor-not-allowed border-border">
      <div className="absolute inset-0 z-10 flex items-center justify-center bg-gray-300 bg-opacity-70">
        <Badge variant="secondary" size="lg">
          Coming soon
        </Badge>
      </div>
      {children}
    </Card>
  );

export default function TaskCard({ name, slug, available }: { name: string; slug: string; available: boolean }) {
  return (
    <CardWithHref key={slug} href={available ? `/tasks/${slug}` : undefined}>
      <div className="flex flex-col h-full">
        <CardContent className="py-4">
          <CardTitle className="flex justify-between text-base">{name}</CardTitle>
        </CardContent>
        <Image src={`/images/task-${slug}.png`} className="w-full h-auto mt-auto" sizes="100vw" width={220} height={140} alt={name} />
      </div>
    </CardWithHref>
  );
}
