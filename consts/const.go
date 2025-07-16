package consts

type UserAction int

const (
	UserActionCheckin  UserAction = 1
	UserActionCheckout UserAction = 2
)

func (a UserAction) ToText() string {
	switch a {
	case UserActionCheckin:
		return "checkin"
	case UserActionCheckout:
		return "checkout"
	default:
		return "unknown"
	}
}

type CardAction int

const (
	CardActionTopup  CardAction = 1
	CardActionPay    CardAction = 2
	CardActionRefund CardAction = 3
)

func (a CardAction) ToText() string {
	switch a {
	case CardActionTopup:
		return "topup"
	case CardActionPay:
		return "pay"
	case CardActionRefund:
		return "refund"
	default:
		return "unknown"
	}
}
