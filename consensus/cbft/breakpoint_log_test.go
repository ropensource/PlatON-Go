package cbft

import (
	"context"
	"os"
	"testing"
)

func TestLogPrepareBP_ReceiveVote(t *testing.T) {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)
	ctx := context.WithValue(context.TODO(), "peer","xxxxxxxx")
	pvote := makeFakePrepareVote()
	logBP.PrepareBP().ReceiveVote(ctx, pvote, engine)

}
