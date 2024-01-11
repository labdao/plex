import { PlusIcon } from "lucide-react";
import Link from "next/link";

import { Card, CardContent } from "@/components/ui/card";
import { Carousel, CarouselContent, CarouselItem, CarouselNext, CarouselPrevious } from "@/components/ui/carousel";

import TaskCard from "./TaskCard";
import { tasks } from "./taskList";

export default function TaskList() {
  return (
    <>
      <div className="mt-8 ">
        <div className="container">
          <h1 className="mb-12 text-4xl font-bold font-heading">Select a task to get started</h1>
        </div>
        <Carousel className="w-full" opts={{ slidesToScroll: "auto" }}>
          <CarouselContent className="pl-8 xl:pl-[calc((100vw-(1176px))/2)]">
            {tasks.map((task) => (
              <CarouselItem className="pr-4 sm:basis-[40%] md:basis-[30%] lg:basis-[22%]" key={task.slug}>
                <TaskCard {...task} />
              </CarouselItem>
            ))}
            <CarouselItem className="pr-4 sm:basis-[40%] md:basis-[30%] lg:basis-[22%]">
              <Link href="https://forms.gle/A5NYYveDiq7gjrMe9" target="_blank" rel="noopener">
                <Card className="flex items-center justify-center h-full bg-gray-200 hover:border-ring hover:text-accent">
                  <CardContent className="flex flex-col items-center text-sm">
                    <PlusIcon size={48} />
                    Suggest a task
                  </CardContent>
                </Card>
              </Link>
            </CarouselItem>
          </CarouselContent>
          <div className="container flex gap-2 mt-4">
            <CarouselPrevious className="static translate-y-0" />
            <CarouselNext className="static translate-y-0" />
          </div>
        </Carousel>
      </div>
    </>
  );
}
