variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "main_instance_type" {
  type    = string
  default = "t2.micro"
}

variable "compute_instance_type" {
  type    = string
  default = "g5.xlarge"
}

variable "ssh_key" {
  type    = string
  default = "steward-dev"
}

variable "environment" {
  type    = string
  default = "staging"
}

variable "cloudflare_zone_id" {
  type    = string
  default = "858fe9f16ace6df3deefd366cb7defd6"
}

variable "availability_zones" {
  type    = list(string)
  default = ["us-east-1c", "us-east-1d"]
}
