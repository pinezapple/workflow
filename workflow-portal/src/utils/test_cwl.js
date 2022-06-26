const { Document, YAMLMap, YAMLSeq, stringify } = require("yaml");

const toolLinks = [
  {
    "from_tool_index": 0,
    "from_tool": {
      "tool_index": 0,
      "id": 2,
      "name": "Fastp",
      "brief": "fast all-in-one preprocessing for FASTQ files",
      "description": "fastp is a tool designed to provide fast all-in-one preprocessing for FASTQ files. This tool is developed in C++ with multithreading supported to afford high performance.",
      "command": "fastp",
      "category": "FASTQ/FASTA",
      "default_version": {},
      "versions": [
        {
          "id": 2,
          "name": "Fastp 0.20.0",
          "semver": "0.20.0",
          "description": "Read Preprocessing Tool",
          "image": "vinbdi/fastp:0.20.0",
          "arguments": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "-i",
              "description": "First reads input file",
              "position": 1,
              "required": true,
              "linked_params": "First reads file"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "-I",
              "description": "Second reads input file",
              "position": 2,
              "linked_params": "Second reads file"
            },
            {
              "id": 5,
              "name": "First reads output file",
              "abbrev": "-o",
              "description": "First reads output file",
              "position": 3,
              "required": true,
              "linked_params": "First output reads file"
            },
            {
              "id": 6,
              "name": "Second reads output file",
              "abbrev": "-I",
              "description": "Second reads output file",
              "position": 4,
              "linked_params": "Second output reads file"
            },
            {
              "id": 7,
              "name": "Qualified filter",
              "abbrev": "-q",
              "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
              "position": 5
            },
            {
              "id": 8,
              "name": "Average qualified filter",
              "abbrev": "-e",
              "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
              "position": 6
            },
            {
              "id": 9,
              "name": "Unqualified filter",
              "abbrev": "-u",
              "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
              "position": 7
            }
          ],
          "inputs": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 3,
              "name": "First output reads file",
              "abbrev": "output_fastq1",
              "description": "First output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            },
            {
              "id": 4,
              "name": "Second output reads file",
              "abbrev": "output_fastq2",
              "description": "Second output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            }
          ],
          "tool_id": 2,
          "runtime_id": 1
        }
      ],
      "selected_version": {
        "id": 2,
        "name": "Fastp 0.20.0",
        "semver": "0.20.0",
        "description": "Read Preprocessing Tool",
        "image": "vinbdi/fastp:0.20.0",
        "arguments": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "-i",
            "description": "First reads input file",
            "position": 1,
            "required": true,
            "linked_params": "First reads file"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "-I",
            "description": "Second reads input file",
            "position": 2,
            "linked_params": "Second reads file"
          },
          {
            "id": 5,
            "name": "First reads output file",
            "abbrev": "-o",
            "description": "First reads output file",
            "position": 3,
            "required": true,
            "linked_params": "First output reads file"
          },
          {
            "id": 6,
            "name": "Second reads output file",
            "abbrev": "-I",
            "description": "Second reads output file",
            "position": 4,
            "linked_params": "Second output reads file"
          },
          {
            "id": 7,
            "name": "Qualified filter",
            "abbrev": "-q",
            "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
            "position": 5
          },
          {
            "id": 8,
            "name": "Average qualified filter",
            "abbrev": "-e",
            "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
            "position": 6
          },
          {
            "id": 9,
            "name": "Unqualified filter",
            "abbrev": "-u",
            "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
            "position": 7
          }
        ],
        "inputs": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 3,
            "name": "First output reads file",
            "abbrev": "output_fastq1",
            "description": "First output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          },
          {
            "id": 4,
            "name": "Second output reads file",
            "abbrev": "output_fastq2",
            "description": "Second output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          }
        ],
        "tool_id": 2,
        "runtime_id": 1
      }
    },
    "from_version": {
      "id": 2,
      "name": "Fastp 0.20.0",
      "semver": "0.20.0",
      "description": "Read Preprocessing Tool",
      "image": "vinbdi/fastp:0.20.0",
      "arguments": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "-i",
          "description": "First reads input file",
          "position": 1,
          "required": true,
          "linked_params": "First reads file"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "-I",
          "description": "Second reads input file",
          "position": 2,
          "linked_params": "Second reads file"
        },
        {
          "id": 5,
          "name": "First reads output file",
          "abbrev": "-o",
          "description": "First reads output file",
          "position": 3,
          "required": true,
          "linked_params": "First output reads file"
        },
        {
          "id": 6,
          "name": "Second reads output file",
          "abbrev": "-I",
          "description": "Second reads output file",
          "position": 4,
          "linked_params": "Second output reads file"
        },
        {
          "id": 7,
          "name": "Qualified filter",
          "abbrev": "-q",
          "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
          "position": 5
        },
        {
          "id": 8,
          "name": "Average qualified filter",
          "abbrev": "-e",
          "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
          "position": 6
        },
        {
          "id": 9,
          "name": "Unqualified filter",
          "abbrev": "-u",
          "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
          "position": 7
        }
      ],
      "inputs": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 3,
          "name": "First output reads file",
          "abbrev": "output_fastq1",
          "description": "First output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        },
        {
          "id": 4,
          "name": "Second output reads file",
          "abbrev": "output_fastq2",
          "description": "Second output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        }
      ],
      "tool_id": 2,
      "runtime_id": 1
    },
    "from_output": {
      "id": 3,
      "name": "First output reads file",
      "abbrev": "output_fastq1",
      "description": "First output reads file",
      "type": "File",
      "extension": ".fastq.gz"
    },
    "to_tool_index": 1,
    "to_tool": {
      "tool_index": 1,
      "id": 1,
      "name": "FastQC",
      "brief": "Read Quality reports",
      "description": "FastQC aims to provide a simple way to do some quality control checks on raw sequence data coming from high throughput sequencing pipelines. It provides a set of analyses which you can use to get a quick impression of whether your data has any problems of which you should be aware before doing any further analysis.",
      "command": "fastqc",
      "category": "FASTQ/FASTA",
      "default_version": {},
      "versions": [
        {
          "id": 1,
          "name": "FastQC 0.72",
          "semver": "0.72",
          "description": "Read Quality Reports",
          "image": "vinbdi/fastqc:0.72",
          "arguments": [
            {
              "id": 1,
              "name": "Adapters",
              "abbrev": "-a",
              "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
              "position": 1
            },
            {
              "id": 2,
              "name": "Kmers",
              "abbrev": "-k",
              "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
              "position": 2
            }
          ],
          "inputs": [
            {
              "id": 1,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 2,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 1,
              "name": "FastQC report file",
              "abbrev": "FastQC_report_html",
              "description": "FastQC report file html",
              "type": "File",
              "extension": ".html"
            },
            {
              "id": 2,
              "name": "FastQC report detail",
              "abbrev": "FastQC_report_zip",
              "description": "FastQC report file zip",
              "type": "File",
              "extension": ".zip"
            }
          ],
          "tool_id": 1,
          "runtime_id": 1
        }
      ],
      "selected_version": {
        "id": 1,
        "name": "FastQC 0.72",
        "semver": "0.72",
        "description": "Read Quality Reports",
        "image": "vinbdi/fastqc:0.72",
        "arguments": [
          {
            "id": 1,
            "name": "Adapters",
            "abbrev": "-a",
            "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
            "position": 1
          },
          {
            "id": 2,
            "name": "Kmers",
            "abbrev": "-k",
            "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
            "position": 2
          }
        ],
        "inputs": [
          {
            "id": 1,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 2,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 1,
            "name": "FastQC report file",
            "abbrev": "FastQC_report_html",
            "description": "FastQC report file html",
            "type": "File",
            "extension": ".html"
          },
          {
            "id": 2,
            "name": "FastQC report detail",
            "abbrev": "FastQC_report_zip",
            "description": "FastQC report file zip",
            "type": "File",
            "extension": ".zip"
          }
        ],
        "tool_id": 1,
        "runtime_id": 1
      }
    },
    "to_version": {
      "id": 1,
      "name": "FastQC 0.72",
      "semver": "0.72",
      "description": "Read Quality Reports",
      "image": "vinbdi/fastqc:0.72",
      "arguments": [
        {
          "id": 1,
          "name": "Adapters",
          "abbrev": "-a",
          "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
          "position": 1
        },
        {
          "id": 2,
          "name": "Kmers",
          "abbrev": "-k",
          "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
          "position": 2
        }
      ],
      "inputs": [
        {
          "id": 1,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 2,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 1,
          "name": "FastQC report file",
          "abbrev": "FastQC_report_html",
          "description": "FastQC report file html",
          "type": "File",
          "extension": ".html"
        },
        {
          "id": 2,
          "name": "FastQC report detail",
          "abbrev": "FastQC_report_zip",
          "description": "FastQC report file zip",
          "type": "File",
          "extension": ".zip"
        }
      ],
      "tool_id": 1,
      "runtime_id": 1
    },
    "to_input": {
      "id": 1,
      "name": "First reads file",
      "abbrev": "input_fastq1",
      "description": "First read",
      "type": "File",
      "extension": ".fastq|.fastq.gz"
    }
  },
  {
    "from_tool_index": 0,
    "from_tool": {
      "tool_index": 0,
      "id": 2,
      "name": "Fastp",
      "brief": "fast all-in-one preprocessing for FASTQ files",
      "description": "fastp is a tool designed to provide fast all-in-one preprocessing for FASTQ files. This tool is developed in C++ with multithreading supported to afford high performance.",
      "command": "fastp",
      "category": "FASTQ/FASTA",
      "default_version": {},
      "versions": [
        {
          "id": 2,
          "name": "Fastp 0.20.0",
          "semver": "0.20.0",
          "description": "Read Preprocessing Tool",
          "image": "vinbdi/fastp:0.20.0",
          "arguments": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "-i",
              "description": "First reads input file",
              "position": 1,
              "required": true,
              "linked_params": "First reads file"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "-I",
              "description": "Second reads input file",
              "position": 2,
              "linked_params": "Second reads file"
            },
            {
              "id": 5,
              "name": "First reads output file",
              "abbrev": "-o",
              "description": "First reads output file",
              "position": 3,
              "required": true,
              "linked_params": "First output reads file"
            },
            {
              "id": 6,
              "name": "Second reads output file",
              "abbrev": "-I",
              "description": "Second reads output file",
              "position": 4,
              "linked_params": "Second output reads file"
            },
            {
              "id": 7,
              "name": "Qualified filter",
              "abbrev": "-q",
              "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
              "position": 5
            },
            {
              "id": 8,
              "name": "Average qualified filter",
              "abbrev": "-e",
              "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
              "position": 6
            },
            {
              "id": 9,
              "name": "Unqualified filter",
              "abbrev": "-u",
              "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
              "position": 7
            }
          ],
          "inputs": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 3,
              "name": "First output reads file",
              "abbrev": "output_fastq1",
              "description": "First output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            },
            {
              "id": 4,
              "name": "Second output reads file",
              "abbrev": "output_fastq2",
              "description": "Second output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            }
          ],
          "tool_id": 2,
          "runtime_id": 1
        }
      ],
      "selected_version": {
        "id": 2,
        "name": "Fastp 0.20.0",
        "semver": "0.20.0",
        "description": "Read Preprocessing Tool",
        "image": "vinbdi/fastp:0.20.0",
        "arguments": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "-i",
            "description": "First reads input file",
            "position": 1,
            "required": true,
            "linked_params": "First reads file"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "-I",
            "description": "Second reads input file",
            "position": 2,
            "linked_params": "Second reads file"
          },
          {
            "id": 5,
            "name": "First reads output file",
            "abbrev": "-o",
            "description": "First reads output file",
            "position": 3,
            "required": true,
            "linked_params": "First output reads file"
          },
          {
            "id": 6,
            "name": "Second reads output file",
            "abbrev": "-I",
            "description": "Second reads output file",
            "position": 4,
            "linked_params": "Second output reads file"
          },
          {
            "id": 7,
            "name": "Qualified filter",
            "abbrev": "-q",
            "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
            "position": 5
          },
          {
            "id": 8,
            "name": "Average qualified filter",
            "abbrev": "-e",
            "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
            "position": 6
          },
          {
            "id": 9,
            "name": "Unqualified filter",
            "abbrev": "-u",
            "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
            "position": 7
          }
        ],
        "inputs": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 3,
            "name": "First output reads file",
            "abbrev": "output_fastq1",
            "description": "First output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          },
          {
            "id": 4,
            "name": "Second output reads file",
            "abbrev": "output_fastq2",
            "description": "Second output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          }
        ],
        "tool_id": 2,
        "runtime_id": 1
      }
    },
    "from_version": {
      "id": 2,
      "name": "Fastp 0.20.0",
      "semver": "0.20.0",
      "description": "Read Preprocessing Tool",
      "image": "vinbdi/fastp:0.20.0",
      "arguments": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "-i",
          "description": "First reads input file",
          "position": 1,
          "required": true,
          "linked_params": "First reads file"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "-I",
          "description": "Second reads input file",
          "position": 2,
          "linked_params": "Second reads file"
        },
        {
          "id": 5,
          "name": "First reads output file",
          "abbrev": "-o",
          "description": "First reads output file",
          "position": 3,
          "required": true,
          "linked_params": "First output reads file"
        },
        {
          "id": 6,
          "name": "Second reads output file",
          "abbrev": "-I",
          "description": "Second reads output file",
          "position": 4,
          "linked_params": "Second output reads file"
        },
        {
          "id": 7,
          "name": "Qualified filter",
          "abbrev": "-q",
          "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
          "position": 5
        },
        {
          "id": 8,
          "name": "Average qualified filter",
          "abbrev": "-e",
          "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
          "position": 6
        },
        {
          "id": 9,
          "name": "Unqualified filter",
          "abbrev": "-u",
          "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
          "position": 7
        }
      ],
      "inputs": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 3,
          "name": "First output reads file",
          "abbrev": "output_fastq1",
          "description": "First output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        },
        {
          "id": 4,
          "name": "Second output reads file",
          "abbrev": "output_fastq2",
          "description": "Second output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        }
      ],
      "tool_id": 2,
      "runtime_id": 1
    },
    "from_output": {
      "id": 4,
      "name": "Second output reads file",
      "abbrev": "output_fastq2",
      "description": "Second output reads file",
      "type": "File",
      "extension": ".fastq.gz"
    },
    "to_tool_index": 1,
    "to_tool": {
      "tool_index": 1,
      "id": 1,
      "name": "FastQC",
      "brief": "Read Quality reports",
      "description": "FastQC aims to provide a simple way to do some quality control checks on raw sequence data coming from high throughput sequencing pipelines. It provides a set of analyses which you can use to get a quick impression of whether your data has any problems of which you should be aware before doing any further analysis.",
      "command": "fastqc",
      "category": "FASTQ/FASTA",
      "default_version": {},
      "versions": [
        {
          "id": 1,
          "name": "FastQC 0.72",
          "semver": "0.72",
          "description": "Read Quality Reports",
          "image": "vinbdi/fastqc:0.72",
          "arguments": [
            {
              "id": 1,
              "name": "Adapters",
              "abbrev": "-a",
              "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
              "position": 1
            },
            {
              "id": 2,
              "name": "Kmers",
              "abbrev": "-k",
              "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
              "position": 2
            }
          ],
          "inputs": [
            {
              "id": 1,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 2,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 1,
              "name": "FastQC report file",
              "abbrev": "FastQC_report_html",
              "description": "FastQC report file html",
              "type": "File",
              "extension": ".html"
            },
            {
              "id": 2,
              "name": "FastQC report detail",
              "abbrev": "FastQC_report_zip",
              "description": "FastQC report file zip",
              "type": "File",
              "extension": ".zip"
            }
          ],
          "tool_id": 1,
          "runtime_id": 1
        }
      ],
      "selected_version": {
        "id": 1,
        "name": "FastQC 0.72",
        "semver": "0.72",
        "description": "Read Quality Reports",
        "image": "vinbdi/fastqc:0.72",
        "arguments": [
          {
            "id": 1,
            "name": "Adapters",
            "abbrev": "-a",
            "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
            "position": 1
          },
          {
            "id": 2,
            "name": "Kmers",
            "abbrev": "-k",
            "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
            "position": 2
          }
        ],
        "inputs": [
          {
            "id": 1,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 2,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 1,
            "name": "FastQC report file",
            "abbrev": "FastQC_report_html",
            "description": "FastQC report file html",
            "type": "File",
            "extension": ".html"
          },
          {
            "id": 2,
            "name": "FastQC report detail",
            "abbrev": "FastQC_report_zip",
            "description": "FastQC report file zip",
            "type": "File",
            "extension": ".zip"
          }
        ],
        "tool_id": 1,
        "runtime_id": 1
      }
    },
    "to_version": {
      "id": 1,
      "name": "FastQC 0.72",
      "semver": "0.72",
      "description": "Read Quality Reports",
      "image": "vinbdi/fastqc:0.72",
      "arguments": [
        {
          "id": 1,
          "name": "Adapters",
          "abbrev": "-a",
          "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
          "position": 1
        },
        {
          "id": 2,
          "name": "Kmers",
          "abbrev": "-k",
          "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
          "position": 2
        }
      ],
      "inputs": [
        {
          "id": 1,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 2,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 1,
          "name": "FastQC report file",
          "abbrev": "FastQC_report_html",
          "description": "FastQC report file html",
          "type": "File",
          "extension": ".html"
        },
        {
          "id": 2,
          "name": "FastQC report detail",
          "abbrev": "FastQC_report_zip",
          "description": "FastQC report file zip",
          "type": "File",
          "extension": ".zip"
        }
      ],
      "tool_id": 1,
      "runtime_id": 1
    },
    "to_input": {
      "id": 2,
      "name": "Second reads file",
      "abbrev": "input_fastq2",
      "description": "Second read",
      "type": "File",
      "extension": ".fastq|.fastq.gz"
    }
  },
  {
    "from_tool_index": 0,
    "from_tool": {
      "tool_index": 0,
      "id": 2,
      "name": "Fastp",
      "brief": "fast all-in-one preprocessing for FASTQ files",
      "description": "fastp is a tool designed to provide fast all-in-one preprocessing for FASTQ files. This tool is developed in C++ with multithreading supported to afford high performance.",
      "command": "fastp",
      "category": "FASTQ/FASTA",
      "default_version": {},
      "versions": [
        {
          "id": 2,
          "name": "Fastp 0.20.0",
          "semver": "0.20.0",
          "description": "Read Preprocessing Tool",
          "image": "vinbdi/fastp:0.20.0",
          "arguments": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "-i",
              "description": "First reads input file",
              "position": 1,
              "required": true,
              "linked_params": "First reads file"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "-I",
              "description": "Second reads input file",
              "position": 2,
              "linked_params": "Second reads file"
            },
            {
              "id": 5,
              "name": "First reads output file",
              "abbrev": "-o",
              "description": "First reads output file",
              "position": 3,
              "required": true,
              "linked_params": "First output reads file"
            },
            {
              "id": 6,
              "name": "Second reads output file",
              "abbrev": "-I",
              "description": "Second reads output file",
              "position": 4,
              "linked_params": "Second output reads file"
            },
            {
              "id": 7,
              "name": "Qualified filter",
              "abbrev": "-q",
              "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
              "position": 5
            },
            {
              "id": 8,
              "name": "Average qualified filter",
              "abbrev": "-e",
              "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
              "position": 6
            },
            {
              "id": 9,
              "name": "Unqualified filter",
              "abbrev": "-u",
              "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
              "position": 7
            }
          ],
          "inputs": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 3,
              "name": "First output reads file",
              "abbrev": "output_fastq1",
              "description": "First output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            },
            {
              "id": 4,
              "name": "Second output reads file",
              "abbrev": "output_fastq2",
              "description": "Second output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            }
          ],
          "tool_id": 2,
          "runtime_id": 1
        }
      ],
      "selected_version": {
        "id": 2,
        "name": "Fastp 0.20.0",
        "semver": "0.20.0",
        "description": "Read Preprocessing Tool",
        "image": "vinbdi/fastp:0.20.0",
        "arguments": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "-i",
            "description": "First reads input file",
            "position": 1,
            "required": true,
            "linked_params": "First reads file"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "-I",
            "description": "Second reads input file",
            "position": 2,
            "linked_params": "Second reads file"
          },
          {
            "id": 5,
            "name": "First reads output file",
            "abbrev": "-o",
            "description": "First reads output file",
            "position": 3,
            "required": true,
            "linked_params": "First output reads file"
          },
          {
            "id": 6,
            "name": "Second reads output file",
            "abbrev": "-I",
            "description": "Second reads output file",
            "position": 4,
            "linked_params": "Second output reads file"
          },
          {
            "id": 7,
            "name": "Qualified filter",
            "abbrev": "-q",
            "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
            "position": 5
          },
          {
            "id": 8,
            "name": "Average qualified filter",
            "abbrev": "-e",
            "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
            "position": 6
          },
          {
            "id": 9,
            "name": "Unqualified filter",
            "abbrev": "-u",
            "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
            "position": 7
          }
        ],
        "inputs": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 3,
            "name": "First output reads file",
            "abbrev": "output_fastq1",
            "description": "First output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          },
          {
            "id": 4,
            "name": "Second output reads file",
            "abbrev": "output_fastq2",
            "description": "Second output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          }
        ],
        "tool_id": 2,
        "runtime_id": 1
      }
    },
    "from_version": {
      "id": 2,
      "name": "Fastp 0.20.0",
      "semver": "0.20.0",
      "description": "Read Preprocessing Tool",
      "image": "vinbdi/fastp:0.20.0",
      "arguments": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "-i",
          "description": "First reads input file",
          "position": 1,
          "required": true,
          "linked_params": "First reads file"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "-I",
          "description": "Second reads input file",
          "position": 2,
          "linked_params": "Second reads file"
        },
        {
          "id": 5,
          "name": "First reads output file",
          "abbrev": "-o",
          "description": "First reads output file",
          "position": 3,
          "required": true,
          "linked_params": "First output reads file"
        },
        {
          "id": 6,
          "name": "Second reads output file",
          "abbrev": "-I",
          "description": "Second reads output file",
          "position": 4,
          "linked_params": "Second output reads file"
        },
        {
          "id": 7,
          "name": "Qualified filter",
          "abbrev": "-q",
          "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
          "position": 5
        },
        {
          "id": 8,
          "name": "Average qualified filter",
          "abbrev": "-e",
          "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
          "position": 6
        },
        {
          "id": 9,
          "name": "Unqualified filter",
          "abbrev": "-u",
          "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
          "position": 7
        }
      ],
      "inputs": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 3,
          "name": "First output reads file",
          "abbrev": "output_fastq1",
          "description": "First output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        },
        {
          "id": 4,
          "name": "Second output reads file",
          "abbrev": "output_fastq2",
          "description": "Second output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        }
      ],
      "tool_id": 2,
      "runtime_id": 1
    },
    "from_output": {
      "id": 3,
      "name": "First output reads file",
      "abbrev": "output_fastq1",
      "description": "First output reads file",
      "type": "File",
      "extension": ".fastq.gz"
    },
    "to_tool_index": 2,
    "to_tool": {
      "tool_index": 2,
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
              "type": "File",
              "extension": ".fasta|.fasta.gz"
            },
            {
              "id": 6,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 7,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 5,
              "name": "Aligned reads file",
              "abbrev": "aligned_read_file",
              "description": "Aligned reads file",
              "type": "File",
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
              "type": "File",
              "extension": ".fasta|.fasta.gz"
            },
            {
              "id": 9,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First reads",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 10,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second reads",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 6,
              "name": "Aligned reads file",
              "abbrev": "aligned_read_file",
              "description": "Aligned reads file",
              "type": "File",
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
              "type": "File",
              "extension": ".fasta|.fasta.gz"
            },
            {
              "id": 12,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 13,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 7,
              "name": "Aligned reads file",
              "abbrev": "aligned_read_file",
              "description": "Aligned reads file",
              "type": "File",
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
            "type": "File",
            "extension": ".fasta|.fasta.gz"
          },
          {
            "id": 6,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 7,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 5,
            "name": "Aligned reads file",
            "abbrev": "aligned_read_file",
            "description": "Aligned reads file",
            "type": "File",
            "extension": ".sam"
          }
        ],
        "tool_id": 6,
        "runtime_id": 1
      }
    },
    "to_version": {
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
          "type": "File",
          "extension": ".fasta|.fasta.gz"
        },
        {
          "id": 6,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 7,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 5,
          "name": "Aligned reads file",
          "abbrev": "aligned_read_file",
          "description": "Aligned reads file",
          "type": "File",
          "extension": ".sam"
        }
      ],
      "tool_id": 6,
      "runtime_id": 1
    },
    "to_input": {
      "id": 6,
      "name": "First reads file",
      "abbrev": "input_fastq1",
      "description": "First read",
      "type": "File",
      "extension": ".fastq|.fastq.gz"
    }
  },
  {
    "from_tool_index": 0,
    "from_tool": {
      "tool_index": 0,
      "id": 2,
      "name": "Fastp",
      "brief": "fast all-in-one preprocessing for FASTQ files",
      "description": "fastp is a tool designed to provide fast all-in-one preprocessing for FASTQ files. This tool is developed in C++ with multithreading supported to afford high performance.",
      "command": "fastp",
      "category": "FASTQ/FASTA",
      "default_version": {},
      "versions": [
        {
          "id": 2,
          "name": "Fastp 0.20.0",
          "semver": "0.20.0",
          "description": "Read Preprocessing Tool",
          "image": "vinbdi/fastp:0.20.0",
          "arguments": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "-i",
              "description": "First reads input file",
              "position": 1,
              "required": true,
              "linked_params": "First reads file"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "-I",
              "description": "Second reads input file",
              "position": 2,
              "linked_params": "Second reads file"
            },
            {
              "id": 5,
              "name": "First reads output file",
              "abbrev": "-o",
              "description": "First reads output file",
              "position": 3,
              "required": true,
              "linked_params": "First output reads file"
            },
            {
              "id": 6,
              "name": "Second reads output file",
              "abbrev": "-I",
              "description": "Second reads output file",
              "position": 4,
              "linked_params": "Second output reads file"
            },
            {
              "id": 7,
              "name": "Qualified filter",
              "abbrev": "-q",
              "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
              "position": 5
            },
            {
              "id": 8,
              "name": "Average qualified filter",
              "abbrev": "-e",
              "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
              "position": 6
            },
            {
              "id": 9,
              "name": "Unqualified filter",
              "abbrev": "-u",
              "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
              "position": 7
            }
          ],
          "inputs": [
            {
              "id": 3,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 4,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 3,
              "name": "First output reads file",
              "abbrev": "output_fastq1",
              "description": "First output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            },
            {
              "id": 4,
              "name": "Second output reads file",
              "abbrev": "output_fastq2",
              "description": "Second output reads file",
              "type": "File",
              "extension": ".fastq.gz"
            }
          ],
          "tool_id": 2,
          "runtime_id": 1
        }
      ],
      "selected_version": {
        "id": 2,
        "name": "Fastp 0.20.0",
        "semver": "0.20.0",
        "description": "Read Preprocessing Tool",
        "image": "vinbdi/fastp:0.20.0",
        "arguments": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "-i",
            "description": "First reads input file",
            "position": 1,
            "required": true,
            "linked_params": "First reads file"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "-I",
            "description": "Second reads input file",
            "position": 2,
            "linked_params": "Second reads file"
          },
          {
            "id": 5,
            "name": "First reads output file",
            "abbrev": "-o",
            "description": "First reads output file",
            "position": 3,
            "required": true,
            "linked_params": "First output reads file"
          },
          {
            "id": 6,
            "name": "Second reads output file",
            "abbrev": "-I",
            "description": "Second reads output file",
            "position": 4,
            "linked_params": "Second output reads file"
          },
          {
            "id": 7,
            "name": "Qualified filter",
            "abbrev": "-q",
            "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
            "position": 5
          },
          {
            "id": 8,
            "name": "Average qualified filter",
            "abbrev": "-e",
            "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
            "position": 6
          },
          {
            "id": 9,
            "name": "Unqualified filter",
            "abbrev": "-u",
            "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
            "position": 7
          }
        ],
        "inputs": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 3,
            "name": "First output reads file",
            "abbrev": "output_fastq1",
            "description": "First output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          },
          {
            "id": 4,
            "name": "Second output reads file",
            "abbrev": "output_fastq2",
            "description": "Second output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          }
        ],
        "tool_id": 2,
        "runtime_id": 1
      }
    },
    "from_version": {
      "id": 2,
      "name": "Fastp 0.20.0",
      "semver": "0.20.0",
      "description": "Read Preprocessing Tool",
      "image": "vinbdi/fastp:0.20.0",
      "arguments": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "-i",
          "description": "First reads input file",
          "position": 1,
          "required": true,
          "linked_params": "First reads file"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "-I",
          "description": "Second reads input file",
          "position": 2,
          "linked_params": "Second reads file"
        },
        {
          "id": 5,
          "name": "First reads output file",
          "abbrev": "-o",
          "description": "First reads output file",
          "position": 3,
          "required": true,
          "linked_params": "First output reads file"
        },
        {
          "id": 6,
          "name": "Second reads output file",
          "abbrev": "-I",
          "description": "Second reads output file",
          "position": 4,
          "linked_params": "Second output reads file"
        },
        {
          "id": 7,
          "name": "Qualified filter",
          "abbrev": "-q",
          "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
          "position": 5
        },
        {
          "id": 8,
          "name": "Average qualified filter",
          "abbrev": "-e",
          "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
          "position": 6
        },
        {
          "id": 9,
          "name": "Unqualified filter",
          "abbrev": "-u",
          "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
          "position": 7
        }
      ],
      "inputs": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 3,
          "name": "First output reads file",
          "abbrev": "output_fastq1",
          "description": "First output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        },
        {
          "id": 4,
          "name": "Second output reads file",
          "abbrev": "output_fastq2",
          "description": "Second output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        }
      ],
      "tool_id": 2,
      "runtime_id": 1
    },
    "from_output": {
      "id": 4,
      "name": "Second output reads file",
      "abbrev": "output_fastq2",
      "description": "Second output reads file",
      "type": "File",
      "extension": ".fastq.gz"
    },
    "to_tool_index": 2,
    "to_tool": {
      "tool_index": 2,
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
              "type": "File",
              "extension": ".fasta|.fasta.gz"
            },
            {
              "id": 6,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 7,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 5,
              "name": "Aligned reads file",
              "abbrev": "aligned_read_file",
              "description": "Aligned reads file",
              "type": "File",
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
              "type": "File",
              "extension": ".fasta|.fasta.gz"
            },
            {
              "id": 9,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First reads",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 10,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second reads",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 6,
              "name": "Aligned reads file",
              "abbrev": "aligned_read_file",
              "description": "Aligned reads file",
              "type": "File",
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
              "type": "File",
              "extension": ".fasta|.fasta.gz"
            },
            {
              "id": 12,
              "name": "First reads file",
              "abbrev": "input_fastq1",
              "description": "First read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            },
            {
              "id": 13,
              "name": "Second reads file",
              "abbrev": "input_fastq2",
              "description": "Second read",
              "type": "File",
              "extension": ".fastq|.fastq.gz"
            }
          ],
          "outputs": [
            {
              "id": 7,
              "name": "Aligned reads file",
              "abbrev": "aligned_read_file",
              "description": "Aligned reads file",
              "type": "File",
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
            "type": "File",
            "extension": ".fasta|.fasta.gz"
          },
          {
            "id": 6,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 7,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 5,
            "name": "Aligned reads file",
            "abbrev": "aligned_read_file",
            "description": "Aligned reads file",
            "type": "File",
            "extension": ".sam"
          }
        ],
        "tool_id": 6,
        "runtime_id": 1
      }
    },
    "to_version": {
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
          "type": "File",
          "extension": ".fasta|.fasta.gz"
        },
        {
          "id": 6,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 7,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 5,
          "name": "Aligned reads file",
          "abbrev": "aligned_read_file",
          "description": "Aligned reads file",
          "type": "File",
          "extension": ".sam"
        }
      ],
      "tool_id": 6,
      "runtime_id": 1
    },
    "to_input": {
      "id": 7,
      "name": "Second reads file",
      "abbrev": "input_fastq2",
      "description": "Second read",
      "type": "File",
      "extension": ".fastq|.fastq.gz"
    }
  }
]

const tools = [
  {
    "tool_index": 0,
    "id": 2,
    "name": "Fastp",
    "brief": "fast all-in-one preprocessing for FASTQ files",
    "description": "fastp is a tool designed to provide fast all-in-one preprocessing for FASTQ files. This tool is developed in C++ with multithreading supported to afford high performance.",
    "command": "fastp",
    "category": "FASTQ/FASTA",
    "default_version": {},
    "versions": [
      {
        "id": 2,
        "name": "Fastp 0.20.0",
        "semver": "0.20.0",
        "description": "Read Preprocessing Tool",
        "image": "vinbdi/fastp:0.20.0",
        "arguments": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "-i",
            "description": "First reads input file",
            "position": 1,
            "required": true,
            "linked_params": "First reads file"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "-I",
            "description": "Second reads input file",
            "position": 2,
            "linked_params": "Second reads file"
          },
          {
            "id": 5,
            "name": "First reads output file",
            "abbrev": "-o",
            "description": "First reads output file",
            "position": 3,
            "required": true,
            "linked_params": "First output reads file"
          },
          {
            "id": 6,
            "name": "Second reads output file",
            "abbrev": "-I",
            "description": "Second reads output file",
            "position": 4,
            "linked_params": "Second output reads file"
          },
          {
            "id": 7,
            "name": "Qualified filter",
            "abbrev": "-q",
            "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
            "position": 5
          },
          {
            "id": 8,
            "name": "Average qualified filter",
            "abbrev": "-e",
            "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
            "position": 6
          },
          {
            "id": 9,
            "name": "Unqualified filter",
            "abbrev": "-u",
            "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
            "position": 7
          }
        ],
        "inputs": [
          {
            "id": 3,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 4,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 3,
            "name": "First output reads file",
            "abbrev": "output_fastq1",
            "description": "First output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          },
          {
            "id": 4,
            "name": "Second output reads file",
            "abbrev": "output_fastq2",
            "description": "Second output reads file",
            "type": "File",
            "extension": ".fastq.gz"
          }
        ],
        "tool_id": 2,
        "runtime_id": 1
      }
    ],
    "selected_version": {
      "id": 2,
      "name": "Fastp 0.20.0",
      "semver": "0.20.0",
      "description": "Read Preprocessing Tool",
      "image": "vinbdi/fastp:0.20.0",
      "arguments": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "-i",
          "description": "First reads input file",
          "position": 1,
          "required": true,
          "linked_params": "First reads file"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "-I",
          "description": "Second reads input file",
          "position": 2,
          "linked_params": "Second reads file"
        },
        {
          "id": 5,
          "name": "First reads output file",
          "abbrev": "-o",
          "description": "First reads output file",
          "position": 3,
          "required": true,
          "linked_params": "First output reads file"
        },
        {
          "id": 6,
          "name": "Second reads output file",
          "abbrev": "-I",
          "description": "Second reads output file",
          "position": 4,
          "linked_params": "Second output reads file"
        },
        {
          "id": 7,
          "name": "Qualified filter",
          "abbrev": "-q",
          "description": "the quality value that a base is qualified. Default 15 means phred quality >=Q15 is qualified.",
          "position": 5
        },
        {
          "id": 8,
          "name": "Average qualified filter",
          "abbrev": "-e",
          "description": " if one read's average quality score <avg_qual, then this read/pair is discarded. Default 0 means no requirement (int [=0])",
          "position": 6
        },
        {
          "id": 9,
          "name": "Unqualified filter",
          "abbrev": "-u",
          "description": "how many percents of bases are allowed to be unqualified (0~100). Default 40 means 40%",
          "position": 7
        }
      ],
      "inputs": [
        {
          "id": 3,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 4,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 3,
          "name": "First output reads file",
          "abbrev": "output_fastq1",
          "description": "First output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        },
        {
          "id": 4,
          "name": "Second output reads file",
          "abbrev": "output_fastq2",
          "description": "Second output reads file",
          "type": "File",
          "extension": ".fastq.gz"
        }
      ],
      "tool_id": 2,
      "runtime_id": 1
    }
  },
  {
    "tool_index": 1,
    "id": 1,
    "name": "FastQC",
    "brief": "Read Quality reports",
    "description": "FastQC aims to provide a simple way to do some quality control checks on raw sequence data coming from high throughput sequencing pipelines. It provides a set of analyses which you can use to get a quick impression of whether your data has any problems of which you should be aware before doing any further analysis.",
    "command": "fastqc",
    "category": "FASTQ/FASTA",
    "default_version": {},
    "versions": [
      {
        "id": 1,
        "name": "FastQC 0.72",
        "semver": "0.72",
        "description": "Read Quality Reports",
        "image": "vinbdi/fastqc:0.72",
        "arguments": [
          {
            "id": 1,
            "name": "Adapters",
            "abbrev": "-a",
            "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
            "position": 1
          },
          {
            "id": 2,
            "name": "Kmers",
            "abbrev": "-k",
            "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
            "position": 2
          }
        ],
        "inputs": [
          {
            "id": 1,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 2,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 1,
            "name": "FastQC report file",
            "abbrev": "FastQC_report_html",
            "description": "FastQC report file html",
            "type": "File",
            "extension": ".html"
          },
          {
            "id": 2,
            "name": "FastQC report detail",
            "abbrev": "FastQC_report_zip",
            "description": "FastQC report file zip",
            "type": "File",
            "extension": ".zip"
          }
        ],
        "tool_id": 1,
        "runtime_id": 1
      }
    ],
    "selected_version": {
      "id": 1,
      "name": "FastQC 0.72",
      "semver": "0.72",
      "description": "Read Quality Reports",
      "image": "vinbdi/fastqc:0.72",
      "arguments": [
        {
          "id": 1,
          "name": "Adapters",
          "abbrev": "-a",
          "description": "Specifies a non-default file which contains the list of adapter sequences which will be explicity searched against the library. The file must contain sets of named adapters in the form name[tab]sequence.  Lines prefixed with a hash will be ignored.",
          "position": 1
        },
        {
          "id": 2,
          "name": "Kmers",
          "abbrev": "-k",
          "description": "Specifies the length of Kmer to look for in the Kmer content module. Specified Kmer length must be between 2 and 10. Default length is 7 if not specified.",
          "position": 2
        }
      ],
      "inputs": [
        {
          "id": 1,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 2,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 1,
          "name": "FastQC report file",
          "abbrev": "FastQC_report_html",
          "description": "FastQC report file html",
          "type": "File",
          "extension": ".html"
        },
        {
          "id": 2,
          "name": "FastQC report detail",
          "abbrev": "FastQC_report_zip",
          "description": "FastQC report file zip",
          "type": "File",
          "extension": ".zip"
        }
      ],
      "tool_id": 1,
      "runtime_id": 1
    }
  },
  {
    "tool_index": 2,
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
            "type": "File",
            "extension": ".fasta|.fasta.gz"
          },
          {
            "id": 6,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 7,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 5,
            "name": "Aligned reads file",
            "abbrev": "aligned_read_file",
            "description": "Aligned reads file",
            "type": "File",
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
            "type": "File",
            "extension": ".fasta|.fasta.gz"
          },
          {
            "id": 9,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First reads",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 10,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second reads",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 6,
            "name": "Aligned reads file",
            "abbrev": "aligned_read_file",
            "description": "Aligned reads file",
            "type": "File",
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
            "type": "File",
            "extension": ".fasta|.fasta.gz"
          },
          {
            "id": 12,
            "name": "First reads file",
            "abbrev": "input_fastq1",
            "description": "First read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          },
          {
            "id": 13,
            "name": "Second reads file",
            "abbrev": "input_fastq2",
            "description": "Second read",
            "type": "File",
            "extension": ".fastq|.fastq.gz"
          }
        ],
        "outputs": [
          {
            "id": 7,
            "name": "Aligned reads file",
            "abbrev": "aligned_read_file",
            "description": "Aligned reads file",
            "type": "File",
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
          "type": "File",
          "extension": ".fasta|.fasta.gz"
        },
        {
          "id": 6,
          "name": "First reads file",
          "abbrev": "input_fastq1",
          "description": "First read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        },
        {
          "id": 7,
          "name": "Second reads file",
          "abbrev": "input_fastq2",
          "description": "Second read",
          "type": "File",
          "extension": ".fastq|.fastq.gz"
        }
      ],
      "outputs": [
        {
          "id": 5,
          "name": "Aligned reads file",
          "abbrev": "aligned_read_file",
          "description": "Aligned reads file",
          "type": "File",
          "extension": ".sam"
        }
      ],
      "tool_id": 6,
      "runtime_id": 1
    }
  }
]


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
  }).length === 0);
  // console.log(notLinkedInputs);
  return notLinkedInputs.map(input => Object.assign({ tool_index: tool.tool_index, tool_name: tool.name }, input));
  // return notLinkedInputs;
});
console.log(externalInputs.flat(1));

const inputs = new YAMLMap();
for (const input of externalInputs.flat(1)) {
  const inputPair = cwlWorkflow.createPair([input.tool_name, input.tool_index, input.abbrev].join("_"), { type: input.type });
  inputs.add(inputPair);
}
cwlWorkflow.set("inputs", inputs);

const steps = new YAMLMap();
for (const tool of tools) {

  const stepIn = new YAMLMap();
  for (const input of tool.selected_version.inputs) {
    const links = toolLinks.filter(link => link.to_tool_index === tool.tool_index && link.to_input.id === input.id);
    let inputValue = [tool.name, tool.tool_index, input.abbrev].join("_");
    if (links.length > 0) {
      const link = links[0];
      inputValue = [link.from_tool.name, link.from_tool_index].join("_") + "/" + link.from_output.abbrev;
    }
    const inputPair = cwlWorkflow.createPair(input.abbrev, inputValue);
    stepIn.add(inputPair);
  }

  const stepOut = new YAMLSeq();
  for (const output of tool.selected_version.outputs) {
    stepOut.add(cwlWorkflow.createNode(output.abbrev));
  }

  const toolPair = cwlWorkflow.createPair(tool.name + "-" + tool.tool_index, {
    id: tool.name + "-" + tool.tool_index + "_" + tool.selected_version.semver,
    in: stepIn,
    out: stepOut,
    run: tool.name + "_" + tool.tool_index
  });
  steps.add(toolPair);
}
cwlWorkflow.set("steps", steps);


// console.log(cwlWorkflow.toString());
// console.log(stringify(cwlWorkflow));