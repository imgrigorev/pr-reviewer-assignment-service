package user

import (
	"context"
	"errors"
	"time"

	"pr-reviewer-assigment-service/internal/client/db"
	innerError "pr-reviewer-assigment-service/internal/errors"
	"pr-reviewer-assigment-service/internal/model"
	"pr-reviewer-assigment-service/internal/repository"
	"pr-reviewer-assigment-service/internal/repository/user/converter"
	modelRepo "pr-reviewer-assigment-service/internal/repository/user/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	tableName = `"user"`

	idColumn        = "id"
	nameColumn      = "name"
	isActiveColumn  = "is_active"
	teamNameColumn  = "team_name"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, u *model.User, teamName string) (*model.User, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			idColumn,
			nameColumn,
			isActiveColumn,
			teamNameColumn,
			createdAtColumn,
		).
		Values(
			u.ID,
			u.Name,
			u.IsActive,
			teamName,
			time.Now().UTC(),
		).
		Suffix("RETURNING id, name, is_active, team_name, created_at, updated_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Get(ctx context.Context, id string) (*model.User, error) {
	builder := sq.
		Select(
			idColumn,
			nameColumn,
			isActiveColumn,
			teamNameColumn,
			createdAtColumn,
			updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, u *model.User) (*model.User, error) {
	builder := sq.
		Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(isActiveColumn, u.IsActive).
		Set(updatedAtColumn, time.Now().UTC()).
		Where(sq.Eq{idColumn: u.ID}).
		Suffix("RETURNING id, name, is_active, team_name, created_at, updated_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.IsActive, &user.TeamName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, innerError.ErrUserNotFound
		}
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) GetList(ctx context.Context, teamName string, isActive *bool) ([]*model.User, error) {
	builder := sq.
		Select(
			idColumn,
			nameColumn,
			isActiveColumn,
			teamNameColumn,
			createdAtColumn,
			updatedAtColumn,
		).
		From(tableName).
		PlaceholderFormat(sq.Dollar)

	if isActive != nil {
		builder = builder.Where(
			sq.And{
				sq.Eq{teamNameColumn: teamName},
				sq.Eq{isActiveColumn: isActive},
			})
	} else {
		builder = builder.Where(sq.Eq{teamNameColumn: teamName})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var u modelRepo.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.IsActive,
			&u.TeamName,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, converter.ToUserFromRepo(&u))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
