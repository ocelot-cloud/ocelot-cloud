# Ocelot-Cloud

## Introduction

Ocelot-Cloud is an open source digital infrastructure management platform that aims to make self-hosting as easy as possible. Information about the project can be found on our [website](https://ocelot-cloud.org). Visit the 'Contact' page if you want to get in touch, find out more about our other communication channels or join our community. You might also want to take a look at the [announcement video](https://youtu.be/WQwBYjMbe8I), which includes a technical demo of the prototype.

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

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
