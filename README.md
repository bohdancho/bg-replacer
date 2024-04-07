# An image processing fullstack app built with Go + Angular.

## Development:
```
cd frontend; pnpm i; pnpm start
cd backend; air
```
 
## Deployment:
Frontend is built locally and deployed to S3, backend is deployed via a docker container running on AWS EC2. Frontend and backend are then reverse proxied via Cloudfront. 

Backend:
```
cd backend
chmod +x ./deploy.sh
./deploy.sh
```
Then log into the server via ssh, clone (pull) the repo,
if applicable stop the old container & delete the old image:
```
docker stop imaginaer-backend
docker rm imaginaer-backend
docker rmi bohdancho/imaginaer-backend
```
pull the new image:
```
docker pull bohdancho/imaginaer-backend
```
and run it:
```
docker run -dp 80:8080 --volume ./data:/data --volume ./static:/static --name imaginaer-backend bohdancho/imaginaer-backend 
```

Frontend:
```
chmod +x frontend/deploy.sh
frontend/deploy.sh
```
