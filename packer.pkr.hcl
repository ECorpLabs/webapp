packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "aws_region" {}
variable "subnet_id" {}
variable ami_users {}
variable "instance_type" {}
variable "profile" {}
variable "volume_size" {}
variable device_name {}
variable volume_type {}
variable script_path {}
variable "file_paths" {}
variable destination_path {}
variable ssh_username {}

source "amazon-ebs" "debian12" {
  ami_name      = "ECorp-debian-12-ami_${formatdate("YYYY-MM-DD-hh-mm-ss", timestamp())}"
  ami_users     = var.ami_users
  profile       = var.profile
  instance_type = var.instance_type
  region        = var.aws_region
  source_ami_filter {
    filters = {
      name                = "debian-12-amd64-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["amazon"]
  }
  subnet_id    = var.subnet_id
  ssh_username = var.ssh_username

  launch_block_device_mappings {
    device_name           = var.device_name
    delete_on_termination = true
    volume_size           = var.volume_size
    volume_type           = var.volume_type
  }
}

build {
  sources = ["source.amazon-ebs.debian12"]

  provisioner "file" {
    sources     = var.file_paths
    destination = var.destination_path
  }

  provisioner "shell" {
    script = var.script_path
  }
}
