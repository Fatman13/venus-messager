package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/venus-messager/config"
	"github.com/filecoin-project/venus-messager/filestore"
	"github.com/filecoin-project/venus-messager/log"
	"github.com/filecoin-project/venus-messager/models"
	"github.com/filecoin-project/venus-messager/models/sqlite"
	types "github.com/filecoin-project/venus/venus-shared/types/messager"
)

func TestMessageStateCache(t *testing.T) {
	fs := filestore.NewMockFileStore(nil)
	db, err := sqlite.OpenSqlite(fs)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.Remove("message.db"))
		assert.NoError(t, os.Remove("message.db-shm"))
		assert.NoError(t, os.Remove("message.db-wal"))
	}()
	assert.NoError(t, db.AutoMigrate())

	msgs := models.NewSignedMessages(10)
	for _, msg := range msgs {
		err := db.MessageRepo().CreateMessage(msg)
		assert.NoError(t, err)
	}

	msgState, err := NewMessageState(db, log.New(), &config.MessageStateConfig{
		BackTime:          60,
		CleanupInterval:   3,
		DefaultExpiration: 2,
	})
	assert.NoError(t, err)

	msgList, err := msgState.repo.MessageRepo().ListMessage()
	assert.NoError(t, err)
	assert.Equal(t, 10, len(msgList))

	assert.NoError(t, msgState.loadRecentMessage())
	assert.Equal(t, 10, len(msgState.idCids.cache))

	state, flag := msgState.GetMessageStateByCid(msgs[0].Cid().String())
	assert.True(t, flag)
	assert.Equal(t, msgs[0].State, state)

	err = msgState.UpdateMessageByCid(msgs[1].Cid(), func(message *types.Message) error {
		message.State = types.OnChainMsg
		return nil
	})
	assert.NoError(t, err)
	state, flag = msgState.GetMessageStateByCid(msgs[1].Cid().String())
	assert.True(t, flag)
	assert.Equal(t, types.OnChainMsg, state)
}
