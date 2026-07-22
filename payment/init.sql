-- payment/ comparte el mismo Postgres que api/ (ver users/assets ahi) --
-- estas son solo las tablas propias de este servicio.

CREATE TABLE IF NOT EXISTS subscriptions (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	plan_id character varying NOT NULL,
	status character varying NOT NULL,
	stripe_customer_id character varying,
	stripe_subscription_id character varying,
	current_period_start timestamp without time zone,
	current_period_end timestamp without time zone,
	canceled_at timestamp without time zone,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT subscriptions_pkey PRIMARY KEY (id),
	CONSTRAINT subscriptions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ads (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	subscription_id uuid NOT NULL,
	title character varying NOT NULL,
	description text,
	image_url character varying,
	target_section character varying NOT NULL,
	target_id character varying,
	status character varying NOT NULL,
	impressions bigint DEFAULT 0,
	clicks bigint DEFAULT 0,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT ads_pkey PRIMARY KEY (id),
	CONSTRAINT ads_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT ads_subscription_id_fkey FOREIGN KEY (subscription_id) REFERENCES subscriptions(id)
);

CREATE TABLE IF NOT EXISTS connected_accounts (
	user_id uuid NOT NULL,
	stripe_account_id character varying NOT NULL,
	charges_enabled boolean DEFAULT false,
	created_at timestamp without time zone DEFAULT now(),
	CONSTRAINT connected_accounts_pkey PRIMARY KEY (user_id),
	CONSTRAINT connected_accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS orders (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	buyer_id uuid NOT NULL,
	seller_id uuid NOT NULL,
	asset_id uuid NOT NULL,
	amount_cents bigint NOT NULL,
	commission_cents bigint NOT NULL,
	seller_amount_cents bigint NOT NULL,
	currency character varying NOT NULL,
	status character varying NOT NULL,
	stripe_customer_id character varying,
	stripe_payment_intent_id character varying,
	stripe_transfer_id character varying,
	created_at timestamp without time zone DEFAULT now(),
	confirmed_at timestamp without time zone,
	CONSTRAINT orders_pkey PRIMARY KEY (id),
	CONSTRAINT orders_buyer_id_fkey FOREIGN KEY (buyer_id) REFERENCES users(id),
	CONSTRAINT orders_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES users(id),
	CONSTRAINT orders_asset_id_fkey FOREIGN KEY (asset_id) REFERENCES assets(id)
);
