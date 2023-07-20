provider "aws" {
  default_tags {
    tags = {
      Env     = "staging"
      Owner   = "Ops"
      Project = "LabDAO"
    }
  }
  region = var.aws_region
}
