#!/bin/bash

curl -X POST http://localhost:17000 -d $'reset\ngreen\nbgrect 0.05 0.05 0.95 0.95\nupdate'
