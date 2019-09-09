#!/bin/bash

cleanup ()
{
    exit 0;
}

trap cleanup SIGINT SIGTERM

while true; do sleep 1 ; echo ""; done
