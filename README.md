# top_phrases

# Deploy
Edit docker-compose.yml to set the image version 
```image: cheojeg/top-quotes-api:<version>```

Command: ```docker buildx build --platform linux/amd64,linux/arm64 -t registry.digitalocean.com/top-quotes-registry/top-quotes-api:<version> --push .```

