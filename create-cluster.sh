#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

scriptDir="$(cd "$(dirname "${0}")"; pwd)"

## Major constants

region="us-west-2"
nodeAMI="ami-993141e1"

keyName="ilya" # TODO - find a way to either upload a key from file or make it optional?

## Arguments

case "$#" in
  0)
    clusterName="cluster-1"
    numberOfNodes="2"
    nodeType="m4.large"
  ;;
  1)
    clusterName="${1}"
    numberOfNodes="2"
    nodeType="m4.large"
    ;;
  2)
    clusterName="${1}"
    numberOfNodes="${2}"
    nodeType="m4.large"
    ;;
  3)
    clusterName="${1}"
    numberOfNodes="${2}"
    nodeType="${3}"
    ;;
  *)
    echo "Usage: ${0} [<clusterName> [<numberOfNodes> [<nodeType>]]]"
    exit 1
    ;;
esac

## Setup

p="${scriptDir}/vendor/1.10.0/2018-05-09"
binariesDir="${p}/bin/$(uname | tr '[:upper:]' '[:lower:]')/amd64"

chmod -R +x "${binariesDir}"

export PATH="${binariesDir}:${PATH}"

aws configure add-model \
  --service-model "file://${p}/eks-2017-11-01.normal.json" \
  --service-name eks

## CloudFormation templates

serviceRoleTemplate="amazon-eks-service-role.yaml"
vpcTemplate="amazon-eks-vpc-sample.yaml"
nodeGroupTemplate="amazon-eks-nodegroup.yaml"

stackNamePrefix="EKS-${clusterName}-"

## Functions

createCluster() {
  echo "Creating cluster ${clusterName}"
  aws eks create-cluster \
    --region "${region}" \
    --role-arn "${clusterRoleARN}" \
    --subnets "${subnets[@]}" \
    --security-groups "${securityGroup}" \
    --cluster-name "${clusterName}"
}

describeCluster() {
  aws eks describe-cluster \
    --region "${region}" \
    --cluster-name "${clusterName}" \
    "$@"
}

createStack() {
  name="${stackNamePrefix}${1}"
  templateBody="file://${p}/${2}"
  shift 2
  aws cloudformation create-stack \
    --region "${region}" \
    --stack-name "${name}" \
    --template-body "${templateBody}" \
    "$@"
}

describeStacks() {
  aws cloudformation describe-stacks \
    --region "${region}"
}

checkStacksReadyCount() {
  ## TODO this is very naive (like most things here)
  ## it assumes stacks won't produce errors, e.g. rollback due to quota
  describeStacks \
    | jq -r \
      --arg prefix "${stackNamePrefix}" \
        '[.Stacks[] | select(.StackName | startswith($prefix)) | select(.StackStatus == "CREATE_COMPLETE")] | length'
}

getStackOutput() {
  name="${stackNamePrefix}${1}"
  outputKey="${2}"
  stacks="$(describeStacks)"
  echo "${stacks}" \
    | jq -r \
      --arg name "${name}" \
      --arg outputKey "${outputKey}" \
        '.Stacks[] | select(.StackName == $name) | .Outputs[] | select(.OutputKey == $outputKey) | .OutputValue'
}

## Action

stacksText="${stackNamePrefix}ServiceRole and ${stackNamePrefix}VPC stacks"

echo "Creating ${stacksText} we need first"

createStack "ServiceRole" "${serviceRoleTemplate}" \
  --capabilities "CAPABILITY_IAM"

createStack "VPC" "${vpcTemplate}" \
  --parameters \
    ParameterKey=ClusterName,ParameterValue="${clusterName}"

echo "Waiting until the ${stacksText} are ready"

until test "$(checkStacksReadyCount)" -ge 2 ; do sleep 20 ; done

echo "Collect outputs from the ${stacksText}"

securityGroup=($(getStackOutput "VPC" "SecurityGroup"))
subnets=($(getStackOutput "VPC" "Subnets"))
subnetsList=($(getStackOutput "VPC" "SubnetsList"))
clusterVPC=($(getStackOutput "VPC" "ClusterVPC"))
clusterRoleARN="$(getStackOutput "ServiceRole" "RoleArn")"

createCluster

echo "Creating ${stackNamePrefix}DefaultNodeGroup stack"

createStack "DefaultNodeGroup" "${nodeGroupTemplate}" \
  --capabilities "CAPABILITY_IAM" \
  --parameters \
    ParameterKey=ClusterName,ParameterValue="${clusterName}" \
    ParameterKey=NodeGroupName,ParameterValue="${clusterName}-DefaultNodeGroup" \
    ParameterKey=KeyName,ParameterValue="${keyName}" \
    ParameterKey=NodeImageId,ParameterValue="${nodeAMI}" \
    ParameterKey=NodeInstanceType,ParameterValue="${nodeType}" \
    ParameterKey=NodeAutoScalingGroupMinSize,ParameterValue="${numberOfNodes}" \
    ParameterKey=NodeAutoScalingGroupMaxSize,ParameterValue="${numberOfNodes}" \
    ParameterKey=ClusterControlPlaneSecurityGroup,ParameterValue="${securityGroup}" \
    ParameterKey=Subnets,ParameterValue=\"${subnetsList}\" \
    ParameterKey=VpcId,ParameterValue="${clusterVPC}"

echo "Waiting until cluster is ready"

until test "$(describeCluster --query "cluster.status")" = '"ACTIVE"' ; do sleep 20 ; done

export KUBECONFIG="${scriptDir}/${clusterName}.${region}.yaml"
echo "Saving cluster credentials in ${KUBECONFIG}"

masterEndpoint="$(describeCluster --query "cluster.masterEndpoint")"
certificateAuthorityData="$(describeCluster --query "cluster.certificateAuthority.data")"

kubeconfig="
apiVersion: v1
kind: Config
preferences: {}
current-context: '${clusterName}.${region}.eks.amazonaws.com'
clusters:
- name: '${clusterName}.${region}.eks.amazonaws.com'
  cluster:
    server: ${masterEndpoint}
    certificate-authority-data: ${certificateAuthorityData}
contexts:
- name: '${clusterName}.${region}.eks.amazonaws.com'
  context:
    cluster: '${clusterName}.${region}.eks.amazonaws.com'
    user: aws-eks-user
users:
- name: aws-eks-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: heptio-authenticator-aws
      args:
        - 'token'
        - '-i'
        - '${clusterName}'
        # - '-r'
        # - '<role-arn>'
"

echo "${kubeconfig}" > "${KUBECONFIG}"

echo "Waiting until ${stackNamePrefix}DefaultNodeGroup stack is ready"

until test "$(checkStacksReadyCount)" -eq 3 ; do sleep 20 ; done

nodeInstanceRoleARN=($(getStackOutput "DefaultNodeGroup" "NodeInstanceRole"))

nodeAuthConfigMap="
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: default
data:
  mapRoles: |
    - rolearn: "${nodeInstanceRoleARN}"
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
        - system:node-proxier
"

echo "${nodeAuthConfigMap}" | kubectl apply --filename='-'

echo "Cluster is ready, nodes will be added soon"
echo "Use the following command to monitor the nodes"
echo "$ kubectl --kubeconfig='${KUBECONFIG}' get nodes --watch"
