# An image processing fullstack app built with Go + Angular.

## Development:
```
cd frontend; pnpm i; pnpm start
cd backend; air
```
 
## Deployment:
First build and push the image locally:
```
docker run -d -p 80:8080 bohdancho/imaginaer
docker push bohdancho/imaginaer
```
Then clone the repo, pull and run them it the server:
```
docker pull bohdancho/imaginaer
docker run -dp 80:8080 bohdancho/imaginaer
```
