let authCookie = "";

const hubWipePath = "http://localhost:8082/wipe-data"
const cloudBaseUrl = "http://localhost:8081"
const homePath = `${cloudBaseUrl}/hub`
const loginPath = `${cloudBaseUrl}/hub/login`
const registrationPath = `${cloudBaseUrl}/hub/registration`
const changePasswordPath = `${cloudBaseUrl}/hub/change-password`
const tagManagementPath = `${cloudBaseUrl}/hub/tag-management`

function createApp() {
    cy.get('#input-app').type('myapp');
    cy.get('#button-create-app').click();
}

function cancelAccountDeletion() {
    cy.get('#dropdown').click();
    cy.get('#button-delete-account').click();
    cy.get('#button-delete-cancel').click();
    cy.get('#button-delete-cancel').should('not.exist');
    cy.get('#button-delete-confirmation').should('not.exist');
    cy.url().should('eq', homePath);
}

function executeAccountDeletion() {
    cy.get('#dropdown').click();
    cy.get('#button-delete-account').click();
    cy.get('#button-delete-confirmation').click();
    cy.get('#button-delete-cancel').should('not.exist');
    cy.get('#button-delete-confirmation').should('not.exist');
    cy.url().should('eq', loginPath);
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
        if (url != homePath)
        cy.visit(homePath)
    });
    cy.get('#app-list').find('li').click()
    cy.get('#button-delete-app').click();
    cy.get('#button-delete-confirmation').click();
    cy.get('#button-delete-cancel').should('not.exist');
    cy.get('#button-delete-confirmation').should('not.exist');
}

function tryToDeleteAppButCancelInConfirmationPopup() {
    cy.get('#app-list').find('li').click()
    cy.get('#button-delete-app').click();
    cy.get('#button-delete-cancel').click();
    cy.get('#button-delete-cancel').should('not.exist');
    cy.get('#button-delete-confirmation').should('not.exist');
    cy.get('#app-list').find('li').click()
}

function checkInputValidationOnLoginPage() {
    cy.visit(loginPath);
    cy.get('#input-username').type('ad');
    cy.get('#input-password').type('pass');
    cy.get('#input-username').should('not.have.class', 'is-invalid');
    cy.get('#input-password').should('not.have.class', 'is-invalid');

    cy.get('#button-login').click();
    cy.get('#input-username').should('have.class', 'is-invalid');
    cy.get('#input-password').should('have.class', 'is-invalid');
    cy.get('body').should('contain.text', 'Invalid username,');
    cy.get('body').should('contain.text', 'Invalid password,');

    cy.get('#input-username').clear().type('admin');
    cy.get('#input-password').clear().type('password');
    cy.get('#input-username').should('not.have.class', 'is-invalid');
    cy.get('#input-password').should('not.have.class', 'is-invalid');
}

function checkInputValidationOnRegistrationPage() {
    cy.visit(registrationPath);
    cy.get('#input-username').type('ad');
    cy.get('#input-password').type('pass');
    cy.get('#input-email').type('a@a');
    cy.get('#input-username').should('not.have.class', 'is-invalid');
    cy.get('#input-password').should('not.have.class', 'is-invalid');
    cy.get('#input-email').should('not.have.class', 'is-invalid')

    cy.get('#button-register').click();
    cy.get('#input-username').should('have.class', 'is-invalid');
    cy.get('#input-password').should('have.class', 'is-invalid');
    cy.get('#input-email').should('have.class', 'is-invalid');
    cy.get('body').should('contain.text', 'Invalid username,');
    cy.get('body').should('contain.text', 'Invalid password,');
    cy.get('body').should('contain.text', 'Invalid email,');

    cy.get('#input-username').clear().type('admin');
    cy.get('#input-password').clear().type('password');
    cy.get('#input-email').clear().type('admin@admin.de');
    cy.get('#input-username').should('not.have.class', 'is-invalid');
    cy.get('#input-password').should('not.have.class', 'is-invalid');
    cy.get('#input-email').should('not.have.class', 'is-invalid');
}

function checkInputValidationOnAppPage() {
    login()
    cy.get('#input-app').type('ad');
    cy.get('#input-app').should('not.have.class', 'is-invalid');
    cy.get('#button-create-app').click();
    cy.get('#input-app').should('have.class', 'is-invalid');
    cy.get('body').should('contain.text', 'Invalid app,');
    cy.get('#input-app').clear().type('asdf');
    cy.get('#input-app').should('not.have.class', 'is-invalid');
}

function checkInputValidationOnTagPage() {
    cy.reload()
    createApp()
    clickOnApp()
    cy.get('#button-edit-tags').click()
    cy.get('input[type="file"]').selectFile({
        contents: Cypress.Buffer.from(''),
        fileName: 'as.tar.gz',
    }, {force: true})
    cy.contains('Invalid tag,').should('be.visible');
    cy.get('input[type="file"]').selectFile({
        contents: Cypress.Buffer.from(''),
        fileName: 'asdf.tar.gz',
    }, {force: true})
    cy.should('not.contain', 'Invalid tag,');
    cy.get('#tag-list').find('li').click()
    cy.get('#button-delete-tag').click()
    cy.get('#button-delete-confirmation').click()
}

describe('Hub Operations', () => {
    it('register and login', () => {
        cy.request(hubWipePath)
        cy.visit(homePath);
        cy.url().should('eq', loginPath);
        cy.get('#registration-redirect').click();
        cy.url().should('eq', registrationPath);
        cy.get('#input-username').type('admin');
        cy.get('#input-password').type('password');
        cy.get('#input-email').type('admin@admin.com');
        cy.get('#button-register').click();
        cy.url().should('eq', loginPath);
        login()
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
        tryToDeleteAppButCancelInConfirmationPopup()
        deleteApp()
        assertEmptyAppList();
    });

    it('check input validation', () => {
        checkInputValidationOnLoginPage();
        checkInputValidationOnRegistrationPage();
        checkInputValidationOnAppPage();
        checkInputValidationOnTagPage();
    });

    it('check upload', () => {
        login()
        createApp()
        cy.get('#app-list').find('li').click()
        cy.get('#button-edit-tags').click()

        cy.get('#tag-list').find('li').should('have.length', 0);
        cy.get('#button-download-tag').should('not.exist')
        cy.get('#button-delete-tag').should('not.exist')

        cy.url().should('contain', tagManagementPath)
        cy.get('input[type="file"]').selectFile({
            contents: Cypress.Buffer.from(''),
            fileName: '1.4.tar.gz',
        }, { force: true }) // "force" is necessary, since the actual <input> is invisible for beauty reasons.

        cy.get('#tag-list').find('li').should('have.length', 1);
        cy.get('#button-download-tag').should('not.exist')
        cy.get('#button-delete-tag').should('not.exist')
        cy.get('#button-delete-cancel').should('not.exist');
        cy.get('#button-delete-confirmation').should('not.exist');

        cy.get('#tag-list').find("li").should('have.text', '1.4').click()

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
        cy.url().should('eq', changePasswordPath);
    });

    it('check logout', () => {
        login()
        cy.get('#dropdown').click();
        cy.get('#button-logout').click();
        cy.visit(homePath)
        cy.url().should('eq', loginPath);

        authCookie = ""
    });

    it('check wrong password prevents login', () => {
        cy.visit(loginPath);
        cy.get('#input-username').type('admin');
        cy.get('#input-password').type('password+x');
        cy.get('#button-login').click();
        cy.url().should('eq', loginPath);
    });

    it('test delete account', () => {
        login()
        cancelAccountDeletion();
        executeAccountDeletion();
    });
});

function login() {
    if(authCookie == "") {
        cy.visit(loginPath);
        cy.get('#input-username').clear().type('admin');
        cy.get('#input-password').clear().type('password');
        cy.get('#button-login').click();
        cy.url().should('eq', loginPath)
        cy.get('#user-label').should('contain', 'admin');
        cy.getCookie("auth").should('exist').then((cookie) => {
            authCookie = cookie.value
        })
    } else {
        cy.setCookie("auth", authCookie)
        cy.visit(homePath)
    }
}

