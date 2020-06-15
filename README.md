<!-- TABLE OF CONTENTS -->
<!-- omit in toc -->
## Table of Contents

- [About The Project](#about-the-project)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Installation](#installation)
- [Usage](#usage)
  - [Requirements](#Requirements)
  - [Running](#Running)

<!-- ABOUT THE PROJECT -->
## About The Project

This project exports build information from Azure DevOps for an organization and places it as JSON files in an AWS s3 bucket.

### Built With

`docker build -t .` in the root of the project.

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Installation
 
1. Clone the repo
```sh
git clone git@github.com:dfds/azure-devops-exporter.git
```
2. Install dependencies 
```sh
go mod download
```

## Usage

### Requirements
The go code requires an environment variable set named `ADO_PERSONAL_ACCESS_TOKEN`.
The variable should be set to your personal access token from Azure DevOps.

To generate a new personal access token, follow [this guide from Microsoft][ms-pa-token].
The go code requires an environment variable set named `ADO_PERSONAL_ACCESS_TOKEN`.

The go code also requires a valid set of AWS credentials will a scope to write to a s3 bucket under the prefix: `azure-devops/`
### Running
```shell script
go run main.go
```
[ms-pa-token]: https://docs.microsoft.com/en-us/azure/devops/organizations/accounts/use-personal-access-tokens-to-authenticate?view=azure-devops&tabs=preview-page
dr