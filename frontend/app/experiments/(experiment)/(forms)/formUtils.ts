export function groupInputs(inputs: any) {
  // Order and group the inputs by their position and grouping value
  const sortedInputs = Object.entries(inputs || {})
    // @ts-ignore
    .sort(([, a], [, b]) => a.position - b.position);

  const groupedInputs = sortedInputs.reduce((acc: { [key: string]: any }, [key, input]: [string, any]) => {
    // _advanced and any others with _ get added to collapsible
    const sectionName = input.grouping?.startsWith("_") ? "collapsible" : "standard";
    const groupName = input.grouping || "Options";
    if (!acc[sectionName]) {
      acc[sectionName] = {};
    }
    if (!acc[sectionName][groupName]) {
      acc[sectionName][groupName] = {};
    }
    acc[sectionName][groupName][key] = input;
    return acc;
  }, {});

  return groupedInputs;
}

type TransformedJSON = {
  name: string;
  modelId: string;
  walletAddress: string;
  scatteringMethod: string;
  kwargs: { [key: string]: any[] }; // Define the type for kwargs where each key is an array
};

export function transformJson(model: any, originalJson: any, walletAddress: string): TransformedJSON {
  const { name, model: modelId, ...dynamicKeys } = originalJson;

  const modelJsonInputs = model.ModelJson.inputs;

  const kwargs = Object.fromEntries(
    Object.entries(dynamicKeys).map(([key, valueArray]) => {
      // Check if the 'array' property for this key is true
      // @ts-ignore
      if (modelJsonInputs[key] && modelJsonInputs[key]["array"]) {
        // Group the entire array as a single element in another array
        // @ts-ignore
        return [key, [valueArray.map((valueObject) => valueObject.value)]];
      } else {
        // Process normally but check if the field is optional and the value is 0
        // @ts-ignore
        const isOptional = !modelJsonInputs[key]["required"];
        // @ts-ignore
        const value = valueArray.map((valueObject) => {
          if (isOptional && valueObject.value === 0) {
            return null;
          }
          return valueObject.value;
        });
        return [key, value];
      }
    })
  );

  // Return the transformed JSON
  return {
    name: name,
    modelId: model.ID,
    walletAddress: walletAddress,
    scatteringMethod: "crossProduct",
    kwargs: kwargs,
  };
}
