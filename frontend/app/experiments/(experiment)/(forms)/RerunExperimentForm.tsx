"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { usePrivy } from "@privy-io/react-auth";
import { ChevronsUpDownIcon, RefreshCcwIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import * as z from "zod";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form } from "@/components/ui/form";
import { AppDispatch, selectFlowDetail, selectToolDetail } from "@/lib/redux";

import ContinuousSwitch from "./ContinuousSwitch";
import { DynamicArrayField } from "./DynamicArrayField";
import { generateDefaultValues, generateRerunSchema, generateValues } from "./formGenerator";
import { groupInputs, transformJson } from "./formUtils";
import { toast } from "sonner";

export default function RerunExperimentForm() {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const { user } = usePrivy();

  const tool = useSelector(selectToolDetail);
  const flow = useSelector(selectFlowDetail);

  const lastJob = flow?.Jobs?.[flow?.Jobs?.length - 1];

  const walletAddress = user?.wallet?.address;

  const groupedInputs = groupInputs(tool.ToolJson?.inputs);
  const formSchema = generateRerunSchema(tool.ToolJson?.inputs);
  const defaultValues = generateValues(lastJob?.Inputs);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: defaultValues,
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("===== Form Submitted =====", values);

    if (!walletAddress) {
      console.error("Wallet address missing");
      return;
    }
    const transformedPayload = transformJson(tool, values, walletAddress);

    console.log("Placeholder for adding run to experiment:", transformedPayload);
    toast.warning("Re-running experiments is coming soon!", { position: "top-center" });
  }

  return flow && lastJob ? (
    <>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          {!!groupedInputs?.standard && (
            <>
              <Card>
                {Object.keys(groupedInputs?.standard || {}).map((groupKey) => {
                  return (
                    <CardContent key={groupKey} className="pt-0 first:pt-2">
                      <div className="space-y-4">
                        {Object.keys(groupedInputs?.standard[groupKey] || {}).map((key) => {
                          const input = groupedInputs?.standard?.[groupKey]?.[key];
                          return <DynamicArrayField key={key} inputKey={key} form={form} input={input} />;
                        })}
                      </div>
                    </CardContent>
                  );
                })}

                {Object.keys(groupedInputs?.collapsible || {}).map((groupKey) => {
                  return (
                    <CardContent key={groupKey} className="border-t">
                      <Collapsible>
                        <CollapsibleTrigger className="flex items-center justify-between w-full gap-2 text-left uppercase font-heading">
                          {groupKey.replace("_", "")}
                          <ChevronsUpDownIcon />
                        </CollapsibleTrigger>
                        <CollapsibleContent>
                          <div className="pt-0 space-y-4">
                            {Object.keys(groupedInputs?.collapsible[groupKey] || {}).map((key) => {
                              const input = groupedInputs?.collapsible?.[groupKey]?.[key];
                              return <DynamicArrayField key={key} inputKey={key} form={form} input={input} />;
                            })}
                          </div>
                        </CollapsibleContent>
                      </Collapsible>
                    </CardContent>
                  );
                })}
                <CardContent>
                  <Button type="submit" className="flex-wrap w-full h-auto">
                    <RefreshCcwIcon /> Re-run Experiment
                  </Button>
                </CardContent>
              </Card>
              <Card className="mt-3">
                <CardContent>
                  <ContinuousSwitch />
                </CardContent>
              </Card>
            </>
          )}
        </form>
      </Form>
    </>
  ) : null;
}
