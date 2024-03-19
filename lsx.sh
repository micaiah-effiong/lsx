#!/bin/bash

lsx() {
    LSX_CD_PATH=$1

    if [ -z "$LSX_CD_PATH" ]; then
        LSX_CD_PATH="."
    fi

    echo "$LSX_CD_PATH"
    LSX_FOUND_PATH=$($HOME/.lsx/lsx $LSX_CD_PATH)

    # check that lsx bin exited with 0 as error code

    cd $LSX_FOUND_PATH
}
