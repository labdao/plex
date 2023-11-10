"use client";
import { zodResolver } from "@hookform/resolvers/zod";
import { BookOpenIcon, GithubIcon, PlusIcon } from "lucide-react";
import { notFound } from "next/navigation";
import React, { useEffect } from "react";
import { useForm } from "react-hook-form";
import { useDispatch, useSelector } from "react-redux";
import * as z from "zod";

import { PageLoader } from "@/components/shared/PageLoader";
import { ToolSelect } from "@/components/shared/ToolSelect";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { badgeVariants } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { LabelDescription } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { AppDispatch, selectToolDetail, selectToolDetailError, selectToolDetailLoading, toolDetailThunk, toolListThunk } from "@/lib/redux";

import { DynamicArrayField } from "./DynamicArrayField";
import formGenerator from "./formGenerator";
import mockToolJson from "./mockToolJson";

export default function TaskDetail({ params }: { params: { slug: string } }) {
  const dispatch = useDispatch<AppDispatch>();

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

  const tool = useSelector(selectToolDetail);
  const toolDetailLoading = useSelector(selectToolDetailLoading);
  const toolDetailError = useSelector(selectToolDetailError);

  // On page load fetch the default tool details
  useEffect(() => {
    const defaultToolCID = task.default_tool?.CID;
    if (defaultToolCID) {
      dispatch(toolDetailThunk(defaultToolCID));
    }
    dispatch(toolListThunk());
  }, [dispatch, task.default_tool?.CID]);

  const { author, name, description, github, paper, inputs } = mockToolJson; //@TODO: replace with tool.ToolJson
  const { generatedFormSchema, generatedDefaultValues } = formGenerator(inputs, task, tool);

  const form = useForm<z.infer<typeof generatedFormSchema>>({
    resolver: zodResolver(generatedFormSchema),
    defaultValues: generatedDefaultValues,
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

  function onSubmit(values: z.infer<typeof generatedFormSchema>) {
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
                  <Card>
                    <CardContent className="space-y-4">
                      {!toolDetailLoading && (
                        <>
                          {Object.keys(inputs || {}).map((key) => {
                            // @ts-ignore
                            const input = inputs[key];
                            return <DynamicArrayField key={key} inputKey={key} form={form} input={input} />;
                          })}
                        </>
                      )}
                    </CardContent>
                  </Card>
                </form>
              </Form>
            </div>
            <div>
              <Card className="sticky top-4">
                <CardContent>
                  <Button type="submit" form="task-form" size="lg" className="w-full">
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
