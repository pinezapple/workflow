FROM ubuntu:18.04
ENV DEBIAN_FRONTEND noninteractive


# Add user biodocker with password biodocker
RUN groupadd fuse && \
  useradd --create-home --shell /bin/bash --user-group --uid 1000 --groups sudo,fuse biodocker && \
  echo `echo "biodocker\nbiodocker\n" | passwd biodocker` 

RUN apt update && apt install python3-pip python3-dev git -y && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN rm -r /var/cache/apt/archives

# ENV GIT_URL=https://github.com/a-xavier/tapes.git
WORKDIR /
COPY tapes tapes
# RUN git clone ${GIT_URL}
RUN pip3 install -r  /tapes/requirements.txt

USER biodocker
