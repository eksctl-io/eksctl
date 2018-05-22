#!/bin/bash -x

set -o errexit
set -o pipefail
set -o nounset

scriptDir="$(cd "$(dirname "${0}")"; pwd)"

region="us-west-2"
nodeAMI="ami-993141e1"

keyName="ilya" # TODO - find a way to either upload a key from file or make it optional?

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

createCluster() {
  aws eks create-cluster \
    --region "${region}" \
    --role-arn "${clusterRoleARN}" \
    --subnets "${subnets[@]}" \
    --security-groups "${securityGroups}" \
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

## TODO this is very naive, it assumes stacks won't produce errors, e.g. rollback due to quota
checkStacksReadyCount() {
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

## Create two stacks we need first
createStack "ServiceRole" "${serviceRoleTemplate}" \
  --capabilities "CAPABILITY_IAM"

createStack "VPC" "${vpcTemplate}" \
  --parameters \
    ParameterKey=ClusterName,ParameterValue="${clusterName}"

## Wait until the two stack are ready
until test "$(checkStacksReadyCount)" -eq 2 ; do sleep 20 ; done

## Obtain outputs from each of the stacks
securityGroups=($(getStackOutput "VPC" "SecurityGroups"))
subnets=($(getStackOutput "VPC" "SubnetIds"))
clusterRoleARN="$(getStackOutput "ServiceRole" "RoleArn")"

## Now, create the actual cluster
createCluster

#### Next, create the nodes
##createStack "DefaultNodeGroup" "${nodeGroupTemplate}" \
##  --capabilities "CAPABILITY_IAM" \
##  --parameters \
##    ParameterKey=ClusterName,ParameterValue="${clusterName}" \
##    ParameterKey=NodeGroupName,ParameterValue="${clusterName}-DefaultNodeGroup" \
##    ParameterKey=KeyName,ParameterValue="${keyName}" \
##    ParameterKey=NodeImageId,ParameterValue="${nodeAMI}" \
##    ParameterKey=NodeInstanceType,ParameterValue="${nodeType}" \
##    ParameterKey=NodeAutoScalingGroupMinSize,ParameterValue="${numberOfNodes}" \
##    ParameterKey=NodeAutoScalingGroupMaxSize,ParameterValue="${numberOfNodes}" \
##    ParameterKey=ControlPlaneSecurityGroup,ParameterValue="${securityGroups[1]}" \
##    ParameterKey=Subnets,ParameterValue="${subnets}" \ # TODO - comma-separated
##    ParameterKey=VpcId,ParameterValue="${clusterVPC}" \ # TODO - not in outputs

## Wait until cluster is ready
until test "$(describeCluster --query "cluster.status")" = '"ACTIVE"' ; do sleep 20 ; done

## Obtain cluster credentials
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
    server: '${masterEndpoint}'
    certificate-authority-data: '${certificateAuthorityData}'
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

nodeAuthConfigMap="
apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: default
data:
  mapRoles: |
    - rolearn: <ARN of instance role (not instance profile)> ## TODO
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
        - system:node-proxier
"

## Write kubeconfig file
export KUBECONFIG="${scriptDir}/../${clusterName}.${region}.yaml"
echo "${kubeconfig}" > "${KUBECONFIG}"

## Authorise the nodes to join the cluster
echo "${nodeAuthConfigMap}" | kubectl apply --filename='-'

##kubectl get nodes --watch
