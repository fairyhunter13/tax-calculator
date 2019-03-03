#!/bin/bash
until $(curl --output /dev/null --silent --get --fail 'http://'"$SERVICE_HOST"':'"$SERVICE_PORT$SERVICE_ENDPOINT"'');
do
  echo 'Waiting for the TaxCalculator to start'
  sleep 1
done

nginx -g 'daemon off;'