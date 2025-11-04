package domain

import "context"

type ImageCommander interface {
	CreateImage (ctx context.Context, url string) (*Image, error)
}