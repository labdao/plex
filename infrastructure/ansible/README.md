Ansible is used to provision and configure infrastructure on AWS.

# Prerequisites

1. [Install Anible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) Most likely this is just `python3 -m pip install --user ansible`
1. This setup assumes the ssh key to access all ec2 instances is in ~/.ssh/setward-dev.pem.
1. Set your aws access token environment variables. These variables are used Ansible's dynamic inventory mechanism.
```
export AWS_ACCESS_KEY_ID="anaccesskey"
export AWS_SECRET_ACCESS_KEY="asecretkey"
```

# Usage

Ansilbe configuration consists of playbooks you run using `ansible-playbook [playbook]`. In theory ansible playbooks are supposed to be idempotent and running them multiple times should not be dangerous. In practice this isn't always the case. Know what you are running before you run it.

# Dynamic Inventory

We are using an aws plug in that provides dynamic ec2 inventory. Through configuration in `inventory.aws_ec2.yaml` we tell the dynamic inventory mechanism to group instances together based on AWS instance tags. Setting (for example) a Type tag on your aws instances of the same type (in Terraform) will allow you to target the instances you wish you perform tasks on.

# Playbooks

* `provision_jupyter.yaml` Targets `jupyter_notebook` instances and installs [The Littlest Jupyter Hub](https://tljh.jupyter.org/). It does not do any configuration beyond installation.
* `set_jupyter_users.yaml` Sets the admins, users, and defines access permissions to a team folder 

# The Teams Task

`set_jupyter_users.yaml` has a Create teams task. It loops over a list of hashes of the form to create the juptyter and unix users and set up the home directories. To add a new user simply add a new line to the appropriate team, or add a enw team.

```yaml
team: TEAM_NAME
users:
  - user1
  - user2
  - user3
```

