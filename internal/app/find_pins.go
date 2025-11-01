package app

import "context"

const FindPingsQueryName = "findPingsQueryName"

type FindPingsQuery struct{}

func (f FindPingsQuery) Name() string {
	return FindPingsQueryName
}

type FindPings struct {
	repo PingRepository
}

func (f FindPings) Handle(ctx context.Context, _ Query) (any, error) {
	return f.repo.Find(ctx)
}

func NewFindPings(repo PingRepository) *FindPings {
	return &FindPings{repo: repo}
}
