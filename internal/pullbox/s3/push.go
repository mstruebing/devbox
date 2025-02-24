// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package s3

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.jetpack.io/devbox/internal/auth"
	"go.jetpack.io/devbox/internal/pullbox/tar"
	"go.jetpack.io/devbox/internal/ux"
)

func Push(ctx context.Context, user *auth.User, dir, profile string) error {
	archivePath, err := tar.Compress(dir)
	if err != nil {
		return err
	}

	config, err := assumeRole(ctx, user)
	if err != nil {
		return err
	}

	s3Client := manager.NewUploader(s3.NewFromConfig(*config))
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}

	_, err = s3Client.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(
			fmt.Sprintf(
				"profiles/%s/%s.tar.gz",
				user.ID(),
				profile,
			),
		),
		Body: io.Reader(file),
	})

	if err != nil {
		return err
	}

	ux.Fsuccess(
		os.Stderr,
		"Profile successfully pushed (profile: %s)\n",
		profile,
	)

	return nil
}
