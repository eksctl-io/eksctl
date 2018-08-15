---
name: Bug report
about: Create a bug report to help us improve eksctl

---

**What happened?**
A description of actual behvaior (with error messages).

**What you expected to happen?**
A clear and concise description of what the bug is.

**How to repoduce it?**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Anything else we need to know?**
Hardware, OS and screenshots (if applicable)

**Versions**
Please paste in the output of these commands:
```
$ eksctl version
$ uname -a
$ kubectl version
```
Also include your version of `heptio-authenticator-aws`

**Logs**
Include the output of the command line when running eksctl. If possible, eksctl should be run with debug logs. For example:
`eksctl get clusters -v 4`
Make sure you redact any sensitive information before posting.
If the output is long, please consider a Gist.
