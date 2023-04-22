const appUser = {
  server: Cypress.env("CYPRESS_SERVER"),
  username: Cypress.env("CYPRESS_USERNAME"),
  password: Cypress.env("CYPRESS_PASSWORD"),
};

describe("Login view", () => {
  it("successful 'login' navigates to home page", () => {
    cy.visit("/login");

    // Enter valid login details
    cy.get("#hostname").type(appUser.server);
    cy.get("#username").type(appUser.username);
    cy.get("#password").type(appUser.password);

    cy.get("#login").click();

    // If redirected, path should be `hostname/`
    cy.location().should((loc) => {
      expect(loc.pathname).to.eq("/");
    });
  });

  it("successful 'login' displays logout button", () => {
    cy.visit("/login");

    // Enter valid login details
    cy.get("#hostname").type(appUser.server);
    cy.get("#username").type(appUser.username);
    cy.get("#password").type(appUser.password);

    cy.get("#login").click();
    cy.get("#log_out").should("be.visible");
  });

  it("incorrect name gives authentication error", () => {
    cy.visit("/login"); // Visit the login page

    // Enter only an invalid username
    cy.get("#hostname").type(appUser.server);
    cy.get("#username").type("user1234");
    cy.get("#password").type(appUser.password);

    cy.get("#login").click();
    cy.get(".error-msg").should(
      "contain",
      "Username should be the name of our main application user"
    );
  });

  it("incorrect password gives authentication error", () => {
    cy.visit("/login");

    // Enter only an invalid password
    cy.get("#hostname").type(appUser.server);
    cy.get("#username").type(appUser.username);
    cy.get("#password").type("password1234");

    cy.get("#login").click();
    cy.get(".error-msg").should("contain", "Incorrect user password");
  });

  it("incorrect server gives authentication network error", () => {
    cy.visit("/login");

    // Enter only an invalid password
    cy.get("#hostname").type("example.hostname.com");
    cy.get("#username").type(appUser.username);
    cy.get("#password").type(appUser.password);

    cy.get("#login").click();
    cy.get(".error-msg").should("be.visible");
  });
});
