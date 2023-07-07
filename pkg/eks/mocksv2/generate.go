package mocksv2

//go:generate "${GOBIN}/mockery" --tags netgo --dir=../../awsapi --all --outpkg=mocksv2 --output=./
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_V2_GO_DIR}/aws --name=CredentialsProvider --outpkg=mocksv2 --output=./
