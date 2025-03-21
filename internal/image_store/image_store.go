package imagestore

import "context"

type ImageStore interface {
	UploadImage(context context.Context, imageData string) (string, error)
}
