package bloodresults

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
	convertapi "github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/convert-api"
)

type Parser struct {
	s3Client      *s3.Client
	s3Bucket      string
	lambdaTmpPath string
}

func (p *Parser) New(s3Client *s3.Client, s3Bucket, lambdaTmpPath string) *Parser {
	return &Parser{
		s3Client:      s3Client,
		s3Bucket:      s3Bucket,
		lambdaTmpPath: lambdaTmpPath,
	}
}

func (p *Parser) Parse(ctx context.Context, s3ObjKey string) error {
	function := "Parse"
	s3OjectBody, err := iAws.GetS3Object(ctx, p.s3Client, p.s3Bucket, s3ObjKey)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(p.lambdaTmpPath, "blood-results-*.pdf")
	if err != nil {
		iLog.Error("failed to create pdf file", nil, err, configuration.ServiceName, function)
		return err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, *s3OjectBody)
	if err != nil {
		iLog.Error("failed to save pdf file to tmp path", nil, err, configuration.ServiceName, function)
		return err
	}
	iLog.Info("pdf saved to", tmpFile.Name(), nil, configuration.ServiceName, function)

	outputFiles, err := convertapi.ConvertPDFToJPG(tmpFile.Name(), p.lambdaTmpPath)
	if err != nil {
		iLog.Error("failed to convert PDF to JPG", err, nil, configuration.ServiceName, function)
		return err
	}

	for _, file := range outputFiles {
		defer os.Remove(file)
		iLog.Info("converted file", file, nil, configuration.ServiceName, function)
	}
	return nil
}
