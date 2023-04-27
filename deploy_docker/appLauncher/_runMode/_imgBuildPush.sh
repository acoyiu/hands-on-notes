#!/bin/bash

# ./_runMode/_imgBuildPush.sh ./project1 localhost:34000 project1 d-May-02

ProjectDirectory=$1
Registry=$2
ImageName=$3
ImageTag=$4
noCache=$5

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

if [[ -z "$ProjectDirectory" ]]; then
  echo -e "${RED}Project Directory is empty"
  exit 2
elif [[ -z "$Registry" ]]; then
  echo -e "${RED}Registry is empty"
  exit 2
elif [[ -z "$ImageName" ]]; then
  echo -e "${RED}ImageName is empty"
  exit 2
fi

if [[ -z "$ImageTag" ]]; then
  ImageTag="latest"
  echo -e "${YELLOW}No image tag existed, using 'latest' instead ${NC}"
fi

ImageId="$ImageName:$ImageTag"
echo -e "buinding image: $ImageId ..."

useCache="--no-cache"
if [[ -n "$ImageName" ]]; then
  useCache=""
fi

cd $ProjectDirectory
time docker build -t $ImageId $useCache .
echo -e "${GREEN}Docker image with ID:($ImageId) built ${NC}"

time docker tag $ImageId "$Registry/$ImageId"
echo -e "${GREEN}Image tagged as: $Registry/$ImageId ${NC}"

time docker push "$Registry/$ImageId"
echo -e "${GREEN}Image ${$Registry}/${ImageId} pushed${NC}"

cd ../
