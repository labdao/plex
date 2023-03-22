# Local example with testdata


```
# from the ligase directory
docker build -t ligase .
docker run -v $(pwd)/testdata/:/inputs -v $(pwd)/outputs:/outputs ligase
```
