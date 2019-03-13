# Contributor Meetings

## Weekly meetings

- 08:00 - 09:00 UTC (Europe & Asia)
- _TBC_ (Americas & Europe)

### 06/02/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- 0.1.20 release update
- testing improvements
    - better integration tests are needed
    - more complete unit tests with some kind cfn mocking
- docs and website
    - need to separate basic usage, getting-started vs advanced config file examples

### 13/02/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- 0.1.20 release update - it's out
- next release and `v1alpha5`
- separation of `NodeGroupConfig` object, so that nodegroups can be managed with a config file

### 20/02/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- 0.1.22 release
- update on latest contributions
- progress update on `kubectl drain` refactoring
- `eksctl create nodegroup` with `--config-file`

### 27/02/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- 0.1.23 release
- update on latest contributions
- progress update on `kubectl drain` refactoring - [k/k#72827](https://github.com/kubernetes/kubernetes/pull/72827) merged
- `eksctl utils describe-stacks`
  - additon of `--trail`
  - call for help on improving output [#585](https://github.com/weaveworks/eksctl/issues/585)

### 06/03/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- upgrades - [#608](https://github.com/weaveworks/eksctl/issues/608)

### 13/03/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- upgrades - [#608](https://github.com/weaveworks/eksctl/issues/608)
   - `k8s.io/client-go` & co
   - `RawClient` and `RawResource` - [#624](https://github.com/weaveworks/eksctl/pull/624)
- update on challenges with going to production from @mumoshu
   - node-local DNS and chaching
   - gitops with helmfile - need better way to manage `aws-auth`
