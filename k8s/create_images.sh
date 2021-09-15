#!/bin/bash
docker build -t kuber/userProfile ..
docker save  kuber/userProfile > userProfile.tar
microk8s ctr image import userProfile.tar
rm userProfile.tar