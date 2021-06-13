package imageService

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"google-cloud-task-processor/config"
	"google-cloud-task-processor/dto"
	"google.golang.org/api/option"
	"io"
	"net/url"
	"os"
	"strings"
	"time"
)

type ImageService struct {
	config *config.Config
}

func New(config *config.Config) *ImageService {
	return &ImageService{config: config}
}

func (s *ImageService) UploadImage(imageRequest *dto.ImageRequest) (*dto.ImageResponse, error) {
	img, err := s.decodeImage(imageRequest)
	if err != nil {
		return nil, err
	}

	f, err := s.createFile()
	if err != nil {
		return nil, err
	}
	err = s.writeImageToFile(f, img)
	if err != nil {
		return nil, err
	}
	name, err := s.uploadFileToBucket(f)
	if err != nil {
		return nil, err
	}
	err = os.Remove(f.Name())
	if err != nil {
		return nil, err
	}
	u, err := url.Parse("/" + s.config.Storage.Bucket + "/" + name)
	if err != nil {
		return nil, err
	}

	return &dto.ImageResponse{Name: imageRequest.Name, Url: s.config.Storage.Url + u.EscapedPath()}, nil
}

func (s *ImageService) decodeImage(imageRequest *dto.ImageRequest) ([]byte, error) {
	b64data := imageRequest.Payload[strings.IndexByte(imageRequest.Payload, ',')+1:]
	img, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (s *ImageService) createFile() (*os.File, error) {
	newUUID, _ := uuid.NewUUID()
	filename := fmt.Sprintf("%s.jpg", newUUID)
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *ImageService) writeImageToFile(f *os.File, img []byte) error {
	if _, err := f.Write(img); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	err := f.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageService) uploadFileToBucket(f *os.File) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("secrets/google-secret.json"))
	if err != nil {
		return "", err
	}
	f, err = os.Open(f.Name())
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(s.config.Storage.Bucket).Object(f.Name()).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	err = client.Close()
	if err != nil {
		return "", err
	}

	return wc.Attrs().Name, nil
}
