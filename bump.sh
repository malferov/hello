#!/bin/bash
echo build = \"$(git rev-parse --short HEAD)\" > build.auto.tfvars
