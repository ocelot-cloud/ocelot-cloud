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
        cy.get('#input-password').type('password');
        cy.get('#button-login').click();
        cy.url().should('eq', 'http://localhost:8081/hub');
        cy.get('#button-logout').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
    });
});
