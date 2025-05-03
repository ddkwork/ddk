mkdir -p .devcontainer
cd .devcontainer
cat >devcontainer.json<<-'EOFJSON'
// For format details, see https://aka.ms/devcontainer.json.
{
    "name": "Manjaro",
    "dockerFile": "Dockerfile",
    "runArgs": [
        "--cap-add=SYS_PTRACE",
        "--security-opt",
        "seccomp=unconfined"
    ],
    // "mounts": [
    //     "source=dind-var-lib-docker,target=/var/lib/docker,type=volume"
    // ],
    "mounts": [
        "source=/var/run/docker.sock,target=/var/run/docker.sock,type=bind"
    ],
    "overrideCommand": false,
    // Configure tool-specific properties.
    "customizations": {
        // Configure properties specific to VS Code.
        "vscode": {
            // Add the IDs of extensions you want installed when the container is created.
            "extensions": [
                // "MS-CEINTL.vscode-language-pack-zh-hans",
                "ms-azuretools.vscode-docker"
            ]
        }
    },
    // Use 'forwardPorts' to make a list of ports inside the container available locally.
    "forwardPorts": [
        5902
    ],
    // Use 'postCreateCommand' to run commands after the container is created.
    // "postCreateCommand": "docker --version",
    // Comment out to connect as root instead. More info: https://aka.ms/vscode-remote/containers/non-root.
    // "build": {
    //     "args": {
    //         "INSTALL_ZSH": "false",
    //         "ENABLE_NONROOT_DOCKER": "false"
    //     }
    // },
    "remoteUser": "ddk"
}
EOFJSON

cat > Dockerfile<<-'EOFDKF'
# syntax=docker/dockerfile:1
#---------------------------
# FROM cake233/manjaro-zsh-amd64

FROM cake233/manjaro-xfce-amd64

# set username & group
ARG USERNAME=ddk
ARG GROUPNAME=ddk
# ARG USER_UID=1001
# ARG USER_GID=$USER_UID

# rm cn mirrorlist
# RUN sed -e '/bfsu.edu.cn/d' \
#     -e '/tuna.tsinghua.edu.cn/d' \
#     -e '/opentuna.cn/d' \
#     -i /etc/pacman.conf

# install dependencies
# live server: https://docs.microsoft.com/en-us/visualstudio/liveshare/reference/linux#install-linux-prerequisites
RUN pacman -Syu \
    --noconfirm \
    --needed \
    base \
    base-devel \
    git \
    lib32-gcc-libs \
    lib32-glibc \
    gcr \
    liburcu \
    openssl-1.0 \
    krb5 \
    icu \
    zlib \
    gnome-keyring \
    libsecret \
    desktop-file-utils \
    xorg-xprop \
    xdg-utils \
    jdk17-openjdk \
    yay \
    xarchiver \
    go \
    gdb lldb cmake clang gcc \
    firefox \
    gradle \
    android-sdk \
    android-ndk

RUN pacman -R --noconfirm  xfce4-screensaver  chromium
RUN yay -Sy --noconfirm goland clion 
 
##########################################################################
# WORKDIR /opt/vlang
# RUN mkdir -p /opt/vlang
# RUN git clone https://github.com/vlang/v /opt/vlang && make && v -version
# RUN mkdir -p /opt/vlang && ln -s /opt/vlang/v /usr/bin/v
 
# /*
# git clone https://github.com/vlang/v
# cd v
# make
# # HINT: Using Windows?: run make.bat in the cmd.exe shell
#     sudo ./v symlink
#     git clone https://github.com/vlang/vls && cd vls
#     v run build.vsh clang
# */

# # v up
# # v install ui
# # sudo pacman -S libxi libxcursor mesa
# # v install --git https://github.com/IsaiahPatton/Vide
# # go install ​module​ ​github.com/jmigpin/editor@latest
##########################################################################

# locale: Chinese Simplified (China)
ENV LANG=zh_CN.UTF-8
 
# add new user
RUN groupadd --force ${GROUPNAME} \
    && useradd --create-home --gid ${GROUPNAME} ${USERNAME} \
    && mkdir -p /etc/sudoers.d \
    && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/ddk \
    && chmod 400 /etc/sudoers.d/ddk

WORKDIR ["/home/$USERNAME"]


# # git clone --recursive https://github.com/goplus/gop.git
# # cd gop
# # # On mac/linux run:
# # ./all.bash

# RUN cd $WORKDIR \
#     git clone --recursive https://github.com/goplus/c2go \
#     go install -v ./...

# # c源码目录
# # c2go  .


# RUN git clone --recursive https://github.com/xmake-io/xmake.git \
#     cd ./xmake \
#     ./scripts/get.sh __local__ \
#     source ~/.xmake/profile

# # RUN  [[ $(command -v curl) ]] || sudo pacman -Syu curl \
# #      bash -c "$(curl -L l.tmoe.me/ee/zsh)"

# RUN git clone https://github.com/goreleaser/goreleaser \
#     cd goreleaser \
#     go install ./... 
# # goreleaser --version
# # goreleaser init
# # goreleaser --snapshot --skip-publish --rm-dist




# clean cache
RUN yes | pacman -Scc; \
    rm -rf /var/cache/pacman/pkg/* \
    /tmp/* \
    2>/dev/null

# command: sleep infinity
CMD [ "sleep", "inf" ]
EOFDKF