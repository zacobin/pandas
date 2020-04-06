#!/bin/bash
CGO_ENABLED=0 swagger generate spec -o ../models.json -w . -m
