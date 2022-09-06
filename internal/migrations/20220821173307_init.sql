-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app_user (
    id serial 			PRIMARY KEY,
    login 				VARCHAR (255) NOT NULL,
    password 			VARCHAR (255) NOT NULL,
    accrual_current 	FLOAT default 0,
    accrual_total 		FLOAT default 0);
CREATE TABLE IF NOT EXISTS app_order (
    id serial 	PRIMARY KEY,
    number 		VARCHAR (255) NOT NULL,
    user_id 	INT NOT NULL,
    status 		VARCHAR (255) NOT NULL,
    accrual 	FLOAT,
    uploaded_at TIMESTAMP default NOW());
CREATE TABLE IF NOT EXISTS withdraw_list (
    id serial 		PRIMARY KEY,
    order_num	 	VARCHAR (255) NOT NULL,
    sum_points	 	FLOAT,
    user_id 		INT NOT NULL,
    processed_at 	TIMESTAMP default NOW());


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE app_user IF EXISTS;
DROP TABLE app_order IF EXISTS;
DROP TABLE withdraw_list IF EXISTS;
-- +goose StatementEnd
