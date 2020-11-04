# bb

![bb logo](.github/bb-logo.png)

---

`bb` is an inoffical bitbucket.org command line tool deeply inspired by the 
official [GitHub CLI](https://github.com/cli/cli/). It brings pull requests, 
downloads, and other bitbucket concepts to you the terminal.

![screenshot showing ](.github/screenshot_create_pr.png)

# Installation

Currently, the only option to install bb is by compiling it from source.

Make sure you have a working Go environment. Follow the 
[Go install instructions](https://golang.org/doc/install).

## Using `go get`

```bash
go get github.com/craftamap/bb
```

## Building from source

```bash
git clone https://github.com/craftamap/bb.git
go build
```

# Set-Up

You need to authenticate with your credentials first. You should generate a
[app password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/)
for that. Make sure to grant read and write access to the features you want to use.
(**Recommended**:Repositories: Read/Write, Pull Requests: Read/Write, Pipelines: Read/Write)

Run the following command to enter your username and password:

```bash
bb auth login
```

Your credentials will be stored to `~/.config/bb/configuration.toml`.

# Usage

To see all available commands, use `bb` without any subcommand.

## Pull Requests

You can use `bb pr` to list, view and merge existing pull requests and see how
their pipelines run. You can also use `bb pr create` to create new ones.

## Downloads

Manage downloads by listing, downloading or uploading them.

# Development
## Used Libraries

We use two multiple different bitbucket libaries:

 - https://github.com/ktrysmt/go-bitbucket
 - https://github.com/jsdidierlaurent/go-bitbucket 
   
Thanks to both of these!
