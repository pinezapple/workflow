################## BASE IMAGE ######################
FROM biocontainers/biocontainers:v1.0.0_cv4

################## METADATA ######################

LABEL base_image="biocontainers:v1.0.0_cv4"
LABEL version="3"
LABEL software="ensembl-vep"
LABEL software.version="95.0"
LABEL about.summary="Variant Prediction"
LABEL about.home="https://ensembl.org/"
LABEL about.documentation="https://ensembl.org/"
LABEL about.license_file="https://ensembl.org/"
LABEL about.license="SPDX:GPL-3.0"
LABEL extra.identifiers.biotools="vep"
LABEL about.tags="Genomics"

################## MAINTAINER ######################
MAINTAINER Saulo Alves Aflitos <sauloal@gmail.com>

################## INSTALLATION ######################
ENV PERL5LIB="/opt/conda/lib/perl5/site_perl/5.22.0/x86_64-linux-thread-multi:/opt/conda/lib/perl5/site_perl/5.22.0"

RUN conda install ensembl-vep=95.0
RUN conda init bash
WORKDIR /data/

CMD ["vep"]
