#!/usr/bin/env bash

clustername="balance-proxy-cluster"
kind delete cluster -n $clustername
