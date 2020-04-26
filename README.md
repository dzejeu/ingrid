# ingrid

In order to run already built docker image fetch it from dockerhub<br/>
```docker run -p 8080:8080 dzejeu/ingrid:latest```<br/>
Then service can be requestes in following way<br/>
```curl 'http://localhost:8080/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219'```<br/>
<br/>
In order to create fresh docker image run:<br/>
```docker build -t myapp:latest .```<br/>
and then:<br/>
```docker run -p 8080:8080 myapp:latest```
