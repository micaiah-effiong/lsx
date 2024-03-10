#!/bin/bash

lsx() {
    echo "message"
    if [ -z "$1" ]; then
        echo "No directory path provided"
        exit 2
    else
        echo "$1"
        p=$(./main $1)
        cd $p
    fi
}
