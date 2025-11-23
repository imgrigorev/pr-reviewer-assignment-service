package pull_request_reviewers

import (
	"context"
	"time"

	"pr-reviewer-assigment-service/internal/client/db"
	"pr-reviewer-assigment-service/internal/model"
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/repository/pull_request_reviewers/converter"
	modelRepo "pr-reviewer-assigment-service/internal/repository/pull_request_reviewers/model"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "pr_reviewers"

	prIDColumn       = "pr_id"
	reviewerIDColumn = "reviewer_id"
	assignedAtColumn = "assigned_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.PullRequestReviewersRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, assignment *model.PullRequestReviewer) (*model.PullRequestReviewer, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			prIDColumn,
			reviewerIDColumn,
			assignedAtColumn,
		).
		Values(
			assignment.ID,
			assignment.ReviewerID,
			time.Now().UTC(),
		).
		Suffix("RETURNING pr_id, reviewer_id, assigned_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "pull_request_reviewers_repository.Create",
		QueryRaw: query,
	}

	var reviewerAssignment modelRepo.PullRequestReviewer
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(
			&reviewerAssignment.PrID,
			&reviewerAssignment.ReviewerID,
			&reviewerAssignment.AssignedAt,
		)
	if err != nil {
		return nil, err
	}

	return converter.ToPullRequestReviewerFromRepo(&reviewerAssignment), nil
}

func (r *repo) GetByUserID(ctx context.Context, userID string) ([]*model.PullRequestReviewer, error) {
	builder := sq.Select(
		prIDColumn,
		reviewerIDColumn,
		assignedAtColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{reviewerIDColumn: userID}).
		OrderBy(assignedAtColumn + " DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "pull_request_reviewers_repository.GetByUserID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*model.PullRequestReviewer
	for rows.Next() {
		var assignment modelRepo.PullRequestReviewer
		err := rows.Scan(
			&assignment.PrID,
			&assignment.ReviewerID,
			&assignment.AssignedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, converter.ToPullRequestReviewerFromRepo(&assignment))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (r *repo) GetByPrID(ctx context.Context, prID string) ([]*model.PullRequestReviewer, error) {
	builder := sq.Select(
		prIDColumn,
		reviewerIDColumn,
		assignedAtColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{prIDColumn: prID}).
		OrderBy(assignedAtColumn + " DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "pull_request_reviewers_repository.GetByPrID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*model.PullRequestReviewer
	for rows.Next() {
		var assignment modelRepo.PullRequestReviewer
		err := rows.Scan(
			&assignment.PrID,
			&assignment.ReviewerID,
			&assignment.AssignedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, converter.ToPullRequestReviewerFromRepo(&assignment))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (r *repo) Delete(ctx context.Context, prID, reviewerID string) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{
			prIDColumn:       prID,
			reviewerIDColumn: reviewerID,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "pull_request_reviewers_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) DeleteByPrID(ctx context.Context, prID string) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{prIDColumn: prID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "pull_request_reviewers_repository.DeleteByPrID",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) DeleteByUserID(ctx context.Context, userID string) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{reviewerIDColumn: userID})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "pull_request_reviewers_repository.DeleteByUserID",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
