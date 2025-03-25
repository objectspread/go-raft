#!/bin/bash

conf_path=$1


config=$1

    case "$conf_path" in
        1)
            config=$(dirname $0)/configs/srv-1.json
            ;;
        2)
            config=$(dirname $0)/configs/srv-2.json
            ;;
        3)
            config=$(dirname $0)/configs/srv-3.json
            ;;
        3)
            config=$(dirname $0)/configs/srv-4.json
            ;;
        *)
            config=$(dirname $0)/configs/srv-0.json
            ;;
    esac

echo "Starting Raft Server with config file ${config}"
go run ./cmd/server/main.go --config-file=${config}