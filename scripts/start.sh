# make sure to run this from the root of the project directory
#
# ./scripts/start.sh NUMBER_OF_NODES

BLUE='\033[1;34m'
NO_COLOR='\033[0m'

log () {
	echo -e "${BLUE}Launcher ==> $1${NO_COLOR}"
}

log "Flushing logs directory"
rm -rf logs/*.txt

log "Creating hosts.txt"
rm -rf hosts.txt
touch hosts.txt
for (( i=0; i<=$(($1-1)); i++ ))
do
    echo "$i,localhost,127.0.0.1" >> hosts.txt
done



log "Starting $1 node(s) locally"
for (( i=0; i<=$(($1-1)); i++ ))
do
    log "Starting node $i"
    ID=$i go run main.go &
done

while true; do sleep 2; done