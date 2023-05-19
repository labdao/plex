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

resource "aws_eip" "plex_prod" {
  instance = aws_instance.plex_prod.id
  vpc      = true

  tags = {
    Name = "plex-prod-gateway"
  }
}

resource "aws_instance" "plex_jupyter" {
  for_each      = toset(["t6"])
  ami           = var.ami_main
  instance_type = "t3.micro"

  vpc_security_group_ids = [aws_security_group.plex.id]
  key_name               = var.key_main
  availability_zone      = var.availability_zones[0]

  root_block_device {
    volume_size = 10
  }

  tags = {
    Name        = "plex-jupyter-${each.key}"
    InstanceKey = each.key
    Type        = "jupyter_notebook"
  }

}

resource "aws_eip" "plex_jupyter" {
  instance = aws_instance.plex_jupyter["t6"].id
  vpc      = true

  tags = {
    Name = "plex-jupyter-eip"
  }
}

resource "cloudflare_record" "jupyter" {
  zone_id = var.cloudflare_zone_id
  name    = "jupyter"
  value   = aws_eip.plex_jupyter.public_dns
  type    = "CNAME"
  ttl     = 3600
}

