# Tasks:
# Check if I can set a fixed version on that OIDC plugin using the github repo url.

# version 3.1.4 is buggy: https://meta.discourse.org/t/after-rebuilding-app-cant-sign-in-with-error-500/290731/21
FROM docker.io/bitnami/discourse:3.1.3
RUN apt-get update && apt-get install -y curl
RUN git config --global --add safe.directory /opt/bitnami/discourse ; cd /opt/bitnami/discourse; RAILS_ENV=production bundle exec rake plugin:install repo=https://github.com/discourse/discourse-openid-connect

# Not necessary: docker exec discourse-discourse bash -c 'git config --global --add safe.directory /opt/bitnami/discourse ; cd /opt/bitnami/discourse; RAILS_ENV=production bundle exec rake plugin:install repo=https://github.com/discourse/discourse-openid-connect; RAILS_ENV=production bundle exec rake assets:precompile'

# https://meta.discourse.org/t/discourse-openid-connect/103632
# https://meta.discourse.org/t/install-plugins-in-discourse/19157
# https://meta.discourse.org/t/install-discourse-in-production-with-the-official-supported-instructions/142537
# https://docs.bitnami.com/aws/apps/discourse/configuration/install-plugins/
# https://stackoverflow.com/questions/73485958/how-to-correct-git-reporting-detected-dubious-ownership-in-repository-withou