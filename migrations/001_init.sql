CREATE TABLE rooms (
    id BIGINT PRIMARY KEY CHECK (
        id >= 10000000000000 AND
        id <= 99999999999999  
    ),
    name VARCHAR(255),
    owner_id INTEGER,
    is_active BOOLEAN NOT NULL

    -- Время жизни
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NOT NULL
)

CREATE TABLE peers (
    id INTEGER PRIMARY KEY CHECK (
        id >= 10000000 AND
        id <= 99999999
    ),
    room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    name VARCHAR(255),
    role VARCHAR(10) NOT NULL,
    join_time TIMESTAMPTZ NOT NULL
)

CREATE TABLE tracks (
    id BIGINT PRIMARY KEY CHECK (
        id >= 1000000000 AND
        id <= 9999999999
    ),
    peer_id NOT NULL INTEGER REFERENCES peers(id) ON DELETE CASCADE,
    room_id NOT NULL BIGINT REFERENCES rooms(id) ON DELETE CASCADE,
    kind VARCHAR(10) NOT NULL,
    is_active BOOLEAN NOT NULL
)

CREATE TABLE peer_connections (
        id BIGINT PRIMARY KEY CHECK (
        id >= 1000000000 AND
        id <= 9999999999
    ),
    peer_id NOT NULL INTEGER REFERENCES peers(id) ON DELETE CASCADE,
    room_id NOT NULL BIGINT REFERENCES rooms(id) ON DELETE CASCADE,
    signaling_state VARCHAR(50) NOT NULL,
    ice_state VARCHAR(50) NOT NULL,
    connection_state VARCHAR(50) NOT NULL,
    gathering_state VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    connected_at TIMESTAMPTZ NULL,
    last_activity TIMESTAMPTZ NOT NULL
)

