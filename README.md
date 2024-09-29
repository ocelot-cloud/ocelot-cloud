# Ocelot-Cloud

## Introduction

Ocelot-Cloud is an open-source digital infrastructure provisioning tool that aims to make self-hosting as easy as possible. In our [announcement video](https://youtu.be/WQwBYjMbe8I) you can get an idea of why this project was started, see a technical demo and see features that will be implemented in the future.

## Attention!

At its current prototype stage, Ocelot-Cloud is not ready for general use until the first official production release.

## Getting Started

This section explains how to build and run Ocelot-Cloud from scratch. 

**Requirements**

* git
* Docker (including Docker Compose)
* Linux shell

### Demo Version

```bash
apt update
apt install -y git docker.io docker-compose
git clone --depth=1 --branch=0.1.0 https://github.com/ocelot-cloud/ocelot-cloud
cd ocelot-cloud/scripts
bash run-demo.sh
```

Visit the application at `http://ocelot-cloud.localhost`. The login credentials for the login page are `admin` and `password`. On the home page, click the `Start` button in the `gitea` or `nocodb` row, for example, as these apps start up quite quickly. Other apps may take a few minutes to download and set up.

### Build From Scratch

Alternatively, you can build the image directly from source using these scripts instead:

```bash
cd ocelot-cloud/scripts
bash build.sh
bash run-demo.sh
```

## Further Information

The Contribution Guide and Code of Conduct can be found in the [docs](docs) folder. To find out more about the project, please visit our [website](https://ocelot-cloud.org). To get in touch or join our community, please visit the contact page.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
