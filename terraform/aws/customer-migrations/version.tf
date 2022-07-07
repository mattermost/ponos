terraform {
  required_version = "~> 1.1.7"
  backend "s3" {
    region = "us-east-1"
  }
  required_providers {
    aws = "~> 4.8.0"
  }
}
