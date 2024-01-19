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
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3">
            {tasks.map((task) => (
              <TaskCard {...task} key={task.slug} />
            ))}
            <Link href="https://forms.gle/A5NYYveDiq7gjrMe9" target="_blank" rel="noopener">
              <Card className="flex items-center justify-center h-full bg-gray-200 hover:border-ring hover:text-accent">
                <CardContent className="flex flex-col items-center text-sm">
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
