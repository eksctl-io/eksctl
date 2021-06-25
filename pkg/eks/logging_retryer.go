/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
)

const (
	maxRetries          = 13
	cfnMinThrottleDelay = 5 * time.Second
)

// LoggingRetryer adds some logging when we are retrying, so we have some idea what is happening
// Right now it is very basic - e.g. it only logs when we retry (so doesn't log when we fail due to too many retries)
// It was copied from k8s.io/kops/upup/pkg/fi/cloudup/awsup/logging_retryer.go; the original version used glog, and
// didn't export the constructor
type LoggingRetryer struct {
	client.DefaultRetryer
	cfnRetryer client.DefaultRetryer
}

var _ request.Retryer = &LoggingRetryer{}

func newLoggingRetryer() *LoggingRetryer {
	return &LoggingRetryer{
		DefaultRetryer: client.DefaultRetryer{
			NumMaxRetries: maxRetries,
		},
		cfnRetryer: client.DefaultRetryer{
			NumMaxRetries:    maxRetries,
			MinThrottleDelay: cfnMinThrottleDelay,
		},
	}
}

// ShouldRetry uses DefaultRetryer.ShouldRetry but also checks for non-retryable
// EC2MetadataError (see #2564)
func (l LoggingRetryer) ShouldRetry(r *request.Request) bool {
	shouldRetry := l.DefaultRetryer.ShouldRetry(r)
	if !shouldRetry {
		return false
	}
	if aerr, ok := r.Error.(awserr.RequestFailure); ok && aerr != nil && aerr.Code() == "EC2MetadataError" {
		switch aerr.StatusCode() {
		case http.StatusForbidden, http.StatusNotFound, http.StatusMethodNotAllowed:
			return false
		}
	}
	return true
}

// RetryRules extends on DefaultRetryer.RetryRules
func (l LoggingRetryer) RetryRules(r *request.Request) time.Duration {
	var (
		duration time.Duration
		service  = r.ClientInfo.ServiceName
	)

	if r.IsErrorThrottle() && service == cloudformation.ServiceName && r.Operation.Name == "DescribeStacks" {
		duration = l.cfnRetryer.RetryRules(r)
	} else {
		duration = l.DefaultRetryer.RetryRules(r)
	}

	name := "?"
	if r.Operation != nil {
		name = r.Operation.Name
	}
	methodDescription := service + "/" + name

	var errorDescription string
	if r.Error != nil {
		// We could check aws error Code & Message, but we expect them to be in the string
		errorDescription = fmt.Sprintf("%v", r.Error)
	} else {
		errorDescription = fmt.Sprintf("%d %s", r.HTTPResponse.StatusCode, r.HTTPResponse.Status)
	}

	logger.Warning("retryable error (%s) from %s - will retry after delay of %v", errorDescription, methodDescription, duration)

	return duration
}
