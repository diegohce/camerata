#!/bin/bash

export OLD_PATH="$PATH"
export PATH="$PATH:$HOME/liteide/bin"

export GOPATH=$(pwd)
export GOBIN=$GOPATH/bin

export OLD_PS1="$PS1"
export PS1="(go $(basename $(pwd)))$PS1"

alias deactivate="unset GOPATH; unset GOBIN; unalias deactivate; export PS1=\"$OLD_PS1\"; unset OLD_PS1; export PATH=\"$OLD_PATH\"; unset OLD_PATH"



