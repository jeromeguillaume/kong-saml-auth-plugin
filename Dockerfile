# docker build -t kong-gateway-saml .
FROM kong/kong-gateway:3.0.0.0-alpine
USER root

RUN apk update && apk add git nodejs npm go musl-dev libffi-dev gcc g++ file make \
&& npm install kong-pdk -g 

# Example for GO:
WORKDIR /saml-go

# Download Go modules
COPY /saml-go/go.mod .
COPY /saml-go/go.sum .
RUN go mod download

WORKDIR /saml-go/plugins
COPY /saml-go/plugins/saml-auth.go .
#RUN go build -o /usr/local/bin/ saml-auth.go
RUN go build saml-auth.go
RUN mv /saml-go/plugins/saml-auth /usr/local/bin/

COPY kong.conf /etc/kong/.

# reset back the defaults
USER kong
ENTRYPOINT ["/docker-entrypoint.sh"]
STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health
CMD ["kong", "docker-start"]
