FROM plugins/base:multiarch
MAINTAINER Jonas Franz <info@jonasfranz.de>

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/JonasFranzDEV/drone-crowdin.git"
LABEL org.label-schema.name="Drone Crowdin"
LABEL org.label-schema.vendor="Jonas Franz"
LABEL org.label-schema.schema-version="1.0"

ADD release/linux/amd64/drone-crowdin /bin/
ENTRYPOINT ["/bin/drone-crowdin"]