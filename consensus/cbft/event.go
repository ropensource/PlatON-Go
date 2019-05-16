package cbft

type NewPrepareBlockEvent struct {
	PrepareBlock *prepareBlock
}

type NewPrepareVoteEvent struct {
	PrepareVote *prepareVote
}

type NewConfirmedPrepareBlockEvent struct {
	ConfirmedPrepareBlock *confirmedPrepareBlock
}

type NewPrepareBlockHashEvent struct {
	PrepareBlockHash *prepareBlockHash
}

type NewViewChangeEvent struct {
	ViewChange *viewChange
}

type NewViewChangeVoteEvent struct {
	ViewChangeVote *viewChangeVote
}
