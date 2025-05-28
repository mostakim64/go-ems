CREATE DATABASE IF NOT EXISTS event_management;
USE event_management;


CREATE TABLE roles
(
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL
);


INSERT INTO roles(name)
VALUES ('ADMIN'),
       ('MANAGER'),
       ('ATTENDEE');



CREATE TABLE users
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    email      VARCHAR(50)  NOT NULL,
    password   VARCHAR(100) NOT NULL,
    first_name VARCHAR(50)  NOT NULL,
    last_name  VARCHAR(50)  NOT NULL,
    role_id    INT          NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_users_role_id FOREIGN KEY (role_id) REFERENCES roles (id)
);

-- admin's password -> 12345

INSERT INTO users(email, password, first_name, last_name, role_id)
VALUES ('admin@vivasoftltd.com', '$2a$10$2fBRiXac/mWv9m1n891zv.K1ooO1ItZtArxGqpO5qFEX6xgtgrDzu', 'Abdul', 'Mukit', 1);


CREATE TABLE permissions
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    permission  VARCHAR(20)  NOT NULL,
    description VARCHAR(100) NOT NULL DEFAULT '',
    created_at  DATETIME              DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT permission UNIQUE (permission)
);


CREATE TABLE role_permissions
(
    role_id       INT,
    permission_id INT,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT fk_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions (id)
);


INSERT INTO permissions (permission, description)
VALUES ('user.create', 'Permission to create a new user'),
       ('user.update', 'Permission to update an existing user'),
       ('user.fetch', 'Permission to fetch a specific user'),
       ('user.list', 'Permission to list users'),
       ('user.delete', 'Permission to delete a user'),
       ('event.create', 'Permission to create a new event'),
       ('event.update', 'Permission to update an existing event'),
       ('event.fetch', 'Permission to fetch a specific event'),
       ('event.list', 'Permission to list events'),
       ('event.delete', 'Permission to delete an event');


SELECT * FROM permissions;

INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 1);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 2);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 3);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 4);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 5);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 3);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 4);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 6);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 7);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 8);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 9);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (1, 10);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 6);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 7);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 8);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 9);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (2, 10);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (3, 8);
INSERT INTO event_management.role_permissions (role_id, permission_id) VALUES (3, 9);



SELECT permissions.id, permissions.permission, r.name
FROM role_permissions
         JOIN permissions on role_permissions.permission_id = permissions.id
         JOIN event_management.roles r on r.id = role_permissions.role_id


SELECT * FROM users;



CREATE TABLE events
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    location    VARCHAR(250),
    start_time  DATETIME,
    end_time    DATETIME,
    created_by  INT          NOT NULL,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_events_created_by FOREIGN KEY (created_by) REFERENCES users (id)
);