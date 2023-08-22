package core

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

func ConnectAws() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(viper.GetString("S3_CREDENTIALS.REGION")),
			Credentials: credentials.NewStaticCredentials(
				viper.GetString("S3_CREDENTIALS.ACCESS_ID"),
				viper.GetString("S3_CREDENTIALS.SECRET_KEY"),
				"", // a token will be created when the session it's used.
			),
			Endpoint: aws.String(viper.GetString("S3_CREDENTIALS.URL")),
		})
	if err != nil {
		panic(err)
	}
	return sess
}
