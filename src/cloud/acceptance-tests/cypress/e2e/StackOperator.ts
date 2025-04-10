import {isFrontendMocked, ocelotUrl, isSecurityEnabled, rootDomain, scheme} from "./Config";

type ButtonType = 'Open' | 'Start' | 'Stop'
type StackOperation = 'start' | 'stop'
type ExpectedState = 'Uninitialized' | 'Starting' | 'Available' | 'Stopping'  | 'Downloading'

const tableSelector = 'table#stack-table';

export class StackOperator {
    private readonly stackName: string;

    constructor(stackName: string) {
        cy.visit(ocelotUrl);
        this.stackName = stackName;
    }

    private getStackRow() {
        return cy.get(tableSelector).contains('td', this.stackName).parent('tr');
    }

    shouldStackBeListed(shouldBeListed: boolean): StackOperator {
        if (shouldBeListed) {
            cy.get(tableSelector).should('contain', this.stackName);
        } else {
            cy.get(tableSelector).should('not.contain', this.stackName);
        }
        return this
    }

    operate(operation: StackOperation): StackOperator {
        this.getStackRow().within(() => {
            cy.get(`button.${operation}-button`).click();
        });
        return this;
    }

    shouldButtonBeEnabled(buttonType: ButtonType, enabled: boolean): StackOperator {
        this.getStackRow().find('button').filter((i, btn) => btn.textContent === buttonType)
            .should(enabled ? 'not.be.disabled' : 'be.disabled');
        return this;
    }

    shouldProcessingAnimationBeVisible(isVisible: boolean): StackOperator {
        if (isVisible) {
            cy.get('.spinner-border').should('exist');
        } else {
            cy.get('.spinner-border').should('not.exist');
        }
        return this;
    }

    assertState(expectedState: ExpectedState): StackOperator {
        const stateCell = this.getStackRow().find('.state-column').should('have.text', expectedState);
        switch (expectedState) {
            case 'Starting':
                stateCell.should('have.class', 'bg-warning');
                break;
            case 'Available':
                stateCell.should('have.class', 'bg-success');
                break;
            case 'Uninitialized':
                stateCell.should('have.class', 'bg-dark');
                break;
        }
        return this;
    }

    waitSeconds(seconds: number): StackOperator {
        if (isFrontendMocked) {
            cy.wait(seconds * 100);
        } else {
            cy.wait(seconds * 1000);
        }
        return this;
    }

    assertWebsiteContent(expectedContent: string): StackOperator {
        let stackUrl = `http://${this.stackName}.` + rootDomain;
        if (isFrontendMocked) {
            cy.exec(`curl ${stackUrl}`).then((response) => {
                expect(response.stdout).to.include(expectedContent);
            });
        }
        return this;
    }

    assertStackNameAlphabeticalOrder(): StackOperator {
        let stackNames = [];
        cy.get(`${tableSelector} tbody tr`).each(($row) => {
            cy.wrap($row).find('td:first').invoke('text').then((text) => {
                stackNames.push(text.trim());
            });
        }).then(() => {
            expect(stackNames).to.deep.equal([...stackNames].sort());
        });
        return this;
    }

    assertOpenButtonUrlPath(urlPath: string): StackOperator {
        cy.get(`#open-button-${this.stackName}`).should('have.attr', 'data-stack-url',  `${scheme}://${this.stackName}.${rootDomain}${urlPath}`);
        return this;
    }
}

export function assertColumnTitles() {
    const checkIfColumnExists = (title: string) => {
        cy.get(`${tableSelector} th`).contains(title);
    };
    ['Name', 'State', 'Link', 'Actions'].forEach(checkIfColumnExists);
}

export function VisitHomePage() {
    cy.visit(ocelotUrl);
    if (isSecurityEnabled) {
        cy.url().should('include', '/login');
        cy.get('#username-field').type('admin');
        cy.get('#password-field').type('password');
        cy.get('#login-button').click();
    }
}

export function stopAllRunningStacks() {
    cy.get(`${tableSelector} tbody tr`).each($row => {
        const state = $row.find('.state-column').text().trim();
        if (state !== 'Uninitialized' && !$row.find('.stop-button').is(':disabled')) {
            cy.wrap($row).find('.stop-button').click();
        }
    });
}
