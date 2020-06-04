package main

import (
	"encoding/json"
	"testing"
)

func TestConvertAgentcloudsRequestsResponseToBuildStatisticsCanConvert(t *testing.T) {
	jsonString := `{"count":95873,"value":[{"agentCloudId":1,"requestId":"df8bb7f6-cd13-473e-8501-0002e6339140","pool":{"id":169,"isHosted":false,"poolType":"automation","size":0,"isLegacy":null,"options":null},"agent":{"id":609,"name":null,"version":null,"status":0,"provisioningState":null},"agentSpecification":{"VMImage":"Ubuntu16"},"agentData":{"RequestId":53253677},"provisionRequestTime":"2019-12-02T08:56:48.3006434Z","provisionedTime":null,"agentConnectedTime":"2019-12-02T08:56:51.727Z","releaseRequestTime":"2019-12-02T09:01:06.9160742Z"},{"agentCloudId":1,"requestId":"cc57ac12-237d-4d24-a55e-00033369ad75","pool":{"id":169,"isHosted":false,"poolType":"automation","size":0,"isLegacy":null,"options":null},"agent":{"id":607,"name":null,"version":null,"status":0,"provisioningState":null},"agentSpecification":{"VMImage":"vs2017-win2016"},"agentData":{"RequestId":57205591},"provisionRequestTime":"2020-01-14T12:12:35.1830738Z","provisionedTime":null,"agentConnectedTime":"2020-01-14T12:12:37.93Z","releaseRequestTime":"2020-01-14T12:13:29.7384754Z"},{"agentCloudId":1,"requestId":"5cae837a-94bc-473f-b0c8-00044baab3b3","pool":{"id":169,"isHosted":false,"poolType":"automation","size":0,"isLegacy":null,"options":null},"agent":{"id":609,"name":null,"version":null,"status":0,"provisioningState":null},"agentSpecification":{"vmImage":"ubuntu-latest"},"agentData":{"RequestId":59276306},"provisionRequestTime":"2020-01-30T14:02:18.3487526Z","provisionedTime":null,"agentConnectedTime":"2020-01-30T14:02:23.327Z","releaseRequestTime":"2020-01-30T14:28:26.4553648Z"}]}`

	agentcloudsResponse := AgentcloudsRequestsResponse{}
	json.Unmarshal([]byte(jsonString), &agentcloudsResponse)

	buildStatistics := ConvertAgentcloudsRequestsResponseToBuildStatistics(agentcloudsResponse)

	expectedBuildStatistics := 3
	if len(buildStatistics) != expectedBuildStatistics {
		t.Error("build statistics count was not the expected: % ")
	}
}
