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

ALTER TABLE users ADD COLUMN IF NOT EXISTS public_key text;

-- roles es el historico acumulado (nunca se quita nada, solo se agrega) --
-- a diferencia de `role`, que sigue siendo el rol "principal/mas reciente"
-- y no se toca para no romper nada que ya lo lea. El backfill solo corre
-- sobre cuentas que todavia no tienen nada acumulado.
ALTER TABLE users ADD COLUMN IF NOT EXISTS roles text[] DEFAULT '{}'::text[];
UPDATE users SET roles = ARRAY[role] WHERE roles = '{}'::text[];

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

ALTER TABLE assets ADD COLUMN IF NOT EXISTS is_for_sale boolean DEFAULT false;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS sale_price numeric;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS sale_description text;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS size character varying;

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

ALTER TABLE businesses ADD COLUMN IF NOT EXISTS specialties text[] DEFAULT '{}'::text[];

-- types reemplaza `type` (permite elegir mas de una categoria a la vez,
-- p.ej. un negocio que hace mantenimiento Y reparacion). Se deja `type` en
-- la tabla sin usarlo -- evita el riesgo de un DROP COLUMN que nadie pidio.
ALTER TABLE businesses ADD COLUMN IF NOT EXISTS types text[] DEFAULT '{}'::text[];
ALTER TABLE businesses ALTER COLUMN type DROP NOT NULL;
UPDATE businesses SET types = ARRAY[type] WHERE types = '{}'::text[] AND type IS NOT NULL;

CREATE TABLE IF NOT EXISTS business_services (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	business_id uuid NOT NULL,
	title character varying NOT NULL,
	description text,
	price numeric NOT NULL,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT business_services_pkey PRIMARY KEY (id),
	CONSTRAINT business_services_business_id_fkey FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS restorer_profiles (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL UNIQUE,
	bio text,
	specialties text[] DEFAULT '{}'::text[],
	created_at timestamp without time zone DEFAULT now(),
	updated_at timestamp without time zone DEFAULT now(),
	CONSTRAINT restorer_profiles_pkey PRIMARY KEY (id),
	CONSTRAINT restorer_profiles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS restorer_services (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	title character varying NOT NULL,
	description text,
	price numeric NOT NULL,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT restorer_services_pkey PRIMARY KEY (id),
	CONSTRAINT restorer_services_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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

CREATE TABLE IF NOT EXISTS saved_posts (
	post_id uuid NOT NULL,
	user_id uuid NOT NULL,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT saved_posts_pkey PRIMARY KEY (post_id, user_id),
	CONSTRAINT saved_posts_post_id_fkey FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
	CONSTRAINT saved_posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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

CREATE TABLE IF NOT EXISTS asset_comments (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	asset_id uuid NOT NULL,
	user_id uuid NOT NULL,
	content text NOT NULL,
	toxicity_score double precision,
	is_visible boolean DEFAULT true,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT asset_comments_pkey PRIMARY KEY (id),
	CONSTRAINT asset_comments_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
	CONSTRAINT asset_comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reviews (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	provider_id uuid NOT NULL,
	content text NOT NULL,
	sentiment_score double precision,
	sentiment_label character varying CHECK (sentiment_label::text = ANY (ARRAY['positivo'::character varying, 'negativo'::character varying, 'neutral'::character varying]::text[])),
	toxicity_score double precision,
	is_visible boolean DEFAULT true,
	likes_count integer DEFAULT 0,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT reviews_pkey PRIMARY KEY (id),
	CONSTRAINT reviews_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT reviews_provider_id_fkey FOREIGN KEY (provider_id) REFERENCES users(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS no agrega columnas a una tabla que ya existe
-- (Railway y el Postgres local ya tienen "reviews" creada) -- este ALTER
-- idempotente sí se aplica en cada arranque vía RunMigrations.
ALTER TABLE reviews ADD COLUMN IF NOT EXISTS sentiment_label character varying
	CHECK (sentiment_label::text = ANY (ARRAY['positivo'::character varying, 'negativo'::character varying, 'neutral'::character varying]::text[]));

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

CREATE TABLE IF NOT EXISTS addresses (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	label character varying NOT NULL,
	recipient character varying NOT NULL,
	phone character varying NOT NULL,
	street character varying NOT NULL,
	city character varying NOT NULL,
	state character varying NOT NULL,
	postal_code character varying NOT NULL,
	reference_notes text,
	is_default boolean NOT NULL DEFAULT false,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT addresses_pkey PRIMARY KEY (id),
	CONSTRAINT addresses_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chat_messages (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	sender_id uuid NOT NULL,
	recipient_id uuid NOT NULL,
	cipher_text text NOT NULL,
	encrypted_aes_key text NOT NULL,
	iv text NOT NULL,
	status character varying NOT NULL DEFAULT 'sent'::character varying CHECK (status::text = ANY (ARRAY['sent'::character varying, 'delivered'::character varying, 'read'::character varying]::text[])),
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT chat_messages_pkey PRIMARY KEY (id),
	CONSTRAINT chat_messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT chat_messages_recipient_id_fkey FOREIGN KEY (recipient_id) REFERENCES users(id) ON DELETE CASCADE
);

-- La llave AES de cada mensaje se cifra dos veces: una con la pública del
-- receptor (encrypted_aes_key, arriba) y otra con la pública del propio
-- emisor -- sin esto, quien envía un mensaje jamás podría releerlo despues
-- (su privada no destraba una llave cifrada para la pública ajena).
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS encrypted_aes_key_sender text;

CREATE INDEX IF NOT EXISTS idx_chat_messages_conversation
	ON chat_messages (LEAST(sender_id, recipient_id), GREATEST(sender_id, recipient_id), created_at);

CREATE OR REPLACE FUNCTION notify_new_chat_message() RETURNS trigger AS $$
BEGIN
	PERFORM pg_notify('new_chat_message', row_to_json(NEW)::text);
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS chat_messages_notify_trigger ON chat_messages;
CREATE TRIGGER chat_messages_notify_trigger
	AFTER INSERT ON chat_messages
	FOR EACH ROW EXECUTE FUNCTION notify_new_chat_message();

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
