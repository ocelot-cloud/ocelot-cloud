import {
  assertColumnTitles,
  stopAllRunningStacks,
  StackOperator,
  VisitHomePage
} from "./StackOperator";
import {
  isFrontendMocked,
} from "./Config";

describe('template spec', () => {

  beforeEach(() => {
    VisitHomePage()
  })

  it('verify state lifecycle', () => {
    new StackOperator('nginx-default')
        .assertState('Uninitialized')
        .operate('start')
        .assertState('Available')
        .operate('stop')
        .assertState('Uninitialized')
  });

  it('assert column titles', () => {
    assertColumnTitles();
  });

  it('assert core services not listed', () => {
    new StackOperator('ocelot-cloud')
        .shouldStackBeListed(false)
  });

  it('verify operations on stacks', () => {
    if (!isFrontendMocked) {
      new StackOperator('nginx-default')
          .assertState('Uninitialized')
          .operate('start')
          .assertState('Available')
          .assertWebsiteContent('nginx index page')
          .operate('stop')
          .assertState('Uninitialized')
    }
  });

  it('check open-button urls', () => {
    new StackOperator('nginx-custom-path').assertOpenButtonUrlPath('/custom-path')
    new StackOperator('nginx-default').assertOpenButtonUrlPath('/')
  })

  it('should verify that stack names are in alphabetical order', () => {
    new StackOperator('nginx-default')
        .assertStackNameAlphabeticalOrder()
        .operate('start')
        .assertStackNameAlphabeticalOrder()
        .operate('stop')
        .assertStackNameAlphabeticalOrder()
        .waitSeconds(1)
  });

  it('check whether custom ports and proxying to stacks work', () => {
    if (!isFrontendMocked) {
      new StackOperator('nginx-custom-port')
          .operate("start")
          .assertWebsiteContent("nginx custom port")
          .operate('stop')
    }
  });

  it('verify button disabling based on current state', () => {
    let operator = new StackOperator('nginx-default')
    operator
        .assertState('Uninitialized')
        .shouldButtonBeEnabled('Open', false)
        .shouldButtonBeEnabled('Start', true)
        .shouldButtonBeEnabled('Stop', false)
        .operate('start')
        .assertState('Available')
        .shouldButtonBeEnabled('Open', true)
        .shouldButtonBeEnabled('Start', false)
        .shouldButtonBeEnabled('Stop', true)
        .operate("stop")
        .assertState('Uninitialized')
        .shouldButtonBeEnabled('Open', false)
        .shouldButtonBeEnabled('Start', true)
        .shouldButtonBeEnabled('Stop', false)
  });

  afterEach(() => {
    stopAllRunningStacks()
  });
});
