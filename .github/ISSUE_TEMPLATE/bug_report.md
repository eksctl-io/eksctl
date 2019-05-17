---
name: Bug report
about: Create a bug report to help us improve eksctl
title: ''
labels: kind/bug
assignees: ''

---

**What happened?**
A description of actual behavior (with error messages).

**What you expected to happen?**
A clear and concise description of what the bug is.

**How to reproduce it?**
Include the steps to reproduce the bug

**Anything else we need to know?**
What OS are you using, are you using a downloaded binary or did you compile eksctl, what type of AWS credentials are you using (i.e. default/named profile, MFA) - please don't include actual credentials though!

**Versions**
Please paste in the output of these commands:
```
$ eksctl version
$ uname -a
$ kubectl version
```

**Logs**
Include the output of the command line when running eksctl. If possible, eksctl should be run with debug logs. For example:
`eksctl get clusters -v 4`
Make sure you redact any sensitive information before posting.
If the output is long, please consider a Gist.
