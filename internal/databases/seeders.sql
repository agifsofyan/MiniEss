INSERT INTO users (employee_id, email, name, password_hash, role)
VALUES (
  1,
  'agif@company.com',
  'Agif Sofyan',
  '$2a$10$HftxpN9C4P42Y4AkgS401e3TZEPcCBgjw/sJX0WkTvgwhSpIkWhg.',  -- hash dari "password123"
  'employee'
);