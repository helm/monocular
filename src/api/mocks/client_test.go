package mocks

import (
	"testing"

	"github.com/arschles/assert"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
)

func TestMockedClient(t *testing.T) {
	client := NewMockedClient()
	assert.Equal(t, client, &mockedClient{}, "Empty mocked client")
	clientBroken := NewMockedBrokenClient()
	assert.Equal(t, clientBroken, &mockedBrokenClient{}, "Empty mocked brokenclient")

	_, err := client.ListReleases(releasesapi.GetAllReleasesParams{})
	assert.NoErr(t, err)
	_, err = clientBroken.ListReleases(releasesapi.GetAllReleasesParams{})
	assert.ExistsErr(t, err, "Cant list")

	_, err = client.InstallRelease("foo", releasesapi.CreateReleaseParams{})
	assert.NoErr(t, err)
	_, err = clientBroken.InstallRelease("foo", releasesapi.CreateReleaseParams{})
	assert.ExistsErr(t, err, "Cant list")

	_, err = client.DeleteRelease("foo")
	assert.NoErr(t, err)
	_, err = clientBroken.DeleteRelease("foo")
	assert.ExistsErr(t, err, "Cant list")

	_, err = client.GetRelease("foo")
	assert.NoErr(t, err)
	_, err = clientBroken.GetRelease("foo")
	assert.ExistsErr(t, err, "Cant list")
}
