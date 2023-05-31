variable "APT_MIRROR" {
  default = "cdn-fastly.deb.debian.org"
#  default = "repo.huaweicloud.com"
}

group "default" {
  targets = ["binary"]
}

target "_common" {
  args = {
    APT_MIRROR = APT_MIRROR
  }
}

target "binary" {
  inherits = ["_common"]
  target = "binary"
  output = ["bin/release"]
}