// Authentication details to be reused throughout the application
export const loginDetails = {
  server: Cypress.env("CYPRESS_SERVER"),
  username: Cypress.env("CYPRESS_USERNAME"),
  password: Cypress.env("CYPRESS_PASSWORD"),
  port: Cypress.env("CYPRESS_BACKEND_PORT"),
};

// Wait 3 seconds to avoid crashing PostgreSQL instance
export const waitForPostgresql = 3000;
