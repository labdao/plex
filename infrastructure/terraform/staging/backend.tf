terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 3.0"
    }
  }
  required_version = ">= 1.4.6"

  # Backend
  backend "s3" {
    bucket = "labdao-infrastructure"
    region = "us-east-1"
    key    = "staging/terraform-state"
  }
}
