-- +goose Up
-- +goose StatementBegin
CREATE TABLE moderated_requests (
    id serial PRIMARY KEY,
    kind VARCHAR(50) NOT NULL,
    state VARCHAR(50) NOT NULL,
    data json NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE moderated_requests;
-- +goose StatementEnd
