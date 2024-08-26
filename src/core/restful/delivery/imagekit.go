package delivery

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/dwprz/prasorganic-product-service/src/common/log"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/cbreaker"
	"github.com/dwprz/prasorganic-product-service/src/infrastructure/imagekit"
	"github.com/dwprz/prasorganic-product-service/src/interface/delivery"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"github.com/sirupsen/logrus"
)

type ImageKitRESTful struct{}

func NewImageKit() delivery.ImageKitRESTful {
	return &ImageKitRESTful{}
}

func (i *ImageKitRESTful) UploadImage(ctx context.Context, path string, filename string) (*uploader.UploadResult, error) {
	res, err := cbreaker.ImageKit.Execute(func() (any, error) {
		fileData, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		base64String := base64.StdEncoding.EncodeToString(fileData)
		file := "data:image/jpeg;base64," + base64String

		useUniqueFileName := false

		res, err := imagekit.IK.Uploader.Upload(ctx, file, uploader.UploadParam{
			FileName:          filename,
			UseUniqueFileName: &useUniqueFileName,
		})

		return &res.Data, err
	})

	if err != nil {
		return nil, err
	}

	result, ok := res.(*uploader.UploadResult)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T (uploader.UploadResult)", result)
	}

	return result, err
}

func (i *ImageKitRESTful) DeleteFile(ctx context.Context, fileId string) {
	_, err := cbreaker.ImageKit.Execute(func() (any, error) {
		_, err := imagekit.IK.Media.DeleteFile(ctx, fileId)
		return nil, err
	})

	if err != nil {
		log.Logger.WithFields(logrus.Fields{"location": "client.ImageKitRESTful/DeleteFile", "section": "ik.Media.DeleteFile"}).Error(err)
	}
}
