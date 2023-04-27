#!/bin/bash

# ./_runMode/_imgBuildPush.sh ./project1 localhost:34000 user1 user1 project-alpha

ProjectDirectory=$1
Registry=$2
Username=$3
Password=$4
ImageName=$5
ImageTag=$6

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
elif [[ -z "$Username" ]]; then
  echo "${RED}Username is empty"
  exit 2
elif [[ -z "$Password" ]]; then
  echo "${RED}Password is empty"
  exit 2
elif [[ -z "$ImageName" ]]; then
  echo "${RED}ImageName is empty"
  exit 2
fi

if [[ -z "$ImageTag" ]]; then
  ImageTag="latest"
  echo "${YELLOW}No image tag existed, using 'latest' instead ${NC}"
fi

echo "Trying to login $Registry with user $Username ..."
echo "$Password" | docker login "$Registry" --username "$Username" --password-stdin

echo "buinding image: $ImageName ..."

ImageId="$ImageName:$ImageTag"
echo "${GREEN}Image built succeed with image Id: $ImageId ${NC}"

cd $ProjectDirectory
docker build -t $ImageId -f ./_run/_docker/Dockerfile .
echo "${GREEN}Docker image $ImageId built ${NC}"

docker tag $ImageId "$Registry/$ImageId"
echo "${GREEN}Image tagged as: $Registry/$ImageId ${NC}"

docker push "$Registry/$ImageId"
echo "${GREEN}Image $ImageId pushed to $Registry ${NC}"

echo "${BLUE}Image List in Registry:"
curl "http://$Username:$Password@$Registry/v2/_catalog"

echo "tags under image are:"
curl "http://$Username:$Password@$Registry/v2/$ImageName/tags/list"
echo "${GREEN}Image Pushed as: ${Registry}/${ImageId} ${NC}"
