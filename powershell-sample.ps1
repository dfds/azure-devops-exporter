# Generate basic auth token header
$personalToken = 'nthvku62ythn6rwogqn73c3cskl2a6fgmy6n3b56u3b6myrwub2q'
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
$agentclouds = 'https://dev.azure.com/dfds/_apis/distributedtask/agentclouds?api-version=6.0-preview.1'
$agentcloudsout = Invoke-RestMethod -Uri $agentclouds -Method Get -ContentType "application/json" -Headers $header
$agentcloudsout.value 

# Get runs from agent cloud
$agentcloudruns = 'https://dev.azure.com/dfds/_apis/distributedtask/agentclouds/1/requests?api-version=6.0-preview.1'
$agentcloudrunsout = Invoke-RestMethod -Uri $agentcloudruns -Method Get -ContentType "application/json" -Headers $header
$agentcloudrunsout.value.count

$agentcloudrunsout.value[95519]
# Get values from agentcloud runs as buildtime
$buildtimes = $agentcloudrunsout.value 

# Convert to timespan
$secs = foreach ($build in $buildtimes) {
    New-TimeSpan -Start $build.agentConnectedTime -End $build.releaseRequestTime
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



df8bb7f6-cd13-473e-8501-0002e6339140



$pools = 'https://dev.azure.com/dfds/_apis/build/builds/53253677?api-version=5.1'
$poolsout = Invoke-RestMethod -Uri $pools -Method Get -ContentType "application/json" -Headers $header
$poolsout.value | Select-Object name, isHosted, agentCloudId




GET https://dev.azure.com/{organization}/

$pools = 'https://dev.azure.com/dfds/_apis/distributedtask/pools/169/agents/609?api-version=5.1'
$poolsout = Invoke-RestMethod -Uri $pools -Method Get -ContentType "application/json" -Headers $header
$poolsout.value | Select-Object name, isHosted, agentCloudId


$pools = 'https://dev.azure.com/dfds/DevelopmentExcellence/_apis/build/builds?api-version=5.1'
$poolsout = Invoke-RestMethod -Uri $pools -Method Get -ContentType "application/json" -Headers $header
$poolsout.value

217334
Setup ECR repos














_links            : @{self=; web=; sourceVersionDisplayUri=; timeline=; badge=}
properties        : 
tags              : {}
validationResults : {}
plans             : {@{planId=d2ed560a-4bda-4451-bf84-194d1df6f61a}}
triggerInfo       : 
id                : 217334
buildNumber       : 217334
status            : completed
result            : succeeded
queueTime         : 23/01/2020 09.50.11
startTime         : 23/01/2020 09.50.20
finishTime        : 23/01/2020 09.51.15
url               : https://dev.azure.com/dfds/ace5e409-c242-4356-93f4-23c53a3dc87b/_apis/build/Builds/217334
definition        : @{drafts=System.Object[]; id=792; name=Setup ECR repos; url=https://dev.azure.com/dfds/ace5e409-c242-4356-93f4-23c53a3dc87b/_apis/build/Definitions/792?revision=3; uri=vstfs:///Build/Definition/792; path=\; type=build; queueStatus=enabled; revision=3; 
                    project=}
project           : @{id=ace5e409-c242-4356-93f4-23c53a3dc87b; name=DevelopmentExcellence; description=Look @ https://github.com/dfds ; url=https://dev.azure.com/dfds/_apis/projects/ace5e409-c242-4356-93f4-23c53a3dc87b; state=wellFormed; revision=1967; visibility=private; 
                    lastUpdateTime=02/10/2019 12.46.29}
uri               : vstfs:///Build/Build/217334
sourceBranch      : refs/heads/master
sourceVersion     : 2016a05ab2b16ec181bf73e4244fe38e0f529789
queue             : @{id=768; name=Hosted VS2017; pool=}
priority          : normal
reason            : individualCI
requestedFor      : @{displayName=Rune Abrahamsson; url=https://spsprodweu1.vssps.visualstudio.com/Ad61c64a7-5a2c-4ee0-b50e-3972af43fa07/_apis/Identities/9f3c4987-8b2b-67d7-811d-72582869d73c; _links=; id=9f3c4987-8b2b-67d7-811d-72582869d73c; uniqueName=ruabr@DFDS.COM; 
                    imageUrl=https://dev.azure.com/dfds/_apis/GraphProfile/MemberAvatars/aad.OWYzYzQ5ODctOGIyYi03N2Q3LTgxMWQtNzI1ODI4NjlkNzNj; descriptor=aad.OWYzYzQ5ODctOGIyYi03N2Q3LTgxMWQtNzI1ODI4NjlkNzNj}
requestedBy       : @{displayName=Microsoft.VisualStudio.Services.TFS; url=https://spsprodweu1.vssps.visualstudio.com/Ad61c64a7-5a2c-4ee0-b50e-3972af43fa07/_apis/Identities/00000002-0000-8888-8000-000000000000; _links=; id=00000002-0000-8888-8000-000000000000; 
                    uniqueName=00000002-0000-8888-8000-000000000000@2c895908-04e0-4952-89fd-54b0046d6288; 
                    imageUrl=https://dev.azure.com/dfds/_apis/GraphProfile/MemberAvatars/s2s.MDAwMDAwMDItMDAwMC04ODg4LTgwMDAtMDAwMDAwMDAwMDAwQDJjODk1OTA4LTA0ZTAtNDk1Mi04OWZkLTU0YjAwNDZkNjI4OA; 
                    descriptor=s2s.MDAwMDAwMDItMDAwMC04ODg4LTgwMDAtMDAwMDAwMDAwMDAwQDJjODk1OTA4LTA0ZTAtNDk1Mi04OWZkLTU0YjAwNDZkNjI4OA}
lastChangedDate   : 11/02/2020 19.38.34
lastChangedBy     : @{displayName=Microsoft.VisualStudio.Services.ReleaseManagement; url=https://spsprodweu1.vssps.visualstudio.com/Ad61c64a7-5a2c-4ee0-b50e-3972af43fa07/_apis/Identities/0000000d-0000-8888-8000-000000000000; _links=; 
                    id=0000000d-0000-8888-8000-000000000000; uniqueName=0000000d-0000-8888-8000-000000000000@2c895908-04e0-4952-89fd-54b0046d6288; 
                    imageUrl=https://dev.azure.com/dfds/_apis/GraphProfile/MemberAvatars/s2s.MDAwMDAwMGQtMDAwMC04ODg4LTgwMDAtMDAwMDAwMDAwMDAwQDJjODk1OTA4LTA0ZTAtNDk1Mi04OWZkLTU0YjAwNDZkNjI4OA; 
                    descriptor=s2s.MDAwMDAwMGQtMDAwMC04ODg4LTgwMDAtMDAwMDAwMDAwMDAwQDJjODk1OTA4LTA0ZTAtNDk1Mi04OWZkLTU0YjAwNDZkNjI4OA}
orchestrationPlan : @{planId=d2ed560a-4bda-4451-bf84-194d1df6f61a}
logs              : @{id=0; type=Container; url=https://dev.azure.com/dfds/ace5e409-c242-4356-93f4-23c53a3dc87b/_apis/build/builds/217334/logs}
repository        : @{id=22ba5736-9b0a-435c-8f06-969f0f23cf5c; type=TfsGit; name=ECR repos; url=https://dev.azure.com/dfds/DevelopmentExcellence/_git/ECR%20repos; clean=; checkoutSubmodules=False}
keepForever       : True
retainedByRelease : False
triggeredByBuild  : 