#!/usr/bin/env bash
set -x

sdk_path=$1
interface=$2
output_file=$3
import=$4

ifacemaker -f "${sdk_path}/*.go" -s Client -i ${interface} -p aws -y "Auto-generated interface" -c "DONT EDIT: Auto generated" | sed "/\"context\".*/a . \"${import}\"" | gofmt | goimports > $output_file

# //go:generate ./generate.bash ${AWS_SDK_GO_DIR_V2}/service/eks/ EKS eksiface.go
# //go:generate ifacemaker -f "${AWS_SDK_GO_DIR_V2}/service/eks/*.go" -s Client -i EKS -p fakes -y "EKS interface" -c "DONT EDIT: Auto generated" -o eksiface.go
