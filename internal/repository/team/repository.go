package team

import (
	"context"
	"errors"
	"time"

	"pr-reviewer-assigment-service/internal/client/db"
	innerError "pr-reviewer-assigment-service/internal/errors"
	"pr-reviewer-assigment-service/internal/model"
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/repository/team/converter"
	modelRepo "pr-reviewer-assigment-service/internal/repository/team/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	tableName = "team"

	nameColumn      = "name"
	createdAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.TeamRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, name string) (*model.Team, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			nameColumn,
			createdAtColumn,
		).
		Values(
			name,
			time.Now().UTC(),
		).
		Suffix("RETURNING name, created_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "team_repository.Create",
		QueryRaw: query,
	}

	var team modelRepo.Team
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&team.Name, &team.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, innerError.ErrTeamExists
			}
		}
		return nil, err
	}

	return converter.ToTeamFromRepo(&team), nil
}

func (r *repo) GetByName(ctx context.Context, name string) (*model.Team, error) {
	builder := sq.
		Select(
			nameColumn,
			createdAtColumn,
		).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{nameColumn: name})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "team_repository.GetByID",
		QueryRaw: query,
	}

	var team modelRepo.Team
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&team.Name, &team.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, innerError.ErrTeamNotFound
		}
		return nil, err
	}

	return converter.ToTeamFromRepo(&team), nil
}
