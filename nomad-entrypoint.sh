#!/bin/sh

nomad agent -server -data-dir=/opt/nomad -bind=0.0.0.0 -bootstrap-expect=${BOOTSTRAP_EXPECT} -consul-address=${CONSUL_ADDR}