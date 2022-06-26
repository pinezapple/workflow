const { Document, YAMLMap, YAMLSeq } = require("yaml");
const tool = {
  "tool_index": 0,
  "id": 6,
  "name": "BWA-MEM",
  "brief": "map medium and long reads (> 100 bp) against reference genome",
  "description": "BWA-MEM is an alignment algorithm for aligning sequence reads or long query sequences against a large reference genome such as human. It automatically chooses between local and end-to-end alignments, supports paired-end reads and performs chimeric alignment. The algorithm is robust to sequencing errors and applicable to a wide range of sequence lengths from 70bp to a few megabases. This tool wraps bwa-mem module of bwa read mapping tool. The Galaxy implementation takes fastq files as input and produces output in BAM format, which can be further processed using various BAM utilities exiting in Galaxy (BAMTools, SAMTools, Picard).",
  "command": "bwa mem",
  "category": "BAM/SAM",
  "default_version": {},
  "versions": [
    {
      "id": 3,
      "name": "BWA-MEM 0.7.17.1",
      "semver": "0.7.17.1",
      "description": "map medium and long reads (> 100 bp) against reference genome (Galaxy Version 0.7.17.1)",
      "image": "vinbdi/bwa-mem:0.7.17.1",
      "arguments": [
        {
          "id": 10,
          "abbrev": "-t",
          "description": "Number of threads",
          "position": 1
        },
        {
          "id": 11,
          "abbrev": "-P",
          "description": "In the paired-end mode, perform SW to rescue missing hits only but do not try to find hits that fit a proper pair.",
          "position": 2
        }
      ],
      "inputs": [
        {
          "id": 5,
          "name": "Reference genome",
          "abbrev": "indexed_reference_fasta",
          "description": "Reference genome",
          "type": "fasta",
          "extension": ".fasta|.fasta.gz"
        },
        {
          "id": 6,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "fasta",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 7,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "fasta",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 5,
          "name": "Aligned reads file",
          "abbrev": "aligned_read_file",
          "description": "Aligned reads file",
          "type": "sam",
          "extension": ".sam"
        }
      ],
      "tool_id": 6,
      "runtime_id": 1
    },
    {
      "id": 4,
      "name": "BWA-MEM 0.7.12.1",
      "semver": "0.7.12.1",
      "description": "map medium and long reads (> 100 bp) against reference genome (Galaxy Version 0.7.12.1)",
      "image": "vinbdi/bwa-mem:0.7.12.1",
      "arguments": [
        {
          "id": 12,
          "abbrev": "-t",
          "description": "Number of threads",
          "position": 1
        },
        {
          "id": 13,
          "abbrev": "-P",
          "description": "In the paired-end mode, perform SW to rescue missing hits only but do not try to find hits that fit a proper pair.",
          "position": 2
        }
      ],
      "inputs": [
        {
          "id": 8,
          "name": "Reference genome",
          "abbrev": "indexed_reference_fasta",
          "description": "Reference genome",
          "type": "fasta",
          "extension": ".fasta|.fasta.gz"
        },
        {
          "id": 9,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First reads",
          "type": "fasta",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 10,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second reads",
          "type": "fasta",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 6,
          "name": "Aligned reads file",
          "abbrev": "aligned_read_file",
          "description": "Aligned reads file",
          "type": "sam",
          "extension": ".sam"
        }
      ],
      "tool_id": 6,
      "runtime_id": 1
    },
    {
      "id": 5,
      "name": "BWA-MEM 0.4.2",
      "semver": "0.4.2",
      "description": "map medium and long reads (> 100 bp) against reference genome (Galaxy Version 0.4.2)",
      "image": "vinbdi/bwa-mem:0.4.2",
      "arguments": [
        {
          "id": 14,
          "abbrev": "-t",
          "description": "Number of threads",
          "position": 1
        },
        {
          "id": 15,
          "abbrev": "-P",
          "description": "In the paired-end mode, perform SW to rescue missing hits only but do not try to find hits that fit a proper pair.",
          "position": 2
        }
      ],
      "inputs": [
        {
          "id": 11,
          "name": "Reference genome",
          "abbrev": "indexed_reference_fasta",
          "description": "Reference genome",
          "type": "fasta",
          "extension": ".fasta|.fasta.gz"
        },
        {
          "id": 12,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "fasta",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 13,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "fasta",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 7,
          "name": "Aligned reads file",
          "abbrev": "aligned_read_file",
          "description": "Aligned reads file",
          "type": "sam",
          "extension": ".sam"
        }
      ],
      "tool_id": 6,
      "runtime_id": 1
    }
  ],
  "selected_version": {
    "id": 3,
    "name": "BWA-MEM 0.7.17.1",
    "semver": "0.7.17.1",
    "description": "map medium and long reads (> 100 bp) against reference genome (Galaxy Version 0.7.17.1)",
    "image": "vinbdi/bwa-mem:0.7.17.1",
    "arguments": [
      {
        "id": 10,
        "abbrev": "-t",
        "description": "Number of threads",
        "position": 1,
        "default_value": 4
      },
      {
        "id": 11,
        "abbrev": "-P",
        "description": "In the paired-end mode, perform SW to rescue missing hits only but do not try to find hits that fit a proper pair.",
        "position": 2,
        "default_value": ""
      }
    ],
    "inputs": [
      {
        "id": 5,
        "name": "Reference genome",
        "abbrev": "indexed_reference_fasta",
        "description": "Reference genome",
        "type": "fasta",
        "extension": ".fasta|.fasta.gz"
      },
      {
        "id": 6,
        "name": "First reads file",
        "abbrev": "input_fastq1",
        "description": "First read",
        "type": "fasta",
        "extension": ".fastq|.fastq.gz"
      },
      {
        "id": 7,
        "name": "Second reads file",
        "abbrev": "input_fastq2",
        "description": "Second read",
        "type": "fasta",
        "extension": ".fastq|.fastq.gz"
      }
    ],
    "outputs": [
      {
        "id": 5,
        "name": "Aligned reads file",
        "abbrev": "aligned_read_file",
        "description": "Aligned reads file",
        "type": "sam",
        "extension": ".sam"
      }
    ],
    "tool_id": 6,
    "runtime_id": 1
  }
}

const cwlTool = new Document()
cwlTool.set("class", "CommandLineTool");
cwlTool.set("cwlVersion", "v1.0");

const version = tool.selected_version;

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
  const outputPair = cwlTool.createPair(output.abbrev, { type: output.type, outputBinding: { glob: output.extension } });
  mapOutputs.add(outputPair);
}
cwlTool.set("outputs", mapOutputs);

const arguments = new YAMLSeq();
for (const argument of version.arguments) {
  const argumentPair = cwlTool.createNode({ prefix: argument.abbrev, position: argument.position, valueFrom: argument.default_value });
  arguments.add(argumentPair)
}
cwlTool.set("arguments", arguments);

cwlTool.set("baseCommand", tool.command);

console.log(cwlTool.toString());
