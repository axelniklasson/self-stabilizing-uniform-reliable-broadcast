# make sure to run this from the root of the project directory
#
# ./scripts/deploy.sh NUMBER_OF_NODES

BLUE='\033[1;34m'
NO_COLOR='\033[0m'

NODE_COUNT=$1

log () {
	echo -e "${BLUE}Deployment ==> $1${NO_COLOR}"
}

log "Provisioning nodes"
