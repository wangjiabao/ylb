package biz

import (
	"context"
	"crypto/md5"
	v1 "dhb/app/app/api"
	"dhb/app/app/internal/pkg/middleware/auth"
	"encoding/base64"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID        int64
	Address   string
	CreatedAt time.Time
}

type Admin struct {
	ID       int64
	Password string
	Account  string
	Type     string
}

type AdminAuth struct {
	ID      int64
	AdminId int64
	AuthId  int64
}

type Auth struct {
	ID   int64
	Name string
	Path string
	Url  string
}

type UserInfo struct {
	ID               int64
	UserId           int64
	Vip              int64
	HistoryRecommend int64
}

type UserRecommend struct {
	ID            int64
	UserId        int64
	RecommendCode string
	CreatedAt     time.Time
}

type UserCurrentMonthRecommend struct {
	ID              int64
	UserId          int64
	RecommendUserId int64
	Date            time.Time
}

type Config struct {
	ID      int64
	KeyName string
	Name    string
	Value   string
}

type UserBalance struct {
	ID          int64
	UserId      int64
	BalanceUsdt int64
	BalanceDhb  int64
}

type Withdraw struct {
	ID              int64
	UserId          int64
	Amount          int64
	RelAmount       int64
	BalanceRecordId int64
	Status          string
	Type            string
	CreatedAt       time.Time
}

type UserUseCase struct {
	repo                          UserRepo
	urRepo                        UserRecommendRepo
	configRepo                    ConfigRepo
	uiRepo                        UserInfoRepo
	ubRepo                        UserBalanceRepo
	locationRepo                  LocationRepo
	userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo
	tx                            Transaction
	log                           *log.Helper
}

type Reward struct {
	ID               int64
	UserId           int64
	Amount           int64
	BalanceRecordId  int64
	Type             string
	TypeRecordId     int64
	Reason           string
	ReasonLocationId int64
	LocationType     string
	CreatedAt        time.Time
}

type Pagination struct {
	PageNum  int
	PageSize int
}

type UserSortRecommendReward struct {
	UserId int64
	Total  int64
}

type ConfigRepo interface {
	GetConfigByKeys(ctx context.Context, keys ...string) ([]*Config, error)
	GetConfigs(ctx context.Context) ([]*Config, error)
	UpdateConfig(ctx context.Context, id int64, value string) (bool, error)
}

type UserBalanceRepo interface {
	CreateUserBalance(ctx context.Context, u *User) (*UserBalance, error)
	LocationReward(ctx context.Context, userId int64, amount int64, locationId int64, myLocationId int64, locationType string) (int64, error)
	WithdrawReward(ctx context.Context, userId int64, amount int64, locationId int64, myLocationId int64, locationType string) (int64, error)
	RecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	RecommendTopReward(ctx context.Context, userId int64, amount int64, locationId int64, vip int64) (int64, error)
	SystemWithdrawReward(ctx context.Context, amount int64, locationId int64) error
	SystemReward(ctx context.Context, amount int64, locationId int64) error
	SystemDailyReward(ctx context.Context, amount int64, locationId int64) error
	GetSystemYesterdayDailyReward(ctx context.Context) (*Reward, error)
	SystemFee(ctx context.Context, amount int64, locationId int64) error
	UserFee(ctx context.Context, userId int64, amount int64) (int64, error)
	UserDailyFee(ctx context.Context, userId int64, amount int64) (int64, error)
	RecommendWithdrawReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	RecommendWithdrawTopReward(ctx context.Context, userId int64, amount int64, locationId int64, vip int64) (int64, error)
	NormalRecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	NormalRecommendTopReward(ctx context.Context, userId int64, amount int64, locationId int64, reasonId int64) (int64, error)
	NormalWithdrawRecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	NormalWithdrawRecommendTopReward(ctx context.Context, userId int64, amount int64, locationId int64, reasonId int64) (int64, error)
	Deposit(ctx context.Context, userId int64, amount int64, dhbAmount int64) (int64, error)
	DepositLast(ctx context.Context, userId int64, lastAmount int64, locationId int64) (int64, error)
	DepositDhb(ctx context.Context, userId int64, amount int64) (int64, error)
	GetUserBalance(ctx context.Context, userId int64) (*UserBalance, error)
	GetUserRewardByUserId(ctx context.Context, userId int64) ([]*Reward, error)
	GetUserRewards(ctx context.Context, b *Pagination, userId int64) ([]*Reward, error, int64)
	GetUserRewardsLastMonthFee(ctx context.Context) ([]*Reward, error)
	GetUserBalanceByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserBalance, error)
	GetUserBalanceUsdtTotal(ctx context.Context) (int64, error)
	GreateWithdraw(ctx context.Context, userId int64, amount int64, coinType string) (*Withdraw, error)
	WithdrawUsdt(ctx context.Context, userId int64, amount int64) error
	WithdrawDhb(ctx context.Context, userId int64, amount int64) error
	GetWithdrawByUserId(ctx context.Context, userId int64) ([]*Withdraw, error)
	GetWithdraws(ctx context.Context, b *Pagination, userId int64) ([]*Withdraw, error, int64)
	GetWithdrawPassOrRewarded(ctx context.Context) ([]*Withdraw, error)
	UpdateWithdraw(ctx context.Context, id int64, status string) (*Withdraw, error)
	GetWithdrawById(ctx context.Context, id int64) (*Withdraw, error)
	GetWithdrawNotDeal(ctx context.Context) ([]*Withdraw, error)
	GetUserBalanceRecordUsdtTotal(ctx context.Context) (int64, error)
	GetUserBalanceRecordUsdtTotalToday(ctx context.Context) (int64, error)
	GetUserWithdrawUsdtTotalToday(ctx context.Context) (int64, error)
	GetUserWithdrawUsdtTotal(ctx context.Context) (int64, error)
	GetUserRewardUsdtTotal(ctx context.Context) (int64, error)
	GetSystemRewardUsdtTotal(ctx context.Context) (int64, error)
	UpdateWithdrawAmount(ctx context.Context, id int64, status string, amount int64) (*Withdraw, error)
	GetUserRewardRecommendSort(ctx context.Context) ([]*UserSortRecommendReward, error)
	UpdateBalance(ctx context.Context, userId int64, amount int64) (bool, error)
}

type UserRecommendRepo interface {
	GetUserRecommendByUserId(ctx context.Context, userId int64) (*UserRecommend, error)
	CreateUserRecommend(ctx context.Context, u *User, recommendUser *UserRecommend) (*UserRecommend, error)
	GetUserRecommendByCode(ctx context.Context, code string) ([]*UserRecommend, error)
	GetUserRecommendLikeCode(ctx context.Context, code string) ([]*UserRecommend, error)
}

type UserCurrentMonthRecommendRepo interface {
	GetUserCurrentMonthRecommendByUserId(ctx context.Context, userId int64) ([]*UserCurrentMonthRecommend, error)
	GetUserCurrentMonthRecommendGroupByUserId(ctx context.Context, b *Pagination, userId int64) ([]*UserCurrentMonthRecommend, error, int64)
	CreateUserCurrentMonthRecommend(ctx context.Context, u *UserCurrentMonthRecommend) (*UserCurrentMonthRecommend, error)
	GetUserCurrentMonthRecommendCountByUserIds(ctx context.Context, userIds ...int64) (map[int64]int64, error)
	GetUserLastMonthRecommend(ctx context.Context) ([]int64, error)
}

type UserInfoRepo interface {
	CreateUserInfo(ctx context.Context, u *User) (*UserInfo, error)
	GetUserInfoByUserId(ctx context.Context, userId int64) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, u *UserInfo) (*UserInfo, error)
	GetUserInfoByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserInfo, error)
}

type UserRepo interface {
	GetUserById(ctx context.Context, Id int64) (*User, error)
	GetAdminByAccount(ctx context.Context, account string, password string) (*Admin, error)
	GetAdminById(ctx context.Context, id int64) (*Admin, error)
	GetUserByAddresses(ctx context.Context, Addresses ...string) (map[string]*User, error)
	GetUserByAddress(ctx context.Context, address string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	CreateAdmin(ctx context.Context, a *Admin) (*Admin, error)
	GetUserByUserIds(ctx context.Context, userIds ...int64) (map[int64]*User, error)
	GetAdmins(ctx context.Context) ([]*Admin, error)
	GetUsers(ctx context.Context, b *Pagination, address string, isLocation bool, vip int64) ([]*User, error, int64)
	GetUserCount(ctx context.Context) (int64, error)
	GetUserCountToday(ctx context.Context) (int64, error)
	CreateAdminAuth(ctx context.Context, adminId int64, authId int64) (bool, error)
	DeleteAdminAuth(ctx context.Context, adminId int64, authId int64) (bool, error)
	GetAuths(ctx context.Context) ([]*Auth, error)
	GetAuthByIds(ctx context.Context, ids ...int64) (map[int64]*Auth, error)
	GetAdminAuth(ctx context.Context, adminId int64) ([]*AdminAuth, error)
	UpdateAdminPassword(ctx context.Context, account string, password string) (*Admin, error)
}

func NewUserUseCase(repo UserRepo, tx Transaction, configRepo ConfigRepo, uiRepo UserInfoRepo, urRepo UserRecommendRepo, locationRepo LocationRepo, userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo, ubRepo UserBalanceRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:                          repo,
		tx:                            tx,
		configRepo:                    configRepo,
		locationRepo:                  locationRepo,
		userCurrentMonthRecommendRepo: userCurrentMonthRecommendRepo,
		uiRepo:                        uiRepo,
		urRepo:                        urRepo,
		ubRepo:                        ubRepo,
		log:                           log.NewHelper(logger),
	}
}

func (uuc *UserUseCase) GetUserByAddress(ctx context.Context, Addresses ...string) (map[string]*User, error) {
	return uuc.repo.GetUserByAddresses(ctx, Addresses...)
}

func (uuc *UserUseCase) GetDhbConfig(ctx context.Context) ([]*Config, error) {
	return uuc.configRepo.GetConfigByKeys(ctx, "level1Dhb", "level2Dhb", "level3Dhb")
}

func (uuc *UserUseCase) GetExistUserByAddressOrCreate(ctx context.Context, u *User, req *v1.EthAuthorizeRequest) (*User, error) {
	var (
		user          *User
		recommendUser *UserRecommend
		userRecommend *UserRecommend
		userInfo      *UserInfo
		userBalance   *UserBalance
		err           error
		userId        int64
		decodeBytes   []byte
	)

	user, err = uuc.repo.GetUserByAddress(ctx, u.Address) // 查询用户
	if nil == user || nil != err {
		code := req.SendBody.Code // 查询推荐码 abf00dd52c08a9213f225827bc3fb100 md5 dhbmachinefirst
		if "abf00dd52c08a9213f225827bc3fb100" != code {
			decodeBytes, err = base64.StdEncoding.DecodeString(code)
			code = string(decodeBytes)
			if 1 >= len(code) {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}
			if userId, err = strconv.ParseInt(code[1:], 10, 64); 0 >= userId || nil != err {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}

			// 查询推荐人的相关信息
			recommendUser, err = uuc.urRepo.GetUserRecommendByUserId(ctx, userId)
			if err != nil {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}
		}

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			user, err = uuc.repo.CreateUser(ctx, u) // 用户创建
			if err != nil {
				return err
			}

			userInfo, err = uuc.uiRepo.CreateUserInfo(ctx, user) // 创建用户信息
			if err != nil {
				return err
			}

			userRecommend, err = uuc.urRepo.CreateUserRecommend(ctx, user, recommendUser) // 创建用户信息
			if err != nil {
				return err
			}

			userBalance, err = uuc.ubRepo.CreateUserBalance(ctx, user) // 创建余额信息
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (uuc *UserUseCase) UserInfo(ctx context.Context, user *User) (*v1.UserInfoReply, error) {
	var (
		myUser                     *User
		userInfo                   *UserInfo
		locations                  []*Location
		userBalance                *UserBalance
		userRecommend              *UserRecommend
		userRecommends             []*UserRecommend
		userRewards                []*Reward
		rewardLocations            []*Location
		userCurrentMonthRecommends []*UserCurrentMonthRecommend
		userRewardTotal            int64
		encodeString               string
		myUserRecommendUserId      int64
		myRecommendUser            *User
		myRow                      int64
		rowNum                     int64
		colNum                     int64
		myCol                      int64
		recommendTeamNum           int64
		recommendTotal             int64
		locationTotal              int64
		feeTotal                   int64
		myCode                     string
		inviteUserAddress          string
		amount                     string
		status                     = "no"
		currentMonthRecommendNum   int64
		configs                    []*Config
		myLastStopLocation         *Location
		myLastLocationCurrent      int64
		hasRunningLocation         bool
		level1Dhb                  string
		level2Dhb                  string
		level3Dhb                  string
		userCount                  string
		err                        error
	)

	myUser, err = uuc.repo.GetUserById(ctx, user.ID)
	if nil != err {
		return nil, err
	}

	userInfo, err = uuc.uiRepo.GetUserInfoByUserId(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}

	locations, err = uuc.locationRepo.GetLocationsByUserId(ctx, myUser.ID)
	if nil != locations && 0 < len(locations) {
		status = "stop"
		for _, v := range locations {
			if "running" == v.Status {
				status = "running"
				if 0 == v.Current {
					status = "yes"
				}
				hasRunningLocation = true
				amount = fmt.Sprintf("%.2f", float64(v.CurrentMax-v.Current)/float64(10000000000))
				myCol = v.Col
				myRow = v.Row
				break
			}
		}
	}

	now := time.Now().UTC().Add(8 * time.Hour)
	myLastStopLocation, err = uuc.locationRepo.GetMyStopLocationLast(ctx, myUser.ID) // 冻结
	if nil != myLastStopLocation && now.Before(myLastStopLocation.StopDate.Add(24*time.Hour)) && !hasRunningLocation {
		myLastLocationCurrent = myLastStopLocation.Current - myLastStopLocation.CurrentMax // 补上
	}

	userBalance, err = uuc.ubRepo.GetUserBalance(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}

	userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, myUser.ID)
	if nil == userRecommend {
		return nil, err
	}

	myCode = "D" + strconv.FormatInt(myUser.ID, 10)
	codeByte := []byte(myCode)
	encodeString = base64.StdEncoding.EncodeToString(codeByte)

	if "" != userRecommend.RecommendCode {
		tmpRecommendUserIds := strings.Split(userRecommend.RecommendCode, "D")
		if 2 <= len(tmpRecommendUserIds) {
			myUserRecommendUserId, _ = strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-1], 10, 64) // 最后一位是直推人
		}
		myRecommendUser, err = uuc.repo.GetUserById(ctx, myUserRecommendUserId)
		if nil != err {
			return nil, err
		}
		inviteUserAddress = myRecommendUser.Address
		myCode = userRecommend.RecommendCode + myCode
	}

	// 团队
	userRecommends, err = uuc.urRepo.GetUserRecommendLikeCode(ctx, myCode)
	if nil != userRecommends {
		recommendTeamNum = int64(len(userRecommends))
	}

	// 累计奖励
	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, myUser.ID)
	if nil != userRewards {
		for _, vUserReward := range userRewards {
			userRewardTotal += vUserReward.Amount
			if "recommend" == vUserReward.Reason || "recommend_vip" == vUserReward.Reason {
				recommendTotal += vUserReward.Amount
			} else if "location" == vUserReward.Reason {
				locationTotal += vUserReward.Amount
			} else if "fee" == vUserReward.Reason {
				feeTotal += vUserReward.Amount
			}
		}
	}

	// 位置
	if 0 < myRow && 0 < myCol {
		rewardLocations, err = uuc.locationRepo.GetRewardLocationByRowOrCol(ctx, myRow, myCol)
		if nil != rewardLocations {
			for _, vRewardLocation := range rewardLocations {
				if myRow == vRewardLocation.Row && myCol == vRewardLocation.Col { // 跳过自己
					continue
				}
				if myRow == vRewardLocation.Row {
					colNum++
				}
				if myCol == vRewardLocation.Col {
					rowNum++
				}
			}
		}
	}

	// 当月推荐人数
	userCurrentMonthRecommends, err = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendByUserId(ctx, myUser.ID)
	if nil != userCurrentMonthRecommends {
		for _, vUserCurrentMonthRecommend := range userCurrentMonthRecommends {
			if vUserCurrentMonthRecommend.Date.After(time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)) {
				currentMonthRecommendNum++
			}
		}
	}

	// 配置
	configs, err = uuc.configRepo.GetConfigByKeys(ctx, "user_count")
	if nil != configs {
		for _, vConfig := range configs {
			if "user_count" == vConfig.KeyName {
				userCount = vConfig.Value
			}
		}
	}

	return &v1.UserInfoReply{
		Address:                  myUser.Address,
		Level:                    userInfo.Vip,
		Status:                   status,
		Amount:                   amount,
		BalanceUsdt:              fmt.Sprintf("%.2f", float64(userBalance.BalanceUsdt)/float64(10000000000)),
		BalanceDhb:               fmt.Sprintf("%.2f", float64(userBalance.BalanceDhb)/float64(10000000000)),
		InviteUrl:                encodeString,
		InviteUserAddress:        inviteUserAddress,
		RecommendNum:             userInfo.HistoryRecommend,
		RecommendTeamNum:         recommendTeamNum,
		Total:                    fmt.Sprintf("%.2f", float64(userRewardTotal)/float64(10000000000)),
		Row:                      rowNum,
		Col:                      colNum,
		CurrentMonthRecommendNum: currentMonthRecommendNum,
		RecommendTotal:           fmt.Sprintf("%.2f", float64(recommendTotal)/float64(10000000000)),
		FeeTotal:                 fmt.Sprintf("%.2f", float64(feeTotal)/float64(10000000000)),
		LocationTotal:            fmt.Sprintf("%.2f", float64(locationTotal)/float64(10000000000)),
		Level1Dhb:                level1Dhb,
		Level2Dhb:                level2Dhb,
		Level3Dhb:                level3Dhb,
		Usdt:                     "0x55d398326f99059fF775485246999027B3197955",
		Dhb:                      "0xb7864be857e00796e6f79e057b3ef1032cbe4a06",
		Account:                  "0x636F2deAAb4C9A8F3c808D23F16f456009C4e9Fd",
		//Usdt:                     "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd",
		//Dhb:                      "0x96BD81715c69eE013405B4005Ba97eA1f420fd87",
		//Account:                  "0xe865f2e5ff04b8b7952d1c0d9163a91f313b158f",
		AmountB:   fmt.Sprintf("%.2f", float64(myLastLocationCurrent)/float64(10000000000)),
		UserCount: userCount,
	}, nil
}

func (uuc *UserUseCase) RewardList(ctx context.Context, req *v1.RewardListRequest, user *User) (*v1.RewardListReply, error) {
	var (
		userRewards    []*Reward
		locationIdsMap map[int64]int64
		locations      map[int64]*Location
		err            error
	)
	res := &v1.RewardListReply{
		Rewards: make([]*v1.RewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	locationIdsMap = make(map[int64]int64, 0)
	if nil != userRewards {
		for _, vUserReward := range userRewards {
			if "location" == vUserReward.Reason && req.Type == vUserReward.LocationType && 1 <= vUserReward.ReasonLocationId {
				locationIdsMap[vUserReward.ReasonLocationId] = vUserReward.ReasonLocationId
			}
		}

		var tmpLocationIds []int64
		for _, v := range locationIdsMap {
			tmpLocationIds = append(tmpLocationIds, v)
		}
		if 0 >= len(tmpLocationIds) {
			return res, nil
		}

		locations, err = uuc.locationRepo.GetRewardLocationByIds(ctx, tmpLocationIds...)

		for _, vUserReward := range userRewards {
			if "location" == vUserReward.Reason && req.Type == vUserReward.LocationType {
				if _, ok := locations[vUserReward.ReasonLocationId]; !ok {
					continue
				}

				res.Rewards = append(res.Rewards, &v1.RewardListReply_List{
					CreatedAt:      vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
					Amount:         fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
					LocationStatus: locations[vUserReward.ReasonLocationId].Status,
					Type:           vUserReward.Type,
				})
			}
		}
	}

	return res, nil
}

func (uuc *UserUseCase) RecommendRewardList(ctx context.Context, user *User) (*v1.RecommendRewardListReply, error) {
	var (
		userRewards []*Reward
		err         error
	)
	res := &v1.RecommendRewardListReply{
		Rewards: make([]*v1.RecommendRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if "recommend" == vUserReward.Reason || "recommend_vip" == vUserReward.Reason {
			res.Rewards = append(res.Rewards, &v1.RecommendRewardListReply_List{
				CreatedAt: vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
				Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
				Type:      vUserReward.Type,
				Reason:    vUserReward.Reason,
			})
		}
	}

	return res, nil
}

func (uuc *UserUseCase) FeeRewardList(ctx context.Context, user *User) (*v1.FeeRewardListReply, error) {
	var (
		userRewards []*Reward
		err         error
	)
	res := &v1.FeeRewardListReply{
		Rewards: make([]*v1.FeeRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if "fee" == vUserReward.Reason {
			res.Rewards = append(res.Rewards, &v1.FeeRewardListReply_List{
				CreatedAt: vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
				Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
			})
		}
	}

	return res, nil
}

func (uuc *UserUseCase) WithdrawList(ctx context.Context, user *User) (*v1.WithdrawListReply, error) {

	var (
		withdraws []*Withdraw
		err       error
	)

	res := &v1.WithdrawListReply{
		Withdraw: make([]*v1.WithdrawListReply_List, 0),
	}

	withdraws, err = uuc.ubRepo.GetWithdrawByUserId(ctx, user.ID)
	if nil != err {
		return res, err
	}

	for _, v := range withdraws {
		res.Withdraw = append(res.Withdraw, &v1.WithdrawListReply_List{
			CreatedAt: v.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(v.Amount)/float64(10000000000)),
			Status:    v.Status,
			Type:      v.Type,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) Withdraw(ctx context.Context, req *v1.WithdrawRequest, user *User) (*v1.WithdrawReply, error) {
	var (
		err         error
		userBalance *UserBalance
	)

	if "dhb" != req.SendBody.Type && "usdt" != req.SendBody.Type {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	userBalance, err = uuc.ubRepo.GetUserBalance(ctx, user.ID)
	if nil != err {
		return nil, err
	}

	amountFloat, _ := strconv.ParseFloat(req.SendBody.Amount, 10)
	amountFloat *= 10000000000
	amount, _ := strconv.ParseInt(strconv.FormatFloat(amountFloat, 'f', -1, 64), 10, 64)
	if 0 >= amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	if "dhb" == req.SendBody.Type && userBalance.BalanceDhb < amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	if "usdt" == req.SendBody.Type && userBalance.BalanceUsdt < amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}
	if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务

		if "usdt" == req.SendBody.Type {
			err = uuc.ubRepo.WithdrawUsdt(ctx, user.ID, amount) // 提现
			if nil != err {
				return err
			}
			_, err = uuc.ubRepo.GreateWithdraw(ctx, user.ID, amount, req.SendBody.Type)
			if nil != err {
				return err
			}

		} else if "dhb" == req.SendBody.Type {
			err = uuc.ubRepo.WithdrawDhb(ctx, user.ID, amount) // 提现
			if nil != err {
				return err
			}
			_, err = uuc.ubRepo.GreateWithdraw(ctx, user.ID, amount, req.SendBody.Type)
			if nil != err {
				return err
			}
		}

		return nil
	}); nil != err {
		return nil, err
	}

	return &v1.WithdrawReply{
		Status: "ok",
	}, nil
}

func (uuc *UserUseCase) AdminRewardList(ctx context.Context, req *v1.AdminRewardListRequest) (*v1.AdminRewardListReply, error) {
	var (
		userSearch  *User
		userId      int64 = 0
		userRewards []*Reward
		users       map[int64]*User
		userIdsMap  map[int64]int64
		userIds     []int64
		err         error
		count       int64
	)
	res := &v1.AdminRewardListReply{
		Rewards: make([]*v1.AdminRewardListReply_List, 0),
	}

	// 地址查询
	if "" != req.Address {
		userSearch, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil != err {
			return res, nil
		}
		userId = userSearch.ID
	}

	userRewards, err, count = uuc.ubRepo.GetUserRewards(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, userId)
	if nil != err {
		return res, nil
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vUserReward := range userRewards {
		userIdsMap[vUserReward.UserId] = vUserReward.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	for _, vUserReward := range userRewards {
		tmpUser := ""
		if nil != users {
			if _, ok := users[vUserReward.UserId]; ok {
				tmpUser = users[vUserReward.UserId].Address
			}
		}

		res.Rewards = append(res.Rewards, &v1.AdminRewardListReply_List{
			CreatedAt: vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
			Type:      vUserReward.Type,
			Address:   tmpUser,
			Reason:    vUserReward.Reason,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminUserList(ctx context.Context, req *v1.AdminUserListRequest) (*v1.AdminUserListReply, error) {
	var (
		users                          []*User
		userIds                        []int64
		userBalances                   map[int64]*UserBalance
		userInfos                      map[int64]*UserInfo
		userCurrentMonthRecommendCount map[int64]int64
		count                          int64
		err                            error
	)

	res := &v1.AdminUserListReply{
		Users: make([]*v1.AdminUserListReply_UserList, 0),
	}

	users, err, count = uuc.repo.GetUsers(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, req.Address, req.IsLocation, req.Vip)
	if nil != err {
		return res, nil
	}
	res.Count = count

	for _, vUsers := range users {
		userIds = append(userIds, vUsers.ID)
	}

	userBalances, err = uuc.ubRepo.GetUserBalanceByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	userInfos, err = uuc.uiRepo.GetUserInfoByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	userCurrentMonthRecommendCount, err = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendCountByUserIds(ctx, userIds...)

	for _, v := range users {
		if _, ok := userBalances[v.ID]; !ok {
			continue
		}
		if _, ok := userInfos[v.ID]; !ok {
			continue
		}

		var tmpCount int64
		if nil != userCurrentMonthRecommendCount {
			if _, ok := userCurrentMonthRecommendCount[v.ID]; ok {
				tmpCount = userCurrentMonthRecommendCount[v.ID]
			}
		}

		res.Users = append(res.Users, &v1.AdminUserListReply_UserList{
			UserId:           v.ID,
			CreatedAt:        v.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Address:          v.Address,
			BalanceUsdt:      fmt.Sprintf("%.2f", float64(userBalances[v.ID].BalanceUsdt)/float64(10000000000)),
			BalanceDhb:       fmt.Sprintf("%.2f", float64(userBalances[v.ID].BalanceDhb)/float64(10000000000)),
			Vip:              userInfos[v.ID].Vip,
			MonthRecommend:   tmpCount,
			HistoryRecommend: userInfos[v.ID].HistoryRecommend,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) GetUserByUserIds(ctx context.Context, userIds ...int64) (map[int64]*User, error) {
	return uuc.repo.GetUserByUserIds(ctx, userIds...)
}

func (uuc *UserUseCase) AdminLocationList(ctx context.Context, req *v1.AdminLocationListRequest) (*v1.AdminLocationListReply, error) {
	var (
		locations  []*Location
		userSearch *User
		userId     int64
		userIds    []int64
		userIdsMap map[int64]int64
		users      map[int64]*User
		count      int64
		err        error
	)

	res := &v1.AdminLocationListReply{
		Locations: make([]*v1.AdminLocationListReply_LocationList, 0),
	}

	// 地址查询
	if "" != req.Address {
		userSearch, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil != err {
			return res, nil
		}
		userId = userSearch.ID
	}

	locations, err, count = uuc.locationRepo.GetLocations(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, userId)
	if nil != err {
		return res, nil
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vLocations := range locations {
		userIdsMap[vLocations.UserId] = vLocations.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range locations {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Locations = append(res.Locations, &v1.AdminLocationListReply_LocationList{
			CreatedAt:    v.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Address:      users[v.UserId].Address,
			Row:          v.Row,
			Col:          v.Col,
			Status:       v.Status,
			CurrentLevel: v.CurrentLevel,
			Current:      fmt.Sprintf("%.2f", float64(v.Current)/float64(10000000000)),
			CurrentMax:   fmt.Sprintf("%.2f", float64(v.CurrentMax)/float64(10000000000)),
		})
	}

	return res, nil

}

func (uuc *UserUseCase) AdminRecommendList(ctx context.Context, req *v1.AdminUserRecommendRequest) (*v1.AdminUserRecommendReply, error) {
	var (
		userRecommends []*UserRecommend
		userRecommend  *UserRecommend
		userIdsMap     map[int64]int64
		userIds        []int64
		users          map[int64]*User
		err            error
	)

	res := &v1.AdminUserRecommendReply{
		Users: make([]*v1.AdminUserRecommendReply_List, 0),
	}

	// 地址查询
	if 0 < req.UserId {
		userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, req.UserId)
		if nil == userRecommend {
			return res, nil
		}

		userRecommends, err = uuc.urRepo.GetUserRecommendByCode(ctx, userRecommend.RecommendCode+"D"+strconv.FormatInt(userRecommend.UserId, 10))
		if nil != err {
			return res, nil
		}
	}

	userIdsMap = make(map[int64]int64, 0)
	for _, vLocations := range userRecommends {
		userIdsMap[vLocations.UserId] = vLocations.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range userRecommends {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Users = append(res.Users, &v1.AdminUserRecommendReply_List{
			Address:   users[v.UserId].Address,
			Id:        v.ID,
			UserId:    v.UserId,
			CreatedAt: v.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminMonthRecommend(ctx context.Context, req *v1.AdminMonthRecommendRequest) (*v1.AdminMonthRecommendReply, error) {
	var (
		userCurrentMonthRecommends []*UserCurrentMonthRecommend
		searchUser                 *User
		userIdsMap                 map[int64]int64
		userIds                    []int64
		searchUserId               int64
		users                      map[int64]*User
		count                      int64
		err                        error
	)

	res := &v1.AdminMonthRecommendReply{
		Users: make([]*v1.AdminMonthRecommendReply_List, 0),
	}

	// 地址查询
	if "" != req.Address {
		searchUser, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil == searchUser {
			return res, nil
		}
		searchUserId = searchUser.ID
	}

	userCurrentMonthRecommends, err, count = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendGroupByUserId(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, searchUserId)
	if nil != err {
		return res, nil
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vRecommend := range userCurrentMonthRecommends {
		userIdsMap[vRecommend.UserId] = vRecommend.UserId
		userIdsMap[vRecommend.RecommendUserId] = vRecommend.RecommendUserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range userCurrentMonthRecommends {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Users = append(res.Users, &v1.AdminMonthRecommendReply_List{
			Address:          users[v.UserId].Address,
			Id:               v.ID,
			RecommendAddress: users[v.RecommendUserId].Address,
			CreatedAt:        v.Date.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminConfig(ctx context.Context, req *v1.AdminConfigRequest) (*v1.AdminConfigReply, error) {
	var (
		configs []*Config
	)

	res := &v1.AdminConfigReply{
		Config: make([]*v1.AdminConfigReply_List, 0),
	}

	configs, _ = uuc.configRepo.GetConfigs(ctx)
	if nil == configs {
		return res, nil
	}

	for _, v := range configs {
		res.Config = append(res.Config, &v1.AdminConfigReply_List{
			Id:    v.ID,
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminConfigUpdate(ctx context.Context, req *v1.AdminConfigUpdateRequest) (*v1.AdminConfigUpdateReply, error) {
	var (
		err error
	)

	res := &v1.AdminConfigUpdateReply{}

	_, err = uuc.configRepo.UpdateConfig(ctx, req.SendBody.Id, req.SendBody.Value)
	if nil != err {
		return res, err
	}

	return res, nil
}

func (uuc *UserUseCase) AdminVipUpdate(ctx context.Context, req *v1.AdminVipUpdateRequest) (*v1.AdminVipUpdateReply, error) {
	var (
		userInfo *UserInfo
		err      error
	)

	userInfo, err = uuc.uiRepo.GetUserInfoByUserId(ctx, req.SendBody.UserId)
	if nil == userInfo {
		return &v1.AdminVipUpdateReply{}, nil
	}

	res := &v1.AdminVipUpdateReply{}

	if 5 == req.SendBody.Vip {
		userInfo.Vip = 5
		userInfo.HistoryRecommend = 10
	} else if 4 == req.SendBody.Vip {
		userInfo.Vip = 4
		userInfo.HistoryRecommend = 8
	} else if 3 == req.SendBody.Vip {
		userInfo.Vip = 3
		userInfo.HistoryRecommend = 6
	} else if 2 == req.SendBody.Vip {
		userInfo.Vip = 2
		userInfo.HistoryRecommend = 4
	} else if 1 == req.SendBody.Vip {
		userInfo.Vip = 1
		userInfo.HistoryRecommend = 2
	}

	_, err = uuc.uiRepo.UpdateUserInfo(ctx, userInfo) // 推荐人信息修改
	if nil != err {
		return res, err
	}

	return res, nil
}

func (uuc *UserUseCase) AdminBalanceUpdate(ctx context.Context, req *v1.AdminBalanceUpdateRequest) (*v1.AdminBalanceUpdateReply, error) {
	var (
		err error
	)
	res := &v1.AdminBalanceUpdateReply{}

	amountFloat, _ := strconv.ParseFloat(req.SendBody.Amount, 10)
	amountFloat *= 10000000000
	amount, _ := strconv.ParseInt(strconv.FormatFloat(amountFloat, 'f', -1, 64), 10, 64)

	_, err = uuc.ubRepo.UpdateBalance(ctx, req.SendBody.UserId, amount) // 推荐人信息修改
	if nil != err {
		return res, err
	}

	return res, nil
}

func (uuc *UserUseCase) AdminLogin(ctx context.Context, req *v1.AdminLoginRequest, ca string) (*v1.AdminLoginReply, error) {
	var (
		admin *Admin
		err   error
	)

	res := &v1.AdminLoginReply{}
	password := fmt.Sprintf("%x", md5.Sum([]byte(req.SendBody.Password)))
	fmt.Println(password)
	admin, err = uuc.repo.GetAdminByAccount(ctx, req.SendBody.Account, password)
	if nil != err {
		return res, err
	}

	claims := auth.CustomClaims{
		UserId:   admin.ID,
		UserType: "admin",
		StandardClaims: jwt2.StandardClaims{
			NotBefore: time.Now().Unix(),              // 签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 7天过期
			Issuer:    "DHB",
		},
	}
	token, err := auth.CreateToken(claims, ca)
	if err != nil {
		return nil, errors.New(500, "AUTHORIZE_ERROR", "生成token失败")
	}
	res.Token = token
	return res, nil
}

func (uuc *UserUseCase) AdminCreateAccount(ctx context.Context, req *v1.AdminCreateAccountRequest) (*v1.AdminCreateAccountReply, error) {
	var (
		admin    *Admin
		myAdmin  *Admin
		newAdmin *Admin
		err      error
	)

	res := &v1.AdminCreateAccountReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" != myAdmin.Type {
		return nil, errors.New(500, "ERROR_TOKEN", "非超管")
	}

	password := fmt.Sprintf("%x", md5.Sum([]byte(req.SendBody.Password)))
	admin, err = uuc.repo.GetAdminByAccount(ctx, req.SendBody.Account, password)
	if nil != admin {
		return nil, errors.New(500, "ERROR_TOKEN", "已存在账户")
	}

	newAdmin, err = uuc.repo.CreateAdmin(ctx, &Admin{
		Password: password,
		Account:  req.SendBody.Account,
	})

	if nil != newAdmin {
		return res, err
	}

	return res, nil
}

func (uuc *UserUseCase) AdminList(ctx context.Context, req *v1.AdminListRequest) (*v1.AdminListReply, error) {
	var (
		admins []*Admin
	)

	res := &v1.AdminListReply{Account: make([]*v1.AdminListReply_List, 0)}

	admins, _ = uuc.repo.GetAdmins(ctx)
	if nil == admins {
		return res, nil
	}

	for _, v := range admins {
		res.Account = append(res.Account, &v1.AdminListReply_List{
			Id:      v.ID,
			Account: v.Account,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminChangePassword(ctx context.Context, req *v1.AdminChangePasswordRequest) (*v1.AdminChangePasswordReply, error) {
	var (
		myAdmin *Admin
		admin   *Admin
		err     error
	)

	res := &v1.AdminChangePasswordReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" != myAdmin.Type {
		return nil, errors.New(500, "ERROR_TOKEN", "非超管")
	}

	password := fmt.Sprintf("%x", md5.Sum([]byte(req.SendBody.Password)))
	admin, err = uuc.repo.UpdateAdminPassword(ctx, req.SendBody.Account, password)
	if nil == admin {
		return res, err
	}

	return res, nil
}

func (uuc *UserUseCase) AuthList(ctx context.Context, req *v1.AuthListRequest) (*v1.AuthListReply, error) {
	var (
		myAdmin *Admin
		Auths   []*Auth
		err     error
	)

	res := &v1.AuthListReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" != myAdmin.Type {
		return nil, errors.New(500, "ERROR_TOKEN", "非超管")
	}

	Auths, err = uuc.repo.GetAuths(ctx)
	if nil == Auths {
		return res, err
	}

	for _, v := range Auths {
		res.Auth = append(res.Auth, &v1.AuthListReply_List{
			Id:   v.ID,
			Name: v.Name,
			Path: v.Path,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) MyAuthList(ctx context.Context, req *v1.MyAuthListRequest) (*v1.MyAuthListReply, error) {
	var (
		myAdmin   *Admin
		adminAuth []*AdminAuth
		auths     map[int64]*Auth
		authIds   []int64
		err       error
	)

	res := &v1.MyAuthListReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" == myAdmin.Type {
		res.Super = int64(1)
		return res, nil
	}

	adminAuth, err = uuc.repo.GetAdminAuth(ctx, adminId)
	if nil == adminAuth {
		return res, err
	}

	for _, v := range adminAuth {
		authIds = append(authIds, v.AuthId)
	}

	if 0 >= len(authIds) {
		return res, nil
	}

	auths, err = uuc.repo.GetAuthByIds(ctx, authIds...)
	for _, v := range adminAuth {
		if _, ok := auths[v.AuthId]; !ok {
			continue
		}
		res.Auth = append(res.Auth, &v1.MyAuthListReply_List{
			Id:   v.ID,
			Name: auths[v.AuthId].Name,
			Path: auths[v.AuthId].Path,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) UserAuthList(ctx context.Context, req *v1.UserAuthListRequest) (*v1.UserAuthListReply, error) {
	var (
		myAdmin   *Admin
		adminAuth []*AdminAuth
		auths     map[int64]*Auth
		authIds   []int64
		err       error
	)

	res := &v1.UserAuthListReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" != myAdmin.Type {
		return nil, errors.New(500, "ERROR_TOKEN", "非超管")
	}

	adminAuth, err = uuc.repo.GetAdminAuth(ctx, req.AdminId)
	if nil == adminAuth {
		return res, err
	}

	for _, v := range adminAuth {
		authIds = append(authIds, v.AuthId)
	}

	if 0 >= len(authIds) {
		return res, nil
	}

	auths, err = uuc.repo.GetAuthByIds(ctx, authIds...)
	for _, v := range adminAuth {
		if _, ok := auths[v.AuthId]; !ok {
			continue
		}
		res.Auth = append(res.Auth, &v1.UserAuthListReply_List{
			Id:   v.ID,
			Name: auths[v.AuthId].Name,
			Path: auths[v.AuthId].Path,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AuthAdminCreate(ctx context.Context, req *v1.AuthAdminCreateRequest) (*v1.AuthAdminCreateReply, error) {
	var (
		myAdmin *Admin
		err     error
	)

	res := &v1.AuthAdminCreateReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" != myAdmin.Type {
		return nil, errors.New(500, "ERROR_TOKEN", "非超管")
	}

	_, err = uuc.repo.CreateAdminAuth(ctx, req.SendBody.AdminId, req.SendBody.AuthId)
	if nil != err {
		return nil, errors.New(500, "ERROR_TOKEN", "创建失败")
	}

	return res, err
}

func (uuc *UserUseCase) AuthAdminDelete(ctx context.Context, req *v1.AuthAdminDeleteRequest) (*v1.AuthAdminDeleteReply, error) {
	var (
		myAdmin *Admin
		err     error
	)

	res := &v1.AuthAdminDeleteReply{}

	// 在上下文 context 中取出 claims 对象
	var adminId int64
	if claims, ok := jwt.FromContext(ctx); ok {
		c := claims.(jwt2.MapClaims)
		if c["UserId"] == nil {
			return nil, errors.New(500, "ERROR_TOKEN", "无效TOKEN")
		}
		adminId = int64(c["UserId"].(float64))
	}
	myAdmin, err = uuc.repo.GetAdminById(ctx, adminId)
	if nil == myAdmin {
		return res, err
	}
	if "super" != myAdmin.Type {
		return nil, errors.New(500, "ERROR_TOKEN", "非超管")
	}

	_, err = uuc.repo.DeleteAdminAuth(ctx, req.SendBody.AdminId, req.SendBody.AuthId)
	if nil != err {
		return nil, errors.New(500, "ERROR_TOKEN", "删除失败")
	}

	return res, err
}

func (uuc *UserUseCase) GetWithdrawPassOrRewardedList(ctx context.Context) ([]*Withdraw, error) {
	return uuc.ubRepo.GetWithdrawPassOrRewarded(ctx)
}

func (uuc *UserUseCase) UpdateWithdrawDoing(ctx context.Context, id int64) (*Withdraw, error) {
	return uuc.ubRepo.UpdateWithdraw(ctx, id, "doing")
}

func (uuc *UserUseCase) UpdateWithdrawSuccess(ctx context.Context, id int64) (*Withdraw, error) {
	return uuc.ubRepo.UpdateWithdraw(ctx, id, "success")
}

func (uuc *UserUseCase) AdminWithdrawList(ctx context.Context, req *v1.AdminWithdrawListRequest) (*v1.AdminWithdrawListReply, error) {
	var (
		withdraws  []*Withdraw
		userIds    []int64
		userSearch *User
		userId     int64
		userIdsMap map[int64]int64
		users      map[int64]*User
		count      int64
		err        error
	)

	res := &v1.AdminWithdrawListReply{
		Withdraw: make([]*v1.AdminWithdrawListReply_List, 0),
	}

	// 地址查询
	if "" != req.Address {
		userSearch, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil != err {
			return res, nil
		}
		userId = userSearch.ID
	}

	withdraws, err, count = uuc.ubRepo.GetWithdraws(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, userId)
	if nil != err {
		return res, err
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vWithdraws := range withdraws {
		userIdsMap[vWithdraws.UserId] = vWithdraws.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range withdraws {
		if _, ok := users[v.UserId]; !ok {
			continue
		}
		res.Withdraw = append(res.Withdraw, &v1.AdminWithdrawListReply_List{
			Id:        v.ID,
			CreatedAt: v.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(v.Amount)/float64(10000000000)),
			Status:    v.Status,
			Type:      v.Type,
			Address:   users[v.UserId].Address,
			RelAmount: fmt.Sprintf("%.2f", float64(v.RelAmount)/float64(10000000000)),
		})
	}

	return res, nil

}

func (uuc *UserUseCase) AdminFee(ctx context.Context, req *v1.AdminFeeRequest) (*v1.AdminFeeReply, error) {

	var (
		userIds        []int64
		userRewardFees []*Reward
		userCount      int64
		fee            int64
		myLocationLast *Location
		err            error
	)

	userIds, err = uuc.userCurrentMonthRecommendRepo.GetUserLastMonthRecommend(ctx)
	if nil != err {
		return nil, err
	}

	if 0 >= len(userIds) {
		return &v1.AdminFeeReply{}, err
	}

	// 全网手续费
	userRewardFees, err = uuc.ubRepo.GetUserRewardsLastMonthFee(ctx)
	if nil != err {
		return nil, err
	}

	for _, vUserRewardFee := range userRewardFees {
		fee += vUserRewardFee.Amount
	}

	if 0 >= fee {
		return &v1.AdminFeeReply{}, err
	}

	userCount = int64(len(userIds))
	fee = fee / 100 / userCount

	for _, v := range userIds {
		// 获取当前用户的占位信息，已经有运行中的跳过
		myLocationLast, err = uuc.locationRepo.GetMyLocationRunningLast(ctx, v)
		if nil == myLocationLast { // 无占位信息
			continue
		}

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			tmpCurrentStatus := myLocationLast.Status // 现在还在运行中
			tmpCurrent := myLocationLast.Current
			tmpBalanceAmount := fee
			myLocationLast.Status = "running"
			myLocationLast.Current += fee
			if myLocationLast.Current >= myLocationLast.CurrentMax { // 占位分红人分满停止
				if "running" == tmpCurrentStatus {
					myLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
				}
				myLocationLast.Status = "stop"
			}

			if 0 < tmpBalanceAmount {
				err = uuc.locationRepo.UpdateLocation(ctx, myLocationLast.ID, myLocationLast.Status, tmpBalanceAmount, myLocationLast.StopDate) // 分红占位数据修改
				if nil != err {
					return err
				}

				if 0 < tmpBalanceAmount && "running" == tmpCurrentStatus && tmpCurrent < myLocationLast.CurrentMax { // 这次还能分红
					tmpCurrentAmount := myLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
					rewardAmount := tmpBalanceAmount
					if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
						rewardAmount = tmpCurrentAmount
					}

					_, err = uuc.ubRepo.UserFee(ctx, v, rewardAmount)
					if nil != err {
						return err
					}
				}
			}

			return nil
		}); nil != err {
			return nil, err
		}
	}

	return &v1.AdminFeeReply{}, err
}

func (uuc *UserUseCase) AdminFeeDaily(ctx context.Context, req *v1.AdminDailyFeeRequest) (*v1.AdminDailyFeeReply, error) {

	var (
		userLocations            []*Location
		userSortRecommendRewards []*UserSortRecommendReward
		fee                      int64
		reward                   *Reward
		myLocationLast           *Location
		err                      error
	)

	userSortRecommendRewards, err = uuc.ubRepo.GetUserRewardRecommendSort(ctx)
	if nil != err {
		return nil, err
	}

	// 全网手续费
	userLocations, err = uuc.locationRepo.GetLocationDailyYesterday(ctx)
	if nil != err {
		return nil, err
	}

	for _, userLocation := range userLocations {
		fee += userLocation.CurrentMax / 5
	}

	if 0 >= fee {
		return &v1.AdminDailyFeeReply{}, err
	}

	// 昨日剩余全网手续费
	reward, _ = uuc.ubRepo.GetSystemYesterdayDailyReward(ctx)
	if nil != reward {
		fmt.Println(reward.Amount, fee)
		fee += reward.Amount
	}
	systemFee := fee / 10000 * 3 * 30
	fee = fee / 10000 * 3 * 70
	if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
		for k, v := range userSortRecommendRewards {
			// 获取当前用户的占位信息，已经有运行中的跳过
			myLocationLast, err = uuc.locationRepo.GetMyLocationLast(ctx, v.UserId)
			if nil == myLocationLast { // 无占位信息
				continue
			}

			var tmpFee int64
			if 0 == k {
				tmpFee = fee / 100 * 40
			} else if 1 == k {
				tmpFee = fee / 100 * 30
			} else if 2 == k {
				tmpFee = fee / 100 * 20
			} else if 3 == k {
				tmpFee = fee / 100 * 10
			} else {
				continue
			}

			tmpCurrentStatus := myLocationLast.Status // 现在还在运行中
			tmpCurrent := myLocationLast.Current
			tmpBalanceAmount := tmpFee
			myLocationLast.Status = "running"
			myLocationLast.Current += tmpFee
			if myLocationLast.Current >= myLocationLast.CurrentMax { // 占位分红人分满停止
				if "running" == tmpCurrentStatus {
					myLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
				}
				myLocationLast.Status = "stop"
			}

			if 0 < tmpBalanceAmount {
				err = uuc.locationRepo.UpdateLocation(ctx, myLocationLast.ID, myLocationLast.Status, tmpBalanceAmount, myLocationLast.StopDate) // 分红占位数据修改
				if nil != err {
					return err
				}

				if 0 < tmpBalanceAmount && "running" == tmpCurrentStatus && tmpCurrent < myLocationLast.CurrentMax { // 这次还能分红
					tmpCurrentAmount := myLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
					rewardAmount := tmpBalanceAmount
					if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
						rewardAmount = tmpCurrentAmount
					}

					_, err = uuc.ubRepo.UserDailyFee(ctx, v.UserId, rewardAmount)
					if nil != err {
						return err
					}
				}
			}
		}

		err = uuc.ubRepo.SystemDailyReward(ctx, systemFee, 0)
		if nil != err {
			return err
		}
		return nil
	}); nil != err {
		return nil, err
	}

	return &v1.AdminDailyFeeReply{}, err
}

func (uuc *UserUseCase) AdminAll(ctx context.Context, req *v1.AdminAllRequest) (*v1.AdminAllReply, error) {

	var (
		userCount                       int64
		userTodayCount                  int64
		userBalanceUsdtTotal            int64
		userBalanceRecordUsdtTotal      int64
		userBalanceRecordUsdtTotalToday int64
		userWithdrawUsdtTotalToday      int64
		userWithdrawUsdtTotal           int64
		userRewardUsdtTotal             int64
		systemRewardUsdtTotal           int64
		userLocationCount               int64
	)
	userCount, _ = uuc.repo.GetUserCount(ctx)
	userTodayCount, _ = uuc.repo.GetUserCountToday(ctx)
	userBalanceUsdtTotal, _ = uuc.ubRepo.GetUserBalanceUsdtTotal(ctx)
	userBalanceRecordUsdtTotal, _ = uuc.ubRepo.GetUserBalanceRecordUsdtTotal(ctx)
	userBalanceRecordUsdtTotalToday, _ = uuc.ubRepo.GetUserBalanceRecordUsdtTotalToday(ctx)
	userWithdrawUsdtTotalToday, _ = uuc.ubRepo.GetUserWithdrawUsdtTotalToday(ctx)
	userWithdrawUsdtTotal, _ = uuc.ubRepo.GetUserWithdrawUsdtTotal(ctx)
	userRewardUsdtTotal, _ = uuc.ubRepo.GetUserRewardUsdtTotal(ctx)
	systemRewardUsdtTotal, _ = uuc.ubRepo.GetSystemRewardUsdtTotal(ctx)
	userLocationCount = uuc.locationRepo.GetLocationUserCount(ctx)

	return &v1.AdminAllReply{
		TodayTotalUser:        userTodayCount,
		TotalUser:             userCount,
		LocationCount:         userLocationCount,
		AllBalance:            fmt.Sprintf("%.2f", float64(userBalanceUsdtTotal)/float64(10000000000)),
		TodayLocation:         fmt.Sprintf("%.2f", float64(userBalanceRecordUsdtTotalToday)/float64(10000000000)),
		AllLocation:           fmt.Sprintf("%.2f", float64(userBalanceRecordUsdtTotal)/float64(10000000000)),
		TodayWithdraw:         fmt.Sprintf("%.2f", float64(userWithdrawUsdtTotalToday)/float64(10000000000)),
		AllWithdraw:           fmt.Sprintf("%.2f", float64(userWithdrawUsdtTotal)/float64(10000000000)),
		AllReward:             fmt.Sprintf("%.2f", float64(userRewardUsdtTotal)/float64(10000000000)),
		AllSystemRewardAndFee: fmt.Sprintf("%.2f", float64(systemRewardUsdtTotal)/float64(10000000000)),
	}, nil
}

func (uuc *UserUseCase) AdminWithdraw(ctx context.Context, req *v1.AdminWithdrawRequest) (*v1.AdminWithdrawReply, error) {
	//time.Sleep(30 * time.Second) // 错开时间和充值
	var (
		currentValue                    int64
		systemAmount                    int64
		rewardLocations                 []*Location
		userRecommend                   *UserRecommend
		myLocationLast                  *Location
		myUserRecommendUserLocationLast *Location
		myUserRecommendUserId           int64
		myUserRecommendUserInfo         *UserInfo
		withdrawAmount                  int64
		stopLocations                   []*Location
		//lock                            bool
		withdrawNotDeal     []*Withdraw
		configs             []*Config
		recommendNeed       int64
		recommendNeedVip1   int64
		recommendNeedVip2   int64
		recommendNeedVip3   int64
		recommendNeedVip4   int64
		recommendNeedVip5   int64
		recommendNeedTwo    int64
		recommendNeedThree  int64
		recommendNeedFour   int64
		recommendNeedFive   int64
		recommendNeedSix    int64
		tmpRecommendUserIds []string
		err                 error
	)
	// 配置
	configs, _ = uuc.configRepo.GetConfigByKeys(ctx, "recommend_need", "recommend_need_one",
		"recommend_need_two", "recommend_need_three", "recommend_need_four", "recommend_need_five", "recommend_need_six",
		"recommend_need_vip1", "recommend_need_vip2",
		"recommend_need_vip3", "recommend_need_vip4", "recommend_need_vip5", "time_again")
	if nil != configs {
		for _, vConfig := range configs {
			if "recommend_need" == vConfig.KeyName {
				recommendNeed, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_two" == vConfig.KeyName {
				recommendNeedTwo, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_three" == vConfig.KeyName {
				recommendNeedThree, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_four" == vConfig.KeyName {
				recommendNeedFour, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_five" == vConfig.KeyName {
				recommendNeedFive, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_six" == vConfig.KeyName {
				recommendNeedSix, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip1" == vConfig.KeyName {
				recommendNeedVip1, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip2" == vConfig.KeyName {
				recommendNeedVip2, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip3" == vConfig.KeyName {
				recommendNeedVip3, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip4" == vConfig.KeyName {
				recommendNeedVip4, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip5" == vConfig.KeyName {
				recommendNeedVip5, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			}
		}
	}

	// todo 全局锁
	//for i := 0; i < 3; i++ {
	//	lock, _ = uuc.locationRepo.LockGlobalWithdraw(ctx)
	//	if !lock {
	//		time.Sleep(12 * time.Second)
	//		continue
	//	}
	//	break
	//}
	//if !lock {
	//	return &v1.AdminWithdrawReply{}, nil
	//}

	withdrawNotDeal, err = uuc.ubRepo.GetWithdrawNotDeal(ctx)
	if nil == withdrawNotDeal {
		//_, _ = uuc.locationRepo.UnLockGlobalWithdraw(ctx)
		return &v1.AdminWithdrawReply{}, nil
	}

	for _, withdraw := range withdrawNotDeal {
		if "" != withdraw.Status {
			continue
		}

		currentValue = withdraw.Amount

		if "dhb" == withdraw.Type { // 提现dhb
			//if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			//	_, err = uuc.ubRepo.UpdateWithdraw(ctx, withdraw.ID, "pass")
			//	if nil != err {
			//		return err
			//	}
			//
			//	return nil
			//}); nil != err {
			//
			//	return nil, err
			//}

			continue
		}

		// 先紧缩一次位置
		stopLocations, err = uuc.locationRepo.GetLocationsStopNotUpdate(ctx)
		if nil != stopLocations {
			// 调整位置紧缩
			for _, vStopLocations := range stopLocations {

				if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
					err = uuc.locationRepo.UpdateLocationRowAndCol(ctx, vStopLocations.ID)
					if nil != err {
						return err
					}
					return nil
				}); nil != err {
					continue
				}
			}
		}

		// 获取当前用户的占位信息，已经有运行中的跳过
		myLocationLast, err = uuc.locationRepo.GetMyLocationLast(ctx, withdraw.UserId)
		if nil == myLocationLast { // 无占位信息
			return nil, err
		}
		// 占位分红人
		rewardLocations, err = uuc.locationRepo.GetRewardLocationByRowOrCol(ctx, myLocationLast.Row, myLocationLast.Col)

		// 推荐人
		userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, withdraw.UserId)
		if nil != err {
			return nil, err
		}
		if "" != userRecommend.RecommendCode {
			tmpRecommendUserIds = strings.Split(userRecommend.RecommendCode, "D")
			if 2 <= len(tmpRecommendUserIds) {
				myUserRecommendUserId, _ = strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-1], 10, 64) // 最后一位是直推人
			}
		}
		if 0 < myUserRecommendUserId {
			myUserRecommendUserInfo, err = uuc.uiRepo.GetUserInfoByUserId(ctx, myUserRecommendUserId)
		}

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			fmt.Println(withdraw.Amount)
			currentValue -= withdraw.Amount / 100 * 5 // 手续费

			// 手续费记录
			err = uuc.ubRepo.SystemFee(ctx, withdraw.Amount/100*5, myLocationLast.ID) // 推荐人奖励
			if nil != err {
				return err
			}

			currentValue = currentValue / 100 * 50 // 百分之50重新分配
			withdrawAmount = currentValue
			systemAmount = currentValue
			fmt.Println(withdrawAmount)
			// 占位分红人分红
			if nil != rewardLocations {
				for _, vRewardLocations := range rewardLocations {
					if "running" != vRewardLocations.Status {
						continue
					}
					if myLocationLast.Row == vRewardLocations.Row && myLocationLast.Col == vRewardLocations.Col { // 跳过自己
						continue
					}

					var locationType string
					var tmpAmount int64
					if myLocationLast.Row == vRewardLocations.Row { // 同行的人
						tmpAmount = currentValue / 100 * 5
						locationType = "row"
					} else if myLocationLast.Col == vRewardLocations.Col { // 同列的人
						tmpAmount = currentValue / 100
						locationType = "col"
					} else {
						continue
					}

					tmpCurrentStatus := vRewardLocations.Status // 现在还在运行中
					tmpCurrent := vRewardLocations.Current

					tmpBalanceAmount := tmpAmount
					vRewardLocations.Status = "running"
					vRewardLocations.Current += tmpAmount
					if vRewardLocations.Current >= vRewardLocations.CurrentMax { // 占位分红人分满停止
						vRewardLocations.Status = "stop"
						if "running" == tmpCurrentStatus {
							vRewardLocations.StopDate = time.Now().UTC().Add(8 * time.Hour)
						}
					}
					fmt.Println(vRewardLocations.StopDate)
					if 0 < tmpBalanceAmount {
						err = uuc.locationRepo.UpdateLocation(ctx, vRewardLocations.ID, vRewardLocations.Status, tmpBalanceAmount, vRewardLocations.StopDate) // 分红占位数据修改
						if nil != err {
							return err
						}
						systemAmount -= tmpBalanceAmount // 占位分红后剩余金额

						if 0 < tmpBalanceAmount && "running" == tmpCurrentStatus && tmpCurrent < vRewardLocations.CurrentMax { // 这次还能分红
							tmpCurrentAmount := vRewardLocations.CurrentMax - tmpCurrent // 最大可分红额度
							rewardAmount := tmpBalanceAmount
							if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
								rewardAmount = tmpCurrentAmount
							}

							_, err = uuc.ubRepo.WithdrawReward(ctx, vRewardLocations.UserId, rewardAmount, myLocationLast.ID, vRewardLocations.ID, locationType) // 分红信息修改
							if nil != err {
								return err
							}
						}
					}
				}
			}

			// 获取当前用户的占位信息，已经有运行中的跳过
			if nil != myUserRecommendUserInfo {
				// 有占位信息
				myUserRecommendUserLocationLast, err = uuc.locationRepo.GetMyLocationLast(ctx, myUserRecommendUserInfo.UserId)
				if nil != myUserRecommendUserLocationLast {
					tmpStatus := myUserRecommendUserLocationLast.Status // 现在还在运行中
					tmpCurrent := myUserRecommendUserLocationLast.Current

					tmpBalanceAmount := currentValue / 100 * recommendNeed // 记录下一次
					myUserRecommendUserLocationLast.Status = "running"
					myUserRecommendUserLocationLast.Current += tmpBalanceAmount
					if myUserRecommendUserLocationLast.Current >= myUserRecommendUserLocationLast.CurrentMax { // 占位分红人分满停止
						myUserRecommendUserLocationLast.Status = "stop"
						if "running" == tmpStatus {
							myUserRecommendUserLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
						}
					}

					fmt.Println(myUserRecommendUserLocationLast.StopDate)
					if 0 < tmpBalanceAmount {
						err = uuc.locationRepo.UpdateLocation(ctx, myUserRecommendUserLocationLast.ID, myUserRecommendUserLocationLast.Status, tmpBalanceAmount, myUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
						if nil != err {
							return err
						}
					}
					systemAmount -= tmpBalanceAmount // 扣除

					if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < myUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
						tmpCurrentAmount := myUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
						rewardAmount := tmpBalanceAmount
						if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
							rewardAmount = tmpCurrentAmount
						}
						_, err = uuc.ubRepo.NormalWithdrawRecommendReward(ctx, myUserRecommendUserId, rewardAmount, myLocationLast.ID) // 直推人奖励
						if nil != err {
							return err
						}

					}
				}

				var recommendNeedLast int64
				var recommendLevel int64
				if nil != myUserRecommendUserLocationLast {
					var tmpMyRecommendAmount int64
					if 5 == myUserRecommendUserInfo.Vip { // 会员等级分红
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip5
						recommendNeedLast = recommendNeedVip5
						recommendLevel = 5
					} else if 4 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip4
						recommendNeedLast = recommendNeedVip4
						recommendLevel = 4
					} else if 3 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip3
						recommendNeedLast = recommendNeedVip3
						recommendLevel = 3
					} else if 2 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip2
						recommendNeedLast = recommendNeedVip2
						recommendLevel = 2
					} else if 1 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip1
						recommendNeedLast = recommendNeedVip1
						recommendLevel = 1
					}
					if 0 < tmpMyRecommendAmount { // 扣除推荐人分红
						tmpStatus := myUserRecommendUserLocationLast.Status // 现在还在运行中
						tmpCurrent := myUserRecommendUserLocationLast.Current

						tmpBalanceAmount := tmpMyRecommendAmount // 记录下一次
						myUserRecommendUserLocationLast.Status = "running"
						myUserRecommendUserLocationLast.Current += tmpBalanceAmount
						if myUserRecommendUserLocationLast.Current >= myUserRecommendUserLocationLast.CurrentMax { // 占位分红人分满停止
							myUserRecommendUserLocationLast.Status = "stop"
							if "running" == tmpStatus {
								myUserRecommendUserLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
							}
						}
						if 0 < tmpBalanceAmount {
							err = uuc.locationRepo.UpdateLocation(ctx, myUserRecommendUserLocationLast.ID, myUserRecommendUserLocationLast.Status, tmpBalanceAmount, myUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
							if nil != err {
								return err
							}
						}
						systemAmount -= tmpBalanceAmount                                                                               // 扣除                                                                                    // 扣除
						if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < myUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
							tmpCurrentAmount := myUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
							rewardAmount := tmpBalanceAmount
							if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
								rewardAmount = tmpCurrentAmount
							}
							_, err = uuc.ubRepo.RecommendWithdrawReward(ctx, myUserRecommendUserId, rewardAmount, myLocationLast.ID) // 推荐人奖励
							if nil != err {
								return err
							}

						}
					}
				}

				// 推荐人的推荐信息，往上找

				if 2 <= len(tmpRecommendUserIds) {
					fmt.Println(tmpRecommendUserIds)
					lasAmount := currentValue / 100 * recommendNeed
					for i := 2; i <= 6; i++ {
						// 有占位信息，推荐人推荐人的上一代
						if len(tmpRecommendUserIds)-i < 1 { // 根据数据第一位是空字符串
							break
						}
						tmpMyTopUserRecommendUserId, _ := strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-i], 10, 64) // 最后一位是直推人

						var tmpMyTopUserRecommendUserLocationLastBalanceAmount int64
						if i == 2 {
							tmpMyTopUserRecommendUserLocationLastBalanceAmount = lasAmount / 100 * recommendNeedTwo // 记录下一次
						} else if i == 3 {
							tmpMyTopUserRecommendUserLocationLastBalanceAmount = lasAmount / 100 * recommendNeedThree // 记录下一次
						} else if i == 4 {
							tmpMyTopUserRecommendUserLocationLastBalanceAmount = lasAmount / 100 * recommendNeedFour // 记录下一次
						} else if i == 5 {
							tmpMyTopUserRecommendUserLocationLastBalanceAmount = lasAmount / 100 * recommendNeedFive // 记录下一次
						} else if i == 6 {
							tmpMyTopUserRecommendUserLocationLastBalanceAmount = lasAmount / 100 * recommendNeedSix // 记录下一次
						} else {
							break
						}

						tmpMyTopUserRecommendUserLocationLast, _ := uuc.locationRepo.GetMyLocationLast(ctx, tmpMyTopUserRecommendUserId)
						if nil != tmpMyTopUserRecommendUserLocationLast {
							tmpMyTopUserRecommendUserLocationLastStatus := tmpMyTopUserRecommendUserLocationLast.Status // 现在还在运行中
							tmpMyTopUserRecommendUserLocationLastCurrent := tmpMyTopUserRecommendUserLocationLast.Current

							tmpMyTopUserRecommendUserLocationLast.Status = "running"
							tmpMyTopUserRecommendUserLocationLast.Current += tmpMyTopUserRecommendUserLocationLastBalanceAmount
							if tmpMyTopUserRecommendUserLocationLast.Current >= tmpMyTopUserRecommendUserLocationLast.CurrentMax { // 占位分红人分满停止
								tmpMyTopUserRecommendUserLocationLast.Status = "stop"
								if "running" == tmpMyTopUserRecommendUserLocationLastStatus {
									tmpMyTopUserRecommendUserLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
								}
							}
							if 0 < tmpMyTopUserRecommendUserLocationLastBalanceAmount {
								err = uuc.locationRepo.UpdateLocation(ctx, tmpMyTopUserRecommendUserLocationLast.ID, tmpMyTopUserRecommendUserLocationLast.Status, tmpMyTopUserRecommendUserLocationLastBalanceAmount, tmpMyTopUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
								if nil != err {
									return err
								}
							}
							systemAmount -= tmpMyTopUserRecommendUserLocationLastBalanceAmount // 扣除

							if 0 < tmpMyTopUserRecommendUserLocationLastBalanceAmount && "running" == tmpMyTopUserRecommendUserLocationLastStatus && tmpMyTopUserRecommendUserLocationLastCurrent < tmpMyTopUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
								tmpCurrentTopAmount := tmpMyTopUserRecommendUserLocationLast.CurrentMax - tmpMyTopUserRecommendUserLocationLastCurrent // 最大可分红额度
								rewardTopAmount := tmpMyTopUserRecommendUserLocationLastBalanceAmount
								if tmpCurrentTopAmount < tmpMyTopUserRecommendUserLocationLastBalanceAmount { // 大于最大可分红额度
									rewardTopAmount = tmpCurrentTopAmount
								}
								_, err = uuc.ubRepo.NormalWithdrawRecommendTopReward(ctx, tmpMyTopUserRecommendUserId, rewardTopAmount, myLocationLast.ID, int64(i)) // 直推人奖励
								if nil != err {
									return err
								}
							}
						}

					}

					fmt.Println(recommendNeedLast)

					for i := 2; i <= len(tmpRecommendUserIds)-1; i++ {
						// 有占位信息，推荐人推荐人的上一代
						if len(tmpRecommendUserIds)-i < 1 { // 根据数据第一位是空字符串
							break
						}
						tmpMyTopUserRecommendUserId, _ := strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-i], 10, 64) // 最后一位是直推人
						if 0 >= tmpMyTopUserRecommendUserId || 0 >= 10-recommendNeedLast {
							break
						}
						fmt.Println(tmpMyTopUserRecommendUserId)

						myUserTopRecommendUserInfo, _ := uuc.uiRepo.GetUserInfoByUserId(ctx, tmpMyTopUserRecommendUserId)
						if nil == myUserTopRecommendUserInfo {
							continue
						}

						if recommendLevel >= myUserTopRecommendUserInfo.Vip {
							continue
						}

						tmpMyTopUserRecommendUserLocationLast, _ := uuc.locationRepo.GetMyLocationLast(ctx, tmpMyTopUserRecommendUserId)
						if nil == tmpMyTopUserRecommendUserLocationLast {
							continue
						}

						var tmpMyRecommendAmount int64
						if 5 == myUserTopRecommendUserInfo.Vip { // 会员等级分红
							tmpMyRecommendAmount = currentValue / 100 * (recommendNeedVip5 - recommendNeedLast)
							recommendNeedLast = recommendNeedVip5
						} else if 4 == myUserTopRecommendUserInfo.Vip {
							tmpMyRecommendAmount = currentValue / 100 * (recommendNeedVip4 - recommendNeedLast)
							recommendNeedLast = recommendNeedVip4
						} else if 3 == myUserTopRecommendUserInfo.Vip {
							tmpMyRecommendAmount = currentValue / 100 * (recommendNeedVip3 - recommendNeedLast)
							recommendNeedLast = recommendNeedVip3
						} else if 2 == myUserTopRecommendUserInfo.Vip {
							tmpMyRecommendAmount = currentValue / 100 * (recommendNeedVip2 - recommendNeedLast)
							recommendNeedLast = recommendNeedVip2
						} else if 1 == myUserTopRecommendUserInfo.Vip {
							tmpMyRecommendAmount = currentValue / 100 * (recommendNeedVip1 - recommendNeedLast)
							recommendNeedLast = recommendNeedVip1
						} else {
							continue
						}

						recommendLevel = myUserTopRecommendUserInfo.Vip
						fmt.Println(tmpMyRecommendAmount)
						if 0 < tmpMyRecommendAmount { // 扣除推荐人分红
							tmpStatus := tmpMyTopUserRecommendUserLocationLast.Status // 现在还在运行中
							tmpCurrent := tmpMyTopUserRecommendUserLocationLast.Current

							tmpBalanceAmount := tmpMyRecommendAmount // 记录下一次
							tmpMyTopUserRecommendUserLocationLast.Status = "running"
							tmpMyTopUserRecommendUserLocationLast.Current += tmpBalanceAmount
							if tmpMyTopUserRecommendUserLocationLast.Current >= tmpMyTopUserRecommendUserLocationLast.CurrentMax { // 占位分红人分满停止
								tmpMyTopUserRecommendUserLocationLast.Status = "stop"
								if "running" == tmpStatus {
									tmpMyTopUserRecommendUserLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
								}
							}
							if 0 < tmpBalanceAmount {
								err = uuc.locationRepo.UpdateLocation(ctx, tmpMyTopUserRecommendUserLocationLast.ID, tmpMyTopUserRecommendUserLocationLast.Status, tmpBalanceAmount, tmpMyTopUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
								if nil != err {
									return err
								}
							}
							systemAmount -= tmpBalanceAmount                                                                                     // 扣除
							if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < tmpMyTopUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
								tmpCurrentAmount := tmpMyTopUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
								rewardAmount := tmpBalanceAmount
								if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
									rewardAmount = tmpCurrentAmount
								}
								_, err = uuc.ubRepo.RecommendWithdrawTopReward(ctx, tmpMyTopUserRecommendUserId, rewardAmount, myLocationLast.ID, recommendLevel) // 推荐人奖励
								if nil != err {
									return err
								}

							}

						}
					}

				}
			}

			err = uuc.ubRepo.SystemWithdrawReward(ctx, systemAmount, myLocationLast.ID)
			if nil != err {
				return err
			}

			_, err = uuc.ubRepo.UpdateWithdrawAmount(ctx, withdraw.ID, "rewarded", withdrawAmount)
			if nil != err {
				return err
			}

			return nil
		}); nil != err {
			continue
		}

		// 调整位置紧缩
		stopLocations, err = uuc.locationRepo.GetLocationsStopNotUpdate(ctx)
		if nil != stopLocations {
			// 调整位置紧缩
			for _, vStopLocations := range stopLocations {

				if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
					err = uuc.locationRepo.UpdateLocationRowAndCol(ctx, vStopLocations.ID)
					if nil != err {
						return err
					}
					return nil
				}); nil != err {
					continue
				}
			}
		}
	}

	//_, _ = uuc.locationRepo.UnLockGlobalWithdraw(ctx)

	return &v1.AdminWithdrawReply{}, nil
}
