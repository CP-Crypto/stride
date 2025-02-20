package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cast"

	"github.com/Stride-Labs/stride/v4/utils"
	recordstypes "github.com/Stride-Labs/stride/v4/x/records/types"
	"github.com/Stride-Labs/stride/v4/x/stakeibc/types"
)

func (k Keeper) CreateEpochUnbondingRecord(ctx sdk.Context, epochNumber uint64) bool {
	k.Logger(ctx).Info(fmt.Sprintf("Creating Epoch Unbonding Records for Epoch %d", epochNumber))

	hostZoneUnbondings := []*recordstypes.HostZoneUnbonding{}

	for _, hostZone := range k.GetAllHostZone(ctx) {
		k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "Creating Epoch Unbonding Record"))

		hostZoneUnbonding := recordstypes.HostZoneUnbonding{
			NativeTokenAmount: sdk.ZeroInt(),
			StTokenAmount:     sdk.ZeroInt(),
			Denom:             hostZone.HostDenom,
			HostZoneId:        hostZone.ChainId,
			Status:            recordstypes.HostZoneUnbonding_UNBONDING_QUEUE,
		}
		hostZoneUnbondings = append(hostZoneUnbondings, &hostZoneUnbonding)
	}

	epochUnbondingRecord := recordstypes.EpochUnbondingRecord{
		EpochNumber:        cast.ToUint64(epochNumber),
		HostZoneUnbondings: hostZoneUnbondings,
	}
	k.RecordsKeeper.SetEpochUnbondingRecord(ctx, epochUnbondingRecord)
	return true
}

func (k Keeper) GetHostZoneUnbondingMsgs(ctx sdk.Context, hostZone types.HostZone) (msgs []sdk.Msg, totalAmtToUnbond sdk.Int, marshalledCallbackArgs []byte, epochUnbondingRecordIds []uint64, err error) {
	totalAmtToUnbond = sdk.ZeroInt()
	// this function goes and processes all unbonded records for this hostZone
	// regardless of what epoch they belong to
	for _, epochUnbonding := range k.RecordsKeeper.GetAllEpochUnbondingRecord(ctx) {
		hostZoneRecord, found := k.RecordsKeeper.GetHostZoneUnbondingByChainId(ctx, epochUnbonding.EpochNumber, hostZone.ChainId)
		if !found {
			errMsg := fmt.Sprintf("Host zone unbonding record not found for hostZoneId %s in epoch %d",
				hostZone.ChainId, epochUnbonding.GetEpochNumber())
			k.Logger(ctx).Error(errMsg)
			continue
		}
		// mark the epoch unbonding record for processing if it's bonded and the host zone unbonding has an amount g.t. zero
		if hostZoneRecord.Status == recordstypes.HostZoneUnbonding_UNBONDING_QUEUE && hostZoneRecord.NativeTokenAmount.GT(sdk.ZeroInt()) {
			totalAmtToUnbond = totalAmtToUnbond.Add(hostZoneRecord.NativeTokenAmount)
			epochUnbondingRecordIds = append(epochUnbondingRecordIds, epochUnbonding.EpochNumber)
			k.Logger(ctx).Info(fmt.Sprintf("[SendHostZoneUnbondings] Sending unbondings, host zone: %s, epochUnbonding: %v", hostZone.ChainId, epochUnbonding))

		}
	}
	delegationAccount := hostZone.GetDelegationAccount()
	if delegationAccount == nil || delegationAccount.GetAddress() == "" {
		errMsg := fmt.Sprintf("Zone %s is missing a delegation address!", hostZone.ChainId)
		k.Logger(ctx).Error(errMsg)
		return nil, sdk.ZeroInt(), nil, nil, sdkerrors.Wrap(types.ErrHostZoneICAAccountNotFound, errMsg)
	}
	validators := hostZone.GetValidators()
	if totalAmtToUnbond.IsZero() {
		return nil, sdk.ZeroInt(), nil, nil, nil
	}
	// we distribute the unbonding based on our target weights
	newUnbondingToValidator, err := k.GetTargetValAmtsForHostZone(ctx, hostZone, totalAmtToUnbond)
	if err != nil {
		errMsg := fmt.Sprintf("Error getting target val amts for host zone %s %d: %s", hostZone.ChainId, totalAmtToUnbond, err)
		k.Logger(ctx).Error(errMsg)
		return nil, sdk.ZeroInt(), nil, nil, sdkerrors.Wrap(types.ErrNoValidatorAmts, errMsg)
	}
	valAddrToUnbondAmt := make(map[string]sdk.Int)
	overflowAmt := sdk.ZeroInt()
	for _, validator := range validators {
		valAddr := validator.GetAddress()
		valUnbondAmt := newUnbondingToValidator[valAddr]
		k.Logger(ctx).Info(fmt.Sprintf("Validator %s - Weight: %d, Expected Unbond Amount: %d, Current Delegations: %d", valAddr, validator.Weight, valUnbondAmt, validator.DelegationAmt))

		if valUnbondAmt.GT(validator.DelegationAmt) { // if we don't have enough assets to unbond
			overflowAmt = overflowAmt.Add(valUnbondAmt).Sub(validator.DelegationAmt)
			valUnbondAmt = validator.DelegationAmt
		}
		if err != nil {
			errMsg := fmt.Sprintf("Error casting validator staked amount %d: %s", validator.DelegationAmt, err.Error())
			k.Logger(ctx).Error(errMsg)
			return nil, sdk.ZeroInt(), nil, nil, sdkerrors.Wrap(types.ErrIntCast, errMsg)
		}
		valAddrToUnbondAmt[valAddr] = valUnbondAmt
	}
	if overflowAmt.GT(sdk.ZeroInt()) { // if we need to reallocate any weights
		k.Logger(ctx).Info(fmt.Sprintf("Expected validator undelegation amount on %s is greater than it's current delegations. Redistributing undelegations accordingly.", hostZone.ChainId))
		for _, validator := range validators {
			valAddr := validator.GetAddress()
			valUnbondAmt := valAddrToUnbondAmt[valAddr]
			// store how many more tokens we could unbond, if needed
			curAmtStaked := validator.DelegationAmt
			amtToPotentiallyUnbond := curAmtStaked.Sub(valUnbondAmt)
			if amtToPotentiallyUnbond.GT(sdk.ZeroInt()) { // if we can afford to unbond more
				if amtToPotentiallyUnbond.GT(overflowAmt) { // we can fully cover the overflow
					valAddrToUnbondAmt[valAddr] = valAddrToUnbondAmt[valAddr].Add(overflowAmt)
					overflowAmt = sdk.ZeroInt()
					break
				} else {
					valAddrToUnbondAmt[valAddr] = valAddrToUnbondAmt[valAddr].Add(amtToPotentiallyUnbond)
					overflowAmt = overflowAmt.Sub(amtToPotentiallyUnbond)
				}
			}
		}
	}
	if overflowAmt.GT(sdk.ZeroInt()) { // what?? we still can't cover the overflow? something is very wrong
		errMsg := fmt.Sprintf("Could not unbond %v on Host Zone %s, unable to balance the unbond amount across validators",
			totalAmtToUnbond, hostZone.ChainId)
		k.Logger(ctx).Error(errMsg)
		return nil, sdk.ZeroInt(), nil, nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, errMsg)
	}
	var splitDelegations []*types.SplitDelegation
	// DO NOT REMOVE: StringMapKeys fixes non-deterministic map iteration
	for _, valAddr := range utils.StringMapKeys(valAddrToUnbondAmt) {
		valUnbondAmt := valAddrToUnbondAmt[valAddr]
		stakeAmt := sdk.NewCoin(hostZone.HostDenom, valUnbondAmt)

		msgs = append(msgs, &stakingtypes.MsgUndelegate{
			DelegatorAddress: delegationAccount.GetAddress(),
			ValidatorAddress: valAddr,
			Amount:           stakeAmt,
		})

		splitDelegations = append(splitDelegations, &types.SplitDelegation{
			Validator: valAddr,
			Amount:    stakeAmt.Amount,
		})
	}

	undelegateCallback := types.UndelegateCallback{
		HostZoneId:              hostZone.ChainId,
		SplitDelegations:        splitDelegations,
		EpochUnbondingRecordIds: epochUnbondingRecordIds,
	}
	k.Logger(ctx).Info(fmt.Sprintf("Marshalling UndelegateCallback args: %v", undelegateCallback))
	marshalledCallbackArgs, err = k.MarshalUndelegateCallbackArgs(ctx, undelegateCallback)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return nil, sdk.ZeroInt(), nil, nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, err.Error())
	}

	return msgs, totalAmtToUnbond, marshalledCallbackArgs, epochUnbondingRecordIds, nil
}

func (k Keeper) SubmitHostZoneUnbondingMsg(ctx sdk.Context, msgs []sdk.Msg, totalAmtToUnbond sdk.Int, marshalledCallbackArgs []byte, hostZone types.HostZone) error {
	delegationAccount := hostZone.GetDelegationAccount()

	// safety check: if msgs is nil, error
	if msgs == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "no msgs to submit for host zone unbondings")
	}

	_, err := k.SubmitTxsDayEpoch(ctx, hostZone.GetConnectionId(), msgs, *delegationAccount, ICACallbackID_Undelegate, marshalledCallbackArgs)
	if err != nil {
		errMsg := fmt.Sprintf("Error submitting unbonding tx: %s", err)
		k.Logger(ctx).Error(errMsg)
		return sdkerrors.Wrap(sdkerrors.ErrNotFound, errMsg)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute("hostZone", hostZone.ChainId),
			sdk.NewAttribute("newAmountUnbonding", totalAmtToUnbond.String()),
		),
	)

	return nil
}

func (k Keeper) InitiateAllHostZoneUnbondings(ctx sdk.Context, dayNumber uint64) (success bool, successfulUnbondings []string, failedUnbondings []string) {
	// this function goes through each host zone, and if it's the right time to
	// initiate an unbonding, it goes and tries to unbond all outstanding records
	// outputs (1) did all chains succeed
	//		   (2) list of strings of successful unbondings
	//		   (3) list of strings of failed unbondings
	success = true
	successfulUnbondings = []string{}
	failedUnbondings = []string{}
	for _, hostZone := range k.GetAllHostZone(ctx) {
		k.Logger(ctx).Info(fmt.Sprintf("Processing epoch unbondings for host zone %s", hostZone.GetChainId()))
		// we only send the ICA call if this hostZone is supposed to be triggered
		if dayNumber%hostZone.UnbondingFrequency == 0 {
			k.Logger(ctx).Info(fmt.Sprintf("Sending unbondings for host zone %s", hostZone.ChainId))
			msgs, totalAmtToUnbond, marshalledCallbackArgs, epochUnbondingRecordIds, err := k.GetHostZoneUnbondingMsgs(ctx, hostZone)
			if err != nil {
				errMsg := fmt.Sprintf("Error getting unbonding msgs for host zone %s: %s", hostZone.ChainId, err.Error())
				k.Logger(ctx).Error(errMsg)
				success = false
				failedUnbondings = append(failedUnbondings, hostZone.ChainId)
				continue
			}
			k.Logger(ctx).Info(fmt.Sprintf("Total unbonded amount for %s: %d%s", hostZone.ChainId, totalAmtToUnbond, hostZone.HostDenom))
			err = k.SubmitHostZoneUnbondingMsg(ctx, msgs, totalAmtToUnbond, marshalledCallbackArgs, hostZone)
			if err != nil {
				errMsg := fmt.Sprintf("Error submitting unbonding tx for host zone %s: %s", hostZone.ChainId, err.Error())
				k.Logger(ctx).Error(errMsg)
				success = false
				failedUnbondings = append(failedUnbondings, hostZone.ChainId)
				continue
			}
			err = k.RecordsKeeper.SetHostZoneUnbondings(ctx, hostZone.ChainId, epochUnbondingRecordIds, recordstypes.HostZoneUnbonding_UNBONDING_IN_PROGRESS)
			if err != nil {
				k.Logger(ctx).Error(err.Error())
				success = false
				continue
			}
			successfulUnbondings = append(successfulUnbondings, hostZone.ChainId)
		}
	}
	return success, successfulUnbondings, failedUnbondings
}

// Deletes any epoch unbonding records that have had all unbondings claimed
func (k Keeper) CleanupEpochUnbondingRecords(ctx sdk.Context, epochNumber uint64) bool {
	k.Logger(ctx).Info("Cleaning Claimed Epoch Unbonding Records...")

	for _, epochUnbondingRecord := range k.RecordsKeeper.GetAllEpochUnbondingRecord(ctx) {
		shouldDeleteEpochUnbondingRecord := true
		hostZoneUnbondings := epochUnbondingRecord.HostZoneUnbondings

		for _, hostZoneUnbonding := range hostZoneUnbondings {
			// if an EpochUnbondingRecord has any HostZoneUnbonding with non-zero balances, we don't delete the EpochUnbondingRecord
			// because it has outstanding tokens that need to be claimed
			if hostZoneUnbonding.NativeTokenAmount != sdk.ZeroInt() {
				shouldDeleteEpochUnbondingRecord = false
				break
			}
		}
		if shouldDeleteEpochUnbondingRecord {
			k.Logger(ctx).Info(fmt.Sprintf("  EpochUnbondingRecord %d - All unbondings claimed, removing record", epochUnbondingRecord.EpochNumber))
			k.RecordsKeeper.RemoveEpochUnbondingRecord(ctx, epochUnbondingRecord.EpochNumber)
		} else {
			k.Logger(ctx).Info(fmt.Sprintf("  EpochUnbondingRecord %d - Has unclaimed unbondings", epochUnbondingRecord.EpochNumber))
		}
	}

	return true
}


	
// Batch transfers any unbonded tokens from the delegation account to the redemption account
func (k Keeper) SweepAllUnbondedTokensForHostZone(ctx sdk.Context, hostZone types.HostZone, epochUnbondingRecords []recordstypes.EpochUnbondingRecord) (success bool, sweepAmount sdk.Int) {
	k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "Sweeping unbonded tokens"))

	// Sum up all host zone unbonding records that have finished unbonding
	totalAmtTransferToRedemptionAcct := sdk.ZeroInt()
	epochUnbondingRecordIds := []uint64{}
	for _, epochUnbondingRecord := range epochUnbondingRecords {

		// Get all the unbondings associated with the epoch + host zone pair
		hostZoneUnbonding, found := k.RecordsKeeper.GetHostZoneUnbondingByChainId(ctx, epochUnbondingRecord.EpochNumber, hostZone.ChainId)
		if !found {
			k.Logger(ctx).Error(fmt.Sprintf("Could not find host zone unbonding %d for host zone %s", epochUnbondingRecord.EpochNumber, hostZone.ChainId))
			continue
		}

		// Get latest blockTime from light client
		blockTime, err := k.GetLightClientTimeSafely(ctx, hostZone.ConnectionId)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("\tCould not find blockTime for host zone %s", hostZone.ChainId))
			continue
		}

		k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "Epoch %d - Status: %s, Amount: %d, Unbonding Time: %d, Block Time: %d",
			epochUnbondingRecord.EpochNumber, hostZoneUnbonding.Status.String(), hostZoneUnbonding.NativeTokenAmount, hostZoneUnbonding.UnbondingTime, blockTime))

		// If the unbonding period has elapsed, then we can send the ICA call to sweep this
		//   hostZone's unbondings to the redemption account (in a batch).
		// Verify:
		//      1. the unbonding time is set (g.t. 0)
		//      2. the unbonding time is less than the current block time
		//      3. the host zone is in the EXIT_TRANSFER_QUEUE state, meaning it's ready to be transferred
		inTransferQueue := hostZoneUnbonding.Status == recordstypes.HostZoneUnbonding_EXIT_TRANSFER_QUEUE
		validUnbondingTime := hostZoneUnbonding.UnbondingTime > 0 && hostZoneUnbonding.UnbondingTime < blockTime
		if inTransferQueue && validUnbondingTime {
			k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "  %d Tokens included in sweep", hostZoneUnbonding.NativeTokenAmount))

			if err != nil {
				errMsg := fmt.Sprintf("Could not convert native token amount to int64 | %s", err.Error())
				k.Logger(ctx).Error(errMsg)
				continue
			}
			totalAmtTransferToRedemptionAcct = totalAmtTransferToRedemptionAcct.Add(hostZoneUnbonding.NativeTokenAmount)
			epochUnbondingRecordIds = append(epochUnbondingRecordIds, epochUnbondingRecord.EpochNumber)
		}
	}

	// If we have any amount to sweep, then we can send the ICA call to sweep them
	if totalAmtTransferToRedemptionAcct.LTE(sdk.ZeroInt()) {
		k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "Record not ready for sweep"))
		return true, totalAmtTransferToRedemptionAcct
	}
	k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "Batch transferring %d to host zone", totalAmtTransferToRedemptionAcct))

	// Get the delegation account and redemption account
	delegationAccount := hostZone.DelegationAccount
	if delegationAccount == nil || delegationAccount.Address == "" {
		k.Logger(ctx).Error(fmt.Sprintf("Zone %s is missing a delegation address!", hostZone.ChainId))
		return false, sdk.ZeroInt()
	}
	redemptionAccount := hostZone.RedemptionAccount
	if redemptionAccount == nil || redemptionAccount.Address == "" {
		k.Logger(ctx).Error(fmt.Sprintf("Zone %s is missing a redemption address!", hostZone.ChainId))
		return false, sdk.ZeroInt()
	}

	// Build transfer message to transfer from the delegation account to redemption account
	sweepCoin := sdk.NewCoin(hostZone.HostDenom, totalAmtTransferToRedemptionAcct)
	msgs := []sdk.Msg{
		&banktypes.MsgSend{
			FromAddress: delegationAccount.Address,
			ToAddress:   redemptionAccount.Address,
			Amount:      sdk.NewCoins(sweepCoin),
		},
	}
	k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "Preparing MsgSend from Delegation Account to Redemption Account"))

	// Store the epoch numbers in the callback to identify the epoch unbonding records
	redemptionCallback := types.RedemptionCallback{
		HostZoneId:              hostZone.ChainId,
		EpochUnbondingRecordIds: epochUnbondingRecordIds,
	}
	marshalledCallbackArgs, err := k.MarshalRedemptionCallbackArgs(ctx, redemptionCallback)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return false, sdk.ZeroInt()
	}

	// Send the transfer ICA
	_, err = k.SubmitTxsDayEpoch(ctx, hostZone.ConnectionId, msgs, *delegationAccount, ICACallbackID_Redemption, marshalledCallbackArgs)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("Failed to SubmitTxs, transfer to redemption account on %s", hostZone.ChainId))
		return false, sdk.ZeroInt()
	}
	k.Logger(ctx).Info(utils.LogWithHostZone(hostZone.ChainId, "ICA MsgSend Successfully Sent"))

	// Update the host zone unbonding records to status IN_PROGRESS
	err = k.RecordsKeeper.SetHostZoneUnbondings(ctx, hostZone.ChainId, epochUnbondingRecordIds, recordstypes.HostZoneUnbonding_EXIT_TRANSFER_IN_PROGRESS)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		return false, sdk.ZeroInt()
	}

	return true, totalAmtTransferToRedemptionAcct
}

// Sends all unbonded tokens to the redemption account
//    returns:
//       * success indicator if all chains succeeded
//       * list of successful chains
//       * list of tokens swept
//       * list of failed chains
func (k Keeper) SweepAllUnbondedTokens(ctx sdk.Context) (success bool, successfulSweeps []string, sweepAmounts []sdk.Int, failedSweeps []string) {
	// this function returns true if all chains succeeded, false otherwise
	// it also returns a list of successful chains (arg 2), tokens swept (arg 3), and failed chains (arg 4)
	k.Logger(ctx).Info("Sweeping All Unbonded Tokens...")

	success = true
	successfulSweeps = []string{}
	sweepAmounts = []sdk.Int{}
	failedSweeps = []string{}
	hostZones := k.GetAllHostZone(ctx)

	epochUnbondingRecords := k.RecordsKeeper.GetAllEpochUnbondingRecord(ctx)
	for _, hostZone := range hostZones {
		hostZoneSuccess, sweepAmount := k.SweepAllUnbondedTokensForHostZone(ctx, hostZone, epochUnbondingRecords)
		if hostZoneSuccess {
			successfulSweeps = append(successfulSweeps, hostZone.ChainId)
			sweepAmounts = append(sweepAmounts, sweepAmount)
		} else {
			success = false
			failedSweeps = append(failedSweeps, hostZone.ChainId)
		}
	}

	return success, successfulSweeps, sweepAmounts, failedSweeps
}
