#!/bin/bash
while ! pg_isready -h ${POSTGRE_HOST} -p ${POSTGRE_PORT} > /dev/null 2> /dev/null; do
    echo "Connecting to postgre failed!"
    sleep 1
done

./taxcalculator
