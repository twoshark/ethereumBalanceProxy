#!/usr/bin/env bash

clustername="balance-proxy-cluster"
kind destroy cluster --name $clustername
