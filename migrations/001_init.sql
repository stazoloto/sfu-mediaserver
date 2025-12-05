CREATE TABLE rooms (
    id BIGINT PRIMARY KEY CHECK (
        id >= 10000000000000 AND
        id <= 99999999999999  
    ),
    name VARCHAR(255),
    owner_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL
)

CREATE TABLE peers (
    id INTEGER PRIMARY KEY CHECK (
        id >= 10000000 AND
        id <= 99999999
    ),
    name VARCHAR(255),
    role VARCHAR(10) NOT NULL,

)

