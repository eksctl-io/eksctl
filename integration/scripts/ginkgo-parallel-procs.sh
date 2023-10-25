#!/bin/sh

PARALLEL_PROCS="-p -procs="

case $1 in
  "crud")
    echo "${PARALLEL_PROCS}5"
    ;;
  "managed")
    echo "${PARALLEL_PROCS}5"
    ;;
  "windows")
    echo "${PARALLEL_PROCS}3"
    ;;
  *)
    echo ""
    ;;
esac