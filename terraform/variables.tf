variable "aws_region" {
  type = string
}

# ========== Lambda ==========
variable "assistant_lambda_name" {
  type    = string
  default = "pr16-assistant"
}

variable "assistant_lambda_role" {
  type    = string
  default = "pr16-assistant-lambda-role"
}

variable "assistant_lambda_policy" {
  type    = string
  default = "pr16-assistant-lambda-policy"
}

# ========== S3 buckets ==========
variable "assistant_bucket_name" {
  type    = string
  default = "pr16-assistant-bucket"
}

# ========== Lambda environment variables ==========
variable "db_host" {
  type = string
}

variable "db_user" {
  type = string
}

variable "db_password" {
  type = string
}

variable "db_name" {
  type = string
}

variable "assistant_api_key" {
  type = string
}

variable "openai_api_key" {
  type = string
}

variable "dub_api_key" {
  type = string
}

variable "airtable_api_key" {
  type = string
}
