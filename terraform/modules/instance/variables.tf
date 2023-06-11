variable "ami" {
  type        = string
  description = "AMI to use for the instance"
}

variable "instance_type" {
  type        = string
  description = "Instance type to use for the instance"
  default     = "t2.micro"
}

variable "subnet_id" {
  type        = string
  description = "VPC subnet ID to launch in"
}

variable "key_name" {
  type        = string
  description = "Key name of the key pair to use for the instance"
}

variable "vpc_security_group_ids" {
  type        = list(string)
  description = "List of the security groups IDs to associate with"
}

variable "disable_api_termination" {
  type        = bool
  description = "If true, enables EC2 termination protection"
  default     = true
}

variable "monitoring" {
  type        = bool
  description = "If true, the launched EC2 instance will have detailed monitoring enable"
  default     = true
}

variable "user_data" {
  type        = string
  description = "Provide a command to a instance at launch time"
}

variable "volume_size" {
  type        = string
  description = "Size of the volume in gibibytes"
}

variable "volume_type" {
  type        = string
  description = "Type of volume"
}

variable "tags" {
  type        = map(string)
  description = "Map of tags to assign to the resource"
}
