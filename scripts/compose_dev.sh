#!/usr/bin/env bash

ROOTDIR=$(dirname $(dirname $(realpath $0)))

MODE=dev
PROJECT_NAME=imageer_$MODE
COMPOSE_FILE=$ROOTDIR/deployments/docker-compose.$MODE.yaml
ENV_FILE=$ROOTDIR/deployments/.$MODE.env

COMPOSE_CMD="docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE -p $PROJECT_NAME"

PS3='Please enter your choice: '
options=("up" "down" "ps" "logs" "Quit")
select opt in "${options[@]}"
do
  case $opt in
    "up")
      CMD="$COMPOSE_CMD up -d"
      echo $CMD && eval $CMD
      break
      ;;
    "down")
      CMD="$COMPOSE_CMD down"
      echo $CMD && eval $CMD
      break
      ;;
    "ps")
      CMD="$COMPOSE_CMD ps"
      echo $CMD && eval $CMD
      break
      ;;
    "logs")
      CMD="$COMPOSE_CMD logs -f"
      echo $CMD && eval $CMD
      break
      ;;
    "Quit")
      exit 0
      ;;
    *) echo "invalid option $REPLY";;
  esac
done
