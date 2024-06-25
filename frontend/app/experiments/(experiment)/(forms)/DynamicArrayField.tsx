"use client";

import { PlusIcon, Trash2Icon } from "lucide-react";
import React from "react";
import { useFieldArray, UseFormReturn } from "react-hook-form";

import { FileSelect } from "@/components/shared/FileSelect";
import { BooleanInput } from "@/components/ui/boolean-input";
import { Button } from "@/components/ui/button";
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { LabelDescription } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

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
    <div className="relative pb-6 space-y-0 border-b">
      {fields.map((field, index) => (
        <FormField
          key={field.id}
          control={form.control}
          name={`${inputKey}.${index}.value`}
          render={({ field }) => (
            <FormItem className="group">
              <FormLabel className="relative flex items-center justify-between gap-2">
                {index === 0 && (
                  <span>
                    <span>{inputKey.replaceAll("_", " ")}</span>{" "}
                    <LabelDescription>
                      {input?.type} {input?.array ? "array" : ""}
                    </LabelDescription>{" "}
                  </span>
                )}
                {fields.length > 1 && (
                  <Button
                    className="absolute z-20 invisible hover:text-destructive hover:bg-destructive-foreground -right-4 -bottom-10 group-hover:visible"
                    size="icon"
                    variant="secondary"
                    onClick={() => remove(index)}
                  >
                    <Trash2Icon size={18} />
                  </Button>
                )}
              </FormLabel>
              <div>
                <FormControl>
                  <>
                    {(input.type === "File" || input.type === "file") && (
                      <FileSelect
                        onValueChange={field.onChange}
                        value={field.value}
                        globPatterns={input.glob}
                        label={input.glob && `${input.glob.join(", ")}`}
                      />
                    )}
                    {input.type === "string" && <Textarea placeholder={input.example ? `e.g. ${input.example}` : ""} {...field} />}
                    {(input.type === "number" || input.type === "int") && <Input type="number" {...field} />}
                    {(input.type === "bool" || input.type === "boolean") && (
                      <BooleanInput {...field} checked={field.value} onChange={(e) => field.onChange(e.target.checked)} />
                    )}
                  </>
                </FormControl>
              </div>
              {index === fields?.length - 1 && (
                <FormDescription>
                  <LabelDescription>{!input?.required && input?.default === "" ? "(Optional)" : ""}</LabelDescription>
                  {input?.description}
                </FormDescription>
              )}
              <FormMessage />
            </FormItem>
          )}
        />
      ))}
      <Button type="button" variant="secondary" size="icon" className="absolute right-0 -bottom-4" onClick={() => append({ value: input.default })}>
        <PlusIcon size={24} />
      </Button>
    </div>
  );
}
