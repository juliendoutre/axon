package worker

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/juliendoutre/axon/internal/extraction"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/structpb"
)

func New(pg *pgxpool.Pool) *Worker {
	return &Worker{pg: pg}
}

type Worker struct {
	pg *pgxpool.Pool
}

func (w *Worker) ExtractAssetsFromObservation(ctx workflow.Context, observationID string) error {
	var serializedAttributes []byte
	if err := workflow.ExecuteActivity(ctx, w.GetObservationAttributes, observationID).
		Get(ctx, &serializedAttributes); err != nil {
		return fmt.Errorf("getting observation attributes: %w", err)
	}

	attributes := &structpb.Struct{}

	if err := attributes.UnmarshalJSON(serializedAttributes); err != nil {
		return fmt.Errorf("unmarshalling observation attributes: %w", err)
	}

	for _, candidate := range extraction.ExtractCandidatesFromStruct(attributes, "$") {
		for _, asset := range extraction.ListMatches(candidate) {
			if err := workflow.ExecuteActivity(ctx, w.InsertExtractedAsset, &InsertExtractedAssetInput{
				ObservationID:  observationID,
				AttributesPath: candidate.Path,
				AssetType:      asset.Type,
				AssetID:        asset.ID,
			}).Get(ctx, &attributes); err != nil {
				return fmt.Errorf("inserting extracted asset: %w", err)
			}
		}
	}

	return nil
}

func (w *Worker) GetObservationAttributes(ctx context.Context, id string) ([]byte, error) {
	var attributes []byte
	if err := w.pg.QueryRow(
		ctx,
		"SELECT attributes FROM axon.observations WHERE id = $1;",
		id,
	).Scan(&attributes); err != nil {
		return nil, fmt.Errorf("getting observation attributes by id: %w", err)
	}

	return attributes, nil
}

type InsertExtractedAssetInput struct {
	ObservationID  string
	AttributesPath string
	AssetType      string
	AssetID        string
}

func (w *Worker) InsertExtractedAsset(ctx context.Context, input *InsertExtractedAssetInput) error {
	if _, err := w.pg.Exec(
		ctx,
		`INSERT INTO axon.extracted_assets (observation_id, attributes_path, asset_type, asset_id)
VALUES ($1, $2, $3, $4);`,
		input.ObservationID,
		input.AttributesPath,
		input.AssetType,
		input.AssetID,
	); err != nil {
		return fmt.Errorf("inserting extracted asset: %w", err)
	}

	return nil
}
