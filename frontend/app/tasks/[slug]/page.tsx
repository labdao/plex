"use client";
import { notFound, usePathname, useRouter, useSearchParams } from "next/navigation";
import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { PageLoader } from "@/components/shared/PageLoader";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardTitle } from "@/components/ui/card";
import { DataTable } from "@/components/ui/data-table";
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

  console.log(tool);

  const { author, name, description, inputs } = tool?.ToolJson;

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
            <div>
              {task.name}/{author}
              {name}
            </div>
          </>
        )}
      </div>
    </>
  );
}
