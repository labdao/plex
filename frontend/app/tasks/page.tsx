import { PlusIcon } from "lucide-react";
import Link from "next/link";

import { Card, CardContent } from "@/components/ui/card";

import TaskCard from "./TaskCard";
import { tasks } from "./taskList";

export default function TaskList() {
  return (
    <>
      <div className="mt-8 ">
        <div className="container">
          <h1 className="mb-12 text-4xl font-bold font-heading">Select a task to get started</h1>
        </div>
        <div className="flex gap-8 overflow-x-auto no-scrollbar pl-8 xl:pl-[calc((100vw-(1176px))/2)]">
          {tasks.map((task) => (
            <div className="w-60 shrink-0" key={task.slug}>
              <TaskCard {...task} />
            </div>
          ))}
          <div className="w-60 shrink-0">
            <Link href="mailto:info@convexity.bio?subject=Task Suggestion">
              <Card className="flex items-center justify-center h-full hover:border-ring hover:text-accent">
                <CardContent className="flex flex-col items-center">
                  <PlusIcon size={48} />
                  Suggest a task
                </CardContent>
              </Card>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}
