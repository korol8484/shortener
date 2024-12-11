#!/bin/bash

baseDir=$(pwd)
protoDir="${baseDir}/proto"
replaceDir="${baseDir}/"
pathArr=()
exclude=("delete/path")

strstr() {
  [ "${1#*$2*}" = "$1" ] && return 1
  return 0
}

function walk_dir () {
    shopt -s nullglob dotglob

    for pathname in "$1"/*; do
        if [ -d "$pathname" ]; then
            walk_dir "$pathname"
        else
            case "$pathname" in
                *.proto)

                  status=0

                  for i in "${exclude[@]}"
                  do
                    if strstr $pathname $i; then
                      status=0
                      break
                    else
                       status=1
                    fi
                  done

                  if [ $status -eq 1 ]; then
                    pathArr+=("${pathname/$replaceDir/''}")
                  fi
            esac
        fi
    done
}

function join_by { local d=$1; shift; echo -n "$1"; shift; printf "%s" "${@/#/$d}"; }

walk_dir $protoDir
resultRun=$(join_by ' ' ${pathArr[@]})

printf "%s\n" "${pathArr[@]}"
mkdir -p /tmp/build

protoc -I$GOPATH/src/ \
-I/usr/local/include/ -I./ \
--go_out=/tmp/build \
--go_opt=paths=source_relative \
--go-grpc_out=/tmp/build \
--go-grpc_opt=paths=source_relative \
$resultRun

cp -R /tmp/build/proto "${baseDir}/internal/app/grpc/"
rm -rf /tmp/build
