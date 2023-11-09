"use client";
import { BookOpenIcon, GithubIcon } from "lucide-react";
import { notFound, usePathname, useRouter, useSearchParams } from "next/navigation";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { PageLoader } from "@/components/shared/PageLoader";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { badgeVariants } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
import { Separator } from "@/components/ui/separator";
import { AppDispatch, selectToolDetail, selectToolDetailError, selectToolDetailLoading, toolDetailThunk, toolPatchDetailThunk } from "@/lib/redux";

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
  const loading = useSelector(selectToolDetailLoading);
  const error = useSelector(selectToolDetailError);

  useEffect(() => {
    const defaultToolCID = task.default_tool?.CID;
    if (defaultToolCID) {
      dispatch(toolDetailThunk(defaultToolCID));
    }
  }, [dispatch, task.default_tool?.CID]);

  // @TODO: GitHub and Paper URLs may be different in actual API
  const { author, name, description, github_url, paper_url, inputs } = tool?.ToolJson;

  return (
    <>
      <div className="container mt-8">
        {loading && <PageLoader />}
        {error && (
          <Alert variant="destructive">
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
        {!loading && (
          <>
            <h1 className="mb-4 text-3xl font-heading">
              <span className="text-muted">
                {task.name}/{author}/
              </span>
              {name}
            </h1>
            <p className="max-w-prose line-clamp-3">{description}</p>
            {github_url ||
              (paper_url && (
                <>
                  <div className="flex gap-2 mt-4">
                    {github_url && (
                      <a href={github_url} target="_blank" className={badgeVariants({ variant: "muted", size: "lg" })}>
                        github <GithubIcon className="ml-2" size={20} />
                      </a>
                    )}
                    {paper_url && (
                      <a href={paper_url} target="_blank" className={badgeVariants({ variant: "muted", size: "lg" })}>
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
        )}
      </div>
    </>
  );
}
