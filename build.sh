#!/bin/bash

echo "Change version number? [Version number/N]";
read x

if [[ $x =~ ^[1-9].[0-9].[0-9]$ ]]
then
  sed -i -e "s/const appVersion = \"[0-9].[0-9].[0-9]\"/const appVersion = \"$x\"/g" main.go
  sed -i -e "s/GOPLOY_VER=v1.3.5/GOPLOY_VER=v$x/g" docker/Dockerfile
fi

echo "Build web? [Y/N]";
read x

if [ "$x" == Y ] || [ "$x" == y ]
then
    cd web
    npm run build
    cd ..
fi


echo "Building goploy-agent";

env GOOS=linux go build -o goploy-agent main.go

env GOOS=darwin go build -o goploy-agent.mac main.go

env GOOS=windows go build -o goploy-agent.exe main.go


