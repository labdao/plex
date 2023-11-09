"use client";
import { BookOpenIcon, GithubIcon } from "lucide-react";
import { notFound } from "next/navigation";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { PageLoader } from "@/components/shared/PageLoader";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { badgeVariants } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
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

  // @TODO: Update GitHub and reference to match API
  const { author, name, description, github_url, reference_url, inputs } = tool?.ToolJson;

  console.log(inputs);

  return (
    <>
      <div className="container mt-8">
        {toolDetailLoading && <PageLoader />}
        {toolDetailError && (
          <Alert variant="destructive">
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{toolDetailError}</AlertDescription>
          </Alert>
        )}
        <>
          <h1 className="mb-4 text-3xl font-heading">
            <span className="text-muted">
              {task.name}/{author}/
            </span>
            {name}
          </h1>
          <p className="max-w-prose line-clamp-3">{description}</p>
          {github_url ||
            (reference_url && (
              <>
                <div className="flex gap-2 mt-4">
                  {github_url && (
                    <a href={github_url} target="_blank" className={badgeVariants({ variant: "muted", size: "lg" })}>
                      github <GithubIcon className="ml-2" size={20} />
                    </a>
                  )}
                  {reference_url && (
                    <a href={reference_url} target="_blank" className={badgeVariants({ variant: "muted", size: "lg" })}>
                      paper <BookOpenIcon className="ml-2" size={20} />
                    </a>
                  )}
                </div>
              </>
            ))}
          <Separator className="my-6" />
          <div className="grid grid-cols-3 gap-8">
            <div className="col-span-2">
              <Card>
                <CardContent></CardContent>
              </Card>
            </div>
            <div>
              <Card>
                <CardContent></CardContent>
              </Card>
            </div>
          </div>
        </>
      </div>
    </>
  );
}
