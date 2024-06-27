"use client";

import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AppDispatch, selectModelDetail, selectModelList, selectModelListError, modelDetailThunk, modelListThunk } from "@/lib/redux";

interface ModelSelectProps {
  onChange: (value: string) => void;
  taskSlug?: string;
}

export function ModelSelect({ onChange, taskSlug }: ModelSelectProps) {
  const dispatch = useDispatch<AppDispatch>();

  const model = useSelector(selectModelDetail);
  const models = useSelector(selectModelList);
  const modelListError = useSelector(selectModelListError);

  useEffect(() => {
    if (taskSlug) {
      dispatch(modelListThunk(taskSlug));
    } else {
      dispatch(modelListThunk());
    }
  }, [dispatch, taskSlug]);

  useEffect(() => {
    const defaultModel = models.find((model) => model.DefaultModel === true);
    const defaultModelCID = defaultModel?.ID;

    if (defaultModelCID) {
      dispatch(modelDetailThunk(defaultModelCID));
    }
  }, [dispatch, models]);

  const handleSelectionChange = (value: string) => {
    onChange(value);
  };

  return (
    <Select onValueChange={handleSelectionChange} value={model?.ID}>
      <SelectTrigger className="font-mono text-sm font-bold uppercase rounded-full border-primary hover:bg-primary/10">
        <SelectValue placeholder="Select a model" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {models.map((model, index) => {
            return (
              <SelectItem key={index} value={model?.ID}>
                {model.Name}
              </SelectItem>
            );
          })}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
