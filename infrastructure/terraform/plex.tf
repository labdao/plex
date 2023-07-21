resource "aws_instance" "plex_prod" {
  ami           = var.ami_main
  instance_type = "g5.2xlarge"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  # Enabling metadata option with instance metadata tags - required for self bootstrapping
  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

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

  vpc_security_group_ids = [aws_security_group.plex.id, aws_security_group.internal.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  # Enabling metadata option with instance metadata tags - required for self bootstrapping
  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

  root_block_device {
    volume_size = 2000
    tags = {
      Name = "plex-prod-${each.key}"
    }
  }

  tags = {
    Name        = "plex-prod-${each.key}"
    Env         = "prod"
    InstanceKey = each.key
    Type        = "compute"
  }
}

resource "aws_instance" "plex_compute_only" {
  for_each = toset([
    "compute_only_1"
  ])
  ami           = "ami-053b0d53c279acc90"
  instance_type = "g5.2xlarge"

  vpc_security_group_ids = [aws_security_group.plex.id, aws_security_group.internal.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  # Enabling metadata option with instance metadata tags - required for self bootstrapping
  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

  root_block_device {
    volume_size = 1000
    tags = {
      Name = "plex-prod-${each.key}"
    }
  }

  tags = {
    Name        = "plex-compute-only-${each.key}"
    Env         = "prod"
    InstanceKey = each.key
    Type        = "compute_only"
  }
}

resource "aws_instance" "plex_requester" {
  ami           = "ami-053b0d53c279acc90"
  instance_type = "m5.large"

  vpc_security_group_ids = [aws_security_group.plex.id, aws_security_group.internal.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  # Enabling metadata option with instance metadata tags - required for self bootstrapping
  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

  root_block_device {
    volume_size = 10
  }

  tags = {
    Name = "plex-requester-prod"
    Env  = "prod"
    Type = "requester"
  }
}

resource "aws_eip" "plex_prod" {
  instance = aws_instance.plex_requester.id
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

# Private DNS record for requester 
resource "cloudflare_record" "plex_compute_prod_private" {
  zone_id = var.cloudflare_zone_id
  name    = "requester"
  value   = aws_eip.plex_prod.private_dns
  type    = "CNAME"
  ttl     = 3600
}

resource "aws_instance" "receptor" {
  for_each      = toset(["judgy"])
  ami           = "ami-053b0d53c279acc90"
  instance_type = "t3.small"

  vpc_security_group_ids = [aws_security_group.external_ssh.id, aws_security_group.internal.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  # Enabling metadata option with instance metadata tags - required for self bootstrapping
  metadata_options {
    http_endpoint          = "enabled"
    http_tokens            = "optional"
    instance_metadata_tags = "enabled"
  }

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
    Env         = "dev"
  }
}

resource "aws_db_instance" "default" {
  allocated_storage           = 10
  db_name                     = "receptor"
  engine                      = "postgres"
  engine_version              = "15.3"
  instance_class              = "db.t3.small"
  username                    = "receptor"
  skip_final_snapshot         = true
  manage_master_user_password = true
  vpc_security_group_ids      = [aws_security_group.internal.id, ]
}
