package party_test

import (
	"github.com/me-next/menext-backend/party"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPreviousQueue(t *testing.T) {

	q := party.NewPreviousStack()
	_, err := q.Pop()

	assert.NotNil(t, err)

	// try adding to teh queue

	q.Push("a")

	sid, err := q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, party.SongUID("a"), sid)

	// check that queue is still empty
	_, err = q.Pop()
	assert.NotNil(t, err)
}
