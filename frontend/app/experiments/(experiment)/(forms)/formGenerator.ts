import dayjs from "dayjs";
import { useEffect, useState } from "react";
import * as z from "zod";

import { ModelDetail } from "@/lib/redux";

type InputType = { [key: string]: any };

const inputsToSchema = (inputs: InputType) => {
  const schema: { [key: string]: z.ZodArray<z.ZodObject<{ value: z.ZodString | z.ZodNumber | z.ZodBoolean }>> } = {};
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
    } else if (input.type === "bool" || input.type === "boolean") {
      schema[key] = z.array(
        z.object({
          value: z.boolean(),
        })
      );
    }
  }
  return schema;
};

// Take the default values indicated on a model and format them for the form
const inputsToDefaultValues = (inputs: InputType) => {
  const defaultValues: { [key: string]: { value: string | number }[] } = {};
  for (var key in inputs) {
    const input = inputs[key];
    defaultValues[key] = [{ value: input.default }];
  }
  return defaultValues;
};

// Take the inputs from a experiment job and format them for the form
const inputsToValues = (inputs: InputType) => {
  const values: { [key: string]: { value: string | number }[] } = {};
  for (var key in inputs) {
    const input = inputs[key];
    values[key] = [{ value: input }];
  }
  return values;
};

export function generateSchema(inputs: InputType) {
  return z.object({
    name: z.string().min(1, { message: "Name is required" }),
    model: z.string().min(1, { message: "Model is required" }),
    ...inputsToSchema(inputs),
  });
}


export function generateRerunSchema(inputs: InputType = {}, jobInputs: any = {}) {
  const schema: { [key: string]: z.ZodArray<z.ZodObject<{ value: z.ZodTypeAny }>> } = {};

  for (const key in inputs) {
    const input = inputs[key] || {}; // Default to an empty object if input is null/undefined
    const originalValue = jobInputs[key]; // Job inputs are likely not nested in an array

    let valueSchema: z.ZodTypeAny;

    switch (input.type) {
      case "File":
      case "file":
        valueSchema = z.string().min(input.required ? 1 : 0, { message: "File is required" });
        break;
      case "string":
        valueSchema = input.required
          ? z.string().min(1, { message: "Field is required" })
          : z.string().optional().default(originalValue as string || "");
        break;
      case "int":
      case "number":
        const min = Number.isFinite(parseInt(input.min)) ? parseInt(input.min) : Number.MIN_SAFE_INTEGER;
        const max = Number.isFinite(parseInt(input.max)) ? parseInt(input.max) : Number.MAX_SAFE_INTEGER;
        valueSchema = input.required
          ? z.coerce.number().int().min(min).max(max)
          : z.union([z.coerce.number().int().min(min).max(max), z.null()]).optional().default(originalValue as number | null ?? null);
        break;
      case "bool":
      case "boolean":
        valueSchema = z.boolean().optional().default(originalValue as boolean ?? false);
        break;
      default:
        // Default case, treat as optional string
        valueSchema = z.string().optional().default(originalValue as string || "");
    }

    schema[key] = z.array(z.object({ value: valueSchema }));
  }

  console.log("Generated schema:", schema); // Add this for debugging

  return z.object(schema).optional().default({});
}

export function generateDefaultValues(inputs: InputType, task: { slug: string }, model: ModelDetail) {
  return {
    name: `${task?.slug}-${dayjs().format("YYYY-MM-DD-hh-mm-ss")}`,
    model: model?.S3URI,
    ...inputsToDefaultValues(inputs),
  };
}

export function generateValues(inputs: InputType = {}, jobInputs: any = {}) {
  const values: { [key: string]: { value: string | number | boolean | null }[] } = {};

  for (const key in inputs) {
    const input = inputs[key] || {}; // Default to an empty object if input is null/undefined
    const originalValue = jobInputs[key]; // Job inputs are likely not nested in an array

    let value: string | number | boolean | null;

    switch (input.type) {
      case "string":
        value = (originalValue as string) || "";
        break;
      case "int":
      case "number":
        value = (originalValue as number) ?? null;
        break;
      case "bool":
      case "boolean":
        value = (originalValue as boolean) ?? false;
        break;
      default:
        // For other types (including file), use the original value or null
        value = originalValue ?? null;
    }

    values[key] = [{ value }];
  }

  console.log("Generated values:", values); // Add this for debugging

  return values;
}
