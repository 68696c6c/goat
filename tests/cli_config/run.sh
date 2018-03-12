#!/bin/bash
set -xe

go build -i -o test_cli_config
./test_cli_config test --config ../config.yml
