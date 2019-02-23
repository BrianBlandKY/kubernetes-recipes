# registry


https://docs.docker.com/registry/deploying/
https://docs.docker.com/registry/deploying/#storage-customization


Add Tag
docker tag nothingsbland-app:1.0.0 flying-dutchman.sea:5000/nothingsbland-app

Push Image to Repo
docker push flying-dutchman.sea:5000/nothingsbland-app 

Review Repo Images
docker image ls -a flying-dutchman.sea:5000/nothingsbland-app