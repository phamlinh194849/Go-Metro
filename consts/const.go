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

type Role int

const (
  AdminRole Role = 1
  StaffRole Role = 2
  UserRole  Role = 3
)

func (r Role) ToText() string {
  switch r {
  case UserRole:
    return "user"
  case AdminRole:
    return "admin"
  case StaffRole:
    return "staff"
  default:
    return "unknown"
  }
}

type CardType int

const (
  StudentCard CardType = 1
  NormalCard  CardType = 2
  VipCard     CardType = 3
)

func (c CardType) ToPrice() float64 {
  switch c {
  case StudentCard:
    return 3000
  case NormalCard:
    return 10000
  case VipCard:
    return 20000
  default:
    return 0
  }
}

func (c CardType) ToDefaultBlance() float64 {
  switch c {
  case StudentCard:
    return 10000
  case NormalCard:
    return 0
  case VipCard:
    return 3000
  default:
    return 0
  }
}

func (c CardType) ToText() string {
  switch c {
  case StudentCard:
    return "student"
  case NormalCard:
    return "normal"
  case VipCard:
    return "vip"
  default:
    return ""
  }
}

type Status string

const (
  ActiveStatus   Status = "active"
  InactiveStatus Status = "inactive"
  BlockedStatus  Status = "blocked"
)
