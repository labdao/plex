variable "ami_main" {
  type    = string
  default = "ami-09cd747c78a9add63"
}

variable "key_main" {
  type    = string
  default = "steward-dev"
}

variable "availability_zones" {
  type    = list(string)
  default = ["us-east-1c"]
}
