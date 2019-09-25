# make sure to run this from the root of the project directory
#
# ./scripts/start.sh NUMBER_OF_NODES ENV

BLUE='\033[1;34m'
NO_COLOR='\033[0m'

INSTANCE_COUNT=$1
ENV=$2

log () {
	echo -e "${BLUE}Launcher ==> $1${NO_COLOR}"
}

log "Flushing logs directory"
rm -rf logs/*.txt

log "Creating hosts.txt"
rm -rf hosts.txt
touch hosts.txt

if [ $ENV == "DEV" ]; then
    IP=127.0.0.1
else
    IP=$(curl ifconfig.me)
fi

PROM_ENDPOINTS=()
for (( i=0; i<=$(($INSTANCE_COUNT-1)); i++ ))
do
    echo "$i,localhost,$IP" >> hosts.txt
    PROM_ENDPOINTS+=("host.docker.internal:$((2112 + $i))")
done

if [ -d "heimdall" ]; then
    log "Generating sd.json for heimdall"

    S="["
    for (( i=0; i<$INSTANCE_COUNT; i++ )); do
        if [ $i == $(($INSTANCE_COUNT - 1)) ]; then
            S+="\"${PROM_ENDPOINTS[$i]}\""
        else
            S+="\"${PROM_ENDPOINTS[$i]}\","
        fi
    done
    S+="]"

    rm heimdall/prometheus/sd.json && touch heimdall/prometheus/sd.json
    echo '[{ "targets": '$S', "labels": { "env": "local", "job": "self-stabilizing-urb" } }]' >> heimdall/prometheus/sd.json
fi


log "Starting $INSTANCE_COUNT node(s) locally"
for (( i=0; i<=$(($INSTANCE_COUNT-1)); i++ ))
do
    log "Starting node $i"
    ID=$i IP=$IP ENV=$ENV go run main.go &
done

while true; do sleep 2; done
