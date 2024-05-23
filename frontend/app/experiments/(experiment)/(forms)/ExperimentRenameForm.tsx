"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";

import { Form, FormControl, FormField, FormItem, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { flowUpdateThunk, useDispatch } from "@/lib/redux";

const formSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
});

export function ExperimentRenameForm({ initialName, inputProps, flowId }: { initialName: string; inputProps?: any; flowId: string }) {
  const dispatch = useDispatch();
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: { name: initialName },
    mode: "onChange",
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    await dispatch(flowUpdateThunk({ flowId, updates: { name: values.name } }));
    toast.success("Experiment renamed successfully!");
  };

  return (
    <Form {...form}>
      <form onBlur={form.handleSubmit(onSubmit)} className="w-full">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input {...inputProps} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}
