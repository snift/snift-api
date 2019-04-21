# Installation

The following steps are for installing go v1.11.4 in an Ubuntu 64-bit machine.

-   Make sure to keep your system libraries updated.

```shell
    sudo apt-get update
    sudo apt-get -y upgrade
```

-   Download the latest Go Package. As of this writing, the version mentioned is v1.11.4.

```shell
wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz
```

-   Extract and Move Go into the `usr/local` location.

```shell
sudo tar -xvf go1.11.4.linux-amd64.tar.gz
sudo mv go /usr/local
```

-   Now, add the Environment Variables in your `~/.profile` or `~/.bashrc` or `~/.zshrc` based on your shell configuration.

```shell
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

-   Update the session so that the changes made will get reflected or open a new terminal session in a different tab.

```shell
    source ~/.profile
```

-   Make sure your installation was successful by trying the `go version` command. Ideally, the output should look something like

```shell
    go1.11.4 linux/amd64
```

## Development Environment

### Visual Studio Code

-   Install the official [vscode-go](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go) extension from Visual Studio Marketplace. Once installed it prompts you to install go-tools such as gocode, go-ved, golint and many more for a smooth development experience.

### Sublime Text

-   Install the [GoSublime](https://github.com/DisposaBoy/GoSublime) extension which provides support for linting, auto-format and other such go-tools.

## Troubleshooting

View the [Troubleshooting Docs](Troubleshooting.md) for any errors during usage or file an issue [here](https://github.com/maruthi-adithya/snift-api/issues) if you think something is wrong.
