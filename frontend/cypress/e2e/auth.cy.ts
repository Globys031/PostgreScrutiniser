import { loginDetails } from "../utils";

beforeEach(() => {
  cy.visit("/login");
});

describe("Login view", () => {
  it("successful 'login' navigates to home page", () => {
    cy.visit("/login");

    // Enter valid login details
    cy.get("#hostname").type(loginDetails.server);
    cy.get("#username").type(loginDetails.username);
    cy.get("#password").type(loginDetails.password);

    cy.get("#login").click();

    // If redirected, path should be `hostname/`
    cy.location().should((loc) => {
      expect(loc.pathname).to.eq("/");
    });
  });

  it("successful 'login' displays logout button", () => {
    // Enter valid login details
    cy.get("#hostname").type(loginDetails.server);
    cy.get("#username").type(loginDetails.username);
    cy.get("#password").type(loginDetails.password);

    cy.get("#login").click();
    cy.get("#log_out").should("be.visible");
  });

  it("incorrect name gives authentication error", () => {
    // Enter only an invalid username
    cy.get("#hostname").type(loginDetails.server);
    cy.get("#username").type("user1234");
    cy.get("#password").type(loginDetails.password);

    cy.get("#login").click();
    cy.get(".error-msg").should(
      "contain",
      "Username should be the name of our main application user"
    );
  });

  it("incorrect password gives authentication error", () => {
    // Enter only an invalid password
    cy.get("#hostname").type(loginDetails.server);
    cy.get("#username").type(loginDetails.username);
    cy.get("#password").type("password1234");

    cy.get("#login").click();
    cy.get(".error-msg").should("contain", "Incorrect user password");
  });

  it("incorrect server gives authentication network error", () => {
    // Enter only an invalid password
    cy.get("#hostname").type("example.hostname.com");
    cy.get("#username").type(loginDetails.username);
    cy.get("#password").type(loginDetails.password);

    cy.get("#login").click();
    cy.get(".error-msg").should("be.visible");
  });
});
