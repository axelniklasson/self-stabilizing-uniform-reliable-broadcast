# this script provisions a planetlab node for running ssurb
# might not be needed later, if binaries are distributed instead

BLUE='\033[1;34m'
NO_COLOR='\033[0m'

log () {
	echo -e "${BLUE}Provision ==> $1${NO_COLOR}"
}

install_go() {
    wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.13.linux-amd64.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.profile
    source $HOME/.profile
    rm -r go1.13.linux-amd64.tar.gz
    mkdir $HOME/go
    mkdir $HOME/go/src
}

if ! [ -x "$(command -v go)" ]; then
	install_go
else
	log "go already installed, skipping"
fi