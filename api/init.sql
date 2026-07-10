CREATE TABLE IF NOT EXISTS users (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	name character varying NOT NULL,
	email character varying NOT NULL UNIQUE,
	password character varying NOT NULL,
	avatar_url character varying,
	role character varying NOT NULL DEFAULT 'usuario'::character varying CHECK (role::text = ANY (ARRAY['usuario'::character varying, 'vendedor'::character varying, 'restaurador'::character varying, 'servicio'::character varying, 'admin'::character varying]::text[])),
	created_at timestamp without time zone DEFAULT now(),
	updated_at timestamp without time zone DEFAULT now(),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS assets (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	name character varying NOT NULL,
	category character varying NOT NULL CHECK (category::text = ANY (ARRAY['sneakers'::character varying, 'gorras'::character varying, 'relojes'::character varying, 'lentes'::character varying, 'carteras'::character varying, 'bolsos'::character varying, 'pulsos'::character varying, 'bisuteria'::character varying, 'coleccionables'::character varying, 'otros'::character varying]::text[])),
	brand character varying,
	purchase_value numeric,
	condition character varying NOT NULL DEFAULT 'nuevo'::character varying CHECK (condition::text = ANY (ARRAY['nuevo'::character varying, 'seminuevo'::character varying, 'usado'::character varying]::text[])),
	purchase_date date,
	store_origin character varying,
	notes text,
	blockchain_tx_id character varying,
	blockchain_hash character varying,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT assets_pkey PRIMARY KEY (id),
	CONSTRAINT assets_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS businesses (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL UNIQUE,
	name character varying NOT NULL,
	type character varying NOT NULL CHECK (type::text = ANY (ARRAY['restaurador'::character varying, 'servicio'::character varying]::text[])),
	description text,
	location character varying,
	is_verified boolean DEFAULT false,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT businesses_pkey PRIMARY KEY (id),
	CONSTRAINT businesses_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS asset_photos (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	asset_id uuid NOT NULL,
	url character varying NOT NULL,
	is_cover boolean DEFAULT false,
	"order" integer DEFAULT 0,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT asset_photos_pkey PRIMARY KEY (id),
	CONSTRAINT asset_photos_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS maintenance_logs (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	asset_id uuid NOT NULL,
	provider_id uuid,
	type character varying NOT NULL CHECK (type::text = ANY (ARRAY['mantenimiento'::character varying, 'restauracion'::character varying]::text[])),
	subtype character varying,
	cost numeric,
	performed_at date,
	notes text,
	blockchain_tx_id character varying,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT maintenance_logs_pkey PRIMARY KEY (id),
	CONSTRAINT maintenance_logs_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES users(id) ON DELETE SET NULL,
	CONSTRAINT maintenance_logs_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	asset_id uuid,
	content text NOT NULL,
	sentiment_score double precision,
	sentiment_label character varying CHECK (sentiment_label::text = ANY (ARRAY['positivo'::character varying, 'negativo'::character varying, 'neutral'::character varying]::text[])),
	toxicity_score double precision,
	is_visible boolean DEFAULT true,
	likes_count integer DEFAULT 0,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT posts_pkey PRIMARY KEY (id),
	CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT posts_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS post_photos (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	post_id uuid NOT NULL,
	url character varying NOT NULL,
	"order" integer DEFAULT 0,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT post_photos_pkey PRIMARY KEY (id),
	CONSTRAINT post_photos_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_likes (
	post_id uuid NOT NULL,
	user_id uuid NOT NULL,
	CONSTRAINT post_likes_pkey PRIMARY KEY (post_id, user_id),
	CONSTRAINT post_likes_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	CONSTRAINT post_likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	post_id uuid NOT NULL,
	user_id uuid NOT NULL,
	content text NOT NULL,
	toxicity_score double precision,
	is_visible boolean DEFAULT true,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT comments_pkey PRIMARY KEY (id),
	CONSTRAINT comments_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	CONSTRAINT comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviews (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	provider_id uuid NOT NULL,
	content text NOT NULL,
	sentiment_score double precision,
	toxicity_score double precision,
	is_visible boolean DEFAULT true,
	likes_count integer DEFAULT 0,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT reviews_pkey PRIMARY KEY (id),
	CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT reviews_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS review_likes (
	review_id uuid NOT NULL,
	user_id uuid NOT NULL,
	CONSTRAINT review_likes_pkey PRIMARY KEY (review_id, user_id),
	CONSTRAINT review_likes_review_id_fkey FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
	CONSTRAINT review_likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS notifications (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	type character varying NOT NULL CHECK (type::text = ANY (ARRAY['servicio'::character varying, 'reparacion'::character varying, 'venta'::character varying, 'blockchain'::character varying, 'comunidad'::character varying]::text[])),
	subtype character varying NOT NULL CHECK (subtype::text = ANY (ARRAY['entro_servicio'::character varying, 'salio_servicio'::character varying, 'entro_reparacion'::character varying, 'salio_reparacion'::character varying, 'pedido_recibido'::character varying, 'pedido_enviado'::character varying, 'nueva_compra'::character varying, 'asset_verificado'::character varying, 'likes_post'::character varying]::text[])),
	title character varying NOT NULL,
	body text NOT NULL,
	data jsonb DEFAULT '{}'::jsonb,
	read boolean DEFAULT false,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT notifications_pkey PRIMARY KEY (id),
	CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION notify_new_notification() RETURNS trigger AS $$
BEGIN
	PERFORM pg_notify('new_notification', row_to_json(NEW)::text);
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS notifications_notify_trigger ON notifications;
CREATE TRIGGER notifications_notify_trigger
	AFTER INSERT ON notifications
	FOR EACH ROW EXECUTE FUNCTION notify_new_notification();

CREATE TABLE IF NOT EXISTS blockchain_certificates (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	asset_id uuid NOT NULL,
	owner_id uuid NOT NULL,
	tx_id character varying NOT NULL UNIQUE,
	asset_hash character varying NOT NULL,
	action character varying NOT NULL CHECK (action::text = ANY (ARRAY['REGISTERED'::character varying, 'MAINTAINED'::character varying, 'RESTORED'::character varying, 'TRANSFERRED'::character varying]::text[])),
	network character varying NOT NULL DEFAULT 'testnet'::character varying CHECK (network::text = ANY (ARRAY['testnet'::character varying, 'mainnet'::character varying]::text[])),
	confirmed_at timestamp without time zone DEFAULT now(),
	CONSTRAINT blockchain_certificates_pkey PRIMARY KEY (id),
	CONSTRAINT blockchain_certificates_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
	CONSTRAINT blockchain_certificates_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);
