FROM scratch

# Expose ports.
#  - 4000: api
EXPOSE 4000

ADD deploy/certs/ca-certificates.crt /etc/ssl/certs/
ADD bin/semver /

WORKDIR /data

ENTRYPOINT ["/semver"]

