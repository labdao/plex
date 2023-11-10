"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import dayjs from "dayjs";
import { BookOpenIcon, GithubIcon, PlusIcon } from "lucide-react";
import { notFound } from "next/navigation";
import React, { useEffect } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import * as z from "zod";

import { DataFileSelect } from "@/components/shared/DataFileSelect";
import { PageLoader } from "@/components/shared/PageLoader";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { badgeVariants } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { LabelDescription } from "@/components/ui/label";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import {
  AppDispatch,
  selectToolDetail,
  selectToolDetailError,
  selectToolDetailLoading,
  selectToolList,
  selectToolListError,
  toolDetailThunk,
  toolListThunk,
} from "@/lib/redux";
import { cn } from "@/lib/utils";

import mockToolJson from "./mockToolJson";

export default function TaskDetail({ params }: { params: { slug: string } }) {
  // Temporarily hardcode the task - we'll fetch this from the API later based on the page slug
  const task = {
    name: "protein design",
    slug: "protein-design", //Could fetch by a slug or ID, whatever you want the url to be
    default_tool: {
      CID: "QmTLQJWhh7iL7sHPpDfve8p1CiW6xu8Mzraamjugaku4j8",
    },
  };

  if (task.slug !== params?.slug) {
    notFound();
  }

  const dispatch = useDispatch<AppDispatch>();

  const tool = useSelector(selectToolDetail);
  const toolDetailLoading = useSelector(selectToolDetailLoading);
  const toolDetailError = useSelector(selectToolDetailError);

  const tools = useSelector(selectToolList);
  const toolListError = useSelector(selectToolListError);

  useEffect(() => {
    const defaultToolCID = task.default_tool?.CID;
    if (defaultToolCID) {
      dispatch(toolDetailThunk(defaultToolCID));
    }
    dispatch(toolListThunk());
  }, [dispatch, task.default_tool?.CID]);

  const { author, name, description, github, paper, inputs } = mockToolJson;

  console.log(inputs);

  const inputsToSchema = (inputs: {}) => {
    const schema = {};
    for (var key in inputs) {
      const input = inputs[key];
      if (input.type === "File") {
        schema[key] = z.array(
          z.object({
            value: z.string().min(input.required ? 1 : 0),
          })
        );
      } else if (input.type === "string") {
        schema[key] = z.array(
          z.object({
            value: z.string().min(input.required ? 1 : 0),
          })
        );
      } else if (input.type === "int") {
        // Hacky
        const min = parseInt(input.min) || Number.MIN_SAFE_INTEGER;
        const max = parseInt(input.max) || Number.MAX_SAFE_INTEGER;
        schema[key] = z.array(
          z.object({
            value: z.coerce.number().int().min(min).max(max),
          })
        );
      }
    }
    return schema;
  };

  const inputsToDefaultValues = (inputs: {}) => {
    const defaultValues = {};
    for (var key in inputs) {
      const input = inputs[key];
      defaultValues[key] = [{ value: input.default }];
    }
    return defaultValues;
  };

  const formSchema = z.object({
    name: z.string().min(1, { message: "Name is required" }),
    tool: z.string().min(1, { message: "Tool is required" }),
    ...inputsToSchema(inputs),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: `${task.slug}-${dayjs().format("YYYY-MM-DD-mm-ss")}`,
      tool: tool?.CID,
      ...inputsToDefaultValues(inputs),
    },
  });

  // If the tool changes, fetch new tool details
  useEffect(() => {
    const subscription = form.watch((value, { name, type }) => {
      if (name === "tool" && value?.tool) {
        dispatch(toolDetailThunk(value.tool));
      }
    });
    return () => subscription.unsubscribe();
  }, [dispatch, form]);

  function onSubmit(values: z.infer<typeof formSchema>) {
    console.log("===== Form Submitted =====", values);
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
          <h1 className="mb-4 text-3xl font-heading">
            <span className="text-muted-foreground">
              {task.name}/{author || "unknown"}/
            </span>
            {name}
          </h1>
          <p className="max-w-prose line-clamp-2 min-h-[3em]">{description}</p>
          {(github || paper) && (
            <div className="flex gap-2 mt-4">
              {github && (
                <a href={github} target="_blank" className={badgeVariants({ variant: "muted", size: "lg" })}>
                  github <GithubIcon className="ml-2" size={20} />
                </a>
              )}
              {paper && (
                <a href={paper} target="_blank" className={badgeVariants({ variant: "muted", size: "lg" })}>
                  paper <BookOpenIcon className="ml-2" size={20} />
                </a>
              )}
            </div>
          )}
          <Separator className="my-6" />
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
                              <Select onValueChange={field.onChange} defaultValue={tool?.CID}>
                                <SelectTrigger>
                                  <SelectValue placeholder="Select a model" />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectGroup>
                                    {tools.map((tool, index) => {
                                      return (
                                        <SelectItem key={index} value={tool?.CID}>
                                          {tool?.ToolJson?.author || "unknown"}/{tool.Name}
                                        </SelectItem>
                                      );
                                    })}
                                  </SelectGroup>
                                </SelectContent>
                              </Select>
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
                      {!toolDetailLoading && (
                        <>
                          {Object.keys(inputs || {}).map((key) => {
                            // @ts-ignore
                            const input = inputs[key];
                            return <DynamicField key={key} inputKey={key} form={form} input={input} />;
                          })}
                        </>
                      )}
                    </CardContent>
                  </Card>
                </form>
              </Form>
            </div>
            <div>
              <Card>
                <CardContent>
                  <Button type="submit" form="task-form" className="w-full">
                    Submit
                  </Button>
                </CardContent>
              </Card>
            </div>
          </div>
        </>
        {toolDetailLoading && <PageLoader />}
      </div>
    </>
  );
}

function DynamicField({ input, inputKey, form }) {
  const { fields, append } = useFieldArray({
    name: inputKey,
    control: form.control,
  });

  return (
    <div className="p-4 border rounded-lg">
      {fields.map((field, index) => (
        <FormField
          key={field.id}
          control={form.control}
          name={`${inputKey}.${index}.value`}
          render={({ field }) => (
            <FormItem>
              <FormLabel>
                {inputKey.replaceAll("_", " ")} <LabelDescription>{input?.type}</LabelDescription>
              </FormLabel>
              <FormControl>
                <>
                  {input.type === "File" && (
                    <DataFileSelect
                      onSelect={field.onChange}
                      value={field.value}
                      globPatterns={input.glob}
                      label={input.glob && `${input.glob.join(", ")}`}
                    />
                  )}
                  {input.type === "string" && <Input {...field} />}
                  {input.type === "int" && <Input type="number" {...field} />}
                </>
              </FormControl>
              <FormDescription>{input?.description}</FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      ))}
      <Button type="button" variant="secondary" size="sm" className="w-full mt-2" onClick={() => append({ value: input.default })}>
        <PlusIcon />
      </Button>
    </div>
  );
}
