
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
* [ ] Add QR code
* [ ] URL shortener
* [ ] Avoid duplicate names
* [ ] Upload max size