#!/bin/sh

set -e

go build
exec ./mc_stats
