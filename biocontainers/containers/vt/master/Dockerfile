################## BASE IMAGE ######################
FROM biocontainers/biocontainers:v1.1.0_cv2

################## METADATA ######################
LABEL base_image="biocontainers:v1.1.0_cv2"
LABEL version="1"
LABEL software="vt"
LABEL software.version="master"
LABEL about.summary="A tool set for short variant discovery in genetic sequence data."
LABEL about.home="https://genome.sph.umich.edu/wiki/Vt"
LABEL about.documentation="https://genome.sph.umich.edu/wiki/Vt"
LABEL about.license_file="https://genome.sph.umich.edu/wiki/Vt"
LABEL about.license="SPDX:GPL-3.0"
LABEL extra.identifiers.biotools="vt"
LABEL about.tags="Genomics"
LABEL extra.binaries="vt"


################## MAINTAINER ######################
MAINTAINER Alvin Chen <ychen@aperiomics.com>

################## INSTALLATION ######################

RUN conda install -c bioconda vt tabix

# WORKDIR /data/

CMD ["vt"]