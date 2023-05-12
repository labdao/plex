Terraform is used to describe and modify our infrastructure on aws.

# Prerequisites

1. [Install terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli).
1. Set your aws access token environment variables. The entity for the keys must have access to the S3 bucket to read and update the state, and access to list and modify any resources that are referenced in the terraform files.
```
export AWS_ACCESS_KEY_ID="anaccesskey"
export AWS_SECRET_ACCESS_KEY="asecretkey"
```

# Getting started

From a command line change directory into the terraform directory and run

```
terraform init -backend-config=backend.conf
```

This initializes your local version of the project and installs any referenced plugins. You only need to run this once, but it is safe to run again.

# Usage

The terraform files should describe the state of the actual infrastructure. To verify this run

```
terraform plan
```

You should see a message saying there is no infrastructure to change.

To update infrastructure modify or add a resource and run

```
terraform plan --out out.tfplan
```

You should see a human readable output of the diff applying this plan will cause to the real configuration.

To apply this change run

***Warning: this will change the real infrastructure in AWS. This command has the power to really break things. Make sure you read the plan before applying it.***

```
terraform apply out.tfplan
```

