cd cmd/go-scribe
build.bat
set /p version= < version.txt
docker build -f Dockerfile -t romanos/scribe:%version% .
cd ../..
docker-compose -f docker-compose.yml -p LogScribe --verbose build