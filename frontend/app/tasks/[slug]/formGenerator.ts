import dayjs from "dayjs";
import { useEffect, useState } from "react";
import * as z from "zod";

import { ToolDetail } from "@/lib/redux";

type InputType = { [key: string]: any };

const inputsToSchema = (inputs: InputType) => {
  const schema: { [key: string]: z.ZodArray<z.ZodObject<{ value: z.ZodString | z.ZodNumber }>> } = {};
  for (var key in inputs) {
    const input = inputs[key];
    if (input.type === "File" || input.type === "file") {
      schema[key] = z.array(
        z.object({
          value: z.string().min(input.required ? 1 : 0, { message: "File is required" }),
        })
      );
    } else if (input.type === "string") {
      schema[key] = z.array(
        z.object({
          value: z.string().min(input.required ? 1 : 0, { message: "Field is required" }),
        })
      );
    } else if (input.type === "int" || input.type === "number") {
      // @TODO: depending on API updates parseInt may not be necessary
      // May be a better way of setting no min/max than SAFE_INTEGER
      const min = parseInt(input.min) || Number.MIN_SAFE_INTEGER;
      const max = parseInt(input.max) || Number.MAX_SAFE_INTEGER;
      schema[key] = z.array(
        z.object({
          value: z.coerce.number().int().min(min).max(max),
        })
      );
    }
  }
  return schema;
};

const inputsToDefaultValues = (inputs: InputType) => {
  const defaultValues: { [key: string]: { value: string | number }[] } = {};
  for (var key in inputs) {
    const input = inputs[key];
    defaultValues[key] = [{ value: input.default }];
  }
  return defaultValues;
};

export function generateSchema(inputs: InputType) {
  return z.object({
    name: z.string().min(1, { message: "Name is required" }),
    tool: z.string().min(1, { message: "Model is required" }),
    ...inputsToSchema(inputs),
  });
}

export function generateDefaultValues(inputs: InputType, task: { slug: string }, tool: ToolDetail) {
  return {
    name: `${task.slug}-${dayjs().format("YYYY-MM-DD-mm-ss")}`,
    tool: tool?.CID,
    ...inputsToDefaultValues(inputs),
  };
}
