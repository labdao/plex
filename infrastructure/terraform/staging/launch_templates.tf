resource "aws_launch_template" "labdao_requester" {
  name = "labdao-${var.environment}-public-lt"

  credit_specification {
    cpu_credits = "standard"
  }

  image_id = data.aws_ami.ubuntu_latest.id

  instance_type = var.main_instance_type

  key_name = var.ssh_key

  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      volume_size = 100
      volume_type = "gp3"

    }
  }

  monitoring {
    enabled = true
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.labdao_public_bacalhau.id, aws_security_group.labdao_public_ipfs.id, aws_security_group.labdao_public_ssh.id, aws_security_group.labdao_egress_all.id, aws_security_group.labdao_private.id]
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = "labdao-${var.environment}-requester"
      Env  = "${var.environment}"
      Type = "requester"
    }
  }
  user_data = base64encode(templatefile("${path.module}/files/userdata.sh", { environment = var.environment }))
}

resource "aws_launch_template" "labdao_compute" {
  name = "labdao-${var.environment}-compute-lt"

  credit_specification {
    cpu_credits = "standard"
  }

  image_id = data.aws_ami.ubuntu_latest.id

  instance_type = var.compute_instance_type

  key_name = var.ssh_key

  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      volume_size = 1000
      volume_type = "gp3"
    }
  }

  monitoring {
    enabled = true
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.labdao_public_ssh.id, aws_security_group.labdao_egress_all.id, aws_security_group.labdao_private.id]
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = "labdao-${var.environment}-compute"
      Env  = "${var.environment}"
      Type = "compute_only"
    }
  }

  user_data = base64encode(templatefile("${path.module}/files/userdata.sh", { environment = var.environment }))
}

resource "aws_launch_template" "labdao_ipfs" {
  name = "labdao-${var.environment}-ipfs-lt"

  credit_specification {
    cpu_credits = "standard"
  }

  image_id = data.aws_ami.ubuntu_latest.id

  instance_type = var.main_instance_type

  key_name = var.ssh_key

  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

  block_device_mappings {
    device_name = "/dev/sda1"

    ebs {
      volume_size = 500
      volume_type = "gp3"
    }
  }

  monitoring {
    enabled = true
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.labdao_public_ssh.id, aws_security_group.labdao_egress_all.id, aws_security_group.labdao_private.id]
  }

  tag_specifications {
    resource_type = "instance"

    tags = {
      Name = "labdao-${var.environment}-ipfs"
      Env  = "${var.environment}"
      Type = "ipfs"
    }
  }

  user_data = base64encode(templatefile("${path.module}/files/userdata.sh", { environment = var.environment }))
}
