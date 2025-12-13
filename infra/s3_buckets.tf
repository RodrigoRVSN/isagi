terraform {
  required_version = "~> 1.14.2"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.33.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}
