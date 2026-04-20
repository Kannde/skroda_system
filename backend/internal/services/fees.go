package services

type FeeSchedule struct {
	FixedFee float64
	Tiers    []FeeTier
}

type FeeTier struct {
	MaxAmount  float64
	Percentage float64
}

var DefaultFeeSchedule = FeeSchedule{
	FixedFee: 5.00,
	Tiers: []FeeTier{
		{MaxAmount: 500.00, Percentage: 5.0},
		{MaxAmount: 5000.00, Percentage: 3.5},
		{MaxAmount: 0, Percentage: 2.0},
	},
}

func CalculateFee(amount float64, schedule FeeSchedule) float64 {
	var percentage float64
	for _, tier := range schedule.Tiers {
		if tier.MaxAmount == 0 || amount <= tier.MaxAmount {
			percentage = tier.Percentage
			break
		}
	}
	return (amount*percentage/100) + schedule.FixedFee
}

func CalculateSellerPayout(amount, feeAmount float64) float64 {
	return amount - feeAmount
}
