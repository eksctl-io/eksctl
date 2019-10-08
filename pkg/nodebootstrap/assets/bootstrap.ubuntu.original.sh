#!/usr/bin/env bash
# This boostrap script is customized for the Ubuntu EKS image

set -o pipefail
set -o nounset
set -o errexit

err_report() {
    echo "Exited with error on line $1"
}
trap 'err_report $LINENO' ERR

IFS=$'\n\t'

function print_help {
    echo "usage: $0 [options] <cluster-name>"
    echo "Bootstraps an instance into an EKS cluster"
    echo ""
    echo "-h,--help print this help"
    echo "--use-max-pods Sets --max-pods for the kubelet when true. (default: true)"
    echo "--b64-cluster-ca The base64 encoded cluster CA content. Only valid when used with --apiserver-endpoint. Bypasses calling \"aws eks describe-cluster\""
    echo "--apiserver-endpoint The EKS cluster API Server endpoint. Only valid when used with --b64-cluster-ca. Bypasses calling \"aws eks describe-cluster\""
    echo "--kubelet-extra-args Extra arguments to add to the kubelet. Useful for adding labels or taints."
    echo "--enable-docker-bridge Restores the docker default bridge network. (default: false)"
    echo "--aws-api-retry-attempts Number of retry attempts for AWS API call (DescribeCluster) (default: 3)"
}

POSITIONAL=()

while [[ $# -gt 0 ]]; do
    key="$1"
    case $key in
        -h|--help)
            print_help
            exit 1
            ;;
        --use-max-pods)
            USE_MAX_PODS="$2"
            shift
            shift
            ;;
        --b64-cluster-ca)
            B64_CLUSTER_CA=$2
            shift
            shift
            ;;
        --apiserver-endpoint)
            APISERVER_ENDPOINT=$2
            shift
            shift
            ;;
        --kubelet-extra-args)
            KUBELET_EXTRA_ARGS=$2
            shift
            shift
            ;;
        --enable-docker-bridge)
            ENABLE_DOCKER_BRIDGE=$2
            shift
            shift
            ;;
        --aws-api-retry-attempts)
            API_RETRY_ATTEMPTS=$2
            shift
            shift
            ;;
        *)    # unknown option
            POSITIONAL+=("$1") # save it in an array for later
            shift # past argument
            ;;
    esac
done

set +u
set -- "${POSITIONAL[@]}" # restore positional parameters
CLUSTER_NAME="$1"
set -u

USE_MAX_PODS="${USE_MAX_PODS:-true}"
B64_CLUSTER_CA="${B64_CLUSTER_CA:-}"
APISERVER_ENDPOINT="${APISERVER_ENDPOINT:-}"
KUBELET_EXTRA_ARGS="${KUBELET_EXTRA_ARGS:-}"
ENABLE_DOCKER_BRIDGE="${ENABLE_DOCKER_BRIDGE:-false}"
API_RETRY_ATTEMPTS="${API_RETRY_ATTEMPTS:-3}"

if [ -z "$CLUSTER_NAME" ]; then
    echo "CLUSTER_NAME is not defined"
    exit  1
fi

echo "Aliasing EKS k8s snap commands"
snap alias kubelet-eks.kubelet kubelet
snap alias kubectl-eks.kubectl kubectl

echo "Stopping k8s daemons until configured"
snap stop kubelet-eks
# Flush the restart-rate for failed starts

ZONE=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)
AWS_DEFAULT_REGION=$(echo $ZONE | awk '{print substr($0, 1, length($0)-1)}')

### kubelet kubeconfig

CA_CERTIFICATE_DIRECTORY=/etc/kubernetes/pki
CA_CERTIFICATE_FILE_PATH=$CA_CERTIFICATE_DIRECTORY/ca.crt
mkdir -p $CA_CERTIFICATE_DIRECTORY
if [[ -z "${B64_CLUSTER_CA}" ]] && [[ -z "${APISERVER_ENDPOINT}" ]]; then
    DESCRIBE_CLUSTER_RESULT="/tmp/describe_cluster_result.txt"
    rc=0
    # Retry the DescribleCluster API for API_RETRY_ATTEMPTS
    for attempt in `seq 0 $API_RETRY_ATTEMPTS`; do
        if [[ $attempt -gt 0 ]]; then
            echo "Attempt $attempt of $API_RETRY_ATTEMPTS"
        fi
        /snap/bin/aws eks describe-cluster \
            --region=${AWS_DEFAULT_REGION} \
            --name=${CLUSTER_NAME} \
            --output=text \
            --query 'cluster.{certificateAuthorityData: certificateAuthority.data, endpoint: endpoint}' > $DESCRIBE_CLUSTER_RESULT || rc=$?
        if [[ $rc -eq 0 ]]; then
            break
        fi
        if [[ $attempt -eq $API_RETRY_ATTEMPTS ]]; then
            exit $rc
        fi
        jitter=$((1 + RANDOM % 10))
        sleep_sec="$(( $(( 5 << $((1+$attempt)) )) + $jitter))"
        sleep $sleep_sec
    done
    B64_CLUSTER_CA=$(cat $DESCRIBE_CLUSTER_RESULT | awk '{print $1}')
    APISERVER_ENDPOINT=$(cat $DESCRIBE_CLUSTER_RESULT | awk '{print $2}')
fi

echo $B64_CLUSTER_CA | base64 -d > $CA_CERTIFICATE_FILE_PATH

sed -i s,CLUSTER_NAME,$CLUSTER_NAME,g /var/lib/kubelet/kubeconfig
/snap/bin/kubectl config \
    --kubeconfig /var/lib/kubelet/kubeconfig \
    set-cluster \
    kubernetes \
    --certificate-authority=/etc/kubernetes/pki/ca.crt \
    --server=$APISERVER_ENDPOINT

### kubelet.service configuration

MAC=$(curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/ -s | head -n 1 | sed 's/\/$//')
TEN_RANGE=$(curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/$MAC/vpc-ipv4-cidr-blocks | grep -c '^10\..*' || true )
DNS_CLUSTER_IP=10.100.0.10
if [[ "$TEN_RANGE" != "0" ]] ; then
    DNS_CLUSTER_IP=172.20.0.10;
fi

snap set kubelet-eks cluster-dns="$DNS_CLUSTER_IP"

INTERNAL_IP=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)
INSTANCE_TYPE=$(curl -s http://169.254.169.254/latest/meta-data/instance-type)

if [[ "$USE_MAX_PODS" = "true" ]]; then
    MAX_PODS_FILE="/etc/eks/eni-max-pods.txt"
    set +o pipefail
    MAX_PODS=$(grep ^$INSTANCE_TYPE $MAX_PODS_FILE | awk '{print $2}')
    set -o pipefail
    if [[ -n "$MAX_PODS" ]]; then
        snap set kubelet-eks max-pods="$MAX_PODS"
    else
        echo "No entry for $INSTANCE_TYPE in $MAX_PODS_FILE. Not setting max pods for kubelet"
    fi
fi

if [[ "$ENABLE_DOCKER_BRIDGE" = "true" ]]; then
    # Enabling the docker bridge network. We have to disable live-restore as it
    # prevents docker from recreating the default bridge network on restart
    echo "$(jq '.bridge="docker0" | ."live-restore"=false' /etc/docker/daemon.json)" > /etc/docker/daemon.json
    systemctl restart docker
fi

echo "Configuring kubelet snap"
snap set kubelet-eks \
    address=0.0.0.0 \
    authentication-token-webhook=true \
    authorization-mode=Webhook \
    allow-privileged=true \
    cloud-provider=aws \
    cluster-domain=cluster.local \
    cni-bin-dir=/opt/cni/bin \
    cni-conf-dir=/etc/cni/net.d \
    config=/etc/kubernetes/kubelet/kubelet-config.json \
    container-runtime=docker \
    node-ip="$INTERNAL_IP" \
    network-plugin=cni \
    pod-infra-container-image="602401143452.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/eks/pause-amd64:3.1" \
    cgroup-driver=cgroupfs \
    register-node=true \
    kubeconfig=/var/lib/kubelet/kubeconfig \
    feature-gates=RotateKubeletServerCertificate=true \
    anonymous-auth=false \
    client-ca-file="$CA_CERTIFICATE_FILE_PATH"

snap set kubelet-eks args="$KUBELET_EXTRA_ARGS"

echo "Starting k8s kubelet daemon"
snap start kubelet-eks
