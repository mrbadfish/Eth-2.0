package spectest

import (
	"path"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

func runProposerSlashingTest(t *testing.T, config string) {
	resetCfg, err := spectest.SetConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	defer resetCfg()

	testFolders, testsFolderPath := testutil.TestFolders(t, config, "operations/proposer_slashing/pyspec_tests")
	for _, folder := range testFolders {
		t.Run(folder.Name(), func(t *testing.T) {
			folderPath := path.Join(testsFolderPath, folder.Name())
			proposerSlashingFile, err := testutil.BazelFileBytes(folderPath, "proposer_slashing.ssz")
			if err != nil {
				t.Fatal(err)
			}
			proposerSlashing := &ethpb.ProposerSlashing{}
			if err := ssz.Unmarshal(proposerSlashingFile, proposerSlashing); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			body := &ethpb.BeaconBlockBody{ProposerSlashings: []*ethpb.ProposerSlashing{proposerSlashing}}
			testutil.RunBlockOperationTest(t, folderPath, body, blocks.ProcessProposerSlashings)
		})
	}
}
