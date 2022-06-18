FROM debian:buster

RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN apt-get update && apt-get install -y apt-utils

# Install and configure locale `en_US.UTF-8`
RUN apt-get update && apt-get install -y locales && \
    sed -i -e "s/# $en_US.*/en_US.UTF-8 UTF-8/" /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8
ENV LANG=en_US.UTF-8

RUN apt-get update && apt-get install -y git python2 python-pip g++ cmake python-ply python-tk tix pkg-config libssl-dev python-setuptools libreadline-dev sudo vim graphviz

# create a user:
RUN useradd -ms /bin/bash user && echo 'user:user' | chpasswd && adduser user sudo
RUN echo "ALL ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers
USER user
WORKDIR /home/user/

RUN git clone https://github.com/kenmcmil/ivy.git
WORKDIR /home/user/ivy/

RUN mkdir -p "/home/user/python/lib/python2.7/site-packages"
ENV PYTHONPATH="/home/user/python/lib/python2.7/site-packages"
# need to install pyparsing manually because otherwise wrong version found
RUN pip install pyparsing
RUN pip install z3-solver ipdb
USER root
USER user
RUN python2.7 setup.py develop --prefix="/home/user/python/"
ENV PATH=$PATH:"/home/user/python/bin/"

# also add poisonivy:
WORKDIR /home/user/
RUN git clone https://github.com/kach/poison-ivy.git
RUN sed -i "s/'ivy_check',/'ivy_check',\ 'complete=fo',/g" /home/user/poison-ivy/poisonivy.py

RUN mkdir /home/user/bin/
ENV PATH=/home/user/.local/bin:/home/user/bin:${PATH}
COPY --chown=user:user plot_dependencies.sh /home/user/bin/plot_dependencies.sh
RUN chmod +x /home/user/bin/plot_dependencies.sh
COPY --chown=user:user check_invariants.sh /home/user/bin/check_invariants.sh
RUN chmod +x /home/user/bin/check_invariants.sh

WORKDIR /home/user/
