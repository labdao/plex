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

resource "aws_instance" "ops_test" {
  for_each      = toset(["opstest"])
  ami           = var.ami_main
  instance_type = "g5.xlarge"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  root_block_device {
    volume_size = 1000
    tags = {
      Name = "plex-prod"
    }
  }

  tags = {
    Name        = "plex-prod-${each.key}"
    InstanceKey = each.key
    Type        = "compute"
  }
}


resource "aws_eip" "plex_prod" {
  instance = aws_instance.plex_prod.id
  vpc      = true

  tags = {
    Name = "plex-prod-gateway"
  }
}
