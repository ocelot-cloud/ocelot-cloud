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

function assertIsAppSelected(isSelected: boolean) {
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
}

function deleteApp() {
    cy.url().then(url => {
        if (url != "http://localhost:8081/hub")
        cy.visit("http://localhost:8081/hub")
    });
    cy.get('#app-list').find('li').click()
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
        assertIsAppSelected(false)
        clickOnApp()
        assertIsAppSelected(true)
        clickOnApp()
        assertIsAppSelected(false)
        deleteApp()
        assertEmptyAppList();
    });

    // TODO abstract urls and paths
    // TODO URL should be determined by the currently used origin/host
    it('check upload', () => {
        login()
        createApp()
        cy.get('#app-list').find('li').click()
        cy.get('#button-edit-tags').click()

        cy.get('#tag-list').find('li').should('have.length', 0);
        cy.get('#button-download-tag').should('not.exist')
        cy.get('#button-delete-tag').should('not.exist')

        cy.url().should('contain', 'http://localhost:8081/hub/tag-management')
        cy.get('input[type="file"]').selectFile({
            contents: Cypress.Buffer.from(''),
            fileName: '1.4.tar.gz',
        }, { force: true }) // "force" is necessary, since the actual <input> is invisible for beauty reasons.

        cy.get('#tag-list').find('li').should('have.length', 1);
        cy.get('#button-download-tag').should('not.exist')
        cy.get('#button-delete-tag').should('not.exist')

        cy.get('#tag-list').find("li").should('have.text', '1) 1.4').click()

        cy.get('#tag-list').find('li').should('have.length', 1);
        cy.get('#button-download-tag').should('exist')
        cy.get('#button-delete-tag').should('exist')

        cy.get('#button-delete-tag').click()
        cy.get('#button-delete-cancel').click()
        cy.get('#tag-list').find('li').should('have.length', 1);

        cy.get('#button-delete-tag').click()
        cy.get('#button-delete-confirmation').click()

        cy.get('#tag-list').find('li').should('have.length', 0);
        deleteApp()
    });

    it('change password', () => {
        login()
        cy.get('#dropdown').click();
        cy.get('#button-change-password').invoke('trigger', 'click');
        cy.url().should('eq', 'http://localhost:8081/hub/change-password');
    });

    it('check logout', () => {
        login()
        cy.get('#dropdown').click();
        cy.get('#button-logout').click();
        cy.visit('http://localhost:8081/hub')
        cy.url().should('eq', 'http://localhost:8081/hub/login');

        authCookie = ""
        login()
        cy.getCookie("auth").should('exist').then((c) => {
            authCookie = c.value
        })
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
        cancelAccountDeletion();
        executeAccountDeletion();
    });
});

function login() {
    if(authCookie == "") {
        cy.visit('http://localhost:8081/hub/login');
        cy.get('#input-username').clear().type('admin');
        cy.get('#input-password').clear().type('password');
        cy.get('#button-login').click();
    } else {
        cy.setCookie("auth", authCookie)
        cy.visit('http://localhost:8081/hub')
    }
    cy.url().should('eq', 'http://localhost:8081/hub')
    cy.get('#user-label').should('contain', 'admin');
}

// TODO When not authenticated on any "/hub", be directed to /hub/login
// TODO deleting apps should be confirmed
// TODO Display error from response as alert message
// TODO Maybe extract App Management Component
// TODO Make GUI pretty
// TODO Input validation