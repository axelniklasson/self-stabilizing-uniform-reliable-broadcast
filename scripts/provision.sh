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
    echo "export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> $HOME/.profile
    source $HOME/.profile
    rm -r go1.13.linux-amd64.tar.gz
    mkdir $HOME/go
    mkdir $HOME/go/src
}

install_dep() {
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
}

if ! [ -x "$(command -v go)" ]; then
	install_go
else
	log "go already installed, skipping"
fi

if ! [ -x "$(command -v dep)" ]; then
	install_dep
else
	log "dep already installed, skipping"
fi