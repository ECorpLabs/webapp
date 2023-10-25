aws_region = "us-east-1"
subnet_id = "subnet-071040365a5f7a5f8"
ami_users = ["042793801071"]
instance_type = "t2.micro"
profile = "default"
volume_size = 25
device_name = "/dev/xvda"
volume_type = "gp2"
file_paths = ["./webapp", "./data/users.csv","scripts/system.service"]
script_path = "./scripts/setup.sh"
destination_path = "/tmp/"
ssh_username = "admin"
