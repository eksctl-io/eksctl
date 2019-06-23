---
title: "Installation"
weight: 20
---
### Shell Completion

To enable bash completion, run the following, or put it in `~/.bashrc` or `~/.profile`:
```
. <(eksctl completion bash)
```

Or for zsh, run:
```
mkdir -p ~/.zsh/completion/
eksctl completion zsh > ~/.zsh/completion/_eksctl
```
and put the following in `~/.zshrc`:
```
fpath=($fpath ~/.zsh/completion)
```
Note if you're not running a distribution like oh-my-zsh you may first have to enable autocompletion:
```
autoload -U compinit
compinit
```

To make the above persistent, run the first two lines, and put the above in `~/.zshrc`.
