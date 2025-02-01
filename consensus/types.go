package consensus

import (
	"fmt"
	"io"
	"strings"

	"github.com/xufeisofly/hotstuff/types"
)

type Hash = []byte
type HashStr = string

// SyncInfo holds the highest known QC or TC.
// Generally, if highQC.View > highTC.View, there is no need to include highTC in the SyncInfo.
// However, if highQC.View < highTC.View, we should still include highQC.
// This can also hold an AggregateQC for Fast-Hotstuff.
type SyncInfo struct {
	qc    *QuorumCert
	tc    *TimeoutCert
	aggQC *AggregateQC
}

// NewSyncInfo returns a new SyncInfo struct.
func NewSyncInfo() SyncInfo {
	return SyncInfo{}
}

// WithQC returns a copy of the SyncInfo struct with the given QC.
func (si SyncInfo) WithQC(qc QuorumCert) SyncInfo {
	si.qc = new(QuorumCert)
	*si.qc = qc
	return si
}

// WithTC returns a copy of the SyncInfo struct with the given TC.
func (si SyncInfo) WithTC(tc TimeoutCert) SyncInfo {
	si.tc = new(TimeoutCert)
	*si.tc = tc
	return si
}

// WithAggQC returns a copy of the SyncInfo struct with the given AggregateQC.
func (si SyncInfo) WithAggQC(aggQC AggregateQC) SyncInfo {
	si.aggQC = new(AggregateQC)
	*si.aggQC = aggQC
	return si
}

// QuorumCert(QC) is certificate for Block created by a quorum of partial certificates.
type QuorumCert struct {
	signature types.QuorumSignature
	view      types.View
	blockHash Hash
}

func NewQuorumCert(signature types.QuorumSignature, view types.View, blockHash Hash) QuorumCert {
	return QuorumCert{signature, view, blockHash}
}

func (qc QuorumCert) Signature() types.QuorumSignature {
	return qc.signature
}

func (qc QuorumCert) View() types.View {
	return qc.view
}

func (qc QuorumCert) BlockHash() Hash {
	return qc.blockHash
}

func (qc QuorumCert) String() string {
	var sb strings.Builder
	if qc.signature != nil {
		_ = writeParticipants(&sb, qc.Signature().Participants())
	}
	return fmt.Sprintf("QC{ hash: %.6s, Addrs: [ %s] }", qc.blockHash, &sb)
}

// AggregateQC is a set of QCs extracted from timeout messages
// and an aggregate signature of the timeout signatures.
// This is used by the Fast-HotStuff consensus protocol.
type AggregateQC struct {
	qcs       map[types.AddressStr]QuorumCert
	signature types.QuorumSignature
	view      types.View
}

func NewAggregateQC(
	qcs map[types.AddressStr]QuorumCert,
	signature types.QuorumSignature,
	view types.View,
) AggregateQC {
	return AggregateQC{qcs, signature, view}
}

func (aggQC AggregateQC) QCs() map[types.AddressStr]QuorumCert {
	return aggQC.qcs
}

func (aggQC AggregateQC) Signature() types.QuorumSignature {
	return aggQC.signature
}

func (aggQC AggregateQC) View() types.View {
	return aggQC.view
}

func (aggQC AggregateQC) String() string {
	var sb strings.Builder
	if aggQC.signature != nil {
		_ = writeParticipants(&sb, aggQC.signature.Participants())
	}
	return fmt.Sprintf("AggQC{ view: %d, Addrs: [ %s] }", aggQC.view, &sb)
}

func writeParticipants(wr io.Writer, participants types.AddressSet) (err error) {
	participants.RangeWhile(func(addr types.Address) bool {
		_, err = fmt.Fprintf(wr, "%s ", addr)
		return err == nil
	})
	return err
}

// TimeoutCert (TC) is a certificate created by a quorum of timeout messsages.
type TimeoutCert struct {
	signature types.QuorumSignature
	view      types.View
}

func NewTimeoutCert(signature types.QuorumSignature, view types.View) TimeoutCert {
	return TimeoutCert{signature, view}
}

func (tc TimeoutCert) Signature() types.QuorumSignature {
	return tc.signature
}

func (tc TimeoutCert) View() types.View {
	return tc.view
}

func (tc TimeoutCert) String() string {
	var sb strings.Builder
	if tc.signature != nil {
		_ = writeParticipants(&sb, tc.Signature().Participants())
	}
	return fmt.Sprintf("TC{ view: %d, Addrs: [ %s] }", tc.view, &sb)
}
