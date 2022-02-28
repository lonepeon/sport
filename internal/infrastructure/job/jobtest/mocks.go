package jobtest

//go:generate go run ../../../../vendor/github.com/golang/mock/mockgen/ -destination=job.go -package jobtest github.com/lonepeon/sport/internal/infrastructure/job Enqueuer
