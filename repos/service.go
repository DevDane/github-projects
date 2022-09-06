package repos

import (
	"context"

	"devdane.com/internal/endpoints"
	pb "devdane.com/repos/proto"
	"google.golang.org/grpc"
)

type service struct {
	pb.UnimplementedReposServiceServer
}

func NewService() *service {
	return &service{}
}

func (svc *service) RegisterService(server *grpc.Server) {
	pb.RegisterReposServiceServer(server, svc)
}

func (svc *service) GetRepos(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
	lim := int(*req.Limit)
	var limit *int = &lim
	if req.Limit == nil {
		limit = nil
	}

	_repos := endpoints.GetRepos(limit)
	var repos []*pb.Repo
	for _, _repo := range _repos {
		langs := endpoints.GetLanguages(_repo.Url)
		readme := endpoints.GetReadme(_repo.Url)
		file := endpoints.GetFile(_repo.Url, endpoints.ImageFromMD(readme))

		repo := pb.Repo{
			Id:              int64(_repo.ID),
			NodeId:          _repo.NodeId,
			Name:            _repo.Name,
			Description:     _repo.Description,
			FullName:        _repo.FullName,
			HtmlUrl:         _repo.HTMLUrl,
			Url:             _repo.Url,
			StargazersCount: int64(_repo.StarGazersCount),
			WatchersCount:   int64(_repo.WatchersCount),
			Visibility:      _repo.Visibility,
			ImageUrl:        file.DownloadUrl,
			Language:        langs,
		}
		repos = append(repos, &repo)
	}

	return &pb.RepoResponse{
		Repos: repos,
	}, nil
}
