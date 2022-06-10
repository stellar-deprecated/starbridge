FROM ubuntu:jammy

ARG DEBIAN_FRONTEND=noninteractive

# always run apt-get update before install to make sure we don't have a stale database
RUN apt-get update && apt-get install -y --no-install-recommends apt-utils ca-certificates

# Install and configure locale `en_US.UTF-8`
RUN apt-get update && apt-get install -y locales && \
    sed -i -e "s/# $en_US.*/en_US.UTF-8 UTF-8/" /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV TZ="America/Los_Angeles"

RUN apt-get update && apt-get install -y git python3 python3-pip g++ cmake python3-ply python3-tk tix pkg-config libssl-dev python3-setuptools libreadline-dev sudo python3-pyparsing build-essential pkg-config wget

RUN useradd -ms /bin/bash user
USER user
RUN pip install z3-solver

# install ivy
WORKDIR /home/user/
RUN git clone https://github.com/kenmcmil/ivy.git
WORKDIR /home/user/ivy/
# checkout Graydon's python3 branch:
RUN git checkout 2to3
RUN python3 setup.py develop --user

RUN mkdir /home/user/bin/
ENV PATH=/home/user/.local/bin:/home/user/bin:${PATH}
COPY --chown=user:user plot_dependencies.sh /home/user/bin/plot_dependencies.sh
RUN chmod +x /home/user/bin/plot_dependencies.sh
COPY --chown=user:user check_invariants.sh /home/user/bin/check_invariants.sh
RUN chmod +x /home/user/bin/check_invariants.sh

CMD ivy_check
