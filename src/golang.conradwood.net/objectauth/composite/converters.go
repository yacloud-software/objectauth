package composite

import (
	"context"
	"golang.conradwood.net/apis/artefact"
)

//given an artefactid, get the git repo id
func get_git_repo_for_artefact(ctx context.Context, artefactid uint64) (uint64, error) {
	/*
		if artefactid == 231 {
			// testing
			return 332, nil
		}
	*/
	afid, err := artefact.GetArtefactServiceClient().GetRepoForArtefact(ctx, &artefact.ID{ID: artefactid})
	if err != nil {
		return 0, err
	}
	return afid.ID, nil
}

