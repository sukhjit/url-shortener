# URL Shortener API
This service uses AWS Api gateway, Lambda and Dynamo DB for url shortening service. For local development, the api uses *inmemory* database; in AWS it uses DynamoDB. The api looks for *LOCAL* env var to start in development/local mode. If this variable is not found, the api start in Lambda mode.

## Develop locally in Docker
#### Copy env file
```bash
cp .env.dist .env
```

#### Start development
```bash
docker-compose up
```


## Usage
#### Create new entry
```bash
curl -X POST 'http://127.0.0.1:3000/add' --data-raw '{ "url": "http://www.example.com" }'
```

Response
```bash
{
    "slug": "5i0w42l4",
    "url": "http://www.example.com",
    "visits": 0
}
```

#### Perform redirection to site: http://www.example.com
```bash
curl -X GET 'http://127.0.0.1:3000/5i0w42l4'
```

#### Get entry info
```bash
curl -X GET 'http://127.0.0.1:3000/info/5i0w42l4'
```

Response
```bash
{
    "slug": "5i0w42l4",
    "url": "http://www.example.com",
    "visits": 1
}
```


## Deploy to AWS

#### Install npm
```bash
npm install
```

#### Export aws key and secret for serverless deployment.
```bash
export AWS_ACCESS_KEY_ID=<your-aws-key>
export AWS_SECRET_ACCESS_KEY=<your-aws-secret-key>
```

#### Deploy using *make* command which uses serverless for deployment to AWS
```bash
make deploy
```
