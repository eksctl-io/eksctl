# Contributor Meetings

## Weekly meetings

- 1st/3rd Wednesday of the month: 09:00-09:30 UK time
- 2nd/4th Wednesday of the month: 16:00-16:30 UK time

### Meeting notes & agenda moved

To make it easier for the whole team to take notes, we moved the
[meeting notes](https://docs.google.com/document/d/1wHfEw7PAUDIFKuqEr7OwkpjgUUYlIDBHFj3I323vLUA/edit).

Please join our [Google
group](https://groups.google.com/forum/#!forum/eksctl)
to get write access to the doc and/or the meeting invite into
your calendar.

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
  - call for help on improving output [#585](https://github.com/eksctl-io/eksctl/issues/585)

### 06/03/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- upgrades - [#608](https://github.com/eksctl-io/eksctl/issues/608)

### 13/03/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- upgrades - [#608](https://github.com/eksctl-io/eksctl/issues/608)
   - `k8s.io/client-go` & co
   - `RawClient` and `RawResource` - [#624](https://github.com/eksctl-io/eksctl/pull/624)
- update on challenges with going to production from @mumoshu
   - node-local DNS and chaching
   - gitops with helmfile - need better way to manage `aws-auth`

### 20/03/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- release updates - 0.1.24
- plan towards v1alpha5
- more discussion around storing config objects [#642](https://github.com/eksctl-io/eksctl/issues/642)

### 27/03/2019

Attendees:
- @errordeveloper
- @mumoshu
- @martina-if

Topics:
- update on latest features that made it to the release
- v1alpha5
- nodegroup deletion [#664](https://github.com/eksctl-io/eksctl/issues/664)

## 03/04/2019

Attendees:
- @errordeveloper
- @mumoshu
- @pawelprazak

Topics:
- cloudformation template export functionality
  - make it explicitly about importing stacks
  - will help with refactoring
  - not a replacement for `eksctl apply`
- update about on-going efforts (#673, #695 etc)
- more discussion of storing config in cluster, and Cluster API

## 10/04/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- 0.1.28 & 0.1.29 releases are out
- @martina's progress on refactoring the SSH config for next release
- progress on task manager refactoring, handling of creation and deletion

## 17/04/2019

Attendees:
- @errordeveloper
- @mumoshu
- @pawelprazak

Topics:
- update on v1alpha5 progress
- @errordeveloper will be absent from next week's meeting 
- progress on other ongoing work (#714, #741)

## 01/04/2019

Attendees:
- @errordeveloper
- @mumoshu
- @martina
- @pawelprazak

Topics:
- 0.1.31 release is out
- next key feature: spot instances
- @pawelprazak's PR #718

## 15/05/2019

Attendees:
- @errordeveloper
- @mumoshu
- @martina-if

Topics:
- recent meeting cancellations
- no meeting for next two weeks due to KubeCon and vocations
- 0.1.32 release from previous week
- next feature release - spot instances
- config file documentation must be our priority

## 05/06/2019

Attendees:
- @errordeveloper
- @mumoshu
- @martina-if

Topics:
- update on releases: 0.1.33 and 0.1.34
- change in `--max-pods`



## 12/06/2019

Attendees:
- @errordeveloper
- @mumoshu
- @martina-if
- @gemagomez

Topics:
- update on releases: 0.1.35 and 0.1.36
- new commands for IAM identity mappings
- progress update on website & docs

## 07/08/2019

Attendees:
- @errordeveloper
- @mumoshu

Topics:
- latest releases, new cadence
- gitops features, add-ons and cluster configuration

## 14/08/2019

Attendees:
- @errordeveloper

Topics:
- 0.4.0 update
