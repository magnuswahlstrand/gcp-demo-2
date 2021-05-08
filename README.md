
**build docker image**
```
gcloud builds submit --project gcp-upload-demo
```

**apply changes**
```
terraform apply
```

## TODO
* [ ] Fix correct redirect url after POST
* [ ] Clean up older uploads
* [x] Add QR code
* [ ] URL shortener
* [x] Avoid duplicate names
* [ ] Make `files` cache safe for concurrency 
* [ ] Upload max size