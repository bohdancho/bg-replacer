# An image processing fullstack app built with Go + Angular.

## Development:
```
cd frontend; pnpm i; pnpm start
cd backend; air
```
 
## Deployment:
Backend:
```
cd backend
docker build -t bohdancho/imaginaer-backend .
docker push bohdancho/imaginaer-backend
```
Then clone the repo, pull and run them it the server:
```
docker pull bohdancho/imaginaer-backend
docker run -dp 80:8080 bohdancho/imaginaer-backend
```

Frontend is deployed via AWS Amplify and autobuilt on push on main, backend is deployed via a docker container running on AWS EC2. Frontend and backend are then served via Cloudfront. 
