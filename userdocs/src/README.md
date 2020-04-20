
# Writing and publishing user docs

The user docs are written in [MkDocs](https://www.mkdocs.org/) using [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/). Install `mkdocs` and then you can do the following:

While writing the docs, to preview locally:

```console
 $ mkdocs serve
INFO    -  Building documentation...
INFO    -  Cleaning site directory
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - index.md
[I 191218 10:16:30 server:296] Serving on http://127.0.0.1:8000
[I 191218 10:16:30 handlers:62] Start watching changes
[I 191218 10:16:30 handlers:64] Start detecting changes
[I 191218 10:16:33 handlers:135] Browser Connected: http://127.0.0.1:8000/intro/
...
```

When you're done and have committed the changes, then you can publish it:

```console
$ mkdocs gh-deploy
INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /Users/hausenbl/Documents/repos/eksctl/site/site
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - index.md
WARNING -  Version check skipped: No version specificed in previous deployment.
INFO    -  Copying '/Users/hausenbl/Documents/repos/eksctl/site/site' to 'gh-pages' branch and pushing to GitHub.
INFO    -  Your documentation should shortly be available at: https://mhausenblas.github.io/eksctl/
```
