variable "region" {
  type        = string
  description = "The region of project"
  default     = "ap-northeast-1"
}

variable "border_security_group_id" {
  type        = string
  description = "The security group ID of the border"
  default     = "sg-0f9a124088814d296"
}
