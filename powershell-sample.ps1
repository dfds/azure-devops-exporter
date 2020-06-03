# Generate basic auth token header
$personalToken = ''
$token = [System.Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes(":$($personalToken)"))
$header = @{authorization = "Basic $token"}

# Get list of pools
$pools = 'https://dev.azure.com/dfds/_apis/distributedtask/pools?api-version=6.0-preview.1'
$poolsout = Invoke-RestMethod -Uri $pools -Method Get -ContentType "application/json" -Headers $header
$poolsout.value | Select-Object name, isHosted, agentCloudId


# Get jobs in pool id
$requests = 'https://dev.azure.com/dfds/_apis/distributedtask/pools/169/jobrequests'
$requestsout = Invoke-RestMethod -Uri $requests -Method Get -ContentType "application/json" -Headers $header
$requestsout.value 

# Get agent clouds available
$agentclouds = 'https://dev.azure.com/dfds/_apis/distributedtask/agentclouds?api-version=5.1-preview.1'
$agentcloudsout = Invoke-RestMethod -Uri $agentclouds -Method Get -ContentType "application/json" -Headers $header
$agentcloudsout.value 

# Get runs from agent cloud
$agentcloudruns = 'https://dev.azure.com/dfds/_apis/distributedtask/agentclouds/1/requests?api-version=5.1-preview.1'
$agentcloudrunsout = Invoke-RestMethod -Uri $agentcloudruns -Method Get -ContentType "application/json" -Headers $header
$agentcloudrunsout.value.count

$agentcloudrunsout.value[95519]
# Get values from agentcloud runs as buildtime
$buildtimes = $agentcloudrunsout.value 

# Convert to timespan
$secs = foreach ($build in $buildtimes) {
    New-TimeSpan -Start $build.provisionRequestTime -End $build.agentConnectedTime
}

# Get min, max and avarage times.
($secs | Measure-Object -Maximum ).Maximum
($secs | Measure-Object -Minimum ).Minimum
#slight rounding error here
[TimeSpan][int]($secs | Measure-Object -Average -Property Ticks ).Average

# name                             id
# ----                             --
# Hosted                            2
# Hosted VS2017                    29
# Hosted macOS                     70
# Hosted Ubuntu 1604              109
# Hosted Windows Container        114
# Hosted Windows 2019 with VS2019 141
# Hosted macOS High Sierra        150
# Azure Pipelines                 169