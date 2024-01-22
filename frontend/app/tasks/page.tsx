import { PlusIcon } from "lucide-react";
import Link from "next/link";

import { Breadcrumbs } from "@/components/global/Breadcrumbs";
import { Card, CardContent } from "@/components/ui/card";

import TaskCard from "./TaskCard";
import { tasks } from "./taskList";

export default function TaskList() {
  return (
    <div className="p-6">
      <Breadcrumbs items={[{ name: "Tasks", href: "/tasks" }]} />
      <h1 className="mb-6 text-4xl font-heading">Select a task to get started</h1>
      <div className="flex flex-wrap gap-8 [&>*]:basis-72 [&>*]:min-h-[15rem]">
        {tasks.map((task) => (
          <TaskCard {...task} key={task.slug} />
        ))}
        <Link href="https://forms.gle/A5NYYveDiq7gjrMe9" target="_blank" rel="noopener">
          <Card className="flex items-center justify-center h-full bg-gray-200 hover:border-ring hover:text-accent">
            <CardContent className="flex flex-col items-center text-sm">
              <PlusIcon size={48} className="mb-1" absoluteStrokeWidth />
              Suggest a task
            </CardContent>
          </Card>
        </Link>
      </div>
    </div>
  );
}
