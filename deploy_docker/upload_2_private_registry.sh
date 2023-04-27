#!/bin/bash
 
Registry=$1
Username=$2
Password=$3
ImagesFromDockerHub=(${@:4}) 

echo "Trying to login $Registry with user $Username ..."
echo "$Password" | docker login "$Registry" --username "$Username" --password-stdin

for imageTag in "${ImagesFromDockerHub[@]}"
do
  docker pull $imageTag --platform linux/amd64

  PrivateTag="$Registry/$imageTag"

  docker tag $imageTag $PrivateTag
  echo "${GREEN}Image tagged as: $Registry/$imageTag ${NC}"
  
  FAILS=0
  for n in {1..3}; do
    docker push $PrivateTag && break;
    FAILS=$((FAILS + 1))
    if [[ $FAILS -ge 3 ]]; then
      echo "ERROR: Failed to push $FAILS times"
      exit 1
    fi
  done
  echo "${GREEN}Image $PrivateTag pushed to $Registry ${NC}"

done
