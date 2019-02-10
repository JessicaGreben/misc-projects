#!/bin/sh

attempt_counter=0
max_attempts=5

until $(curl --output /dev/null --silent --head --fail http://localhost:8080); do
    if [ ${attempt_counter} -eq ${max_attempts} ];then
      echo "Max attempts reached"
      exec "$@"
    fi
    attempt_counter=$(($attempt_counter+1))
    printf '.'
    sleep 1
done

exec "$@"
