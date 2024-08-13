let authCookie = "";

describe('Hub Operations', () => {
    it('register and login', () => {
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
        login()
        cy.getCookie("auth").should('exist').then((c) => {
            authCookie = c.value
        })
    });

    it('create and delete app', () => {
        login()
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
    });

    /* TODO No idea why that fails. Clicking the button does not cause a redirect.
    it('change password', () => {
        login()
        cy.get('#dropdown').click();
        cy.get('#button-change-password').click();
        cy.url().should('eq', 'http://localhost:8081/hub/change-password');
    });
     */

    it('check logout', () => {
        login()
        cy.url().should('eq', 'http://localhost:8081/hub')
        cy.get('#dropdown').click();
        cy.get('#button-logout').click();
        cy.visit('http://localhost:8081/hub')
        cy.url().should('eq', 'http://localhost:8081/hub/login');
    });

    it('check wrong password prevents login', () => {
        cy.visit('http://localhost:8081/hub/login');
        cy.get('#input-username').type('admin');
        cy.get('#input-password').type('password+x');
        cy.get('#button-login').click();
        cy.url().should('eq', 'http://localhost:8081/hub/login');
    });

    it('test delete account', () => {
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
    cy.setCookie("auth", authCookie)
    cy.visit('http://localhost:8081/hub');
    cy.url().then((url) => {
        if (url != 'http://localhost:8081/hub') {
            cy.visit('http://localhost:8081/hub/login');
            cy.get('#input-username').clear().type('admin');
            cy.get('#input-password').clear().type('password');
            cy.get('#button-login').click();
            cy.url().should('eq', 'http://localhost:8081/hub');
        }
    });
}