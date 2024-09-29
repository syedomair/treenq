package cdk

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	tqsdk "github.com/treenq/treenq/pkg/sdk"
	"github.com/treenq/treenq/src/domain"
)

//go:embed testdata/app.yaml
var appYaml string

func TestAppDefinition(t *testing.T) {
	k := NewKube()
	ctx := context.Background()
	res := k.DefineApp(ctx, "id-1234", tqsdk.Space{
		Key: "space",
		Service: tqsdk.Service{
			Name: "simple-app",
			RuntimeEnvs: map[string]string{
				"PORT": "1234",
			},
			HttpPort: 8000,
			Repicas:  1,
			SizeSlug: tqsdk.SizeSlugS,
		},
	}, domain.Image{
		Registry:   "localhost:5000",
		Repository: "treenq",
		Tag:        "1.0.0",
	})

	_ = res
	t.Log(res)
	// assert.Equal(t, appYaml, res)
}

var conf = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJkekNDQVIyZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQWpNU0V3SHdZRFZRUUREQmhyTTNNdGMyVnkKZG1WeUxXTmhRREUzTWpjMk1EazNNakV3SGhjTk1qUXdPVEk1TVRFek5USXhXaGNOTXpRd09USTNNVEV6TlRJeApXakFqTVNFd0h3WURWUVFEREJock0zTXRjMlZ5ZG1WeUxXTmhRREUzTWpjMk1EazNNakV3V1RBVEJnY3Foa2pPClBRSUJCZ2dxaGtqT1BRTUJCd05DQUFSMVA0UlpJSUJRcExEdmtkZWt1enI5SWdHQUt0T1FiMGxZcW9YWVNKc04KSUV2cm03Y0VzSDhabjhmVWx0a0t2SDUzelRyR09DUWNWZmt1aGUyelhuWHZvMEl3UURBT0JnTlZIUThCQWY4RQpCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVU5zK3ZRWEhoQThFSmErMnhZcGVQCjFac0VCK0V3Q2dZSUtvWkl6ajBFQXdJRFNBQXdSUUlnU3Y0c1BCK0pRODU4clVvVjY3NHRKQ3prR1NXME01MkMKeWlaWkdPa2pIYkVDSVFEOHl1Y0hGaytnSmd4bmd5R1BDMXhTYU10UGlRVzBOT0hzNTZnODZtdEoydz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://127.0.0.1:6443
  name: default
contexts:
- context:
    cluster: default
    user: default
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: default
  user:
    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJrVENDQVRlZ0F3SUJBZ0lJZXVhekFXd01Vdm93Q2dZSUtvWkl6ajBFQXdJd0l6RWhNQjhHQTFVRUF3d1kKYXpOekxXTnNhV1Z1ZEMxallVQXhOekkzTmpBNU56SXhNQjRYRFRJME1Ea3lPVEV4TXpVeU1Wb1hEVEkxTURreQpPVEV4TXpVeU1Wb3dNREVYTUJVR0ExVUVDaE1PYzNsemRHVnRPbTFoYzNSbGNuTXhGVEFUQmdOVkJBTVRESE41CmMzUmxiVHBoWkcxcGJqQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJEaGducEI3Z2o4Q3hObHgKRGc2Q0N3a2k1aVFWblk3YklVYnpjZ0hQTlk5bTNsQ0tOOWJIZnU0U1UvUWlsRHNnTE9sZjNBOVhVYVZrVFRpVgpaSkpId2pDalNEQkdNQTRHQTFVZER3RUIvd1FFQXdJRm9EQVRCZ05WSFNVRUREQUtCZ2dyQmdFRkJRY0RBakFmCkJnTlZIU01FR0RBV2dCUnRkWUZvU29SZGZ1c1dxRTJLY1cxcEwzY3BOekFLQmdncWhrak9QUVFEQWdOSUFEQkYKQWlFQS81RGE3b0N4Z3l4Y1BuVG84S2V1N1RkNUduaWYzM3RhaVo4bS9MSlIvN3NDSUdyZ243bVN5TWpQL3BSTgp3WlZmNlJzV2x1dnZ2My93bkkrZHlySWlmUzNOCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJkakNDQVIyZ0F3SUJBZ0lCQURBS0JnZ3Foa2pPUFFRREFqQWpNU0V3SHdZRFZRUUREQmhyTTNNdFkyeHAKWlc1MExXTmhRREUzTWpjMk1EazNNakV3SGhjTk1qUXdPVEk1TVRFek5USXhXaGNOTXpRd09USTNNVEV6TlRJeApXakFqTVNFd0h3WURWUVFEREJock0zTXRZMnhwWlc1MExXTmhRREUzTWpjMk1EazNNakV3V1RBVEJnY3Foa2pPClBRSUJCZ2dxaGtqT1BRTUJCd05DQUFUaDBYUGY3bHhhaGhMbm5FQWk2V2NuYVBpYmp1UlVnLzQ5WllNc1dmQjAKY1ZGL2Irb05ZcnZaK0xxRW9vTHkrV1R4SnBuUGxCaHhSUnBoZXg3d2lwWXBvMEl3UURBT0JnTlZIUThCQWY4RQpCQU1DQXFRd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBZEJnTlZIUTRFRmdRVWJYV0JhRXFFWFg3ckZxaE5pbkZ0CmFTOTNLVGN3Q2dZSUtvWkl6ajBFQXdJRFJ3QXdSQUlnSHh2R1lzOVk1MmJvbytTTFNSTU1jN2RqS1JIS3EwOUUKb0gzakgyU3hKS1VDSURVTFlMVWt4b29kYlV4cEpoeGJqVGx1dm1HMXB2bnl5eVdSdjNsYUFiakYKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
    client-key-data: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1IY0NBUUVFSVBmVDI0Mk91NDhyQWkvUUF3dk4vVjlFK2x2UW5STFBDR3liMnpSbm5tQVhvQW9HQ0NxR1NNNDkKQXdFSG9VUURRZ0FFT0dDZWtIdUNQd0xFMlhFT0RvSUxDU0xtSkJXZGp0c2hSdk55QWM4MWoyYmVVSW8zMXNkKwo3aEpUOUNLVU95QXM2Vi9jRDFkUnBXUk5PSlZra2tmQ01BPT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo=`

//go:embed testdata/app.yaml
var data string

func TestKubeClient(t *testing.T) {
	k := NewKube()
	ctx := context.Background()

	err := k.Apply(ctx, conf, data)
	require.NoError(t, err)
}

func TestInvalidNamespaceName(t *testing.T) {

}
