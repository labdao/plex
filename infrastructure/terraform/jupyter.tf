resource "aws_instance" "plex_jupyter" {
  for_each      = toset(["littlehub"])
  ami           = var.ami_main
  instance_type = "t3.xlarge"

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
  }

  tags = {
    Name        = "plex-jupyter-${each.key}"
    InstanceKey = each.key
    Type        = "jupyter_notebook"
  }

}

resource "aws_eip" "plex_jupyter" {
  instance = aws_instance.plex_jupyter["littlehub"].id
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
