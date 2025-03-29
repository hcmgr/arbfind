## SET THESE
outcomes = [1.29, 5]
totalStake = 10

## CALCULATE THESE
r = sum(1/o for o in outcomes)

# Unrounded stakes
stakes = [(totalStake / o) * (1/r) for o in outcomes]

# Rounded stakes (real bets)
roundedStakes = [round(s, 2) for s in stakes]

# Actual returns per outcome
returns = [roundedStakes[i] * outcomes[i] for i in range(len(outcomes))]

# Find worst-case return
worstReturn = min(returns)

# Actual profit after rounding
actualProfit = worstReturn - sum(roundedStakes)

print("Odds=", outcomes)
print("r =", r)
print("Total stake (S) =", totalStake)
print("Original stakes =", stakes)
print("Rounded stakes =", roundedStakes)
print("Returns =", returns)
print("Worst-case return =", worstReturn)
print("Actual profit =", round(actualProfit, 2))