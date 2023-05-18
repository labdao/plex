resource "aws_security_group" "plex" {
  name = "dev-web"
  description = "SSH with key and HTTP open"
}
