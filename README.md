# Sonarqube statistics bitbucket application

## Table of Contents
1. [Bitbucket ap registration](#Bitbucket ap registration)
2. [Running](#Running)
3. [Building](#Building)



## Bitbucket ap registration
1. Go to your bitbucket organization/account settings
2. OAuth consumers - Add consumer
3. Fill in details, check `This is a private consumer` checkbox, choose scopes: `account: read`, `repositories: read`
4. Save oauth key and secret

## Running
1. Prepare `config.toml` file.
2. `docker pull 1node/bb-sonar-stats:latest`
3. `docker run -itd -v </path/to/config.toml>/config.toml:/app/config.toml -p <app_port>:9001 1node/bb-s
    onar-stats:latest`


## Building
1. git clone repo, cd repo
2. `docker build -t 1node/bb-sonar-stats:latest -f ./docker/Dockerfile .`
 
