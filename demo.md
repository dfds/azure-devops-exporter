# Demo notes

## Why (purpose)
* There have been complaints about long queue times.
  - We want to use data to decide if we should scale up the amount of build agents.
* We want to get experience with data lakes.

## How (Process)
1. We scrape Azure devops endpoints
1. store the scraped information in s3
1. Import data into AWS Athena
1. We transform some of the data to make it queryable
1. We analyse our data

## What (Outcome)

* Decision foundation based on data we can query with sql, r, matlab.
* Dashboards showing current throughput and potential bottlenecks.

## Feedback
