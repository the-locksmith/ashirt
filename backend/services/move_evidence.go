package services

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/theparanoids/ashirt/backend"
	"github.com/theparanoids/ashirt/backend/database"
	"github.com/theparanoids/ashirt/backend/policy"
)

type MoveEvidenceInput struct {
	SourceOperationSlug string
	EvidenceUUID        string
	TargetOperationSlug string
}

func MoveEvidence(ctx context.Context, db *database.Connection, i MoveEvidenceInput) error {
	sourceOperation, evidence, err := lookupOperationEvidence(db, i.SourceOperationSlug, i.EvidenceUUID)
	if err != nil {
		return backend.UnauthorizedReadErr(err)
	}

	destinationOperation, err := lookupOperation(db, i.TargetOperationSlug)
	if err != nil {
		return backend.UnauthorizedReadErr(err)
	}

	if err := policyRequireWithAdminBypass(ctx,
		policy.CanModifyOperation{OperationID: sourceOperation.ID},
		policy.CanModifyOperation{OperationID: destinationOperation.ID},
	); err != nil {
		return backend.UnauthorizedWriteErr(err)
	}

	//Check which tags can be migrated
	tagDifferences, err := ListTagDifferenceForEvidence(ctx, db, ListTagDifferenceForEvidenceInput{
		ListTagsDifferenceInput: ListTagsDifferenceInput{
			SourceOperationSlug:      i.SourceOperationSlug,
			DestinationOperationSlug: i.TargetOperationSlug,
		},
		SourceEvidenceUUID: i.EvidenceUUID,
	})

	if err != nil {
		return err
	}

	err = db.WithTx(ctx, func(tx *database.Transactable) {
		// remove findings
		tx.Delete(sq.Delete("evidence_finding_map").Where(sq.Eq{"evidence_id": evidence.ID}))
		// remove tags
		tx.Delete(sq.Delete("tag_evidence_map").Where(sq.Eq{"evidence_id": evidence.ID}))
		// reassociate evidence with new operation
		tx.Update(sq.Update("evidence").Set("operation_id", destinationOperation.ID).Where(sq.Eq{"id": evidence.ID}))
		// associate with common tags
		tx.BatchInsert("tag_evidence_map", len(tagDifferences.Included), func(idx int) map[string]interface{} {
			pair := tagDifferences.Included[idx]
			return map[string]interface{}{
				"tag_id":      pair.DestinationTag.ID,
				"evidence_id": evidence.ID,
			}
		})
	})

	return err
}
