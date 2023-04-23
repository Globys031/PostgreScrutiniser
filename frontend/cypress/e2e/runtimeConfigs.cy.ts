/* eslint-disable cypress/no-unnecessary-waiting */
import { loginDetails } from "../auth";

// Wait 3 seconds to avoid crashing PostgreSQL instance
const waitForPostgresql = 3000;

describe("RuntimeConfigs view", () => {
  beforeEach(() => {
    /**
      Using `beforeEach` for login instead of `before` because `localStorage`
      doesn't seem to be getting saved in-between tests even with the
      `cacheAcrossSpecs` options
    **/
    cy.login().then(() => {
      cy.request({
        method: "DELETE",
        url: `http://${loginDetails.server}:${loginDetails.port}/api/resource`,
        headers: {
          Authorization: `Bearer ${localStorage.token}`,
        },
      });
    });
    cy.visit("/configurations");
    cy.wait(waitForPostgresql);
  });

  it("'Get checks' mounts two collapsible tables", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div.text", "Suggestions").should("be.visible");
    cy.contains("div.text", "Checks that passed").should("be.visible");
  });

  it("'Get checks' returns a total of 13 parameters within two collapsible tables", () => {
    cy.visit("/configurations");
    cy.contains("div", "Get checks").click();
    cy.get("div.content")
      .find("div.collapsible")
      .its("length")
      .should("eq", 13);
  });

  it("'Reset Configurations' returns success notification", () => {
    cy.contains("div", "Reset Configurations").click();
    cy.contains(".modal .button-container div", "Reset").click();
    cy.get(".notifications-container").should(
      "contain",
      "Operation successful"
    );
    cy.wait(waitForPostgresql);
  });

  it("'Apply all suggestions' empties suggestions table", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div", "Apply all suggestions").click();
    cy.contains("div.text", "Suggestions (0)");
    cy.wait(waitForPostgresql);
  });

  it("'Apply all suggestions' should move suggestions to passed checks table", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div", "Apply all suggestions").click();
    cy.contains("div.text", "Checks that passed (13)");
    cy.wait(waitForPostgresql);
  });

  it("'Apply all suggestions' returns 'suggestions applied' success notification", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div", "Apply all suggestions").click();
    cy.get(".notifications-container").should("contain", "Suggestions applied");
    cy.wait(waitForPostgresql);
  });

  it("'Apply suggestion' should remove suggestions from suggestions table", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div.text", "Suggestions (").click();

    cy.get("div.content div.collapsible")
      .first()
      .invoke("text")
      .then((paramHeader) => {
        // Only suggestions use parameter names with parentheses
        const match = paramHeader.match(/^\w+ \(/);
        if (!match) {
          throw new Error("Could not get parameter name");
        }
        const paramName = match[0];

        cy.contains("div.collapsible", paramName).click();
        cy.contains("div", "Apply suggestion").click();

        // After applying parameter it should no longer be amongst suggestions
        cy.contains("div.collapsible", paramName).click();
        cy.contains("div.collapsible", paramName).should("not.exist");
      });
  });

  it("'Apply suggestion' should move suggestion to passed checks table", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div.text", "Suggestions (").click();

    cy.get("div.content div.collapsible")
      .first()
      .invoke("text")
      .then((paramHeader) => {
        // Only suggestions use parameter names with parentheses
        const match = paramHeader.match(/^\w+/);
        if (!match) {
          throw new Error("Could not get parameter name");
        }
        const paramNameInPassed = match[0];
        const paramNameInSuggestions = `${match[0]} (`;

        cy.contains("div.collapsible", paramNameInSuggestions).click();
        cy.contains("div", "Apply suggestion").click();

        // After applying parameter it should no longer be amongst suggestions
        cy.contains("div.collapsible", paramNameInSuggestions).should(
          "not.exist"
        );
        // Should instead be amongst passed checks
        cy.contains("div.collapsible", paramNameInPassed).should("exist");
      });
  });

  it("'Apply suggestion' returns 'suggestions applied' success notification", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div.text", "Suggestions (").click();

    cy.get("div.content div.collapsible")
      .first()
      .invoke("text")
      .then((paramHeader) => {
        // Only suggestions use parameter names with parentheses
        const match = paramHeader.match(/^\w+ \(/);
        if (!match) {
          throw new Error("Could not get parameter name");
        }
        const paramName = match[0];

        cy.contains("div.collapsible", paramName).click();
        cy.contains("div", "Apply suggestion").click();
        cy.get(".notifications-container").should(
          "contain",
          "Suggestions applied"
        );
        cy.wait(waitForPostgresql);
      });
  });

  it("'Apply %{n} suggestions' returns 'suggestions applied' success notification", () => {
    cy.contains("div", "Get checks").click();
    cy.contains("div.text", "Suggestions (").click();

    cy.get("div.content div.collapsible").first().click();
    cy.get("div.suggestion-container input").first().click();
    cy.contains("div", /Apply \d suggestions/).click();
    cy.get(".notifications-container").should("contain", "Suggestions applied");

    cy.wait(waitForPostgresql);
  });
});
