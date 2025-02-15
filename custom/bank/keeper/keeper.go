package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	banktypes "github.com/furya-official/furya/custom/bank/types"
	furyakeeper "github.com/furya-official/furya/x/furya/keeper"
	furyatypes "github.com/furya-official/furya/x/furya/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Keeper struct {
	bankkeeper.BaseKeeper

	ak   furyakeeper.Keeper
	sk   banktypes.StakingKeeper
	acck accountkeeper.AccountKeeper
}

var (
	_ bankkeeper.Keeper = Keeper{}
)

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	ak accountkeeper.AccountKeeper,
	paramSpace paramtypes.Subspace,
	blockedAddrs map[string]bool,
) Keeper {
	keeper := Keeper{
		BaseKeeper: bankkeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs),
		ak:         furyakeeper.Keeper{},
		sk:         stakingkeeper.Keeper{},
		acck:       ak,
	}
	return keeper
}

func (k *Keeper) RegisterKeepers(ak furyakeeper.Keeper, sk banktypes.StakingKeeper) {
	k.ak = ak
	k.sk = sk
}

// SupplyOf implements the Query/SupplyOf gRPC method
func (k Keeper) SupplyOf(c context.Context, req *types.QuerySupplyOfRequest) (*types.QuerySupplyOfResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid denom")
	}

	ctx := sdk.UnwrapSDKContext(c)
	supply := k.GetSupply(ctx, req.Denom)

	if req.Denom == k.sk.BondDenom(ctx) {
		assets := k.ak.GetAllAssets(ctx)
		totalRewardWeights := sdk.ZeroDec()
		for _, asset := range assets {
			totalRewardWeights = totalRewardWeights.Add(asset.RewardWeight)
		}
		furyaBonded := k.ak.GetFuryaBondedAmount(ctx, k.acck.GetModuleAddress(furyatypes.ModuleName))
		supply.Amount = supply.Amount.Sub(furyaBonded)
	}

	return &types.QuerySupplyOfResponse{Amount: sdk.NewCoin(req.Denom, supply.Amount)}, nil
}

// TotalSupply implements the Query/TotalSupply gRPC method
func (k Keeper) TotalSupply(ctx context.Context, req *types.QueryTotalSupplyRequest) (*types.QueryTotalSupplyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	totalSupply, pageRes, err := k.GetPaginatedTotalSupply(sdkCtx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	furyaBonded := k.ak.GetFuryaBondedAmount(sdkCtx, k.acck.GetModuleAddress(furyatypes.ModuleName))
	bondDenom := k.sk.BondDenom(sdkCtx)
	if totalSupply.AmountOf(bondDenom).IsPositive() {
		totalSupply = totalSupply.Sub(sdk.NewCoin(bondDenom, furyaBonded))
	}

	return &types.QueryTotalSupplyResponse{Supply: totalSupply, Pagination: pageRes}, nil
}
