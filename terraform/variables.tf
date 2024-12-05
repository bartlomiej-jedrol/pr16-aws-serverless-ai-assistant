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

variable "openai_api_key" {
  type = string
}

variable "dub_api_key" {
  type = string
}
