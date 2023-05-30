group "default" {
  targets = ["binary"]
}

target "_common" {
  args = {
  }
}

target "binary" {
  inherits = ["_common"]
  target = "binary"
  output = ["bin/release"]
}