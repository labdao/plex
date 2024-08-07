"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { usePrivy } from "@privy-io/react-auth";
import { ChevronsUpDownIcon, RefreshCcwIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import { toast } from "sonner";
import * as z from "zod";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form } from "@/components/ui/form";
import { AppDispatch, addExperimentWithCheckoutThunk, experimentDetailThunk, experimentListThunk, refreshUserDataThunk, selectExperimentDetail, selectModelDetail, selectIsUserSubscribed, selectUserTier } from "@/lib/redux";
import { addJobToExperiment } from "@/lib/redux/slices/experimentAddSlice/asyncActions";


import ContinuousSwitch from "./ContinuousSwitch";
import { DynamicArrayField } from "./DynamicArrayField";
import { generateRerunSchema, generateValues } from "./formGenerator";
import { groupInputs, transformJson } from "./formUtils";

export default function RerunExperimentForm() {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const { user } = usePrivy();

  const model = useSelector(selectModelDetail);
  const experiment = useSelector(selectExperimentDetail);
  const userTier = useSelector(selectUserTier);
  const isUserSubscribed = useSelector(selectIsUserSubscribed);

  const lastJob = experiment?.Jobs?.[experiment?.Jobs?.length - 1];

  const walletAddress = user?.wallet?.address;

  const groupedInputs = groupInputs(model.ModelJson?.inputs);
  const formSchema = generateRerunSchema(model.ModelJson?.inputs, lastJob?.Inputs);
  const defaultValues = generateValues(model.ModelJson?.inputs, lastJob?.Inputs);
  const experimentID = experiment?.ID || 0;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: defaultValues,
  });

  useEffect(() => {
    if (model && lastJob) {
      form.reset(generateValues(model.ModelJson?.inputs, lastJob?.Inputs));
    }
  }, [model, lastJob, form]);

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("===== Form Submitted =====", values);
  
    if (!walletAddress) {
      console.error("Wallet address missing");
      return;
    }
  
    const transformedPayload = transformJson(model, values, walletAddress);
    console.log("Submitting Payload:", transformedPayload);
  
    try {
      if (userTier === 'Free' || (userTier === 'Paid' && isUserSubscribed)) {
        const response = await addJobToExperiment(experimentID, transformedPayload);
        if (response.redirectUrl) {
          router.push(response.redirectUrl);
        } else if (response && response.ID) {
          console.log("Job added to experiment", response);
          dispatch(experimentDetailThunk(response.ID));
          dispatch(experimentListThunk(walletAddress));
          toast.success("Job added to experiment successfully");
        } else {
          console.log("Something went wrong", response);
          toast.error("Failed to add job to experiment");
        }
      } else if (userTier === 'Paid' && !isUserSubscribed) {
        console.log('Redirecting...');
        router.push('/subscribe');
      } else {
        console.error("Invalid user tier");
        toast.error("Unable to process request due to invalid user tier");
      }
  
      await dispatch(refreshUserDataThunk());
    } catch (error) {
      console.error("Failed to add job to experiment", error);
      toast.error("Failed to add job to experiment");
    }
  }

  return experiment && lastJob ? (
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
