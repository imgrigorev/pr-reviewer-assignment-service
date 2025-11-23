package app

import (
	"context"
	"log"

	"pr-reviewer-assigment-service/internal/api/pullrequest"
	"pr-reviewer-assigment-service/internal/api/team"
	"pr-reviewer-assigment-service/internal/api/user"
	"pr-reviewer-assigment-service/internal/client/db"
	"pr-reviewer-assigment-service/internal/client/db/pg"
	"pr-reviewer-assigment-service/internal/closer"
	"pr-reviewer-assigment-service/internal/config"
	"pr-reviewer-assigment-service/internal/repository"
	pullRequestRepository "pr-reviewer-assigment-service/internal/repository/pull_request"
	pullRequestReviewersRepository "pr-reviewer-assigment-service/internal/repository/pull_request_reviewers"
	teamRepository "pr-reviewer-assigment-service/internal/repository/team"
	userRepository "pr-reviewer-assigment-service/internal/repository/user"
	"pr-reviewer-assigment-service/internal/service"
	pullRequestService "pr-reviewer-assigment-service/internal/service/pull_request"
	teamService "pr-reviewer-assigment-service/internal/service/team"
	userService "pr-reviewer-assigment-service/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	httpConfig config.HTTPConfig

	dbClient db.Client

	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *user.Implementation

	teamRepository repository.TeamRepository
	teamService    service.TeamService
	teamImpl       *team.Implementation

	pullRequestRepository repository.PullRequestRepository
	pullRequestService    service.PullRequestService
	pullRequestImpl       *pullrequest.Implementation

	pullRequestReviewersRepository repository.PullRequestReviewersRepository
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.PullRequestRepository(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) TeamRepository(ctx context.Context) repository.TeamRepository {
	if s.teamRepository == nil {
		s.teamRepository = teamRepository.NewRepository(s.DBClient(ctx))
	}

	return s.teamRepository
}

func (s *serviceProvider) TeamService(ctx context.Context) service.TeamService {
	if s.teamService == nil {
		s.teamService = teamService.NewService(
			s.UserRepository(ctx),
			s.TeamRepository(ctx),
		)
	}

	return s.teamService
}

func (s *serviceProvider) TeamImpl(ctx context.Context) *team.Implementation {
	if s.teamImpl == nil {
		s.teamImpl = team.NewImplementation(s.TeamService(ctx))
	}

	return s.teamImpl
}

func (s *serviceProvider) PullRequestReviewersRepository(ctx context.Context) repository.PullRequestReviewersRepository {
	if s.pullRequestReviewersRepository == nil {
		s.pullRequestReviewersRepository = pullRequestReviewersRepository.NewRepository(s.DBClient(ctx))
	}

	return s.pullRequestReviewersRepository
}

func (s *serviceProvider) PullRequestRepository(ctx context.Context) repository.PullRequestRepository {
	if s.pullRequestRepository == nil {
		s.pullRequestRepository = pullRequestRepository.NewRepository(s.DBClient(ctx))
	}

	return s.pullRequestRepository
}

func (s *serviceProvider) PullRequestService(ctx context.Context) service.PullRequestService {
	if s.pullRequestService == nil {
		s.pullRequestService = pullRequestService.NewService(
			s.PullRequestRepository(ctx),
			s.PullRequestReviewersRepository(ctx),
			s.TeamRepository(ctx),
			s.UserRepository(ctx),
		)
	}

	return s.pullRequestService
}

func (s *serviceProvider) PullRequestImpl(ctx context.Context) *pullrequest.Implementation {
	if s.pullRequestImpl == nil {
		s.pullRequestImpl = pullrequest.NewImplementation(s.PullRequestService(ctx))
	}

	return s.pullRequestImpl
}
