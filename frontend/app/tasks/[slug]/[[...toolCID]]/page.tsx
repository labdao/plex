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

import ProtectedComponent from "@/components/auth/ProtectedComponent";
import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { ToolSelect } from "@/components/shared/ToolSelect";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { LabelDescription } from "@/components/ui/label";
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

import { tasks } from "../../taskList";
import { DynamicArrayField } from "./DynamicArrayField";
import { generateDefaultValues, generateSchema } from "./formGenerator";
import ModelInfo from "./ModelInfo";
import { TaskSummary } from "./TaskSummary";

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

export default function TaskDetail({ params }: { params: { slug: string; toolCID: string } }) {
  const dispatch = useDispatch<AppDispatch>();
  const router = useRouter();
  const { user } = usePrivy();

  // Fetch the task from our static list of tasks
  const task = tasks.find((task) => task.slug === params?.slug);

  if (!task?.slug || !task?.available) {
    notFound();
  }

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
  }, [params.slug, dispatch]);

  const defaultTool = tools.find((tool) => tool.DefaultTool === true);
  const default_tool_cid = defaultTool?.CID;

  useEffect(() => {
    if (default_tool_cid) {
      dispatch(toolDetailThunk(default_tool_cid));
    }
  }, [dispatch, default_tool_cid]);

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
  useEffect(() => {
    const subscription = form.watch((value, { name, type }) => {
      if (name === "tool" && value?.tool) {
        dispatch(toolDetailThunk(value.tool));
      }
    });
    return () => subscription.unsubscribe();
  }, [dispatch, form]);

  function transformJson(originalJson: any, walletAddress: string): TransformedJSON {
    const { name, tool: toolCid, ...dynamicKeys } = originalJson;

    const toolJsonInputs = tool.ToolJson.inputs;

    const kwargs = Object.fromEntries(
      Object.entries(dynamicKeys).map(([key, valueArray]) => {
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
      toolCid: toolCid,
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

  const renderDescriptionParagraphs = (text: string) => {
    const paragraphs = text.split('\n');
    const hasNumberedSteps = paragraphs.some((paragraph) => paragraph.match(/^\d+\. /));
  
    if (hasNumberedSteps) {
      const steps = paragraphs.filter((paragraph) => paragraph.match(/^\d+\. /));
      const nonStepParagraphs = paragraphs.filter((paragraph) => !paragraph.match(/^\d+\. /));
  
      return (
        <>
          {nonStepParagraphs.map((paragraph, index) => (
            <p key={index} className="mt-2">
              {paragraph}
            </p>
          ))}
          <ol className="list-decimal list-inside mt-2">
            {steps.map((step, index) => (
              <li key={index} className="mt-1">
                {step.replace(/^\d+\. /, '')}
              </li>
            ))}
          </ol>
        </>
      );
    } else {
      return paragraphs.map((paragraph, index) => (
        <p key={index} className="mt-2">
          {paragraph}
        </p>
      ));
    }
  };  

  return (
    <>
      <Breadcrumbs
        items={[
          { name: "Tasks", href: "/tasks" },
          { name: task.name, href: `/tasks/${task.slug}` },
          { name: !toolDetailLoading ? `${author || "unknown"}/${name}` : "" },
        ]}
      />
      <div>
        {toolDetailError && (
          <Alert variant="destructive">
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{toolDetailError}</AlertDescription>
          </Alert>
        )}
        <>
          <ProtectedComponent method="hide" message="Log in to run an experiment">
            <div className="grid min-h-screen grid-cols-1 p-6 lg:pr-0 lg:grid-cols-3">
              <div className="col-span-2">
                <Form {...form}>
                  <form id="task-form" onSubmit={form.handleSubmit((values) => onSubmit(values))} className="space-y-8">
                    <Card>
                      <CardHeader>
                        <CardTitle className="uppercase">Model</CardTitle>
                      </CardHeader>
                      <CardContent>
                        <FormField
                          control={form.control}
                          name="name"
                          render={({ field }) => (
                            <FormItem className="hidden">
                              <FormLabel>
                                Name <LabelDescription>string</LabelDescription>
                              </FormLabel>
                              <FormControl>
                                <Input {...field} />
                              </FormControl>
                              <FormDescription>Name your experiment</FormDescription>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={form.control}
                          name="tool"
                          key={tool?.CID}
                          render={({ field }) => (
                            <FormItem>
                              <FormControl>
                                <ToolSelect onChange={field.onChange} taskSlug={params.slug} defaultValue={tool?.CID} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </CardContent>
                    </Card>
                    {tool.ToolJson?.guide && (
                      <Card>
                        <Collapsible defaultOpen={true}>
                          <CollapsibleTrigger className="flex items-center justify-between w-full p-6 text-left uppercase font-heading">
                            <CardTitle className="uppercase">Guide</CardTitle>
                            <ChevronsUpDownIcon />
                          </CollapsibleTrigger>
                          <CollapsibleContent>
                            <CardContent>
                              <div className="mt-2 space-y-2 text-sm text-muted-foreground">
                                {renderDescriptionParagraphs(tool.ToolJson.guide)}
                              </div>
                            </CardContent>
                          </CollapsibleContent>
                        </Collapsible>
                      </Card>
                    )}
                    {!toolDetailLoading && (
                      <>
                        {Object.keys(groupedInputs?.standard || {}).map((groupKey) => {
                          return (
                            <Card key={groupKey}>
                              <CardHeader>
                                <CardTitle className="uppercase">{groupKey}</CardTitle>
                              </CardHeader>
                              <CardContent className="space-y-4">
                                {Object.keys(groupedInputs?.standard[groupKey] || {}).map((key) => {
                                  const input = groupedInputs?.standard?.[groupKey]?.[key];
                                  return <DynamicArrayField key={key} inputKey={key} form={form} input={input} />;
                                })}
                              </CardContent>
                            </Card>
                          );
                        })}
                        {Object.keys(groupedInputs?.collapsible || {}).map((groupKey) => {
                          return (
                            <Card key={groupKey}>
                              <Collapsible>
                                <CollapsibleTrigger className="flex items-center justify-between w-full p-6 text-left uppercase text-bold font-heading">
                                  {groupKey.replace("_", "")}
                                  <ChevronsUpDownIcon />
                                </CollapsibleTrigger>
                                <CollapsibleContent>
                                  <CardContent className="pt-0 space-y-4">
                                    {Object.keys(groupedInputs?.collapsible[groupKey] || {}).map((key) => {
                                      const input = groupedInputs?.collapsible?.[groupKey]?.[key];
                                      return <DynamicArrayField key={key} inputKey={key} form={form} input={input} />;
                                    })}
                                  </CardContent>
                                </CollapsibleContent>
                              </Collapsible>
                            </Card>
                          );
                        })}
                      </>
                    )}
                  </form>
                </Form>
              </div>
              <div>
                <TaskSummary
                  sortedInputs={sortedInputs}
                  form={form}
                  outputs={tool?.ToolJson?.outputs}
                  // description={tool.ToolJson?.description}
                  tool={tool}
                />
              </div>
            </div>
          </ProtectedComponent>
        </>
      </div>
    </>
  );
}