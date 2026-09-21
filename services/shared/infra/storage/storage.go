package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Storage struct {
	cloudinary *cloudinary.Cloudinary
}

func NewStorage(c *cloudinary.Cloudinary) *Storage {
	return &Storage{cloudinary: c}
}

// NOTE: Old and new names must include the file's location. Ex: temp/images/filename
func (s *Storage) Rename(ctx context.Context, oldName, newName string, overwrite, invalidate bool) (*uploader.RenameResult, error) {
	return s.cloudinary.Upload.Rename(
		ctx,
		uploader.RenameParams{
			FromPublicID: oldName,
			ToPublicID:   newName,
			Overwrite:    &overwrite,
			Invalidate:   &invalidate,
		},
	)
}

func (s *Storage) ExtractPublicID(rawURL, expectedPrefix string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsed.Host != "res.cloudinary.com" {
		return "", fmt.Errorf("unauthorized domain: %s", parsed.Host)
	}

	// Path Ex.		: "/{cloud}/image/upload/v12345678/temp/images/filename.jpg"
	// - before		:	"/{cloud}/image"
	// - after		: "v12345678/temp/images/filename.jpg"
	_, after, found := strings.Cut(parsed.Path, "/upload/")
	if !found {
		return "", errors.New("invalid cloudinary path structure")
	}
	if parts := strings.SplitN(after, "/", 2); len(parts) == 2 && s.isVersionSegment(parts[0]) {
		after = parts[1] // "temp/images/filename.jpg"
	}

	publicID := after
	if lastDot := strings.LastIndex(after, "."); lastDot != -1 && lastDot > strings.LastIndex(after, "/") {
		publicID = after[:lastDot] // "temp/images/filename"
	}
	if !strings.HasPrefix(publicID, expectedPrefix) {
		return "", fmt.Errorf("unauthorized file source: must be in %s", expectedPrefix)
	}

	return publicID, nil
}

func (s *Storage) isVersionSegment(value string) bool {
	if len(value) < 2 || value[0] != 'v' {
		return false
	}
	for i := 1; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return false
		}
	}
	return true
}
