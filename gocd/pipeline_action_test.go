package gocd

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

func testPipelineServiceUnPause(t *testing.T) {
	for n, test := range []struct {
		name          string
		v             *ServerVersion
		confirmHeader string
		acceptHeader  string
	}{
		{
			name:          "server-version-14.3.0",
			v:             &ServerVersion{Version: "14.3.0"},
			confirmHeader: "Confirm",
			acceptHeader:  apiV0,
		},
		{
			name:          "server-version-18.3.0",
			v:             &ServerVersion{Version: "18.3.0"},
			confirmHeader: "X-GoCD-Confirm",
			acceptHeader:  apiV1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if runIntegrationTest(t) {

				pipelineName := fmt.Sprintf("test-pipeline-un-pause%d", n)

				err := test.v.parseVersion()
				assert.NoError(t, err)

				cachedServerVersion = test.v

				ctx := context.Background()

				pausePipeline, _, err := intClient.PipelineConfigs.Create(ctx, mockTestingGroup, &Pipeline{
					Name: pipelineName,
				})

				pp, _, err := intClient.Pipelines.Pause(context.Background(), pausePipeline.Name)
				assert.NoError(t, err)
				assert.True(t, pp)

				deleteResponse, _, err := intClient.PipelineConfigs.Delete(ctx, pausePipeline.Name)
				assert.Equal(t, "", deleteResponse)
			}
		})
	}
}
