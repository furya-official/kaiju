package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	test_helpers "github.com/furya-official/furya/app"
	"github.com/furya-official/furya/x/furya/keeper"
	"github.com/furya-official/furya/x/furya/types"
)

var ULUNA_FURYA = "ufury"

func TestQueryFuryas(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:        FURYA_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.NewDec(0),
				TotalTokens:  sdk.ZeroInt(),
			},
			{
				Denom:        FURYA_2_TOKEN_DENOM,
				RewardWeight: sdk.NewDec(10),
				TakeRate:     sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING THE FURYAS LIST
	furyas, err := queryServer.Furyas(ctx, &types.QueryFuryasRequest{})

	// THEN: VALIDATE THAT BOTH FURYAS HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryFuryasResponse{
		Furyas: []types.FuryaAsset{
			{
				Denom:                "furya",
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                "furya2",
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   2,
		},
	}, furyas)
}

func TestQueryAnUniqueFurya(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                FURYA_2_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING THE FURYAS LIST
	furyas, err := queryServer.Furya(ctx, &types.QueryFuryaRequest{
		Denom: "furya2",
	})

	// THEN: VALIDATE THAT BOTH FURYAS HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryFuryaResponse{
		Furya: &types.FuryaAsset{
			Denom:                "furya2",
			RewardWeight:         sdk.NewDec(10),
			TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
			TotalTokens:          sdk.ZeroInt(),
			TotalValidatorShares: sdk.NewDec(0),
			RewardChangeRate:     sdk.NewDec(0),
			RewardChangeInterval: 0,
		},
	}, furyas)
}

func TestQueryAnUniqueIBCFurya(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                "ibc/" + FURYA_2_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING THE FURYAS LIST
	furyas, err := queryServer.IBCFurya(ctx, &types.QueryIBCFuryaRequest{
		Hash: "furya2",
	})

	// THEN: VALIDATE THAT BOTH FURYAS HAVE THE CORRECT MODEL WHEN QUERYING
	require.Nil(t, err)
	require.Equal(t, &types.QueryFuryaResponse{
		Furya: &types.FuryaAsset{
			Denom:                "ibc/furya2",
			RewardWeight:         sdk.NewDec(10),
			TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
			TotalTokens:          sdk.ZeroInt(),
			TotalValidatorShares: sdk.NewDec(0),
			RewardChangeRate:     sdk.NewDec(0),
			RewardChangeInterval: 0,
		},
	}, furyas)
}

func TestQueryFuryaNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING THE FURYA
	_, err := queryServer.Furya(ctx, &types.QueryFuryaRequest{
		Denom: "furya2",
	})

	// THEN: VALIDATE THE ERROR
	require.Equal(t, err.Error(), "furya asset is not whitelisted")
}

func TestQueryAllFuryas(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING THE FURYA
	res, err := queryServer.Furyas(ctx, &types.QueryFuryasRequest{})

	// THEN: VALIDATE THE ERROR
	require.Nil(t, err)
	require.Equal(t, len(res.Furyas), 0)
	require.Equal(t, res.Pagination, &query.PageResponse{
		NextKey: nil,
		Total:   0,
	})
}

func TestQueryParams(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH AN FURYA ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING THE PARAMS...
	queyParams, err := queryServer.Params(ctx, &types.QueryParamsRequest{})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND OUTPUT IS AS WE EXPECT
	require.Nil(t, err)

	require.Equal(t, queyParams.Params.RewardDelayTime, time.Hour)
	require.Equal(t, queyParams.Params.TakeRateClaimInterval, time.Minute*5)

	// there is no way to match the exact time when the module is being instantiated
	// but we know that this time should be older than actually the time when this
	// following two lines are executed
	require.NotNil(t, queyParams.Params.LastTakeRateClaimTime)
	require.LessOrEqual(t, queyParams.Params.LastTakeRateClaimTime, time.Now())
}

func TestClaimQueryReward(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH ACCOUNTS
	app, ctx := createTestContext(t)
	startTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(startTime)
	ctx = ctx.WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.Params{
			RewardDelayTime:       time.Minute * 60,
			TakeRateClaimInterval: time.Minute * 5,
			LastTakeRateClaimTime: startTime,
		},
		Assets: []types.FuryaAsset{
			{
				Denom:                ULUNA_FURYA,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.MustNewDecFromStr("0.00005"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	feeCollectorAddr := app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val1, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)
	delAddr := test_helpers.AddTestAddrsIncremental(app, ctx, 1, sdk.NewCoins(sdk.NewCoin(ULUNA_FURYA, sdk.NewInt(1000_000_000))))[0]

	// WHEN: DELEGATING ...
	delRes, delErr := app.FuryaKeeper.Delegate(ctx, delAddr, val1, sdk.NewCoin(ULUNA_FURYA, sdk.NewInt(1000_000_000)))
	require.Nil(t, delErr)
	require.Equal(t, sdk.NewDec(1000000000), *delRes)
	assets := app.FuryaKeeper.GetAllAssets(ctx)
	err := app.FuryaKeeper.RebalanceBondTokenWeights(ctx, assets)
	require.NoError(t, err)

	// ...and advance block...
	timePassed := time.Minute*5 + time.Second
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(timePassed))
	ctx = ctx.WithBlockHeight(2)
	app.FuryaKeeper.DeductAssetsHook(ctx, assets)
	app.BankKeeper.GetAllBalances(ctx, feeCollectorAddr)
	require.Equal(t, startTime.Add(time.Minute*5), app.FuryaKeeper.LastRewardClaimTime(ctx))
	app.FuryaKeeper.GetAssetByDenom(ctx, ULUNA_FURYA)

	// ... at the next begin block, tokens will be distributed from the fee pool...
	cons, _ := val1.GetConsAddr()
	app.DistrKeeper.AllocateTokens(ctx, 1, 1, cons, []abcitypes.VoteInfo{
		{
			Validator: abcitypes.Validator{
				Address: cons,
				Power:   1,
			},
			SignedLastBlock: true,
		},
	})

	// THEN: Query the delegation rewards ...
	queryDelegation, queryErr := queryServer.FuryaDelegationRewards(ctx, &types.QueryFuryaDelegationRewardsRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: valAddr.String(),
		Denom:         ULUNA_FURYA,
	})

	// ... validate that no error has been produced.
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryFuryaDelegationRewardsResponse{
		Rewards: []sdk.Coin{
			{
				Denom:  ULUNA_FURYA,
				Amount: math.NewInt(32666),
			},
		},
	}, queryDelegation)
}

func TestQueryFuryaDelegation(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)))
	queryDelegation, queryErr := queryServer.FuryaDelegation(ctx, &types.QueryFuryaDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
		Denom:         FURYA_TOKEN_DENOM,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryFuryaDelegationResponse{
		Delegation: types.DelegationResponse{
			Delegation: types.Delegation{
				DelegatorAddress:      delAddr.String(),
				ValidatorAddress:      val.OperatorAddress,
				Denom:                 FURYA_TOKEN_DENOM,
				Shares:                sdk.NewDec(1000_000),
				RewardHistory:         nil,
				LastRewardClaimHeight: uint64(ctx.BlockHeight()),
			},
			Balance: sdk.Coin{
				Denom:  FURYA_TOKEN_DENOM,
				Amount: sdk.NewInt(1000_000),
			},
		},
	}, queryDelegation)
	require.Equal(t, sdk.NewDec(1000000), *delegationTxRes)
}

func TestQueryFuryaDelegationNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.StakingKeeper.GetValidator(ctx, valAddr)
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING ...
	_, err := queryServer.FuryaDelegation(ctx, &types.QueryFuryaDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
		Denom:         FURYA_TOKEN_DENOM,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "FuryaAsset not found by denom furya"))
}

func TestQueryFuryaValidatorNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING ...
	_, err := queryServer.FuryaDelegation(ctx, &types.QueryFuryaDelegationRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: "cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu",
		Denom:         FURYA_TOKEN_DENOM,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "Validator not found by address cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu"))
}

func TestQueryFuryasDelegationByValidator(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)))
	queryDelegation, queryErr := queryServer.FuryasDelegationByValidator(ctx, &types.QueryFuryasDelegationByValidatorRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: val.OperatorAddress,
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryFuryasDelegationsResponse{
		Delegations: []types.DelegationResponse{
			{
				Delegation: types.Delegation{
					DelegatorAddress:      delAddr.String(),
					ValidatorAddress:      val.OperatorAddress,
					Denom:                 FURYA_TOKEN_DENOM,
					Shares:                sdk.NewDec(1000_000),
					RewardHistory:         nil,
					LastRewardClaimHeight: uint64(ctx.BlockHeight()),
				},
				Balance: sdk.Coin{
					Denom:  FURYA_TOKEN_DENOM,
					Amount: sdk.NewInt(1000_000),
				},
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	}, queryDelegation)
	require.Equal(t, sdk.NewDec(1000_000), *delegationTxRes)
}

func TestQueryFuryasDelegationByValidatorNotFound(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)

	// WHEN: QUERYING ...
	_, err := queryServer.FuryasDelegationByValidator(ctx, &types.QueryFuryasDelegationByValidatorRequest{
		DelegatorAddr: delAddr.String(),
		ValidatorAddr: "cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu",
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Equal(t, err, status.Error(codes.NotFound, "Validator not found by address cosmosvaloper19lss6zgdh5vvcpjhfftdghrpsw7a4434elpwpu"))
}

func TestQueryFuryasFuryasDelegation(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                FURYA_2_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	delegationTxRes, txErr := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)))
	delegation2TxRes, tx2Err := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	queryDelegation, queryErr := queryServer.FuryasDelegation(ctx, &types.QueryFuryasDelegationsRequest{
		DelegatorAddr: delAddr.String(),
	})

	// THEN: VALIDATE THAT NO ERRORS HAVE BEEN PRODUCED AND BOTH OUTPUTS ARE AS WE EXPECT
	require.Nil(t, txErr)
	require.Nil(t, tx2Err)
	require.Nil(t, queryErr)
	require.Equal(t, &types.QueryFuryasDelegationsResponse{
		Delegations: []types.DelegationResponse{
			{
				Delegation: types.Delegation{
					DelegatorAddress:      delAddr.String(),
					ValidatorAddress:      val.OperatorAddress,
					Denom:                 FURYA_TOKEN_DENOM,
					Shares:                sdk.NewDec(1000_000),
					RewardHistory:         nil,
					LastRewardClaimHeight: uint64(ctx.BlockHeight()),
				},
				Balance: sdk.Coin{
					Denom:  FURYA_TOKEN_DENOM,
					Amount: sdk.NewInt(1000_000),
				},
			},
			{
				Delegation: types.Delegation{
					DelegatorAddress:      delAddr.String(),
					ValidatorAddress:      val.OperatorAddress,
					Denom:                 FURYA_2_TOKEN_DENOM,
					Shares:                sdk.NewDec(1000_000),
					RewardHistory:         nil,
					LastRewardClaimHeight: uint64(ctx.BlockHeight()),
				},
				Balance: sdk.Coin{
					Denom:  FURYA_2_TOKEN_DENOM,
					Amount: sdk.NewInt(1000_000),
				},
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   2,
		},
	}, queryDelegation)
	require.Equal(t, sdk.NewDec(1000_000), *delegationTxRes)
	require.Equal(t, sdk.NewDec(1000_000), *delegation2TxRes)
}

func TestQueryAllDelegations(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                FURYA_2_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	_, txErr := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, txErr)
	_, tx2Err := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, tx2Err)
	queryDelegations, queryErr := queryServer.AllFuryasDelegations(ctx, &types.QueryAllFuryasDelegationsRequest{
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})
	require.NoError(t, queryErr)
	require.Equal(t, 1, len(queryDelegations.Delegations))

	require.Equal(t, types.DelegationResponse{
		Delegation: types.Delegation{
			DelegatorAddress:      delAddr.String(),
			ValidatorAddress:      val.OperatorAddress,
			Denom:                 FURYA_TOKEN_DENOM,
			Shares:                sdk.NewDec(1000_000),
			RewardHistory:         nil,
			LastRewardClaimHeight: uint64(ctx.BlockHeight()),
		},
		Balance: sdk.Coin{
			Denom:  FURYA_TOKEN_DENOM,
			Amount: sdk.NewInt(1000_000),
		},
	}, queryDelegations.Delegations[0])

	queryDelegations, queryErr = queryServer.AllFuryasDelegations(ctx, &types.QueryAllFuryasDelegationsRequest{
		Pagination: &query.PageRequest{
			Key:        queryDelegations.Pagination.NextKey,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})
	require.NoError(t, queryErr)
	require.Equal(t, types.DelegationResponse{
		Delegation: types.Delegation{
			DelegatorAddress:      delAddr.String(),
			ValidatorAddress:      val.OperatorAddress,
			Denom:                 FURYA_2_TOKEN_DENOM,
			Shares:                sdk.NewDec(1000_000),
			RewardHistory:         nil,
			LastRewardClaimHeight: uint64(ctx.BlockHeight()),
		},
		Balance: sdk.Coin{
			Denom:  FURYA_2_TOKEN_DENOM,
			Amount: sdk.NewInt(1000_000),
		},
	}, queryDelegations.Delegations[0])
}

func TestQueryValidator(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                FURYA_2_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	_, txErr := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, txErr)
	_, tx2Err := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, tx2Err)

	queryVal, queryErr := queryServer.FuryaValidator(ctx, &types.QueryFuryaValidatorRequest{
		ValidatorAddr: val.GetOperator().String(),
	})

	require.NoError(t, queryErr)
	require.Equal(t, &types.QueryFuryaValidatorResponse{
		ValidatorAddr: val.GetOperator().String(),
		TotalDelegationShares: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(FURYA_TOKEN_DENOM, sdk.NewDec(1000000)),
			sdk.NewDecCoinFromDec(FURYA_2_TOKEN_DENOM, sdk.NewDec(1000000)),
		),
		ValidatorShares: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(FURYA_TOKEN_DENOM, sdk.NewDec(1000000)),
			sdk.NewDecCoinFromDec(FURYA_2_TOKEN_DENOM, sdk.NewDec(1000000)),
		),
		TotalStaked: sdk.NewDecCoins(
			sdk.NewDecCoinFromDec(FURYA_TOKEN_DENOM, sdk.NewDec(1000_000)),
			sdk.NewDecCoinFromDec(FURYA_2_TOKEN_DENOM, sdk.NewDec(1000_000)),
		),
	}, queryVal)
}

func TestQueryValidators(t *testing.T) {
	// GIVEN: THE BLOCKCHAIN WITH FURYAS ON GENESIS
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx = ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.FuryaKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.FuryaAsset{
			{
				Denom:                FURYA_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.NewDec(0),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
			{
				Denom:                FURYA_2_TOKEN_DENOM,
				RewardWeight:         sdk.NewDec(10),
				TakeRate:             sdk.MustNewDecFromStr("0.14159265359"),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.FuryaKeeper)
	delegations := app.StakingKeeper.GetAllDelegations(ctx)
	delAddr, _ := sdk.AccAddressFromBech32(delegations[0].DelegatorAddress)
	valAddr, _ := sdk.ValAddressFromBech32(delegations[0].ValidatorAddress)
	val, _ := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr)

	addrs := test_helpers.AddTestAddrsIncremental(app, ctx, 3, sdk.NewCoins(
		sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)),
		sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(1000_000)),
	))
	valAddr2 := sdk.ValAddress(addrs[0])
	_val2 := teststaking.NewValidator(t, valAddr2, test_helpers.CreateTestPubKeys(1)[0])
	test_helpers.RegisterNewValidator(t, app, ctx, _val2)
	val2, err := app.FuryaKeeper.GetFuryaValidator(ctx, valAddr2)
	require.NoError(t, err)

	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, delAddr, sdk.NewCoins(sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(2000_000))))

	// WHEN: DELEGATING AND QUERYING ...
	_, txErr := app.FuryaKeeper.Delegate(ctx, delAddr, val, sdk.NewCoin(FURYA_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, txErr)
	_, tx2Err := app.FuryaKeeper.Delegate(ctx, delAddr, val2, sdk.NewCoin(FURYA_2_TOKEN_DENOM, sdk.NewInt(1000_000)))
	require.NoError(t, tx2Err)

	queryVal, queryErr := queryServer.AllFuryaValidators(ctx, &types.QueryAllFuryaValidatorsRequest{
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})

	require.NoError(t, queryErr)
	require.Equal(t, &types.QueryFuryaValidatorsResponse{
		Validators: []types.QueryFuryaValidatorResponse{
			{
				ValidatorAddr: val.GetOperator().String(),
				TotalDelegationShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(FURYA_TOKEN_DENOM, sdk.NewDec(1000000)),
				),
				ValidatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(FURYA_TOKEN_DENOM, sdk.NewDec(1000000)),
				),
				TotalStaked: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(FURYA_TOKEN_DENOM, sdk.NewDec(1000_000)),
				),
			},
		},
		Pagination: queryVal.Pagination,
	}, queryVal)

	queryVal2, queryErr := queryServer.AllFuryaValidators(ctx, &types.QueryAllFuryaValidatorsRequest{
		Pagination: &query.PageRequest{
			Key:        queryVal.Pagination.NextKey,
			Offset:     0,
			Limit:      1,
			CountTotal: false,
			Reverse:    false,
		},
	})

	require.NoError(t, queryErr)
	require.Equal(t, &types.QueryFuryaValidatorsResponse{
		Validators: []types.QueryFuryaValidatorResponse{
			{
				ValidatorAddr: val2.GetOperator().String(),
				TotalDelegationShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(FURYA_2_TOKEN_DENOM, sdk.NewDec(1000000)),
				),
				ValidatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(FURYA_2_TOKEN_DENOM, sdk.NewDec(1000000)),
				),
				TotalStaked: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec(FURYA_2_TOKEN_DENOM, sdk.NewDec(1000_000)),
				),
			},
		},
		Pagination: queryVal2.Pagination,
	}, queryVal2)
}
