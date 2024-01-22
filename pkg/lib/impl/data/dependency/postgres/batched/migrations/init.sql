CREATE TABLE dependencies
(
    id     varchar(255),
    status varchar(255),
    PRIMARY KEY (id)
);

CREATE TABLE groups
(
    id         varchar(255),
    pending    int,
    status     varchar(16),
    updated_at timestamp,
    PRIMARY KEY (id)
);

CREATE TABLE group_dependencies
(
    dependency_id varchar(255),
    group_id      varchar(255),
    PRIMARY KEY (dependency_id, group_id),
    CONSTRAINT fk_group
        FOREIGN KEY (group_id)
            REFERENCES groups (id)
            ON DELETE CASCADE
);