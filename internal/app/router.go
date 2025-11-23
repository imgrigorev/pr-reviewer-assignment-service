package app

import (
	"context"

	"github.com/go-chi/chi/v5"
)

func (a *App) initRouter(ctx context.Context) error {
	a.router = chi.NewRouter()

	user := a.serviceProvider.UserImpl(ctx)
	team := a.serviceProvider.TeamImpl(ctx)
	pullRequest := a.serviceProvider.PullRequestImpl(ctx)

	// users
	{
		a.router.Post("/users", user.Create)
		a.router.Post("/users/setIsActive", user.Update)
		a.router.Get("/users/getReview", user.GetUserReviews)
	}

	// teams
	{
		a.router.Post("/team", team.Create)
		a.router.Get("/team/get", team.Get)
	}
	// pull request
	{
		a.router.Post("/pullRequest/create", pullRequest.Create)
		a.router.Post("/pullRequest/merge", pullRequest.Merge)
		a.router.Post("/pullRequest/reassign", pullRequest.ReassignReviewer)
	}
	return nil
}
