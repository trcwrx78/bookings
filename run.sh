#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings
./bookings -dbname=bookings -dbuser=postgres -cache=false -production=false