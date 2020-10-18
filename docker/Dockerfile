FROM golang:1.14.3-alpine as build
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build

FROM alpine:latest
LABEL maintainer="Bogdan Kosarevskyi <bogdan.kosarevskyi@gmail.com>"
LABEL vendor="1node"
LABEL lastUpdate="18-10-2020"
LABEL description="Sonarqube stats bitbucket application."
COPY --from=build /app/bb-sonarqube-integration /app/bb-sonarqube-integration
COPY --from=build /app/static /app/static
COPY --from=build /app/templates /app/templates
WORKDIR /app
CMD ["/app/bb-sonarqube-integration"]
