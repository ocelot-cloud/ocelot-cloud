let authCookie = "";

function createApp() {
    cy.get('#input-app').type('myapp');
    cy.get('#button-create-app').click();
}

function cancelAccountDeletion() {
    cy.get('#dropdown').click();
    cy.get('#button-delete-account').click();
    cy.get('#button-delete-cancel').click();
    cy.url().should('eq', 'http://localhost:8081/hub');
}

function executeAccountDeletion() {
    cy.get('#dropdown').click();
    cy.get('#button-delete-account').click();
    cy.get('#button-delete-confirmation').click();
    cy.url().should('eq', 'http://localhost:8081/hub/login');
}

function assertEmptyAppList() {
    cy.get('#app-list').find('li').should('have.length', 0);
    cy.get('#button-delete-app').should('not.exist')
    cy.get('#button-edit-tags').should('not.exist')
}

function assertIsSelected(isSelected: boolean) {
    let prefix = ""
    if (!isSelected) {
        prefix = "not."
    }
    cy.get('#button-delete-app').should(prefix + 'exist')
    cy.get('#button-edit-tags').should(prefix + 'exist')
    cy.get('#app-list').find('li')
        .should('have.length', 1).and('contain', 'myapp').and(prefix + 'have.class', 'active')
}

function clickOnApp() {
    cy.get('#app-list').find('li').click()
    // TODO .should('have.class', 'active') or non.have.class active?
}

function deleteApp() {
    cy.get('#app-list').find('li').click().should('have.class', 'active')
    cy.get('#button-delete-app').click();
}

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
        assertEmptyAppList()
        createApp()
        assertIsSelected(false)
        clickOnApp()
        assertIsSelected(true)
        clickOnApp()
        assertIsSelected(false)
        deleteApp()
        assertEmptyAppList();
    });

    it('change password', () => {
        login()
        cy.get('#dropdown').click();
        cy.get('#button-change-password').invoke('trigger', 'click');
        cy.url().should('eq', 'http://localhost:8081/hub/change-password');
    });

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

    it('check upload', () => {
        login()
        createApp()
        cy.get('#app-list').find('li').click()
        cy.get('#button-edit-tags').click()
        // TODO check URL
        cy.get('input[type="file"]').selectFile({
            contents: Cypress.Buffer.from(''),
            fileName: '1.4.tar.gz',
        }, { force: true }) // "force" is necessary, since the actual <input> is invisible for beauty reasons.
        cy.get('#tag-list').find("li").should('have.text', '1.4')
    });

    it('test delete account', () => {
        login()
        cancelAccountDeletion();
        executeAccountDeletion();
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