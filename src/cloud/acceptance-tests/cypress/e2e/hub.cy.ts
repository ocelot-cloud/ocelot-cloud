describe('Hub Operations', () => {
    it('register, login and delete account', () => {
        cy.request("http://localhost:8082/wipe-data")
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
        cy.get('#input-password').clear().type('password');
        cy.get('#button-login').click();
        cy.url().should('eq', 'http://localhost:8081/hub');

        cy.get('#app-list').find('li').should('have.length', 0);
        cy.get('#button-delete-app').should('not.exist')
        cy.get('#button-edit-tags').should('not.exist')
        cy.get('#input-app').type('myapp');
        cy.get('#button-create-app').click();
        cy.get('#app-list').find('li')
            .should('have.length', 1).and('contain', 'myapp').and('not.have.class', 'active')
            .click().should('have.class', 'active')
        cy.get('#button-delete-app').should('exist')
        cy.get('#button-edit-tags').should('exist')
        cy.get('#app-list').find('li')
            .click().should('not.have.class', 'active')
        cy.get('#button-delete-app').should('not.exist')
        cy.get('#button-edit-tags').should('not.exist')
        cy.get('#app-list').find('li').click().should('have.class', 'active')
        cy.get('#button-delete-app').click();
        cy.get('#app-list').find('li').should('have.length', 0);


        // TODO Check that login fails?
    });

    it('change password', () => {
        console.log("todo");
    });

    it('check logout', () => {
        cy.visit('http://localhost:8081/hub');
        login()
        cy.get('#dropdown').click();
        cy.get('#button-logout').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');

        cy.get('#input-username').type('admin');
        cy.get('#input-password').type('password+x');
        cy.get('#button-login').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');

        login()
    });

    it('test delete account', () => {
        cy.visit('http://localhost:8081/hub');
        login()
        cy.get('#dropdown').click();
        cy.get('#button-delete-account').click();
        cy.get('#button-delete-cancel').click();
        cy.url().should('eq', 'http://localhost:8081/hub');

        cy.get('#dropdown').click();
        cy.get('#button-delete-account').click();
        cy.get('#button-delete-confirmation').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
    });
});

function login() {
    cy.get('#input-username').clear().type('admin');
    cy.get('#input-password').clear().type('password');
    cy.get('#button-login').click();
    cy.url().should('eq', 'http://localhost:8081/hub');
}