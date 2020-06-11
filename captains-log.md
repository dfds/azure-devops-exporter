## 2020-06-11
Does the json object we place in a single file have to be in an array or can we just dump comma separated object?
## 2020-06-10

Could we load the diff from last file to current file in s3. 
Does it cost that much to have separate files for each build?

Is there a correlation between a build to a service running in hellman. Can we attach meta data to the build files to tie this in? 

## 2020-06-09

We added a function to get a list of project ID's.

This will allow us to use the endpoint below to get the builds from each project:
https://dev.azure.com/dfds/projectid/_apis/build/builds?api-version=5.1&top=5000

## 2020-06-04
There are 83 projects in DFDS: https://dev.azure.com/dfds/_apis/projects?api-version=5.1

We can get all builds for a project at https://dev.azure.com/dfds/DevelopmentExcellence/_apis/build/builds/


We could use this one and the next to get the time on job runs:
https://dev.azure.com/dfds/_apis/distributedtask/pools?api-version=6.0-preview.1&properties=isHosted

https://dev.azure.com/dfds/_apis/distributedtask/pools/169/jobrequests

Willi suggests that we keep as much as the original payload as possible.