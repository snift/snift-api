## Steps to Install Latest version of Go in Ubuntu 64 bit

1) First step as always is a common one to update the existing packages.
```
sudo apt-get update
sudo apt-get -y upgrade
```
2) Download the latest Go Package.(Here it is 1.11.4)
```
wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz
```
3) Extract and Move Go into the `usr/local` location.
```
sudo tar -xvf go1.11.4.linux-amd64.tar.gz
sudo mv go /usr/local
```
4) Now, add the Environment Variables. If want them to be permanent, place them in `~/.profile` 

In case, you are not a big fan of vim, you can use Sublime Text Editor using the following commands.
```
subl ~/.profile
```
Add the following lines at the end of the file.

```
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```
5) Update the session so that the changes made will get reflected.
```
source ~/.profile
```
6)When `go version` command is typed, the output must be the following.

```
go version go1.11.4 linux/amd64
```

