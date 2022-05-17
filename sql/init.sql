CREATE TABLE cities
(
    name VARCHAR(25) NOT NULL,
    post_code INT NOT NULL,
    uuid uuid DEFAULT uuid_generate_v4 (),
    PRIMARY KEY (UUID)
);
