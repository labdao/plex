"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { notFound } from "next/navigation";
import React, { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import * as z from "zod";

import { PageLoader } from "@/components/shared/PageLoader";
import { ToolSelect } from "@/components/shared/ToolSelect";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { LabelDescription } from "@/components/ui/label";
import { AppDispatch, selectToolDetail, selectToolDetailError, selectToolDetailLoading, toolDetailThunk } from "@/lib/redux";

import { DynamicArrayField } from "./DynamicArrayField";
import { generateDefaultValues, generateSchema } from "./formGenerator";
import TaskPageHeader from "./TaskPageHeader";
import { ChevronsUpDownIcon } from "lucide-react";
import { VariantSummary } from "./VariantSummary";

import { createFlow } from "@/lib/redux/slices/flowAddSlice/asyncActions";

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

export default function TaskDetail({ params }: { params: { slug: string } }) {
  const dispatch = useDispatch<AppDispatch>();

  // Temporarily hardcode the task - we'll fetch this from the API later based on the page slug
  const task = useMemo(
    () => ({
      name: "protein design",
      slug: "protein-design", //Could fetch by a slug or ID, whatever you want the url to be
      default_tool: {
        //CID: "QmXHQhZG8PSNY9g8MAcsHjGYWSrbZpPjaoPPk9QoRZCT3w",
        CID: "QmbxyLKaZg73PvnREPdVitKziw2xTDjTp268VNy1hMkR5E",
      },
    }),
    []
  );

  if (task.slug !== params?.slug) {
    notFound();
  }

  const tool = useSelector(selectToolDetail);
  const toolDetailLoading = useSelector(selectToolDetailLoading);
  const toolDetailError = useSelector(selectToolDetailError);

  // On page load fetch the default tool details
  useEffect(() => {
    const defaultToolCID = task.default_tool?.CID;
    if (defaultToolCID) {
      dispatch(toolDetailThunk(defaultToolCID));
    }
  }, [dispatch, task.default_tool?.CID]);

  // Order and group the inputs by their position and grouping value
  const sortedInputs = Object.entries(tool.ToolJson?.inputs)
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

  function transformJson(
    originalJson: any,
    walletAddress: string
  ): TransformedJSON {
    const { name, tool, ...dynamicKeys } = originalJson;

    // Define the transformation for the dynamic keys using a more specific type
    // @ts-ignore
    const kwargs = Object.fromEntries(
      // @ts-ignore
      Object.entries<DynamicKeys>(dynamicKeys).map(([key, valueArray]: [string, JsonValueArray]) => [
        key,
        valueArray.map((valueObject) => valueObject.value),
      ])
    );

    // Return the transformed JSON
    return {
      name: name,
      toolCid: tool,
      walletAddress: walletAddress,
      scatteringMethod: "crossProduct",
      kwargs: kwargs,
    };
  }

  function onSubmit(values: z.infer<typeof formSchema>) {
    // submit here
    console.log("===== Form Submitted =====", values);
    /*
    payload='{
          "name": "test equibind new backend",
          "walletAddress": "0xab5801a7d398351b8be11c439e05c5b3259aec9b",
          "toolCid": "QmXt1WFynKUcQCR5ihPz3gFmXnLoq5pFGHFEj32oRxPwxE",
          "scatteringMethod": "crossProduct",
          "kwargs": {
            "protein": ["QmUWCBTqbRaKkPXQ3M14NkUuM4TEwfhVfrqLNoBB7syyyd/7n9g.pdb"],
            "small_molecule": ["QmV6qVzdQLNM6SyEDB3rJ5R5BYJsQwQTn1fjmPzvCCkCYz/ZINC000003986735.sdf", "QmViB4EnKX6PXd77WYSgMDMq9ZMX14peu3ZNoVV1LHUZwS/ZINC000019632618.sdf"]
          }
        }'
    */
   /*
   {
    "name": "protein-design-2023-11-13-07-27",
    "tool": "QmbxyLKaZg73PvnREPdVitKziw2xTDjTp268VNy1hMkR5E",
    "binder_length": [
        {
            "value": 50
        },
        {
            "value": 75
        }
    ],
    "full_prompt": [
        {
            "value": ""
        }
    ],
    "hotspot": [
        {
            "value": "A30,A33"
        }
    ],
    "number_of_binders": [
        {
            "value": 2
        }
    ],
    "target_chain": [
        {
            "value": "A"
        }
    ],
    "target_end_residue": [
        {
            "value": 1
        }
    ],
    "target_protein": [
        {
            "value": "QmVM6i5sHMyUm9qrSsAzLLWRDYAResgMttQj9s2D8Qb8dJ/best.pdb"
        }
    ],
    "target_start_residue": [
        {
            "value": 1
        }
    ]
}
*/

    console.log('transformed payload')
    console.log(transformJson(values, "0xab5801a7d398351b8be11c439e05c5b3259aec9b"));
    console.log('submitting')
    createFlow(transformJson(values, "0xab5801a7d398351b8be11c439e05c5b3259aec9b"))
  }

  return (
    <>
      <div className="container mt-8">
        {toolDetailError && (
          <Alert variant="destructive">
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{toolDetailError}</AlertDescription>
          </Alert>
        )}
        <>
          <TaskPageHeader tool={tool} task={task} />

          <div className="grid grid-cols-3 gap-8">
            <div className="col-span-2">
              <Form {...form}>
                <form id="task-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                  <Card>
                    <CardContent className="space-y-4">
                      <FormField
                        control={form.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem>
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
                            <FormLabel>Model</FormLabel>
                            <FormControl>
                              <ToolSelect onChange={field.onChange} defaultValue={tool?.CID} />
                            </FormControl>
                            <FormDescription>
                              <a
                                className="text-accent hover:underline"
                                target="_blank"
                                href={`${process.env.NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT}${tool?.CID}/`}
                              >
                                View Tool Manifest
                              </a>
                            </FormDescription>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </CardContent>
                  </Card>
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
                                // @ts-ignore
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
                                    // @ts-ignore
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
              <VariantSummary sortedInputs={sortedInputs} form={form} />
            </div>
          </div>
        </>
        {toolDetailLoading && <PageLoader />}
      </div>
    </>
  );
}
