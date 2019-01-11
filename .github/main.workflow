workflow "New workflow" {
  on = "push"
  resolves = ["go t"]
}

action "go t" {
  uses = "go"
  args = "test ./..."
}
