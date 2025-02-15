package furya

import (
	"fmt"
	"time"

	"github.com/furya-official/furya/x/furya/keeper"
	"github.com/furya-official/furya/x/furya/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	k.CompleteRedelegations(ctx)
	if err := k.CompleteUndelegations(ctx); err != nil {
		panic(fmt.Errorf("Failed to complete undelegations from x/furya module: %s", err))
	}

	assets := k.GetAllAssets(ctx)
	if _, err := k.DeductAssetsHook(ctx, assets); err != nil {
		panic(fmt.Errorf("Failed to deduct take rate from furya in x/furya module: %s", err))
	}
	k.RewardWeightChangeHook(ctx, assets)
	if err := k.RebalanceHook(ctx, assets); err != nil {
		panic(fmt.Errorf("Failed to rebalance assets in x/furya module: %s", err))
	}
	return []abci.ValidatorUpdate{}
}
