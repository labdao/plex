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

export function generateRerunSchema(inputs: InputType) {
  return z.object({
    ...inputsToSchema(inputs),
  });
}

export function generateDefaultValues(inputs: InputType, task: { slug: string }, model: ModelDetail) {
  return {
    name: `${task?.slug}-${dayjs().format("YYYY-MM-DD-mm-ss")}`,
    model: model?.s3uri || "s3://convexity/258462036b9903fb51af718933b222b56b6ab49b5235b332b9c60167cddd0256/relay_sample_rfdiffusion_both_ipfs_s3_v6.json",
    ...inputsToDefaultValues(inputs),
  };
}

export function generateValues(inputs: InputType) {
  return {
    ...inputsToValues(inputs),
  };
}
