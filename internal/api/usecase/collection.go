package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zd4r/artblocks-stats/internal/api/entity"
)

type CollectionUseCase struct {
	repo   HoldersRepo
	webAPI CollectionWebAPI
}

// NewCollection creates new collection use case
func NewCollection(r HoldersRepo, w CollectionWebAPI) *CollectionUseCase {
	return &CollectionUseCase{
		repo:   r,
		webAPI: w,
	}
}

// СalculateStats gather holders of provided collection and calculate their distribution based on scores
func (uc *CollectionUseCase) СalculateStats(ctx context.Context, c entity.Collection) (entity.Collection, error) {
	collection, err := uc.GetHolders(ctx, c)
	if err != nil {
		return entity.Collection{}, err
	}

	collection.CountHoldersDistribution()

	return collection, nil
}

// GetHolders gather holders of provided collection
func (uc *CollectionUseCase) GetHolders(ctx context.Context, c entity.Collection) (entity.Collection, error) {
	collection, err := uc.webAPI.GetHoldersCount(c)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - uc.webAPI.GetHoldersCount: %w", err)
	}

	collection.Holders = make([]entity.Holder, collection.HoldersCount)

	collection, err = uc.webAPI.GetHolders(collection)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - uc.webAPI.GetHolders: %w", err)
	}

	collection, err = uc.GatherHoldersScores(ctx, collection)
	if err != nil {
		return entity.Collection{}, err
	}

	return collection, nil
}

// GatherHoldersScores - fills []entity.Holder with scores from repo or from Artacle API if data in repo is outdated
func (uc *CollectionUseCase) GatherHoldersScores(ctx context.Context, c entity.Collection) (entity.Collection, error) {
	for i, h := range c.Holders {
		holder, err := uc.repo.Get(h)
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrHolderNotFound):
				holder, err = uc.webAPI.GetHolderScores(h)
				if err != nil {
					return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - uc.webAPI.GetHolderScores: %w", err)
				}

				holder, err = uc.repo.Insert(holder)
				if err != nil {
					return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - uc.repo.Insert: %w", err)
				}
			default:
				return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - us.repo.Get: %w", err)
			}
		}

		// TODO: Fix: expiration date hardcode
		if time.Now().Sub(holder.UpdatedAt) >= 72*time.Hour {
			holder, err = uc.webAPI.GetHolderScores(holder)
			if err != nil {
				return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - uc.webAPI.GetHolderScores: %w", err)
			}

			holder, err = uc.repo.UpdateScores(holder)
			if err != nil {
				return entity.Collection{}, fmt.Errorf("CollectionUseCase - GetHolders - uc.repo.UpdateScores: %w", err)
			}
		}
		c.Holders[i] = holder
	}

	return c, nil
}
