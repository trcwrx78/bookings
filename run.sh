#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings
./bookings -dbname=bookings -dbuser=postgres -dbpass=BrynAd02 -cache=false -production=false