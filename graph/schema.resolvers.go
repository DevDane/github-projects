package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"devdane.com/graph/generated"
	"devdane.com/graph/model"
	"devdane.com/internal/endpoints"
)

// Repos is the resolver for the repos field.
func (r *queryResolver) Repos(ctx context.Context, limit *int) ([]*model.Repo, error) {
	_repos := endpoints.GetRepos(limit)
	var repos []*model.Repo
	for _, _repo := range _repos {
		langs := endpoints.GetLanguages(_repo.Url)
		readme := endpoints.GetReadme(_repo.Url)
		file := endpoints.GetFile(_repo.Url, endpoints.ImageFromMD(readme))

		repo := model.Repo{
			ID:              _repo.ID,
			NodeID:          _repo.NodeId,
			Name:            _repo.Name,
			Description:     _repo.Description,
			FullName:        _repo.FullName,
			HTMLURL:         _repo.HTMLUrl,
			URL:             _repo.Url,
			StargazersCount: _repo.StarGazersCount,
			WatchersCount:   _repo.WatchersCount,
			Visibility:      _repo.Visibility,
			ImageURL:        file.DownloadUrl,
			Languages:       langs,
		}
		repos = append(repos, &repo)
	}

	return repos, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
