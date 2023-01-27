CREATE TABLE IF NOT EXISTS holders (
    id bigserial PRIMARY KEY,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    address text NOT NULL,
    commitment_score numeric NOT NULL,
    portfolio_score numeric NOT NULL,
    trading_score numeric NOT NULL,
    version integer NOT NULL DEFAULT 1
)