# registry

https://docs.docker.com/registry/deploying/
https://docs.docker.com/registry/deploying/#storage-customization


Add Tag
docker tag nothingsbland-www:1.0.0 black-pearl.sea:5000/nothingsbland-www

Push Image to Repo
docker push black-pearl.sea:5000/nothingsbland-www 

Review Repo Images
docker image ls -a black-pearl.sea:5000/nothingsbland-www

NOTE: Use docker->preferences->daemon to update the insecure registries