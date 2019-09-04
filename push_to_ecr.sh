#!/usr/bin/env bash
docker rmi -f genosha &&\
docker build --network host -t genosha . &&\
docker tag genosha xxx/genosha &&\
COMMAND=`eval aws ecr get-login --no-include-email` &&\
echo `eval $COMMAND` &&\
docker push xxx/genosha &&\
MANIFEST=$(aws ecr batch-get-image --repository-name genosha --image-ids imageTag=latest --query 'images[].imageManifest' --output text) &&\
aws ecr put-image --repository-name genosha --image-tag $1 --image-manifest "$MANIFEST"