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
  toolCid: string;
  walletAddress: string;
  scatteringMethod: string;
  kwargs: { [key: string]: any[] }; // Define the type for kwargs where each key is an array
};

export function transformJson(model: any, originalJson: any, walletAddress: string): TransformedJSON {
  const { name, model: toolCid, ...dynamicKeys } = originalJson;

  const toolJsonInputs = model.ToolJson.inputs;

  const kwargs = Object.fromEntries(
    Object.entries(dynamicKeys).map(([key, valueArray]) => {
      // Check if the 'array' property for this key is true
      // @ts-ignore
      if (toolJsonInputs[key] && toolJsonInputs[key]["array"]) {
        // Group the entire array as a single element in another array
        // @ts-ignore
        return [key, [valueArray.map((valueObject) => valueObject.value)]];
      } else {
        // Process normally
        // @ts-ignore
        return [key, valueArray.map((valueObject) => valueObject.value)];
      }
    })
  );

  // Return the transformed JSON
  return {
    name: name,
    toolCid: model.CID,
    walletAddress: walletAddress,
    scatteringMethod: "crossProduct",
    kwargs: kwargs,
  };
}
