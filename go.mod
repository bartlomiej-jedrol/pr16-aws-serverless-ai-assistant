module github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant

go 1.22.4

replace github.com/bartlomiej-jedrol/go-toolkit => ../go-toolkit

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/bartlomiej-jedrol/go-toolkit v0.0.0-00010101000000-000000000000
)

require (
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
)
