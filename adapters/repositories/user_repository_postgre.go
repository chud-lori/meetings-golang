package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"meeting_service/domain/entities"
	"time"
)

type UserRepositoryPostgre struct {
	db *sql.DB
}

func NewUserRepositoryPostgre(db *sql.DB) (*UserRepositoryPostgre, error) {
	return &UserRepositoryPostgre{
		db: db,
	}, nil
}

func (repository *UserRepositoryPostgre) Save(ctx context.Context, user *entities.User) (*entities.User, error) {
	var id string
	var createdAt time.Time
	query := `
            INSERT INTO users (email, passcode)
            VALUES ($1, $2)
            RETURNING id, created_at`
	err := repository.db.QueryRow(query, user.Email, user.Passcode).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	user.Id = id
	user.Created_at = createdAt

	return user, nil
}

func (repository *UserRepositoryPostgre) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := "UPDATE users SET email = $1, passcode = $2 WHERE id = $3"
	_, err := repository.db.Exec(query, user.Email, user.Passcode, user.Id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := repository.db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}

func (repository *UserRepositoryPostgre) FindById(ctx context.Context, id string) (*entities.User, error) {
	user := &entities.User{}

	query := "SELECT id, email, created_at FROM users WHERE id = $1"

	err := repository.db.QueryRow(query, id).Scan(&user.Id, &user.Email, &user.Created_at)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No data found")
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (repository *UserRepositoryPostgre) FindAll(ctx context.Context) ([]*entities.User, error) {
	query := "SELECT id, email, created_at FROM users"

	rows, err := repository.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var users []*entities.User

	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.Id, &user.Email, &user.Created_at)

		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repository *UserRepositoryPostgre) CheckInteractionStatus(ctx context.Context, user1 string, user2 string) ([]uint8, error) {
	var statuses []uint8

	query := "SELECT status FROM interactions WHERE (from_id = $1 and to_id = $2) or (to_id = $1 and from_id = $2) ORDER BY created_at DESC"

	rows, err := repository.db.QueryContext(ctx, query, user1, user2)
	if err != nil {
		log.Printf("Error in checkinteraction repo psql %#v", err)
		return nil, err
	}

	for rows.Next() {
		var status uint8
		if err := rows.Scan(&status); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error in checkinteraction repo psql %#v", err)
		return nil, err
	}

	return statuses, nil
}

func (repository *UserRepositoryPostgre) Request(ctx context.Context, fromUser string, toUser string) error {
	query := "INSERT into interactions(from_id, to_id, status) VALUES($1, $2, $3)"
	_, err := repository.db.Exec(query, fromUser, toUser, 1)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryPostgre) Receive(ctx context.Context, fromUser string, toUser string) error {
	query := "INSERT into interactions(from_id, to_id, status) VALUES($1, $2, $3)"
	_, err := repository.db.Exec(query, fromUser, toUser, 2)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryPostgre) Decline(ctx context.Context, fromUser string, toUser string) error {
	query := "INSERT into interactions(from_id, to_id, status) VALUES($1, $2, $3)"
	_, err := repository.db.Exec(query, fromUser, toUser, 3)

	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepositoryPostgre) StoreMeet(ctx context.Context, fromUser string, toUser string) error {
	query := "INSERT INTO meets(user_1, user_2) VALUES($1, $2)"
	_, err := repository.db.Exec(query, fromUser, toUser)

	if err != nil {
		return err
	}

	return nil
}
