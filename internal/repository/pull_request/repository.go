package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pr-reviewer-assigment-service/internal/client/db"
	innerError "pr-reviewer-assigment-service/internal/errors"
	"pr-reviewer-assigment-service/internal/model"
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/repository/pull_request/converter"
	modelRepo "pr-reviewer-assigment-service/internal/repository/pull_request/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	tableName = "prs"

	idColumn        = "id"
	titleColumn     = "title"
	authorIDColumn  = "author_id"
	statusColumn    = "status"
	createdAtColumn = "created_at"
	mergedAtColumn  = "merged_at"

	mergedStatus = "MERGED"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.PullRequestRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, pr *model.PullRequest) (*model.PullRequest, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			idColumn,
			titleColumn,
			authorIDColumn,
			statusColumn,
			createdAtColumn,
		).
		Values(
			pr.ID,
			pr.Name,
			pr.AuthorID,
			pr.Status,
			time.Now().UTC(),
		).
		Suffix("RETURNING id, title, author_id, status, created_at, merged_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "pull_request_repository.Create",
		QueryRaw: query,
	}

	var pullRequest modelRepo.PullRequest
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(
			&pullRequest.ID,
			&pullRequest.Name,
			&pullRequest.AuthorID,
			&pullRequest.Status,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt,
		)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, innerError.ErrPRAlreadyExists
			}
		}
		return nil, err
	}

	return converter.ToPullRequestFromRepo(&pullRequest), nil
}

func (r *repo) Merge(ctx context.Context, prID string) (*model.PullRequest, error) {
	query := `
        WITH current_pr AS (
            SELECT *
            FROM prs 
            WHERE id = $1
        ),
        updated_pr AS (
            UPDATE prs 
            SET status = 'MERGED', 
                merged_at = CASE 
                    WHEN merged_at IS NULL THEN NOW() 
                    ELSE merged_at 
                END
            WHERE id = $1 
            AND status != 'MERGED'
            RETURNING id, title, author_id, status, created_at, merged_at
        )
        SELECT 
            COALESCE(u.id, c.id) as id,
            COALESCE(u.title, c.title) as title,
            COALESCE(u.author_id, c.author_id) as author_id,
            COALESCE(u.status, c.status) as status,
            COALESCE(u.created_at, c.created_at) as created_at,
            COALESCE(u.merged_at, c.merged_at) as merged_at
        FROM current_pr c
        LEFT JOIN updated_pr u ON true
        WHERE c.id IS NOT NULL
    `

	q := db.Query{
		Name:     "pull_request_repository.Merge",
		QueryRaw: query,
	}

	var pullRequest modelRepo.PullRequest
	err := r.db.DB().QueryRowContext(ctx, q, prID).Scan(
		&pullRequest.ID,
		&pullRequest.Name,
		&pullRequest.AuthorID,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("pull request not found")
		}
		return nil, err
	}

	return converter.ToPullRequestFromRepo(&pullRequest), nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*model.PullRequest, error) {
	builder := sq.Select(
		idColumn,
		titleColumn,
		authorIDColumn,
		statusColumn,
		createdAtColumn,
		mergedAtColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "pull_request_repository.GetByID",
		QueryRaw: query,
	}

	var pullRequest modelRepo.PullRequest
	err = r.db.DB().QueryRowContext(ctx, q, args...).
		Scan(
			&pullRequest.ID,
			&pullRequest.Name,
			&pullRequest.AuthorID,
			&pullRequest.Status,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt,
		)
	if err != nil {
		return nil, err
	}

	return converter.ToPullRequestFromRepo(&pullRequest), nil
}

func (r *repo) GetByReviewerID(ctx context.Context, reviewerID string) ([]*model.PullRequest, error) {
	builder := sq.Select(
		"p.id",
		"p.title",
		"p.author_id",
		"p.status",
		"p.created_at",
		"p.merged_at",
	).
		From("prs p").
		LeftJoin("pr_reviewers pr ON p.id = pr.pr_id").
		Where(sq.Eq{"pr.reviewer_id": reviewerID}).
		PlaceholderFormat(sq.Dollar).
		OrderBy("p.created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "pull_request_repository.GetByReviewerID",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*model.PullRequest
	for rows.Next() {
		var pullRequest modelRepo.PullRequest
		err := rows.Scan(
			&pullRequest.ID,
			&pullRequest.Name,
			&pullRequest.AuthorID,
			&pullRequest.Status,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt,
		)
		if err != nil {
			return nil, err
		}
		prs = append(prs, converter.ToPullRequestFromRepo(&pullRequest))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
