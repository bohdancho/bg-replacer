source .env

ng build --configuration production
aws s3 cp ./dist/browser s3://imaginaer-static --recursive
aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DISTRIBUTION_ID --paths "/*"
