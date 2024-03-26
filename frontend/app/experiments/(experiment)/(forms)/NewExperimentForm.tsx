"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { usePrivy } from "@privy-io/react-auth";
import { ChevronsUpDownIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import * as z from "zod";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { AppDispatch, resetToolDetail, resetToolList, selectToolDetail } from "@/lib/redux";
import { createFlow } from "@/lib/redux/slices/flowAddSlice/asyncActions";

import { DynamicArrayField } from "./DynamicArrayField";
import { generateDefaultValues, generateSchema } from "./formGenerator";
import { groupInputs, transformJson } from "./formUtils";

export default function NewExperimentForm({ task }: { task: any }) {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const { user } = usePrivy();

  const tool = useSelector(selectToolDetail);
  const walletAddress = user?.wallet?.address;

  useEffect(() => {
    return () => {
      dispatch(resetToolDetail());
      dispatch(resetToolList());
    };
  }, [dispatch]);

  const groupedInputs = groupInputs(tool.ToolJson?.inputs);
  const formSchema = generateSchema(tool.ToolJson?.inputs);
  const defaultValues = generateDefaultValues(tool.ToolJson?.inputs, task, tool);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: defaultValues,
  });

  // If the tool changes, fetch new tool details
  useEffect(() => {
    form.reset(generateDefaultValues(tool.ToolJson?.inputs, task, tool));
  }, [tool, form, task]);

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("===== Form Submitted =====", values);

    if (!walletAddress) {
      console.error("Wallet address missing");
      return;
    }
    const transformedPayload = transformJson(tool, values, walletAddress);
    console.log("Submitting Payload:", transformedPayload);

    try {
      const response = await createFlow(transformedPayload);
      if (response && response.ID) {
        console.log("Flow created", response);
        console.log(response.ID);
        // Redirecting to another page, for example, a success page or dashboard
        router.push(`/experiments/${response.ID}`);
      } else {
        console.log("Something went wrong", response);
      }
    } catch (error) {
      console.error("Failed to create flow", error);
      // Handle error, maybe show message to user
    }
  }

  return (
    <>
      <Form {...form}>
        <form onSubmit={form.handleSubmit((values) => onSubmit(values))}>
          <Card className="mb-2">
            <CardContent>
              <div className="flex items-center gap-1">
                <div className="block w-3 h-3 border border-gray-300 rounded-full" />
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem className="grow">
                      <FormControl>
                        <Input variant="subtle" className="text-xl shrink-0 font-heading" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </CardContent>
          </Card>
          {!!groupedInputs?.standard && (
            <>
              <Card>
                {Object.keys(groupedInputs?.standard || {}).map((groupKey) => {
                  return (
                    <CardContent key={groupKey} className="border-t first:border-0">
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
                    <CardContent key={groupKey} className="pt-0 first:pt-2">
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
                    Start Experiment
                  </Button>
                </CardContent>
              </Card>
            </>
          )}
        </form>
      </Form>
    </>
  );
}