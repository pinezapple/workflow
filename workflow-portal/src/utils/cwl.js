const { Document, YAMLMap, YAMLSeq } = require("yaml");

export function generateCwl(tools, toolLinks) {
  const cwlTools = new Map();
  for (const tool of tools) {
    const cwlTool = generateCWLTool(tool);
    cwlTools.set("tasks/" + tool.name + "_" + tool.tool_index + ".cwl", cwlTool)
  }

  const cwlWorkflow = generateCwlWorkflow(tools, toolLinks)
  return { cwlWorkflow: cwlWorkflow, cwlSteps: cwlTools };
}

function generateCWLTool(tool) {
  const version = tool.selected_version;

  const cwlTool = new Document();
  cwlTool.set("class", "CommandLineTool");
  cwlTool.set("cwlVersion", "v1.2");

  const requirements = new YAMLSeq();
  requirements.add(cwlTool.createNode({ class: "DockerRequirement", dockerPull: version.image }));
  cwlTool.set("requirements", requirements);

  const mapInputs = new YAMLMap();
  for (const input of version.inputs) {
    const inputPair = cwlTool.createPair(input.abbrev, { type: input.type });
    mapInputs.add(inputPair);
  }
  cwlTool.set("inputs", mapInputs);

  const mapOutputs = new YAMLMap();
  for (const output of version.outputs) {
    // const outputPair = cwlTool.createPair(output.abbrev, { type: output.type, outputBinding: { glob: "*" + output.extension } });
    if (output.binding) {
      const outputPair = cwlTool.createPair(output.abbrev, { type: output.type, outputBinding: { glob: output.binding } });
      mapOutputs.add(outputPair);
    } else {
      const extensions = output.extension.split("|");
      const binding = extensions.map(ext => "*" + ext).join("|");
      const outputPair = cwlTool.createPair(output.abbrev, { type: output.type, outputBinding: { glob: binding } });
      mapOutputs.add(outputPair);
    }
  }
  cwlTool.set("outputs", mapOutputs);

  const args = new YAMLSeq();
  for (const argument of version.arguments) {
    const argumentObj = {
      position: argument.position
    };
    if (argument.abbrev) argumentObj["prefix"] = argument.abbrev;

    if (argument.linked_params) {
      const inputBinding = version.inputs.find(input => input.abbrev === argument.linked_params);
      if (inputBinding) argumentObj["valueFrom"] = "$(inputs." + inputBinding.abbrev + ")";
      //   else {
      //     const outputBinding = version.outputs.find(output => output.name === argument.linked_params);
      //     if (outputBinding) argumentObj["valueFrom"] = "$(outputs." + outputBinding.abbrev + ")";
      //   }
      // } else if (argument.binding) {
      //   argumentObj["valueFrom"] = argument.binding;
    } else {
      argumentObj["valueFrom"] = argument.default_value
    }
    console.log("argument: ", argumentObj);
    const argumentPair = cwlTool.createNode(argumentObj);
    args.add(argumentPair)
  }
  cwlTool.set("arguments", args);

  cwlTool.set("baseCommand", tool.command);

  return cwlTool;
}

function generateCwlWorkflow(tools, toolLinks) {
  const cwlWorkflow = new Document();

  cwlWorkflow.set("id", "workflow");
  cwlWorkflow.set("class", "Workflow");
  cwlWorkflow.set("cwlVersion", "v1.2");

  const requirements = new YAMLSeq();
  requirements.add(cwlWorkflow.createNode({ class: "StepInputExpressionRequirement" }));
  cwlWorkflow.set("requirements", requirements);

  const externalInputs = tools.map(tool => {
    const inputs = tool.selected_version.inputs;
    const notLinkedInputs = inputs.filter(input => toolLinks.filter(link => {
      return link.to_input.id === input.id && link.to_tool_index === tool.tool_index;
    }).length === 0 && !input.hasOwnProperty("binding"));
    return notLinkedInputs.map(input => Object.assign({ tool_index: tool.tool_index, tool_name: tool.name }, input));
  });

  console.log("external inputs: ", externalInputs);

  const inputs = new YAMLMap();
  for (const input of externalInputs.flat(1)) {
    const inputPair = cwlWorkflow.createPair([input.tool_name, input.tool_index, input.abbrev].join("_"), { type: input.type });
    inputs.add(inputPair);
  }
  inputs.add(cwlWorkflow.createPair("sample_name", { type: "string" }));
  cwlWorkflow.set("inputs", inputs);

  const steps = new YAMLMap();
  for (const tool of tools) {
    const stepIn = new YAMLMap();
    for (const input of tool.selected_version.inputs) {
      if (input.binding) {
        const inputPair = cwlWorkflow.createPair(input.abbrev, { source: "sample_name", valueFrom: input.binding });
        stepIn.add(inputPair);
      } else {
        const links = toolLinks.filter(link => link.to_tool_index === tool.tool_index && link.to_input.id === input.id);
        let inputValue = [tool.name, tool.tool_index, input.abbrev].join("_");
        if (links.length > 0) {
          const link = links[0];
          inputValue = [link.from_tool.name, link.from_tool_index].join("_") + "/" + link.from_output.abbrev;
        }
        const inputPair = cwlWorkflow.createPair(input.abbrev, inputValue);
        stepIn.add(inputPair);
      }
    }

    const stepOut = new YAMLSeq();
    for (const output of tool.selected_version.outputs) {
      stepOut.add(cwlWorkflow.createNode(output.abbrev));
    }

    const toolPair = cwlWorkflow.createPair(tool.name + "_" + tool.tool_index, {
      id: [tool.name, tool.tool_index, tool.selected_version.semver].join("_"),
      in: stepIn,
      out: stepOut,
      run: "tasks/" + tool.name + "_" + tool.tool_index + ".cwl"
    });
    steps.add(toolPair);
  }
  cwlWorkflow.set("steps", steps);

  return cwlWorkflow;
}