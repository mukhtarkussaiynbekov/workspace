package sampleapp_test

import (
	"context"
	"testing"

	"github.com/mukhtarkv/workspace/sample/sampleapp"
	"github.com/mukhtarkv/workspace/sample/sampleapp/inmem"
	"github.com/stretchr/testify/assert"
)

func TestSampleApp(t *testing.T) {
	storage := inmem.New()
	userSrv := sampleapp.New(storage)

	u := sampleapp.User{
		Name: "Toto",
	}

	ctx := context.TODO()
	err := userSrv.Create(ctx, &u)
	assert.NoError(t, err)

	user, err := userSrv.Fetch(ctx, u.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, &u, user)
}
