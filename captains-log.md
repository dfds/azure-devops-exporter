## 2020-06-04
There are 83 projects in DFDS: https://dev.azure.com/dfds/_apis/projects?api-version=5.1

We can get all builds for a project at https://dev.azure.com/dfds/DevelopmentExcellence/_apis/build/builds/



We could use this one and the next to get the time on job runs:
https://dev.azure.com/dfds/_apis/distributedtask/pools?api-version=6.0-preview.1&properties=isHosted

https://dev.azure.com/dfds/_apis/distributedtask/pools/169/jobrequests