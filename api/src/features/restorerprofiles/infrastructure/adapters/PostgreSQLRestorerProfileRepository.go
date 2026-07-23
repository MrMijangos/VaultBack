package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/restorerprofiles/domain/entities"
	"vault/src/features/restorerprofiles/domain/repositories"
)

const selectRestorerProfileQuery = `
	SELECT rp.user_id, COALESCE(rp.bio, ''), COALESCE(rp.specialties, '{}'), rp.created_at, rp.updated_at,
	       u.name, COALESCE(u.avatar_url, '')
	FROM restorer_profiles rp
	JOIN users u ON u.id = rp.user_id
`

const selectRestorerServicesQuery = `
	SELECT id, user_id, title, COALESCE(description, ''), price::float8, created_at
	FROM restorer_services
`

type PostgreSQLRestorerProfileRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLRestorerProfileRepository(pool *pgxpool.Pool) *PostgreSQLRestorerProfileRepository {
	return &PostgreSQLRestorerProfileRepository{pool: pool}
}

func scanRestorerProfile(row pgx.Row) (entities.RestorerProfile, error) {
	var p entities.RestorerProfile
	err := row.Scan(&p.UserID, &p.Bio, &p.Specialties, &p.CreatedAt, &p.UpdatedAt, &p.Name, &p.AvatarURL)
	return p, err
}

func (r *PostgreSQLRestorerProfileRepository) findServicesByUserID(ctx context.Context, userID string) ([]entities.RestorerService, error) {
	rows, err := r.pool.Query(ctx, selectRestorerServicesQuery+" WHERE user_id = $1 ORDER BY created_at", userID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los servicios: %w", err)
	}
	defer rows.Close()

	var list []entities.RestorerService
	for rows.Next() {
		var s entities.RestorerService
		if err := rows.Scan(&s.ID, &s.UserID, &s.Title, &s.Description, &s.Price, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("no se pudo leer un servicio: %w", err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *PostgreSQLRestorerProfileRepository) Upsert(ctx context.Context, userID string, bio string, specialties []string, services []entities.RestorerService) (entities.RestorerProfile, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return entities.RestorerProfile{}, fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	const upsertProfileQuery = `
		INSERT INTO restorer_profiles (user_id, bio, specialties)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET bio = EXCLUDED.bio, specialties = EXCLUDED.specialties, updated_at = now()
	`
	if _, err := tx.Exec(ctx, upsertProfileQuery, userID, bio, specialties); err != nil {
		return entities.RestorerProfile{}, fmt.Errorf("no se pudo guardar el perfil: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM restorer_services WHERE user_id = $1`, userID); err != nil {
		return entities.RestorerProfile{}, fmt.Errorf("no se pudieron limpiar los servicios anteriores: %w", err)
	}

	for _, s := range services {
		if _, err := tx.Exec(ctx, `
			INSERT INTO restorer_services (user_id, title, description, price)
			VALUES ($1, $2, $3, $4)
		`, userID, s.Title, s.Description, s.Price); err != nil {
			return entities.RestorerProfile{}, fmt.Errorf("no se pudo guardar un servicio: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return entities.RestorerProfile{}, fmt.Errorf("no se pudo confirmar la transaccion: %w", err)
	}

	return r.FindByUserID(ctx, userID)
}

func (r *PostgreSQLRestorerProfileRepository) FindByUserID(ctx context.Context, userID string) (entities.RestorerProfile, error) {
	row := r.pool.QueryRow(ctx, selectRestorerProfileQuery+" WHERE rp.user_id = $1", userID)
	p, err := scanRestorerProfile(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.RestorerProfile{}, repositories.ErrProfileNotFound
	}
	if err != nil {
		return entities.RestorerProfile{}, fmt.Errorf("no se pudo obtener el perfil: %w", err)
	}

	services, err := r.findServicesByUserID(ctx, userID)
	if err != nil {
		return entities.RestorerProfile{}, err
	}
	p.Services = services

	return p, nil
}

func (r *PostgreSQLRestorerProfileRepository) ListWithServices(ctx context.Context) ([]entities.RestorerProfile, error) {
	rows, err := r.pool.Query(ctx, selectRestorerProfileQuery+`
		WHERE EXISTS (SELECT 1 FROM restorer_services rs WHERE rs.user_id = rp.user_id)
		ORDER BY rp.created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los perfiles: %w", err)
	}
	defer rows.Close()

	var list []entities.RestorerProfile
	for rows.Next() {
		p, err := scanRestorerProfile(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer un perfil: %w", err)
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range list {
		services, err := r.findServicesByUserID(ctx, list[i].UserID)
		if err != nil {
			return nil, err
		}
		list[i].Services = services
	}

	return list, nil
}
