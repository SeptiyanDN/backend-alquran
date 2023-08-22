package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetUrlObject(fileName string, sess *session.Session) string {

	svc := s3.New(sess)

	// List objects in the bucket with the specified prefix
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("staging-smartpatrol"),
		Key:    aws.String(fileName),
	})

	urlStr, err := req.Presign(24 * time.Hour) // Sesuaikan waktu kadaluarsa URL
	if err != nil {
		// Handle error
		fmt.Println("Failed to generate URL", err)
	}

	// Hapus parameter URL yang tidak diperlukan
	parsedURL, _ := url.Parse(urlStr)
	parsedURL.RawQuery = ""
	finalURL := parsedURL.String()

	// Gunakan finalURL sesuai kebutuhan Anda

	return finalURL
}

func SaveObjectToS3(sess *session.Session, bucketName, keyName string, file io.Reader) error {
	svc := s3.New(sess)

	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(file)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(keyName),
		Body:        bytes.NewReader(buf.Bytes()),
		ACL:         aws.String("public-read"), // Menambahkan parameter ACL untuk set permission public
		ContentType: aws.String(http.DetectContentType(buf.Bytes())),
		Metadata: map[string]*string{
			"Content-Disposition": aws.String("inline"),
		},
	}

	_, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})

	// if err != nil {
	// 	if reqErr, ok := err.(awserr.RequestFailure); ok && reqErr.StatusCode() == 404 {
	// 		_, err := svc.PutObject(&s3.PutObjectInput{
	// 			Bucket:      aws.String(bucketName),
	// 			Key:         aws.String("/" + keyName),
	// 			ACL:         aws.String("public-read"),
	// 			Body:        bytes.NewReader(buf.Bytes()),
	// 			ContentType: aws.String(http.DetectContentType(buf.Bytes())),
	// 			Metadata: map[string]*string{
	// 				"Content-Disposition": aws.String("inline"),
	// 			}, // Menambahkan parameter ACL untuk set permission public
	// 		})

	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		return err
	// 	}
	// }

	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}
