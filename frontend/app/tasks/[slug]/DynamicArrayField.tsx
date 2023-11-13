"use client";

import { PlusIcon, Trash2Icon } from "lucide-react";
import React from "react";
import { useFieldArray, UseFormReturn } from "react-hook-form";

import { DataFileSelect } from "@/components/shared/DataFileSelect";
import { Button } from "@/components/ui/button";
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { LabelDescription } from "@/components/ui/label";
import { cn } from "@/lib/utils";

interface DynamicFieldProps {
  form: UseFormReturn<any>;
  input: any;
  inputKey: string;
}

export function DynamicArrayField({ input, inputKey, form }: DynamicFieldProps) {
  const { fields, append, remove } = useFieldArray({
    name: inputKey,
    control: form.control,
  });

  const hasError = form.formState.errors?.[inputKey];

  return (
    <div className={cn("p-4 space-y-4 border rounded-lg", hasError && "border-destructive")}>
      {fields.map((field, index) => (
        <FormField
          key={field.id}
          control={form.control}
          name={`${inputKey}.${index}.value`}
          render={({ field }) => (
            <FormItem className="group">
              <FormLabel className="flex items-center justify-between gap-2">
                <span>
                  <span>{inputKey.replaceAll("_", " ")}</span>{" "}
                  <LabelDescription>
                    {input?.type} {input?.array ? "array" : ""}
                  </LabelDescription>{" "}
                </span>
                {fields.length > 1 && (
                  <Button className="invisible -mt-2 -mb-1 group-hover:visible" size="icon" variant="ghost" onClick={() => remove(index)}>
                    <Trash2Icon size={18} />
                  </Button>
                )}
              </FormLabel>
              <div>
                <FormControl>
                  <>
                    {input.type === "File" && (
                      <DataFileSelect
                        onChange={field.onChange}
                        value={field.value}
                        globPatterns={input.glob}
                        label={input.glob && `${input.glob.join(", ")}`}
                      />
                    )}
                    {input.type === "string" && <Input placeholder={input.example ? `e.g. ${input.example}` : ""} {...field} />}
                    {input.type === "int" && <Input type="number" {...field} />}
                  </>
                </FormControl>
              </div>
              <FormDescription>
                <LabelDescription>{!input?.required && input?.default === "" ? "(Optional)" : ""}</LabelDescription>
                {input?.description}
              </FormDescription>
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
