workflow "New workflow" {
  on = "push"
  resolves = ["go t"]
}

action "go t" {
  uses = "docker://golang:1.11"
  args = "test ./..."
}
