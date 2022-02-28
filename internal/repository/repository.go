package repository

import (
	"context"
	"io"

	"github.com/lonepeon/sport/internal/domain"
)

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	GetRunningActivity(context.Context, domain.RunningActivitySlug) (domain.RunningActivity, error)
	ListRunningActivities(context.Context) ([]domain.RunningActivity, error)
}

type Writer interface {
	AnnotateMapWithStats(context.Context, domain.MapFile, domain.Distance, domain.Speed) (domain.ShareableMapFile, error)
	CleanGPXFile(context.Context, io.Reader) (domain.GPXFile, error)
	GenerateMap(context.Context, domain.GPXFile) (domain.MapFile, error)
	DeleteRunningActivity(context.Context, domain.RunningActivitySlug) error
	RecordRunningActivity(context.Context, domain.RunningActivity) error
	StoreAsset(content io.Reader, fileName string) error
	DeleteAsset(fileName string) error
}
