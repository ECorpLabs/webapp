packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "subnet_id" {
  type    = string
  default = "subnet-071040365a5f7a5f8"
}

variable ami_users {
  type    = list(string)
  default = ["042793801071"]
}
variable "instance_type" {
  type    = string
  default = "t2.micro"
}
variable "profile" {
  type    = string
  default = "default"
}


source "amazon-ebs" "debian12" {
  ami_name      = "ECorp-debian-12-ami_${formatdate("YYYY-MM-DD-hh-mm-ss", timestamp())}"
  ami_users     = "${var.ami_users}"
  profile       = "${var.profile}"
  instance_type = "${var.instance_type}"
  region        = "${var.aws_region}"
  source_ami_filter {
    filters = {
      name                = "debian-12-amd64-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["amazon"]
  }
  subnet_id    = "${var.subnet_id}"
  ssh_username = "admin"

  launch_block_device_mappings {
    device_name           = "/dev/xvda"
    delete_on_termination = true
    volume_size           = 25
    volume_type           = "gp2"
  }
}

build {
  sources = ["source.amazon-ebs.debian12"]

  provisioner "file" {
    sources = [
      "./webapp",
      "./data/users.csv"
    ]
    destination = "/tmp/"
  }

  provisioner "shell" {

    script = "./scripts/setup.sh"
  }
}
