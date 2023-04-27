#!/bin/bash

# ./imageBuildPusher_noLogin.sh ./directory <registry>:<port> <image-name> <tag-name>

ProjectDirectory=$1
Registry=$2
ImageName=$3
ImageTag=$4

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

if [[ -z "$ProjectDirectory" ]]; then
  echo "${RED}Project Directory is empty"
  exit 2
elif [[ -z "$Registry" ]]; then
  echo "${RED}Registry is empty"
  exit 2
elif [[ -z "$ImageName" ]]; then
  echo "${RED}ImageName is empty"
  exit 2
fi

if [[ -z "$ImageTag" ]]; then
  ImageTag="latest"
  echo "${YELLOW}No image tag existed, using 'latest' instead ${NC}"
fi

echo "buinding image: $ImageName ..."

ImageId="$ImageName:$ImageTag"
echo "${GREEN}Image building with image Id: $ImageId ${NC}"

cd $ProjectDirectory
docker build -t $ImageId .
echo "${GREEN}Docker image $ImageId built ${NC}"

docker tag $ImageId "$Registry/$ImageId"
echo "${GREEN}Image tagged as: $Registry/$ImageId ${NC}"

FAILS=0
for n in {1..3}; do
  docker push "$Registry/$ImageId" && break;
  FAILS=$((FAILS + 1))
  if [[ $FAILS -ge 3 ]]; then
    echo "ERROR: Failed to push $FAILS times"
    exit 1
  fi
done
echo "${GREEN}Image $ImageId pushed to $Registry ${NC}"

echo "${BLUE}Image List in Registry:"
echo "http://$Username:$Password@$Registry/v2/_catalog"

echo "tags under image are:"
echo "http://$Username:$Password@$Registry/v2/$ImageName/tags/list"

echo ""
echo "${GREEN}Image Pushed as: ${Registry}/${ImageId} ${NC}"
