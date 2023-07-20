resource "aws_autoscaling_group" "labdao_requester_asg" {
  name             = "labdao_requester_asg"
  desired_capacity = 1
  max_size         = 1
  min_size         = 1

  vpc_zone_identifier = data.aws_subnets.default_filtered.ids

  launch_template {
    id      = aws_launch_template.labdao_requester.id
    version = "$Latest"
  }
  target_group_arns = [
    aws_lb_target_group.labdao_requester_bacalhau_tg.arn,
    aws_lb_target_group.labdao_requester_ipfs_tg.arn,
    aws_lb_target_group.labdao_requester_ipfs_swarm_tg.arn,
    aws_lb_target_group.labdao_requester_bacalhau_swarm_tg.arn,
  ]
}

resource "aws_autoscaling_group" "labdao_compute_asg" {
  name             = "labdao_compute_asg"
  desired_capacity = 1
  max_size         = 1
  min_size         = 1

  vpc_zone_identifier = data.aws_subnets.default_filtered.ids

  launch_template {
    id      = aws_launch_template.labdao_compute.id
    version = "$Latest"
  }
}
