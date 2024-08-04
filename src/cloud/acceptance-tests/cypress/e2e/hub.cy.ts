describe('Hub Registration and Login Flow', () => {
    it('should complete the registration and login process', () => {
        cy.visit('http://localhost:8081/hub');
        cy.url().should('eq', 'http://localhost:8081/hub/login');
        cy.get('#registration-redirect').click();
        cy.url().should('eq', 'http://localhost:8081/hub/registration');
        cy.get('#input-username').type('admin');
        cy.get('#input-password').type('password');
        cy.get('#input-email').type('admin@admin.com');
        cy.get('#button-register').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');

        cy.get('#input-username').type('admin');
        cy.get('#input-password').type('password+x');
        cy.get('#button-login').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');

        cy.get('#input-password').clear().type('password');
        cy.get('#button-login').click();
        cy.url().should('eq', 'http://localhost:8081/hub');

        cy.get('#button-logout').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
        login()
        cy.get('#button-delete-account').click();
        cy.get('#button-delete-cancel').click();
        cy.url().should('eq', 'http://localhost:8081/hub');
        cy.get('#button-delete-account').click();
        cy.get('#button-delete-confirmation').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
        // TODO Check that login fails?
    });
});

function login() {
    cy.get('#input-username').type('admin');
    cy.get('#input-password').type('password');
    cy.get('#button-login').click();
    cy.url().should('eq', 'http://localhost:8081/hub');
}