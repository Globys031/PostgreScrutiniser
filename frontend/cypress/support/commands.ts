/// <reference types="cypress" />
// ***********************************************
// This example commands.ts shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
//
// declare global {
//   namespace Cypress {
//     interface Chainable {
//       login(email: string, password: string): Chainable<void>
//       drag(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       dismiss(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       visit(originalFn: CommandOriginalFn, url: string, options: Partial<VisitOptions>): Chainable<Element>
//     }
//   }
// }

import { loginDetails } from "../utils";

const login = () => {
  cy.session(
    "loginID",
    () => {
      cy.visit("/login");

      cy.get("#hostname").type(loginDetails.server);
      cy.get("#username").type(loginDetails.username);
      cy.get("#password").type(loginDetails.password);

      cy.get("#login").click();
      cy.get("#log_out").should("be.visible");
    },
    {
      validate: () => {
        expect(localStorage.token).to.exist;
      },
    }
  );
}

const resetConfigs = () => {
  cy.request({
    method: "DELETE",
    url: `http://${loginDetails.server}:${loginDetails.port}/api/resource`,
    headers: {
      Authorization: `Bearer ${localStorage.token}`,
    },
  });
};

Cypress.Commands.addAll({
  login,
  resetConfigs,
});

declare global {
  namespace Cypress {
    interface Chainable {
      login(): Chainable<void>;
      resetConfigs(): Chainable<void>;
    }
  }
}

export {};
