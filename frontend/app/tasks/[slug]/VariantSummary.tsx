"use client";

import React from "react";
import { UseFormReturn } from "react-hook-form";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

interface VariantSummaryProps {
  form: UseFormReturn<any>;
  sortedInputs: any;
}

type VariantSummaryItem = {
  name: string;
  variantCount: number;
};

export function VariantSummary({ sortedInputs, form }: VariantSummaryProps) {
  const watchAllFields = form.watch();
  const variantSummaryInfo = { items: [] as VariantSummaryItem[], total: 1 };

  for (const [key, input] of sortedInputs) {
    //If the field is required or has a value, we show it in the summary
    //We don't show optional fields with no value
    if (watchAllFields?.[key]?.[0]?.value || input?.required) {
      const count = watchAllFields[key]?.length;
      variantSummaryInfo.items.push({
        name: key.replaceAll("_", " "),
        variantCount: count,
      });
      if (!input?.array) variantSummaryInfo.total *= count;
    }
  }

  return (
    <>
      <Card className="sticky top-4">
        <CardContent>
          <div className="mb-4 font-mono text-sm font-bold uppercase">Variant Summary</div>
          <div className="mb-4 space-y-2 lowercase">
            {(variantSummaryInfo?.items || []).map((item: { name: string; variantCount: number }, index: number) => (
              <div className="flex justify-between text-muted-foreground" key={index}>
                <span>{item.name}</span>{" "}
                <span className="mr-3">
                  x <span className="text-foreground">{item.variantCount}</span>
                </span>
              </div>
            ))}
            <div className="flex justify-between pt-2 border-t text-muted-foreground">
              <span>Total Runs</span>{" "}
              <Badge size="lg" className="text-base" variant="secondary">
                {variantSummaryInfo?.total || 1}
              </Badge>
            </div>
          </div>
          <Button type="submit" form="task-form" className="w-full">
            Submit <Badge className="mx-1">{variantSummaryInfo?.total || 1}</Badge> Experimental Run{variantSummaryInfo?.total > 1 && "s"}
          </Button>
        </CardContent>
      </Card>
    </>
  );
}
