import {
  assertColumnTitles,
  stopAllRunningStacks,
  StackOperator,
  VisitHomePage
} from "./StackOperator";
import {
  areRealDockerContainersUsed,
} from "./Config";

// TODO It would be nice to have the option to run a specific single test. To be implemented in CI runner.
// TODO It would also be faster, to have an option to skip the rebuild during acceptance testing, when the code did not change.
describe('template spec', () => {

  beforeEach(() => {
    VisitHomePage()
  })

  it('assert column titles', () => {
    assertColumnTitles();
  });

  it('assert core services not listed', () => {
    new StackOperator('ocelot-cloud')
        .shouldStackBeListed(false)
  });

  it('verify operations on stacks', () => {
    new StackOperator('nginx-default')
        .operate('start')
        .assertStateWithLongWaitingTime('Available')
        .waitSeconds(2)
        .assertWebsiteContent('nginx index page')
        .operate('stop')
        .assertState('Uninitialized')
  });

  it('check open-button urls', () => {
    new StackOperator('nginx-custom-path').assertOpenButtonUrlPath('/custom-path')
    new StackOperator('nginx-default').assertOpenButtonUrlPath('/')
  })

  it('should verify that stack names are in alphabetical order', () => {
    new StackOperator('nginx-default')
        .assertStackNameAlphabeticalOrder()
        .operate('start')
        .waitSeconds(2)
        .assertStackNameAlphabeticalOrder()
        .operate('stop')
        .waitSeconds(2)
        .assertStackNameAlphabeticalOrder()
        .waitSeconds(2)
  });

  it('verify state lifecycle', () => {
    if (areRealDockerContainersUsed) {
      new StackOperator('nginx-slow-start')
          .assertStateWithLongWaitingTime('Uninitialized')
          .shouldProcessingAnimationBeVisible(false)
          .operate('start')
          .assertStateWithLongWaitingTime('Starting')
          .shouldProcessingAnimationBeVisible(true)
          .assertStateWithLongWaitingTime('Available')
          .shouldProcessingAnimationBeVisible(false)
          .operate('stop')
          .assertStateWithLongWaitingTime('Stopping')
          .shouldProcessingAnimationBeVisible(true)
          .assertStateWithLongWaitingTime('Uninitialized')
          .shouldProcessingAnimationBeVisible(false)
    } else {
      new StackOperator('nginx-slow-start')
          .assertState('Uninitialized')
          .operate('start')
          .assertState('Available')
          .operate('stop')
          .assertState('Uninitialized')
    }
  });

  it('check whether custom ports and proxying to stacks work', () => {
    if (areRealDockerContainersUsed) {
      new StackOperator('nginx-custom-port')
          .operate("start")
          .waitSeconds(2)
          .assertWebsiteContent("nginx custom port")
    }
  });

  // TODO This test failed once on a fresh installation, but worked on seconds run. Flaky test should ot exist. To be researched.
  it('verify button disabling based on current state', () => {
    let operator = new StackOperator('nginx-slow-start')
    operator
        .assertStateWithLongWaitingTime('Uninitialized')
        .shouldButtonBeEnabled('Open', false)
        .shouldButtonBeEnabled('Start', true)
        .shouldButtonBeEnabled('Stop', false)
        .operate('start')

    if (areRealDockerContainersUsed) {
      operator
        .assertStateWithLongWaitingTime('Starting')
        .shouldButtonBeEnabled('Open', false)
        .shouldButtonBeEnabled('Start', false)
        .shouldButtonBeEnabled('Stop', false)
    }

    operator.assertState('Available')
        .shouldButtonBeEnabled('Open', true)
        .shouldButtonBeEnabled('Start', false)
        .shouldButtonBeEnabled('Stop', true)
        .operate("stop")

    if (areRealDockerContainersUsed) {
      operator.assertStateWithLongWaitingTime('Stopping')
          .shouldButtonBeEnabled('Open', false)
          .shouldButtonBeEnabled('Start', false)
          .shouldButtonBeEnabled('Stop', false)
    }

    operator.assertStateWithLongWaitingTime('Uninitialized')
        .shouldButtonBeEnabled('Open', false)
        .shouldButtonBeEnabled('Start', true)
        .shouldButtonBeEnabled('Stop', false)
  });

  afterEach(() => {
    stopAllRunningStacks()
  });
});
