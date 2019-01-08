workflow "Build every commit" {
  on = "push"
  resolves = ["build and test"]
}

action "build and test" {
  uses = "./build"
  runs = "make"
  args = ["github-build"]
}
