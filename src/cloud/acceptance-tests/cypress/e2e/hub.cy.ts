describe('Hub Registration and Login Flow', () => {
    it('should complete the registration and login process', () => {
        cy.visit('http://localhost:8081/hub');
        cy.url().should('eq', 'http://localhost:8081/hub/login');
        cy.get('#registration-redirect').click();
        cy.url().should('eq', 'http://localhost:8081/hub/registration');
        cy.get('input[name="username"]').type('admin');
        cy.get('input[name="password"]').type('password');
        cy.get('input[name="email"]').type('admin@admin.com');
        cy.get('button').contains('register').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
        cy.get('input[name="username"]').type('admin');
        cy.get('input[name="password"]').type('password');
        cy.get('button').contains('login').click();
        cy.url().should('eq', 'http://localhost:8081/hub/');
        cy.get('button').contains('logout').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
    });
});
