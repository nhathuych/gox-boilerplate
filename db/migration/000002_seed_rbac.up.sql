INSERT INTO roles (id, name)
VALUES
  (1, 'admin'),
  (2, 'user');

SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));

INSERT INTO permissions (id, name)
VALUES
  (1, 'article:create'),
  (2, 'article:read'),
  (3, 'article:update'),
  (4, 'article:delete'),
  (5, 'article:publish');

SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));

INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id
FROM permissions;

INSERT INTO role_permissions (role_id, permission_id)
VALUES
  (2, 1),
  (2, 2),
  (2, 3);

-- Default admin: password = "ChangeMe123!" (bcrypt cost 10)
INSERT INTO users (email, password_hash, role_id)
VALUES (
  'admin@example.com',
  '$2a$10$JTjtSsAh.Wr19l/ajigKJOFBL03LJYwQKCmMd/DaqIn8wvdWGu4Z2',
  1
);
