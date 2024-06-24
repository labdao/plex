"use client";

import React from "react";
import { UseFormReturn } from "react-hook-form";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { ModelDetail } from "@/lib/redux";

interface ExperimentSummaryProps {
  form: UseFormReturn<any>;
  sortedInputs: any;
  model: ModelDetail;
  showVariants?: boolean;
}

type VariantSummaryItem = {
  name: string;
  variantCount: number;
};

export function ExperimentSummary({ sortedInputs, form, showVariants }: ExperimentSummaryProps) {
  const watchAllFields = form.watch();

  let variantSummaryInfo = { items: [] as VariantSummaryItem[], total: 1 };
  for (const [key, input] of sortedInputs) {
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
    <div className="p-4 border-t">
      {showVariants && (
        <>
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
        </>
      )}
    </div>
  );
}
