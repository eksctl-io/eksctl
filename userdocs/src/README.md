
# Writing and publishing user docs

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
[I 191218 10:16:33 handlers:135] Browser Connected: http://127.0.0.1:8000/intro/
...
```
