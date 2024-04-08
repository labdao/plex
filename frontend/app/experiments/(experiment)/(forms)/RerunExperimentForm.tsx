"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { usePrivy } from "@privy-io/react-auth";
import { ChevronsUpDownIcon, RefreshCcwIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import React from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import { toast } from "sonner";
import * as z from "zod";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form } from "@/components/ui/form";
import { AppDispatch, flowDetailThunk, flowListThunk, selectFlowDetail, selectToolDetail } from "@/lib/redux";
import { addJobToFlow } from "@/lib/redux/slices/flowAddSlice/asyncActions";


import ContinuousSwitch from "./ContinuousSwitch";
import { DynamicArrayField } from "./DynamicArrayField";
import { generateRerunSchema, generateValues } from "./formGenerator";
import { groupInputs, transformJson } from "./formUtils";

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
  const flowID = flow?.ID || 0;

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

    console.log("Submitting Payload:", transformedPayload);
    try {
      const response = await addJobToFlow(flowID, transformedPayload);
      console.log("Response from addJobToFlow", response);
      if (response && response.ID) {
        console.log("Flow created", response);
        console.log(response.ID);
        dispatch(flowDetailThunk(response.ID));
        dispatch(flowListThunk(walletAddress));
        toast.success("Experiment started successfully");
      } else {
        console.log("Something went wrong", response);
      }
    } catch (error) {
      console.error("Failed to create flow", error);
      // Handle error, maybe show message to user
    }
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
                        <CollapsibleTrigger className="flex items-center w-full gap-2 text-sm text-left lowercase text-muted-foreground font-heading">
                          <ChevronsUpDownIcon />
                          {groupKey.replace("_", "")}
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
