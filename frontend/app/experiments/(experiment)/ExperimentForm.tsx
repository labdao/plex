"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { usePrivy } from "@privy-io/react-auth";
import { ChevronsUpDownIcon } from "lucide-react";
import { notFound } from "next/navigation";
import { useRouter } from "next/navigation";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import * as z from "zod";

import { tasks } from "@/app/tasks/taskList";
import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import TransactionSummaryInfo from "@/components/payment/TransactionSummaryInfo";
import { ToolSelect } from "@/components/shared/ToolSelect";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  AppDispatch,
  resetToolDetail,
  resetToolList,
  selectToolDetail,
  selectToolDetailError,
  selectToolDetailLoading,
  selectToolList,
  toolDetailThunk,
  toolListThunk,
} from "@/lib/redux";
import { createFlow } from "@/lib/redux/slices/flowAddSlice/asyncActions";

import { DynamicArrayField } from "./DynamicArrayField";
import { ExperimentSummary } from "./ExperimentSummary";
import { generateDefaultValues, generateSchema } from "./formGenerator";
import ModelInfo from "./ModelInfo";

type JsonValueArray = Array<{ value: any }>; // Define the type for the array of value objects

type DynamicKeys = {
  [key: string]: Array<{ value: any }>;
};

type TransformedJSON = {
  name: string;
  toolCid: string;
  walletAddress: string;
  scatteringMethod: string;
  kwargs: { [key: string]: any[] }; // Define the type for kwargs where each key is an array
};

export default function ExperimentForm({ task, showNameField }: { task: any; showNameField: boolean }) {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const { user } = usePrivy();

  const tool = useSelector(selectToolDetail);
  const tools = useSelector(selectToolList);
  const toolDetailLoading = useSelector(selectToolDetailLoading);
  const toolDetailError = useSelector(selectToolDetailError);
  const walletAddress = user?.wallet?.address;
  const { author, name } = tool.ToolJson;

  useEffect(() => {
    return () => {
      dispatch(resetToolDetail());
      dispatch(resetToolList());
    };
  }, [dispatch]);

  // Order and group the inputs by their position and grouping value
  const sortedInputs = Object.entries(tool.ToolJson?.inputs || {})
    // @ts-ignore
    .sort(([, a], [, b]) => a.position - b.position);

  const groupedInputs = sortedInputs.reduce((acc: { [key: string]: any }, [key, input]: [string, any]) => {
    // _advanced and any others with _ get added to collapsible
    const sectionName = input.grouping?.startsWith("_") ? "collapsible" : "standard";
    const groupName = input.grouping || "Options";
    if (!acc[sectionName]) {
      acc[sectionName] = {};
    }
    if (!acc[sectionName][groupName]) {
      acc[sectionName][groupName] = {};
    }
    acc[sectionName][groupName][key] = input;
    return acc;
  }, {});

  const formSchema = generateSchema(tool.ToolJson?.inputs);
  const defaultValues = generateDefaultValues(tool.ToolJson?.inputs, task, tool);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: defaultValues,
  });

  useEffect(() => {
    form.reset(generateDefaultValues(tool.ToolJson?.inputs, task, tool));
  }, [tool, form, task]);

  // If the tool changes, fetch new tool details
  function transformJson(originalJson: any, walletAddress: string): TransformedJSON {
    const { name, tool: toolCid, ...dynamicKeys } = originalJson;

    const toolJsonInputs = tool.ToolJson.inputs;

    const kwargs = Object.fromEntries(
      Object.entries(dynamicKeys).map(([key, valueArray]) => {
        console.log(valueArray);
        // Check if the 'array' property for this key is true
        // @ts-ignore
        if (toolJsonInputs[key] && toolJsonInputs[key]["array"]) {
          // Group the entire array as a single element in another array
          // @ts-ignore
          return [key, [valueArray.map((valueObject) => valueObject.value)]];
        } else {
          // Process normally
          // @ts-ignore
          return [key, valueArray.map((valueObject) => valueObject.value)];
        }
      })
    );

    // Return the transformed JSON
    return {
      name: name,
      toolCid: tool.CID,
      walletAddress: walletAddress,
      scatteringMethod: "crossProduct",
      kwargs: kwargs,
    };
  }

  async function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("===== Form Submitted =====", values);

    console.log("transformed payload");
    if (!walletAddress) {
      console.error("Wallet address missing");
      return;
    }
    const transformedPayload = transformJson(values, walletAddress);
    console.log(transformedPayload);
    console.log("submitting");

    try {
      const response = await createFlow(transformedPayload);
      if (response && response.ID) {
        console.log("Flow created", response);
        console.log(response.ID);
        // Redirecting to another page, for example, a success page or dashboard
        router.push(`/experiments/${response.ID}`);
      } else {
        console.log("something went wrong", response);
      }
    } catch (error) {
      console.error("Failed to create flow", error);
      // Handle error, maybe show message to user
    }
  }

  return (
    <>
      <ProtectedComponent method="hide" message="Log in to run an experiment">
        <Form {...form}>
          <form id="task-form" onSubmit={form.handleSubmit((values) => onSubmit(values))}>
            {showNameField && (
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
            )}
            {!toolDetailLoading && (
              <>
                <Card>
                  {Object.keys(groupedInputs?.standard || {}).map((groupKey) => {
                    return (
                      <CardContent key={groupKey} className="border-t first:border-0">
                        <div className="uppercase font-heading">{groupKey}</div>
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
                  <ExperimentSummary sortedInputs={sortedInputs} form={form} tool={tool} />
                </Card>
              </>
            )}
          </form>
        </Form>
      </ProtectedComponent>
    </>
  );
}
