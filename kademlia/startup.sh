#!/bin/sh

attempt_counter=0
max_attempts=3

# Wait for the bootstrap node to be up before allowing any nodes to join.
until $(./kademlia -c); do
    if [ ${attempt_counter} -eq ${max_attempts} ];then
      echo "Max attempts reached"
      exec "$@"
    fi
    attempt_counter=$(($attempt_counter+1))
    printf '.'
    sleep 1
done

exec "$@"
