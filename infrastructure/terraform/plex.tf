resource "aws_instance" "plex_prod" {
  ami           = var.ami_main
  instance_type = "g5.2xlarge"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  root_block_device {
    volume_size = 1000
    tags = {
      Name = "keep-plex-prod"
    }
  }

  tags = {
    Name = "plex-prod"
  }
}

resource "aws_instance" "plex_compute_prod" {
  for_each      = toset(["compute1"])
  ami           = "ami-053b0d53c279acc90"
  instance_type = "g5.2xlarge"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  root_block_device {
    volume_size = 1000
    tags = {
      Name = "plex-prod-${each.key}"
    }
  }

  tags = {
    Name        = "plex-prod-${each.key}"
    InstanceKey = each.key
    Type        = "compute"
  }
}


resource "aws_eip" "plex_prod" {
  instance = aws_instance.plex_compute_prod["compute1"].id
  vpc      = true

  tags = {
    Name = "plex-prod-gateway"
  }
}

resource "cloudflare_record" "plex_compute_prod" {
  zone_id = var.cloudflare_zone_id
  name    = "bacalhau"
  value   = aws_eip.plex_prod.public_dns
  type    = "CNAME"
  ttl     = 3600
}

resource "aws_instance" "plex_compute_dev" {
  for_each      = toset(["compute-dev"])
  ami           = "ami-053b0d53c279acc90"
  instance_type = "m5.large"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  root_block_device {
    volume_size = 10
    tags = {
      Name = "plex-prod-${each.key}"
    }
  }

  tags = {
    Name        = "plex-prod-${each.key}"
    InstanceKey = each.key
    Type        = "compute"
    Env         = "dev"
  }
}

resource "aws_instance" "receptor" {
  for_each      = toset(["judgy"])
  ami           = "ami-053b0d53c279acc90"
  instance_type = "t3.small"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  root_block_device {
    volume_size = 10
    tags = {
      Name = "plex-receptor-${each.key}"
    }
  }

  tags = {
    Name        = "plex-receptor-${each.key}"
    InstanceKey = each.key
    Type        = "receptor"
  }
}
