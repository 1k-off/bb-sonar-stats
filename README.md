# Sonarqube statistics bitbucket application

## Table of Contents
1. [Bitbucket-app-registration](#Bitbucket-app-registration)
2. [Running](#Running)
3. [Building](#Building)



## Bitbucket-app-registration
1. Go to your bitbucket organization/account settings
2. OAuth consumers - Add consumer
3. Fill in details, check `This is a private consumer` checkbox, choose scopes: `account: read`, `repositories: read`
4. Save oauth key and secret
5. Go to Installed apps, enable Development mode
6. Go to Develop apps - Register app, enter app url, register app
7. Click on installation url and grant access to the app

## Running
1. Prepare `config.toml` file.
2. `docker pull 1node/bb-sonar-stats:latest`
3. `docker run -itd -v </path/to/config.toml>/config.toml:/app/config.toml -p <app_port>:9001 1node/bb-sonar-stats:latest`
4. Place `sonar.json` to master branch in your repo

```json
{
    "server": "https://sonar.domain.tld",
    "project_key": "My-Awesome-Project",
    "token": "sonarToken"
}
```

## Building
1. git clone repo, cd repo
2. `docker build -t 1node/bb-sonar-stats:latest -f ./docker/Dockerfile .`
 
