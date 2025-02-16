#!/bin/bash

# This script is run on a daily basis via GitHub Actions. It is responsible for
# automatically updating all of goformation's CloudFormation resources to match
# what is currently supported in AWS CloudFormation.
# 
# It does the following:
# 
# 1. Clone github.com/goformation/goformation to local filesystem
# 2. Switch to aws-goformation-updates branch (creating it if necessary)
# 3. Merge in any changes from github.com/awslabs/goformation (upstream)
# 4. Run 'go generate'
# 5. Check whether any resources changed
# 6. Run 'go test' to ensure all tests pass
# 7. Commit any changes and push to github.com/goformation/goformation
# 8. Create a pull request from github.com/goformation/goformation (branch: aws-goformation-updates) to github.com/awslabs/goformation (master)
#
# It relies on the environment GITHUB_TOKEN containing a Personal Access Token for the github.com/goformation user.

# Bomb out on errors
set -e

# The repo/branch to create the PR against
SRC_REPO="aws-goformation/goformation"
SRC_BRANCH="master"
DST_REPO="awslabs/goformation"
DST_BRANCH="master"

# Git details (for the commit)
COMMIT_NAME="AWS GoFormation"
COMMIT_EMAIL="goformation@amazon.com"
COMMIT_MSG="feat(schema): AWS CloudFormation Update ($(date +%F))"

# Configure the Git user for any commits
git config --global user.name "${COMMIT_NAME}"
git config --global user.email "${COMMIT_EMAIL}"

echo "Checking out github.com/${SRC_REPO}..."
UPSTREAM_DIR=/tmp/upstream/src/github.com/awslabs/goformation
mkdir -p ${UPSTREAM_DIR}
git clone https://github.com/${SRC_REPO}.git ${UPSTREAM_DIR} > /dev/null 2>&1
GOPATH=/tmp/upstream
cd ${UPSTREAM_DIR}

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
REQUEST_BRANCH="aws-goformation-updates"
echo "Creating new branch: ${REQUEST_BRANCH} tracking origin/${SRC_BRANCH}"
git checkout -b ${REQUEST_BRANCH} origin/${SRC_BRANCH}
git pull --rebase origin ${REQUEST_BRANCH} || true

# Merging in any changes from upstream
git remote add upstream https://github.com/${DST_REPO}.git
git pull --rebase upstream ${SRC_BRANCH} || true

 
echo "Auto-generating AWS CloudFormation resources..."
go generate

# Don't proceed if there were no new resources generated/updated
if [[ -z $(git status -s | grep "cloudformation/") ]]; then
    echo "No changes - no pull request necessary."
    exit 0;
fi
  
echo "Changes found:"
git status -s

echo "Running tests..."
go test -v ./...

echo "Committing changes..."
git add cloudformation/*
git add schema/*
git commit -m "${COMMIT_MSG}" 

echo "Pushing changes..."
git remote add origin-push https://x-access-token:${GITHUB_TOKEN}@github.com/${SRC_REPO}.git > /dev/null 2>&1
git push --quiet --set-upstream origin-push ${REQUEST_BRANCH}

echo "Installing GitHub Hub"
git clone https://github.com/github/hub.git /tmp/hub
cd /tmp/hub
go get ./...
./script/build 

echo "Generating Pull Request for merging ${REPO}/${REQUEST_BRANCH} to ${REPO}/${DST_BRANCH}..."
cd ${UPSTREAM_DIR}
echo "${COMMIT_MSG}" | /tmp/hub/bin/hub pull-request -F --head aws-goformation:${REQUEST_BRANCH} - /dev/null 2>&1 || true
