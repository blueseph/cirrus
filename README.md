# Cirrus
Cirrus is a simple CLI tool for deploying and deleting CloudFormation stacks safely.

_Cirrus is under development and there may be display errors. Please report any [issues](github.com/blueseph/cirrus/issues)_

[![asciicast](https://asciinema.org/a/J98ZyQx1lhpV4Ok1UBksNdB1w.svg)](https://asciinema.org/a/J98ZyQx1lhpV4Ok1UBksNdB1w)

## Installation

Cirrus is available for Windows, Mac, and Linux and the i386 and amd64 architectures. You'll find the binaries on the [release page](github.com/blueseph/cirrus/releases)

It's a [Go](https://golang.org/) binary and can be installed with:

`GO111MODULE=on go get github.com/blueseph/cirrus`

## Quick Overview

`cirrus up --stack MySecureVPC`  
`cirrus down --stack MySecureVPC`


Cirrus will follow CloudFormation best practices such as creating a change set before creates/updates, deleting empty (0 resource) stacks, and linting your templates.

A best effort has been made to apply sensible deployment defaults, such as assuming a template.yaml or template.json file in the directory as the intended template, and a parameters.json file as the intended parameters file.

Some CloudFormation options have been disabled as a way of promoting best practices. Ad-hoc parameters are tags are not supported. The only supported options are having a parameters.json and/or a tags.json file. These are config files that can be sourced and vetted -- ad-hoc parameters/tags cannot.

## Commands

```
cirrus up 
    --stack stack-name              - Name of stack to be created/updated
    --template template.yaml        - Template to be uploaded. Default template.yaml
    --tags tags.json                - Tags to be uploaded. Default tags.json
    --parameters parameters.json    - Parameters to be uploaded. Default parameters.json
    --skip-lint                     - Skips linting with cfn-lint. Default false
```

```
cirrus down
    --stack stack-name              - Name of stack to be deleted
````

## Contributing

We'd love your help! See [CONTRIBUTING](github.com/blueseph/cirrus/blob/master/CONTRIBUTING.md) on how to help

## License

Cirrus is released under the MIT License.