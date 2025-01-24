package convertapi

import (
	convertapi "github.com/ConvertAPI/convertapi-go/pkg"
	"github.com/ConvertAPI/convertapi-go/pkg/config"
	"github.com/ConvertAPI/convertapi-go/pkg/param"
	iAws "github.com/bartlomiej-jedrol/go-toolkit/aws"
	iLog "github.com/bartlomiej-jedrol/go-toolkit/log"
	"github.com/bartlomiej-jedrol/pr16-aws-serverless-ai-assistant/configuration"
)

func ConvertPDFToJPG(srcPath, dstPath string) ([]string, error) {
	function := "ConvertPDFToJPG"
	apiKey, err := iAws.GetEnvironmentVariable("CONVERT_API_KEY")
	if err != nil {
		iLog.Error("failed to load env var", "CONVERT_API_KEY", err, configuration.ServiceName, function)
	}

	config.Default = config.NewDefault(apiKey)
	result, errs := convertapi.ConvDef("pdf", "jpg",
		param.NewPath("File", srcPath, nil),
		param.NewString("FileName", "converted_file"),
		param.NewString("ImageResolution", "800"),
		param.NewString("ImageOutputFormat", "jpg"),
	).ToPath(dstPath)

	if err != nil {
		return nil, errs[0]
	}

	var outputFiles []string
	for _, file := range result {
		outputFiles = append(outputFiles, file.Name())
	}

	return outputFiles, nil
}
