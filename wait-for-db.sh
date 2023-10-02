#!/bin/bash
set -e

host="$1"
shift
port="$1"
shift
cmd="$@"

until nc -z -v -w30 "$host" "$port"; do
	echo "Waiting for database connection..."
	# wait for 5 seconds before check again
	sleep 5
done

echo >&2 "Database is up - executing command"
exec $cmd
