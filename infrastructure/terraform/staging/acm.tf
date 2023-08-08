resource "aws_acm_certificate" "wildcard_labdao" {
  domain_name               = "bacalhau.${var.environment}.${var.domain}"
  subject_alternative_names = ["*.${var.environment}.${var.domain}"]
  validation_method         = "DNS"
}
