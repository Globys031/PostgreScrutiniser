export const loginDetails = {
  server: Cypress.env("CYPRESS_SERVER"),
  username: Cypress.env("CYPRESS_USERNAME"),
  password: Cypress.env("CYPRESS_PASSWORD"),
  port: Cypress.env("CYPRESS_BACKEND_PORT"),
};
