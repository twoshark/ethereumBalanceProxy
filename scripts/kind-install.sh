#!/usr/bin/env bash

helm del balance-proxy-test -n balance-proxy-test
helm install balance-proxy-test helm/ethBalanceProxy/. -n balance-proxy-test --create-namespace

