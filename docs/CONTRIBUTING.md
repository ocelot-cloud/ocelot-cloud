# Contributing

### Report bugs

If you find a bug, please open an issue on GitHub with a clear title and detailed description. Include the steps to reproduce the bug, the expected result, and the actual result. If possible, include screenshots or logs. Tag the issue as "bug".

### Feature proposals

If you have an idea for improving the project, you can submit your suggestion as a GitHub issue. Describe the feature, its benefits, and if possible, how it could be implemented. Label the issue as a "feature proposal".

### Create a pull request

To contribute code, fork the repository, create a new branch, make your changes, and then submit a pull request (PR). Branches should have a name of the following form `<issueNumber>-<descriptiveTitle>`. Make sure your PR has a clear title, a brief description of the changes, and any relevant issue numbers. All code contributions should pass existing tests and include new tests for added features.

## Contribution Model

We want to be transparent and clear about our contribution and business model. This section explains them in detail and reflects our philosophy and plans for the future. Currently there is no enterprise edition and the community edition is in prototype stage.

### Why use AGPL + CLA?

The Ocelot-Cloud community edition is open source. But using a permissive open source license for it, such as the MIT license, would allow other companies to take the code, build commercial products from it, and profit from our and the community's efforts without any obligation to contribute money or code in return. This would also undermine our ability to generate revenue. Therefore, we chose AGPL along with a Contributor License Agreement (CLA) to retain the exclusive right to create a commercial proprietary enterprise edition of Ocelot-Cloud, ensuring that revenues are reinvested in the development of both the community and enterprise editions. It also gives us time to implement the enterprise edition without rushing, and to remain independent of private investment that would turn this project into a profit-making venture instead of its originally intended focus on the public good.

**How does this work?** The AGPL license requires that any derivative product also be licensed under the AGPL to ensure that it remains open source. This mechanism is called "copyleft". This rule ensures that only open source editions of Ocelot-Cloud will exist. However, if we own the source code, either by writing it ourselves or by legal transfer of ownership, we have the option to release it under a different license. We can offer both an open source community edition and a proprietary enterprise edition with additional features that generate revenue to fund the project.

### Why is the CLA Important?

When you contribute to Ocelot-Cloud, you retain ownership of your code. Since it is licensed under the AGPL, this prevents us from relicensing your code. To allow us to dual-license and to ensure the financial sustainability of the project, we ask you to sign a CLA. This transfers ownership of your contributions to us, allowing us to offer both an open source community edition and a proprietary enterprise edition. The community edition and your contributions remain open source and can be used as such and forked at any time. We believe this approach strikes a fair balance between the interests of the community and the business interests.

### What does that mean for me as a user of Ocelot-Cloud?

If you use Ocelot-Cloud without modifying its code, you can use it just like any other open source software. However, if you modify the source code, your modifications must also be licensed under the AGPL, and you must make the source code available to all users of the modified version. If you are the only user, you do not need to take any action. Otherwise, there are two ways to make the source code available:

1. Contribute Back: If you enhance the software (e.g., new features, bug fixes), we recommend contributing your changes to this GitHub repository so that everyone can benefit from your efforts.

2. Distribute Yourself: You can distribute the modified code for example through your own source code repository.

### What does that mean for me as contributor of Ocelot-Cloud?

All contributors are required to read the [Ocelot-Cloud Contributor License Agreement](../legal/cla.md) and sign it by copying and renaming the [CLA Signing Template](../legal/cla-signing-template.md) to `legal/contributors/<your-github-name>.md`. Fill out the newly created signing template and submit it as part of your contribution. The CLA signing will automatically apply to all subsequent contributions, so it only needs to be done once for the first contribution. There is also a [sample file](../legal/contributors/sample.md) of what your CLA signing should look like in the end.

## Developer Guide

### Technological Decisions

This project is designed to make development as easy and accessible as possible. Notable decisions in this regard are:

- **Build on Docker**: Without any installations other than Docker, the entire project can be built from scratch and installed as described in the README. This gives users the ability to build their own Ocelot-Cloud images from source code if they prefer not to use the Ocelot-Cloud images uploaded to DockerHub.
- The **ci-runner** component is a tool written in Go that is capable of performing many CI operations such as build, test, deploy, and specifically provides very fine-grained access to the test suite according to your development needs.
- **Easy Tools**: Vue.js and Cypress are intuitive tools for web development and writing acceptance tests. Both technologies are configured to use Typescript because of its type safety. The backend and ci-runner language is Go, which is easy to learn and, in our opinion, pleasant to develop because of its lean design, and its executables have low resource requirements.
- The **CI pipelines**, or "Actions" in GitHub parlance, are enabled by the ci-runner and automatically execute the test suite for new commits pushed to GitHub. This enforces frequent quality assurance and early feedback for developers.

**Installation Requirements**

- Go 1.21
- Node.js 18.10
- Cypress Dependencies
- Docker and Docker Compose v2 (not v1 "docker-compose")
- Add these lines in the `/etc/hosts` file:

```
127.0.0.1       gitea.localhost
127.0.0.1       ocelot-cloud.localhost
```

You may need to add more addresses as tests are executed against them in the test suite. Browsers seem to automatically map `.localhost` addresses to IP 127.0.0.1, but the applications we build do not. Hence the customization of the `hosts` file.

### Automatic Development Preparation

There is a script that can install all the required development dependencies for the project on Debian 12 but may also work for other versions or Ubuntu. The user executing this script must have 'sudo' rights. The instructions for execution are:

```bash
git clone https://github.com/ocelot-cloud/ocelot-cloud
cd ocelot-cloud/scripts
bash install.sh
# Reboot computer
cd ocelot-cloud/src/ci-runner
go build && ./ci-runner test-ci
```

If you work on another system, the script may still be helpful for guiding you through installing some dependencies. In either case, you must manually add the lines to the hosts file.

### ci-runner

The ci-runner is designed to make the testing process as easy as possible. Here is a quick start that prints the available options:

```bash
cd src/ci-runner
go build && ./ci-runner --help
```

There are a lot of options. Feel free to experiment with them to get familiar.

### Development Mode

To facilitate the development of interactions between the frontend and the backend, the 'dev' mode has been introduced. This is a special setup in which the frontend and backend run separately, allowing Vue.js to run its own process that enables hot loading of source code changes, reducing development time. Open a console and run one of these commands to start the setup in development mode:

```bash
cd src/cloud/backend           && bash run-dev.sh
cd src/cloud/frontend          && bash run-dev.sh
cd src/cloud/acceptance-tests  && bash run-dev.sh
```

You can now visit the Vue.js GUI at `http://localhost.localdomain:8081/` and use the Cypress GUI to run acceptance tests if needed.

## Development Policies

### Lean Technology Stack

Stick to the languages Go and Typescript and the tools Vue.js, Cypress and Docker for all development to keep knowledge requirements for this project as low as possible. New technologies should only be introduced when there is a compelling need. Implement DevOps related code in the ci-runner tool which is written in Go and avoid bash scripts.

### Testing Strategy

Only frequently and extensively tested features are stable features. And in this project, only stable features may be merged. Every feature you add should be checked by passing tests. All tests must be included in the ci-runner in a job that is frequently executed by GitHub Actions.

What should be tested?

- 'happy path', if the execution of the production code flow works as expected
- 'unhappy path' tests, where invalid input, unexpected things, or errors occur during execution. The code must be able to handle errors correctly to be robust in production.
- border cases such as nil/null, zero, negative numbers, empty strings, empty collections, etc.
- Security and its testing must be built in from the beginning, where necessary. For example, a new REST API endpoint requires authentication, input validation etc.

Depending on the situation, one or more of the following types of tests may be necessary:

- Unit tests (testing a single unit) or core tests (wiring several units together and testing by interacting with a high-level API). I/O related dependencies should be mocked to make execution very fast.
- Component tests, i.e. running the backend and interacting with it via http requests
- Acceptance tests must be written using Cypress for all GUI features.

### Coding Practices

- The extensive test suite makes it unlikely that code changes will introduce new bugs. Therefore, any kind of code improvement is encouraged. If you accidentally break something, the tests will tell you, so you can fix it right away. Feel free to add any tests that may be missing.

- Contributions must be refactored until they are clean code, i.e.
  
  - readable: easy to understand
  
  - maintainable: easy to apply future changes, e.g. by avoiding duplications

- 3rd party libraries must be wrapped.
