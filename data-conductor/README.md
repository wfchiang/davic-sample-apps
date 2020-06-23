## Clean Build 
#### (1) Ensure Davic Version in glide.yaml 

#### (2) Cleanup and Install Dependencies 
``` 
rm -rf vendor
```

```
glide install 
```

#### (3) Build Go Module 
```
go mod init github.com/wfchiang/davic-sample-apps
``` 