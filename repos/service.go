package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"devdane.com/internal/endpoints"
	pb "devdane.com/repos/proto"
	"github.com/go-redis/redis/v9"
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASS"), // no password set
		DB:       0,                       // use default DB
	})

	lim := int(*req.Limit)
	var limit *int = &lim
	if req.Limit == nil {
		limit = nil
	}

	res := rdb.Get(ctx, "repos")

	if res.Err() != redis.Nil {

		var repos []*pb.Repo
		err := json.Unmarshal([]byte(res.Val()), &repos)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if limit != nil && *limit > 0 {
			return &pb.RepoResponse{
				Repos: repos[:*limit],
			}, nil
		} else {
			return &pb.RepoResponse{
				Repos: repos,
			}, nil
		}

	} else {
		_repos := endpoints.GetRepos()
		repos := make([]*pb.Repo, len(_repos))

		wg := &sync.WaitGroup{}
		for i, _repo := range _repos {
			wg.Add(1)
			go func(i int, _repo endpoints.Repo) {

				langs := make(chan []string)
				go func(langs chan []string) {
					langs <- endpoints.GetLanguages(_repo.Url)
				}(langs)

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
					Language:        <-langs,
				}
				repos[i] = &repo
				wg.Done()
			}(i, _repo)
		}

		wg.Wait()

		b, err := json.Marshal(repos)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		rdb.Set(ctx, "repos", string(b), time.Hour*24*7)

		if limit != nil && *limit > 0 {
			return &pb.RepoResponse{
				Repos: repos[:*limit],
			}, nil
		} else {
			return &pb.RepoResponse{
				Repos: repos,
			}, nil
		}
	}
}
