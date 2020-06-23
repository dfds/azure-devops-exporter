

## 2020-06-23

The Velocity project has gates in their builds resulting in builds takeing up to 14 days. Resulting in skewed picture of running time. We thing that we would need ot look at job running time to get precise running time.

Picking random builds in the FastFwd project shows that our calculated que time is longer than the ui for the agent pool shows queue.
https://dev.azure.com/dfds/_settings/agentpools?poolId=169&view=jobs

We see a deviance in the ui but we feel the numbers are valid and expect that the numbers we see in the ui comes from a different datasource and have been rounded slightly

Is there a init phase that would account for the missing wait time?

I the agent doesn't count the init time, then our number from init to build it more on para.

Current goal:
- Queries within given timelines
    Can we create a view or stored procedure? 
- Do we still put empty files into s3?
- How to demo?
    This is what we are trying to solve
    This is what we can do
    Where would you like to see this go


## 2020-06-11

We need to test if we get all the relevant results when we query:
https://dev.azure.com/dfds/7d9e1da7-15ed-4d53-b567-81d11e7ccec6/_apis/build/builds?api-version=5.1&$top=5000&statusFilter=completed&minTime=2020-06-11T13:00:43Z&maxTime=2020-06-11T13:02:45Z
  

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