import { waitForPostgresql } from "../utils";

/** In some cases API request for listing backups might finish 
  executing before intercept can get to it.
  Cypress doesn't allow exception catching so the best solution
  I could think of is to add an arbitrary timeout instead of using `cy.wait`.
**/
const waitForBackups = 2000;

describe("ConfigBackups view", () => {
  before(() => {
    cy.login().then(() => {
      cy.resetConfigs();
      cy.wait(waitForPostgresql);
    });
  })

  beforeEach(() => {
    /**
      Using `beforeEach` for login instead of `before` because `localStorage`
      doesn't seem to be getting saved in-between tests even with the
      `cacheAcrossSpecs` options
    **/
    cy.login().then(() => {
      cy.visit("/backups");
    });
  });

  it("'List backups' displays at least one backup", () => {
    cy.contains("div", "List backups").click();
    cy.wait(waitForBackups).then(() => {
      cy.get("div.collapsible")
      .its("length")
      .should("be.above", 0);
    })
  /** Left this for later in case I can figure how to deal
   * with interception happening after the api request finishes
   * 
   * 
   const listBackupsUrl = `http://${loginDetails.server}:${loginDetails.port}/api/backup`
    cy.intercept('GET', listBackupsUrl).as('listBackups').wait("@listBackups").then(() => {
    cy.get("div.collapsible")
    .its("length")
    .should("be.above", 0);
  })
  **/
  });

  /** Test checks that:
    1. Old backup name is no longer visible
    2. The amount of collapsible tables (backups) is the same
    This in turn means that a backup has been restored and a new one created
  */
  it("'Restore backup' creates another backup", () => {
    cy.contains("div", "List backups").click();
    cy.wait(waitForBackups).then(() => {
      let originalBackupCount = 0;
      cy.get("div.content").find("div.collapsible").its("length").then((backupCount) => {
        originalBackupCount = backupCount
      }) 

      cy.get("div.content div.collapsible")
      .first()
      .invoke("text")
      .then((backupName) => {
        cy.contains("div.collapsible", backupName).click();
        cy.contains("div", "Restore backup").click();
        cy.contains(".modal .button-container div", "Restore").click();

        // After applying backup it should no longer be amongst backups
        cy.contains("div.collapsible", backupName).click();
        cy.contains("div.collapsible", backupName).should("not.exist");

        // amount of collapsible tables should be the same
        cy.get("div.content").find("div.collapsible").its("length").should("eq", originalBackupCount);
      })
    })
  });

  it("'Restore backup' returns success notification", () => {
    cy.contains("div", "List backups").click();
    cy.wait(waitForBackups).then(() => {
      cy.get("div.content div.collapsible")
      .first()
      .click()

      cy.contains("div", "Restore backup").click();
      cy.contains(".modal .button-container div", "Restore").click();
      cy.get(".notifications-container").should(
        "contain",
        "Backups have been restored"
      );
    })
  });

  // If deleted, the backup should no longer be visible on the page
  it("'Delete backup' removes backup from view", () => {
    cy.resetConfigs(); // ensure a backup exists in the first place

    cy.contains("div", "List backups").click();
    cy.wait(waitForBackups).then(() => {
      cy.get("div.content div.collapsible")
      .first()
      .invoke("text")
      .then((backupName) => {
        cy.contains("div.collapsible", backupName).click();
        cy.contains("div", "Delete backup").click();
        cy.contains(".modal .button-container div", "Delete").click();
        cy.contains("div.collapsible", backupName).should("not.exist");
      })
    })
  });

  it("'Delete backup' returns success notification", () => {
    cy.resetConfigs(); // ensure a backup exists in the first place

    cy.contains("div", "List backups").click();
    cy.wait(waitForBackups).then(() => {
      cy.get("div.content div.collapsible")
      .first()
      .click()

      cy.contains("div", "Delete backup").click();
      cy.contains(".modal .button-container div", "Delete").click();
      cy.get(".notifications-container").should(
        "contain",
        "Backups have been deleted"
      );
    })
  });

  it("'Delete backups' removes all backups", () => {
    cy.resetConfigs(); // ensure a backup exists in the first place

    cy.contains("div", "Delete backups").click();
    cy.contains(".modal .button-container div", "Delete").click();
    cy.get("div.collapsible").should("not.exist")
  });

  it("'Delete backups' returns success notification", () => {
    cy.resetConfigs(); // ensure a backup exists in the first place

    cy.contains("div", "Delete backups").click();
    cy.contains(".modal .button-container div", "Delete").click();
    cy.get(".notifications-container").should(
      "contain",
      "Backups have been deleted"
    );
  });
});
