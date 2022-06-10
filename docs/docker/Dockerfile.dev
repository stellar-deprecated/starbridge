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

RUN apt-get update && apt-get install -y git python3 python3-pip g++ cmake python3-ply python3-tk tix pkg-config libssl-dev python3-setuptools libreadline-dev sudo python3-pyparsing python3-pynvim vim vim-nox aptitude tmux gdb build-essential pkg-config fzf bat ripgrep stow wget curl bear
RUN pip install z3-solver ipdb

# create a user who can use sudo:
RUN apt-get update && apt-get install -y sudo
RUN useradd -ms /bin/bash user && echo 'user:user' | chpasswd && adduser user sudo
RUN echo "ALL ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

USER user

# install septum (tool to grep in big codebases)
WORKDIR /home/user/
RUN git clone https://github.com/pyjarrett/septum.git
RUN wget https://github.com/alire-project/alire/releases/download/v1.1.2/alr-1.1.2-bin-x86_64-linux.zip
RUN unzip /home/user/alr-1.1.2-bin-x86_64-linux.zip
RUN rm LICENSE.txt
WORKDIR /home/user/septum
RUN /home/user/bin/alr toolchain -n --select
RUN /home/user/bin/alr build
RUN cp bin/septum /home/user/bin
WORKDIR /home/user/

# configure stuff:
RUN curl -fLo /home/user/.vim/autoload/plug.vim --create-dirs https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
COPY --chown=user:user dotfiles dotfiles
WORKDIR /home/user/dotfiles/
RUN stow vim
# TODO: why not have stow create the following?
RUN mkdir /home/user/.vim/backupfiles
RUN mkdir /home/user/.vim/swapfiles
RUN mkdir /home/user/.vim/undofiles
RUN stow ctags
RUN stow tmux
# RUN stow zsh
RUN stow gdb
# RUN stow bash
WORKDIR /home/user/
RUN echo '\n\
# disable C-S in case it causes freezing in vim\n\
stty -ixon\n\
# Avoid duplicates\n\
HISTCONTROL=ignoredups:erasedups\n\
# When the shell exits, append to the history file instead of overwriting it\n\
shopt -s histappend\n\
# After each command, append to the history file and reread it\n\
PROMPT_COMMAND="${PROMPT_COMMAND:+$PROMPT_COMMAND$'\n'}history -a; history -c; history -r"\n'\
>> /home/user/.bashrc
RUN vim -E -s -u "/home/user/.vimrc" +PlugInstall +qall | true
ENV BAT_THEME=ansi

# install ivy
RUN git clone https://github.com/kenmcmil/ivy.git
WORKDIR /home/user/ivy/
# checkout Graydon's python3 branch:
RUN git checkout 2to3
RUN python3 setup.py develop --user

ENV PATH=/home/user/.local/bin:/home/user/bin:${PATH}

WORKDIR /home/user/
