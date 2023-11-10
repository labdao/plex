import { ArrowRightIcon } from "lucide-react";
import Link from "next/link";

import { Card, CardContent, CardTitle } from "@/components/ui/card";

export default function TaskList() {
  return (
    <>
      <div className="container mt-8">
        <h1 className="mb-4 text-3xl font-bold font-heading">Tasks</h1>
        <div className="grid grid-cols-3">
          <Link href="/tasks/protein-design">
            <Card className="hover:border hover:border-ring">
              <CardContent>
                <CardTitle className="flex justify-between">
                  Protein Design <ArrowRightIcon />
                </CardTitle>
              </CardContent>
            </Card>
          </Link>
        </div>
      </div>
    </>
  );
}
