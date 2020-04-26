# ingrid

In order to run already built docker image fetch it from dockerhub
sh docker run -p 8080:8080 dzejeu/ingrid:latest
Then service can be requestes in following way
sh curl 'http://localhost:8080/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219'

In order to create fresh docker image run:
sh docker build -t myapp:latest .
and then:
sh docker run -p 8080:8080 myapp:latest
