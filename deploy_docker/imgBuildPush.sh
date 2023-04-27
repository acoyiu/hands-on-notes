#!/bin/bash

# ./imgBuildPush.sh ./directory registry-url username password img-name [tag-name]

ProjectDirectory=$1
Registry=$2
Username=$3
Password=$4
ImageName=$5
ImageTag=$6

if [[ -z "$ProjectDirectory" ]]; then
  echo -e "Error: Project Directory is empty"
  exit 2
elif [[ -z "$Registry" ]]; then
  echo "Error: Registry is empty"
  exit 2
elif [[ -z "$Username" ]]; then
  echo "Error: Username is empty"
  exit 2
elif [[ -z "$Password" ]]; then
  echo "Error: Password is empty"
  exit 2
elif [[ -z "$ImageName" ]]; then
  echo "Error: ImageName is empty"
  exit 2
fi

if [[ -z "$ImageTag" ]]; then
  ImageTag="latest"
  echo "No image tag existed, using 'latest' instead"
fi

echo "Trying to login $Registry with user $Username ..."
echo "$Password" | docker login "$Registry" --username "$Username" --password-stdin

echo "building image: $ImageName ..."

ImageId="$ImageName:$ImageTag"
echo "Image built succeed with image Id: $ImageId"

cd $ProjectDirectory
docker build -t $ImageId .
echo "Docker image $ImageId built"

docker tag $ImageId "$Registry/$ImageId"
echo "Image tagged as: $Registry/$ImageId"

docker push "$Registry/$ImageId"
echo "Image $ImageId pushed to $Registry"

echo "Pls call curl to check is image pushed successfully: curl http://$Username:$Password@$Registry/v2/_catalog"
echo "tags under image are: curl http://$Username:$Password@$Registry/v2/$ImageName/tags/list"
