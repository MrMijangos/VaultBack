package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/businesses/domain/entities"
	"vault/src/features/businesses/domain/repositories"
)

const selectBusinessesQuery = `
	SELECT id, user_id, name, COALESCE(types, '{}'), COALESCE(description, ''), COALESCE(location, ''), is_verified, created_at, COALESCE(specialties, '{}')
	FROM businesses
`

type PostgreSQLBusinessRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLBusinessRepository(pool *pgxpool.Pool) *PostgreSQLBusinessRepository {
	return &PostgreSQLBusinessRepository{pool: pool}
}

func scanBusiness(row pgx.Row) (entities.Business, error) {
	var b entities.Business
	err := row.Scan(&b.ID, &b.UserID, &b.Name, &b.Types, &b.Description, &b.Location, &b.IsVerified, &b.CreatedAt, &b.Specialties)
	return b, err
}

func (r *PostgreSQLBusinessRepository) Create(ctx context.Context, business entities.Business) (entities.Business, error) {
	const query = `
		INSERT INTO businesses (user_id, name, types, description, location, specialties)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, business.UserID, business.Name, business.Types, business.Description, business.Location, business.Specialties).Scan(&business.ID)
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo crear el negocio: %w", err)
	}
	return r.FindByID(ctx, business.ID)
}

func (r *PostgreSQLBusinessRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM businesses WHERE user_id = $1)`, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("no se pudo verificar el negocio: %w", err)
	}
	return exists, nil
}

func (r *PostgreSQLBusinessRepository) FindAll(ctx context.Context) ([]entities.Business, error) {
	rows, err := r.pool.Query(ctx, selectBusinessesQuery+" ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los negocios: %w", err)
	}
	defer rows.Close()

	var list []entities.Business
	for rows.Next() {
		b, err := scanBusiness(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el negocio: %w", err)
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *PostgreSQLBusinessRepository) FindByID(ctx context.Context, id string) (entities.Business, error) {
	row := r.pool.QueryRow(ctx, selectBusinessesQuery+" WHERE id = $1", id)
	b, err := scanBusiness(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Business{}, repositories.ErrBusinessNotFound
	}
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo obtener el negocio: %w", err)
	}
	return b, nil
}

func (r *PostgreSQLBusinessRepository) FindByUserID(ctx context.Context, userID string) (entities.Business, error) {
	row := r.pool.QueryRow(ctx, selectBusinessesQuery+" WHERE user_id = $1", userID)
	b, err := scanBusiness(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Business{}, repositories.ErrBusinessNotFound
	}
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo obtener el negocio: %w", err)
	}
	return b, nil
}

func (r *PostgreSQLBusinessRepository) Update(ctx context.Context, id string, userID string, business entities.Business) (entities.Business, error) {
	const query = `
		UPDATE businesses
		SET name = $1, types = $2, description = $3, location = $4, specialties = $5
		WHERE id = $6 AND user_id = $7
	`
	tag, err := r.pool.Exec(ctx, query, business.Name, business.Types, business.Description, business.Location, business.Specialties, id, userID)
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo actualizar el negocio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.Business{}, repositories.ErrBusinessNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *PostgreSQLBusinessRepository) FindPhotosByBusinessID(ctx context.Context, businessID string) ([]entities.BusinessPhoto, error) {
	const query = `
		SELECT id, business_id, url, is_cover, "order", created_at
		FROM business_photos
		WHERE business_id = $1
		ORDER BY "order"
	`

	rows, err := r.pool.Query(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las fotos: %w", err)
	}
	defer rows.Close()

	var photos []entities.BusinessPhoto
	for rows.Next() {
		var p entities.BusinessPhoto
		if err := rows.Scan(&p.ID, &p.BusinessID, &p.URL, &p.IsCover, &p.Order, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("no se pudo leer la foto: %w", err)
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

func (r *PostgreSQLBusinessRepository) AddPhoto(ctx context.Context, businessID string, url string) (entities.BusinessPhoto, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM business_photos WHERE business_id = $1`, businessID).Scan(&count); err != nil {
		return entities.BusinessPhoto{}, fmt.Errorf("no se pudo verificar las fotos existentes: %w", err)
	}

	const query = `
		INSERT INTO business_photos (business_id, url, is_cover, "order")
		VALUES ($1, $2, $3, $4)
		RETURNING id, business_id, url, is_cover, "order", created_at
	`

	var p entities.BusinessPhoto
	err := r.pool.QueryRow(ctx, query, businessID, url, count == 0, count).Scan(
		&p.ID, &p.BusinessID, &p.URL, &p.IsCover, &p.Order, &p.CreatedAt,
	)
	if err != nil {
		return entities.BusinessPhoto{}, fmt.Errorf("no se pudo guardar la foto: %w", err)
	}

	return p, nil
}

func (r *PostgreSQLBusinessRepository) DeletePhoto(ctx context.Context, photoID string, businessID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM business_photos WHERE id = $1 AND business_id = $2`, photoID, businessID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar la foto: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrBusinessNotFound
	}
	return nil
}

func (r *PostgreSQLBusinessRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM businesses WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el negocio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrBusinessNotFound
	}
	return nil
}
