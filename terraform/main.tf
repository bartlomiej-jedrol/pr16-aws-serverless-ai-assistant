terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket  = "bj-terraform-states"
    key     = "pr16-aws-serverless-ai-assistant/terraform.tfstate"
    region  = "eu-central-1"
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region
}

# ========== Assistant Lambda ==========
resource "aws_lambda_function" "assistant_lambda" {
  function_name = var.assistant_lambda_name
  handler       = "main"
  runtime       = "provided.al2023"
  filename      = "../bin/assistant.zip"
  timeout       = 10

  role = aws_iam_role.assistant_lambda_role.arn

  environment {
    variables = {
      OPENAI_API_KEY = var.openai_api_key
      DUB_API_KEY    = var.dub_api_key
    }
  }
}

# Assistant Lambda Role
resource "aws_iam_role" "assistant_lambda_role" {
  name = var.assistant_lambda_role
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}