################## BASE IMAGE ######################
FROM biocontainers/biocontainers:v1.0.0_cv4

################## METADATA ######################

LABEL base_image="biocontainers:v1.0.0_cv4" \
  version="2" \
  software="gatk" \
  software.version="4.0" \
  about.summary="The Genome Analysis Toolkit or GATK is a software package developed at the Broad Institute to analyze high-throughput sequencing data." \
  about.home="https://www.broadinstitute.org/gatk/|https://github.com/broadgsa/gatk-protected/" \
  about.license_file="https://github.com/broadgsa/gatk-protected/blob/master/licensing/protected_license.txt" \
  about.license="SPDX:BSD-3-Clause" \
  about.tags="Genomics" \
  extra.identifiers.biotools="gatk" \
  extra.binaries="GenomeAnalysisTK" \
  about.provides="gatk 4.0"

################## MAINTAINER ######################
MAINTAINER Saulo Alves Aflitos <sauloal@gmail.com>

################## INSTALLATION ######################

RUN conda install gatk4=4.0.5.1

WORKDIR /data/

CMD ["gatk"]
