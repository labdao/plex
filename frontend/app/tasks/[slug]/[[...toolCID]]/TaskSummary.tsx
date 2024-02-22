"use client";

import React from "react";
import { UseFormReturn } from "react-hook-form";

import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface TaskSummaryProps {
  form: UseFormReturn<any>;
  sortedInputs: any;
  outputs: {
    [key: string]: {
      glob: string[];
      item: string;
      type: string;
    };
  } | null;
}

type VariantSummaryItem = {
  name: string;
  variantCount: number;
};

type OutputSummaryItem = {
  name: string;
  fileExtensions: string;
  fileNames: string;
  multiple: boolean;
};

export function TaskSummary({ sortedInputs, form, outputs }: TaskSummaryProps) {
  const watchAllFields = form.watch();

  let variantSummaryInfo = { items: [] as VariantSummaryItem[], total: 1 };
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

  let outputSummaryInfo = { items: [] as OutputSummaryItem[] };
  for (const key in outputs) {
    outputSummaryInfo.items.push({
      name: key.replaceAll("_", " "),
      fileExtensions: outputs?.[key]?.glob?.map((glob: string) => glob.split(".").pop())?.join(", "),
      fileNames: outputs?.[key]?.glob?.join(", "),
      multiple: outputs?.[key]?.type === "Array",
    });
  }

  return (
    <div className="sticky max-h-screen pb-6 mt-6 overflow-y-auto lg:mt-0 lg:pl-6 top-4">
      <Card>
        <TransactionSummaryInfo className="px-6 rounded-b-none" />
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
          <Button type="submit" form="task-form" className="flex-wrap w-full h-auto">
            Submit <Badge className="mx-1 bg-black/10">{variantSummaryInfo?.total || 1}</Badge> Experimental Run
            {variantSummaryInfo?.total > 1 && "s"}
          </Button>
        </CardContent>

        {outputs && (
          <>
            <Separator className="my-2" />
            <CardContent>
              <div className="mb-4 font-mono text-sm font-bold uppercase">Expected Output</div>
              <div className="space-y-2 lowercase">
                {(outputSummaryInfo?.items || []).map((item, index) => (
                  <div key={index}>
                    {item.multiple ? <div>{item.fileExtensions} files</div> : <div>{item.fileNames} file</div>}
                    <div className="mr-3 text-xs text-muted-foreground">{item.name}</div>
                  </div>
                ))}
              </div>
            </CardContent>
          </>
        )}
      </Card>
    </div>
  );
}
