# GoFormation Versioning Changelog

## [4.15.5](https://github.com/awslabs/goformation/compare/v4.15.4...v4.15.5) (2020-11-06)


### Bug Fixes

* **schema:** CloudFormation Updates ([#331](https://github.com/awslabs/goformation/issues/331)) ([12f9c83](https://github.com/awslabs/goformation/commit/12f9c835cbd9b980e06e30c458f95b14e16a3771))

## [4.15.4](https://github.com/awslabs/goformation/compare/v4.15.3...v4.15.4) (2020-11-01)


### Bug Fixes

* **schema:** CloudFormation Updates ([#330](https://github.com/awslabs/goformation/issues/330)) ([4070319](https://github.com/awslabs/goformation/commit/40703191771425b3519128027478cf740d10f2d9))

## [4.15.3](https://github.com/awslabs/goformation/compare/v4.15.2...v4.15.3) (2020-10-23)


### Bug Fixes

* **schema:** CloudFormation Updates ([#329](https://github.com/awslabs/goformation/issues/329)) ([4c1362b](https://github.com/awslabs/goformation/commit/4c1362bac14a3d6fc2cebd56eddcaf440228a4e9))

## [4.15.2](https://github.com/awslabs/goformation/compare/v4.15.1...v4.15.2) (2020-10-11)


### Bug Fixes

* **schema:** CloudFormation Updates ([#320](https://github.com/awslabs/goformation/issues/320)) ([49879b4](https://github.com/awslabs/goformation/commit/49879b45f79dc9880071facd32d76cc4bf0570eb))

## [4.15.1](https://github.com/awslabs/goformation/compare/v4.15.0...v4.15.1) (2020-10-11)


### Bug Fixes

* **intrinsics:** Join function to allow to use parameters of type `List<>` ([#309](https://github.com/awslabs/goformation/issues/309)) ([6cc1cd3](https://github.com/awslabs/goformation/commit/6cc1cd329047227674caaf1b546066e3043c6616))

# [4.15.0](https://github.com/awslabs/goformation/compare/v4.14.0...v4.15.0) (2020-08-16)


### Features

* **schema:** dummy commit - trigger CI for schema update ([66bc344](https://github.com/awslabs/goformation/commit/66bc344ed6a17a613abbb4d217afcc0b8ea02b48))

# [4.14.0](https://github.com/awslabs/goformation/compare/v4.13.1...v4.14.0) (2020-07-26)


### Features

* **schema:** Add support for Template Outputs ([#291](https://github.com/awslabs/goformation/issues/291)) ([6875c50](https://github.com/awslabs/goformation/commit/6875c50d00d8e1af71d9bad5788446b29ab03513))

## [4.13.1](https://github.com/awslabs/goformation/compare/v4.13.0...v4.13.1) (2020-07-26)


### Bug Fixes

* **schema:** Add Change and Update policies to the Unmarshal method ([#288](https://github.com/awslabs/goformation/issues/288)) ([989b05f](https://github.com/awslabs/goformation/commit/989b05fa78cb9e72f6d59298fb8bb287612f322e))

# [4.13.0](https://github.com/awslabs/goformation/compare/v4.12.0...v4.13.0) (2020-07-26)


### Features

* **schema:** adding AWS::Serverless::StateMachine and FileSystemConfigs to Function ([#284](https://github.com/awslabs/goformation/issues/284)) ([d2d23ca](https://github.com/awslabs/goformation/commit/d2d23cafba606a8ea40649cc666073fa0e2d5ad3))

# [4.12.0](https://github.com/awslabs/goformation/compare/v4.11.0...v4.12.0) (2020-07-21)


### Features

* **schema:** Add new DynamoDBEvent options ([#289](https://github.com/awslabs/goformation/issues/289)) ([741228d](https://github.com/awslabs/goformation/commit/741228d6923ea10f1d22a901bbddf106d5c71cd7))

# [4.11.0](https://github.com/awslabs/goformation/compare/v4.10.1...v4.11.0) (2020-06-28)


### Features

* **schema:** CFN Updates ([#287](https://github.com/awslabs/goformation/issues/287)) ([9778479](https://github.com/awslabs/goformation/commit/97784795e35035b71b946d0ca69ef4d380d3b4a8))

## [4.10.1](https://github.com/awslabs/goformation/compare/v4.10.0...v4.10.1) (2020-06-22)


### Bug Fixes

* **generator:** update the generation making it easier to fix CF schema errors to generate ([#285](https://github.com/awslabs/goformation/issues/285)) ([6751e5b](https://github.com/awslabs/goformation/commit/6751e5b6ecbe1daee45171528cb1300efc6fb300))

# [4.10.0](https://github.com/awslabs/goformation/compare/v4.9.0...v4.10.0) (2020-06-22)


### Features

* **schema:** Serverless eventbridgeruleevent ([#279](https://github.com/awslabs/goformation/issues/279)) ([2a9e572](https://github.com/awslabs/goformation/commit/2a9e572313485023dc4e57cb8facda72a3571307))

# [4.9.0](https://github.com/awslabs/goformation/compare/v4.8.0...v4.9.0) (2020-06-22)


### Features

* **schema:** Add OpenApiVersion field to serverless Api ([#281](https://github.com/awslabs/goformation/issues/281)) ([bccc71b](https://github.com/awslabs/goformation/commit/bccc71b90531fb6bba8465b578fc2accc4dc6e34))

# [4.8.0](https://github.com/awslabs/goformation/compare/v4.7.1...v4.8.0) (2020-04-04)


### Features

* **schema:** Add UpdateReplacePolicy to the templates and the policies so that it is generated for every resource ([#272](https://github.com/awslabs/goformation/issues/272)) ([696c515](https://github.com/awslabs/goformation/commit/696c515bcbb07105683a328ef0e161d62146114b))

## [4.7.1](https://github.com/awslabs/goformation/compare/v4.7.0...v4.7.1) (2020-04-04)


### Bug Fixes

* **intrinsics:** change Fn::Sub to allow AWS pseudo parameters ([#275](https://github.com/awslabs/goformation/issues/275)) ([5a48c27](https://github.com/awslabs/goformation/commit/5a48c27630b945dcdc33133defd0241f898ccc52)), closes [#274](https://github.com/awslabs/goformation/issues/274) [#202](https://github.com/awslabs/goformation/issues/202)

# [4.7.0](https://github.com/awslabs/goformation/compare/v4.6.0...v4.7.0) (2020-02-28)


### Features

* **schema:** Added CloudWatch Logs event for SAM ([#271](https://github.com/awslabs/goformation/issues/271)) ([fedb013](https://github.com/awslabs/goformation/commit/fedb013e3b19ab1242cf8e3ae28a40240103d9b1))

# [4.6.0](https://github.com/awslabs/goformation/compare/v4.5.1...v4.6.0) (2020-02-22)


### Features

* **schema:** CloudFormation Updates (2020-02-22) ([#269](https://github.com/awslabs/goformation/issues/269)) ([ffd88a6](https://github.com/awslabs/goformation/commit/ffd88a6a9b0349853517e811169ee66804d79a2e))

## [4.5.1](https://github.com/awslabs/goformation/compare/v4.5.0...v4.5.1) (2020-02-14)


### Bug Fixes

* **schema, parser:** change Transform json schema to allow multiple macros ([#268](https://github.com/awslabs/goformation/issues/268)) ([072fc74](https://github.com/awslabs/goformation/commit/072fc74628c8ee9a603c2e502ac458af916afc07)), closes [#267](https://github.com/awslabs/goformation/issues/267)

# [4.5.0](https://github.com/awslabs/goformation/compare/v4.4.0...v4.5.0) (2020-02-13)


### Features

* **schema:** CloudFormation Updates (2020-02-13) ([#266](https://github.com/awslabs/goformation/issues/266)) ([bc75922](https://github.com/awslabs/goformation/commit/bc75922eb604d6e43f290912234a644c4d7584b5))

# [4.4.0](https://github.com/awslabs/goformation/compare/v4.3.0...v4.4.0) (2020-01-30)


### Features

* **schema:** CloudFormation Updates (2020-01-30) ([#263](https://github.com/awslabs/goformation/issues/263)) ([fda2d31](https://github.com/awslabs/goformation/commit/fda2d31f384eabbbf432ad1ee77ff8db6d0f2e73))

# [4.3.0](https://github.com/awslabs/goformation/compare/v4.2.0...v4.3.0) (2020-01-30)


### Features

* **schema:** add CloudFormation parameter type ([#259](https://github.com/awslabs/goformation/issues/259)) ([27fe204](https://github.com/awslabs/goformation/commit/27fe204f7addb8cb1bd6e977b0f717c04b09364a))

# [4.2.0](https://github.com/awslabs/goformation/compare/v4.1.0...v4.2.0) (2020-01-29)


### Features

* **parser:** Add support for Conditions ([#260](https://github.com/awslabs/goformation/issues/260)) ([1b00f17](https://github.com/awslabs/goformation/commit/1b00f17a33109023ad8a4471812448dc1d0db776))

# [4.1.0](https://github.com/awslabs/goformation/compare/v4.0.3...v4.1.0) (2019-12-09)


### Features

* **schema:** CloudFormation Updates (2019-12-09) ([#251](https://github.com/awslabs/goformation/issues/251)) ([a23ba41](https://github.com/awslabs/goformation/commit/a23ba416a24649c7296a0bc507c7940d9082ea30))

## [4.0.3](https://github.com/awslabs/goformation/compare/v4.0.2...v4.0.3) (2019-11-30)


### Bug Fixes

* **schema:** AWS::Serverless::Function S3 notification filters ([#249](https://github.com/awslabs/goformation/issues/249)) ([a50ef92](https://github.com/awslabs/goformation/commit/a50ef9291026420ea8a5e74790fc49b8a9c7fd85)), closes [#74](https://github.com/awslabs/goformation/issues/74)

## [4.0.2](https://github.com/awslabs/goformation/compare/v4.0.1...v4.0.2) (2019-11-30)


### Bug Fixes

* **schema:** AWS::Serverless:Api.Cors ([#246](https://github.com/awslabs/goformation/issues/246)) ([62fd56a](https://github.com/awslabs/goformation/commit/62fd56a62586c65722f99dbd4c8308ab42fcfc1d)), closes [#244](https://github.com/awslabs/goformation/issues/244)

## [4.0.1](https://github.com/awslabs/goformation/compare/v4.0.0...v4.0.1) (2019-11-30)


### Bug Fixes

* **schema:** AWS::Serverless::Api.MethodSettings should be a list ([a1f340a](https://github.com/awslabs/goformation/commit/a1f340a07e0ba4f21b8655da2c4d608849278901)), closes [#242](https://github.com/awslabs/goformation/issues/242)

# [4.0.0](https://github.com/awslabs/goformation/compare/v3.1.0...v4.0.0) (2019-11-30)


* Fix method conflicts (#245) ([d0b0a8b](https://github.com/awslabs/goformation/commit/d0b0a8bc322e27f72e840c9847f3c822d4efa933)), closes [#245](https://github.com/awslabs/goformation/issues/245) [#241](https://github.com/awslabs/goformation/issues/241) [#294](https://github.com/awslabs/goformation/issues/294)


### BREAKING CHANGES

* This change refactors the DependsOn, Metadata, CreationPolicy,
UpdatePolicy and DeletionPolicy methods on each resource to a new
name. This is required, as some CloudFormation resources use these
keywords as properties (AWS::AppMesh::Route.GrpcRouteMatch has a
Metadata field for example), which causes a conflict.

`resource.DependsOn()` method is refactored to `resource.AWSCloudFormationDependsOn` field.
`resource.SetDependsOn()` method is refactored to `resource.AWSCloudFormationDependsOn` field.
`resource.Metadata()` method is refactored to `resource.AWSCloudFormationMetadata` field.
`resource.SetMetadata()` method is refactored to `resource.AWSCloudFormationMetadata` field.
`resource.CreationPolicy()` method is refactored to `resource.AWSCloudFormationCreationPolicy` field.
`resource.SetCreationPolicy()` method is refactored to `resource.AWSCloudFormationCreationPolicy` field.
`resource.UpdatePolicy()` method is refactored to `resource.AWSCloudFormationUpdatePolicy` field.
`resource.SetUpdatePolicy()` method is refactored to `resource.AWSCloudFormationUpdatePolicy` field.
`resource.DeletionPolicy()` method is refactored to `resource.AWSCloudFormationDeletionPolicy` field.
`resource.SetDeletionPolicy()` method is refactored to `resource.AWSCloudFormationDeletionPolicy` field.

# [3.1.0](https://github.com/awslabs/goformation/compare/v3.0.1...v3.1.0) (2019-10-29)


### Features

* **schema:** AWS CloudFormation Update (2019-10-29) ([#239](https://github.com/awslabs/goformation/issues/239)) ([7ff8499](https://github.com/awslabs/goformation/commit/7ff84990c89e11815d22e06d377e110ae422cc17))

## [3.0.1](https://github.com/awslabs/goformation/compare/v3.0.0...v3.0.1) (2019-10-29)


### Bug Fixes

* **schema:** Ordered cloudformation/all.go file ([#238](https://github.com/awslabs/goformation/issues/238)) ([91254f3](https://github.com/awslabs/goformation/commit/91254f30925b89db5e79604d812a1ee9279267bd))

# [3.0.0](https://github.com/awslabs/goformation/compare/v2.3.1...v3.0.0) (2019-10-27)


* Group CloudFormation resources by AWS service name (#234) ([d0749e6](https://github.com/awslabs/goformation/commit/d0749e6a8fc5e7b0ddc301aef0170e12c7dc459c)), closes [#234](https://github.com/awslabs/goformation/issues/234)


### BREAKING CHANGES

* this change moves all Cloudformation resources to
packages based on the AWS service name. The main motivation for this is
that building goformation on some platforms (Windows) failed due to too
many files in the old cloudformation/resources package. This new package
style has a nice benefit of slightly nicer to use API, but is a breaking
change and will require refactoring existing codebases to update to v3.

Old usage:

```go
import "github.com/awslabs/goformation/v2/cloudformation/resources"

... snip ...

topic := &resources.AWSSNSTopic{}
```

New usage:

```go
import "github.com/awslabs/goformation/v4/cloudformation/sns"

...snip...

topic := &sns.Topic{}
```

Most tests are still failing at this point and need refactoring.

* fix(schema): Tag handling

Fixed tag handling for new grouped resources style (via new tags.Tag
struct).

* fix(schema): SAM specification

SAM Specification now generates nicely with new grouped resources
format. Also all tests are now passing \o/

# [2.3.0](https://github.com/awslabs/goformation/compare/v2.2.2...v2.3.0) (2019-03-20)


### Bug Fixes

* **parser:** Unmarshalling of resources with polymorphic properties (like S3 events) now works ([#188](https://github.com/awslabs/goformation/issues/188)) ([8eff90a](https://github.com/awslabs/goformation/commit/8eff90a))


### Features

* **sam:** Add support for `AWS::Serverless::Api.TracingEnabled`, `AWS::Serverless::Function.PermissionsBoundary`, `AWS::Serverless::Function.DynamoEvent.Enabled`, `AWS::Serverless::Function.KinesisEvent.Enabled`, and `AWS::Serverless::Function.SQSEvent.Enabled` ([#191](https://github.com/awslabs/goformation/issues/191)) ([38f0187](https://github.com/awslabs/goformation/commit/38f0187))
* **schema:** AWS CloudFormation Update (2019-03-15) ([#189](https://github.com/awslabs/goformation/issues/189)) ([8b332a4](https://github.com/awslabs/goformation/commit/8b332a4))

## [2.2.2](https://github.com/awslabs/goformation/compare/v2.2.1...v2.2.2) (2019-03-13)


### Bug Fixes

* **parser:** Select the correct AWS CloudFormation resource type based on similarity ([#183](https://github.com/awslabs/goformation/issues/183)) ([5749b23](https://github.com/awslabs/goformation/commit/5749b23))

## [2.2.1](https://github.com/awslabs/goformation/compare/v2.2.0...v2.2.1) (2019-03-10)


### Bug Fixes

* **parser:** fix invalid YAML template error for custom tag marshaler ([#177](https://github.com/awslabs/goformation/issues/177)) ([035d438](https://github.com/awslabs/goformation/commit/035d438))

# [2.2.0](https://github.com/awslabs/goformation/compare/v2.1.5...v2.2.0) (2019-03-10)


### Features

* **schema:** regenerated resources to apply SAM schema fixes from previous PR ([b30c019](https://github.com/awslabs/goformation/commit/b30c019))

## [2.1.5](https://github.com/awslabs/goformation/compare/v2.1.4...v2.1.5) (2019-03-10)


### Bug Fixes

* **parser:** do not break if a non-intrinsic `Condition` statement is found in a YAML template ([#169](https://github.com/awslabs/goformation/issues/169)) ([e4671e3](https://github.com/awslabs/goformation/commit/e4671e3))

## [2.1.4](https://github.com/awslabs/goformation/compare/v2.1.3...v2.1.4) (2019-03-10)


### Bug Fixes

* **schema:** fixed incorrect field type for AWS::Serverless::Application.Location ([#167](https://github.com/awslabs/goformation/issues/167)) ([3f1817b](https://github.com/awslabs/goformation/commit/3f1817b))

## [2.1.3](https://github.com/awslabs/goformation/compare/v2.1.2...v2.1.3) (2019-03-10)


### Bug Fixes

* **schema:** maps within YAML templates should allow unknown fields/properties ([3b6e359](https://github.com/awslabs/goformation/commit/3b6e359))

## [2.1.2](https://github.com/awslabs/goformation/compare/v2.1.1...v2.1.2) (2019-03-10)


### Bug Fixes

* **CI:** fix broken GitHub PR integration ([#185](https://github.com/awslabs/goformation/issues/185)) ([d42d00a](https://github.com/awslabs/goformation/commit/d42d00a))

## [2.1.1](https://github.com/awslabs/goformation/compare/v2.1.0...v2.1.1) (2019-03-10)


### Bug Fixes

* **CI:** only run semantic-release on push-to-master (not on pull requests) ([#184](https://github.com/awslabs/goformation/issues/184)) ([c83945a](https://github.com/awslabs/goformation/commit/c83945a))

# [2.1.0](https://github.com/awslabs/goformation/compare/v2.0.0...v2.1.0) (2019-03-10)


### Features

* **CI:** auto-generate AUTHORS.md file ([b37af7b](https://github.com/awslabs/goformation/commit/b37af7b))

# Semantic Versioning Changelog

# [2.0.0](https://github.com/awslabs/goformation/compare/v1.4.1...v2.0.0) (2019-03-10)


### Code Refactoring

* **generator:** moving resources and policies into their own packages ([#161](https://github.com/awslabs/goformation/issues/161)) ([03a0123](https://github.com/awslabs/goformation/commit/03a0123))


### BREAKING CHANGES

* **generator:** this PR refactors the auto-generated CloudFormation resources out of the cloudformation package and into a dedicated package (resources). This helps keep the auto generated files separate from others.

E.g. cloudformation.AWSSnsTopic{} becomes resources.AWSSnsTopic{}

## [1.4.1](https://github.com/awslabs/goformation/compare/v1.4.0...v1.4.1) (2019-03-10)


### Bug Fixes

* **spec:** corrected AWS::Serverless::Api.Auth.Authorizers to be of type JSON rather than string  ([#164](https://github.com/awslabs/goformation/issues/164)) ([4cf1bee](https://github.com/awslabs/goformation/commit/4cf1bee))

# [1.4.0](https://github.com/awslabs/goformation/compare/v1.3.0...v1.4.0) (2019-03-09)


### Features

* **parser:** Default to parsing as YAML unless the filename ends in .json ([#176](https://github.com/awslabs/goformation/issues/176)) ([42e7146](https://github.com/awslabs/goformation/commit/42e7146))

# [1.3.0](https://github.com/awslabs/goformation/compare/v1.2.1...v1.3.0) (2019-03-09)


### Bug Fixes

* **CI:** speed up PR builds by only downloading the cfn spec and regenerating resources on cron schedule (not on every build) ([7ae2a32](https://github.com/awslabs/goformation/commit/7ae2a32))
* **CI:** Update TravisCI configuration based on https://github.com/se… ([#180](https://github.com/awslabs/goformation/issues/180)) ([88e1e85](https://github.com/awslabs/goformation/commit/88e1e85))
* **CI:** Update TravisCI configuration for semantic-release to use jobs ([f6c2fee](https://github.com/awslabs/goformation/commit/f6c2fee))


### Features

* Added semantic-release CI setup ([a9b368a](https://github.com/awslabs/goformation/commit/a9b368a))
* Added semantic-release configuration file ([3b25fdb](https://github.com/awslabs/goformation/commit/3b25fdb))
