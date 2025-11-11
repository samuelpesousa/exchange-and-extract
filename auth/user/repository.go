package user

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("usuário não encontrado")
	ErrEmailAlreadyExists = errors.New("email já cadastrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create cria um novo usuário
func (r *Repository) Create(email, passwordHash, nome string) (*User, error) {
	query := `
		INSERT INTO users (email, senha, nome)
		VALUES ($1, $2, $3)
		RETURNING id, email, senha, nome, created_at, updated_at
	`

	user := &User{}
	err := r.db.QueryRow(query, email, passwordHash, nome).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Nome,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

// FindByEmail busca um usuário por email
func (r *Repository) FindByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, senha, nome, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Nome,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID busca um usuário por ID
func (r *Repository) FindByID(id int) (*User, error) {
	query := `
		SELECT id, email, senha, nome, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Nome,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update atualiza os dados do usuário
func (r *Repository) Update(user *User) error {
	query := `
		UPDATE users
		SET email = $1, nome = $2, updated_at = $3
		WHERE id = $4
	`

	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.Email, user.Nome, user.UpdatedAt, user.ID)
	return err
}

// Delete remove um usuário
func (r *Repository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
