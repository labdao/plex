resource "aws_autoscaling_group" "labdao_requester_asg" {
  name             = "labdao_requester_asg"
  desired_capacity = 1
  max_size         = 1
  min_size         = 1

  termination_policies = ["OldestInstance"]

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
  max_size         = 10
  min_size         = 0

  termination_policies = ["OldestInstance"]

  vpc_zone_identifier = data.aws_subnets.default_filtered.ids

  mixed_instances_policy {
    instances_distribution {
      spot_allocation_strategy                 = "lowest-price"
      on_demand_base_capacity                  = 0
      on_demand_percentage_above_base_capacity = 0
      spot_instance_pools                      = 10
    }
    launch_template {
      launch_template_specification {
        launch_template_id = aws_launch_template.labdao_compute.id
        version            = "$Latest"
      }
      override {
        instance_type = "g5.xlarge"
      }
      override {
        instance_type = "g5.2xlarge"
      }
    }
  }
}

# NOTE: autoscaling to stop instances at Friday 8pm EST
resource "aws_autoscaling_schedule" "labdao_compute_asg_schedule_0" {
  scheduled_action_name  = "labdao-${var.environment}-compute-asg-count-0"
  autoscaling_group_name = aws_autoscaling_group.labdao_compute_asg.name
  recurrence             = "00 20 * * FRI"
  time_zone              = "America/Toronto"

  # NOT Adjusting
  min_size = -1

  # NOT Adjusting
  max_size = -1

  # NOTE: Dropping to 0
  desired_capacity = 0
}

# NOTE: autoscaling to start single instance on Monday 8am CEST
resource "aws_autoscaling_schedule" "labdao_compute_asg_schedule_1" {
  scheduled_action_name  = "labdao-${var.environment}-compute-asg-count-1"
  autoscaling_group_name = aws_autoscaling_group.labdao_compute_asg.name
  recurrence             = "00 8 * * MON"
  time_zone              = "Europe/Berlin"

  # NOT Adjusting
  min_size = -1

  # NOT Adjusting
  max_size = -1

  # NOTE: Upping to 1
  desired_capacity = 1
}
