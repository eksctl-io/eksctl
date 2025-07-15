# Important - Content Migration to AWS Docs

The eksctl docs are live on the AWS Docs platform

* New site: https://docs.aws.amazon.com/eks/latest/eksctl/what-is-eksctl.html
* New GitHub source: https://github.com/eksctl-io/eksctl-docs
* For information about editing the new repo in AsciiDoc format, see https://docs.aws.amazon.com/eks/latest/userguide/contribute.html


This directory is for the `eksctl.io` site. 

The `eksctl.io` site is still live, but in the future it will redirect to the AWS Docs.


# Writing and publishing docs on eksctl.io 

As of July 2025, the `eksctl.io` domain is owned by AWS Legal. The DNS has been migrated to Route53. The site hosting has been migrated to AWS Amplify. 

The user docs are written in [MkDocs](https://www.mkdocs.org/) using [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/). Install `mkdocs` and then you can do the following:

While writing the docs, to preview locally:

```console
 $ make serve-pages
INFO    -  Building documentation...
INFO    -  Cleaning site directory
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - index.md
[I 191218 10:16:30 server:296] Serving on http://127.0.0.1:8000
[I 191218 10:16:30 handlers:62] Start watching changes
[I 191218 10:16:30 handlers:64] Start detecting changes
[I 191218 10:16:33 handlers:135] Browser Connected: http://127.0.0.1:8000/introduction/
...
```
