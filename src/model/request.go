package model

// TeamJoinRequest 请求入队
type TeamJoinRequest struct {
	TeamId   int    `json:"teamId"`
	Password string `json:"password"`
}

// TeamQuitRequest 请求离队
type TeamQuitRequest struct {
	TeamId int `json:"teamId"`
}

/**
 * 创建队伍请求体
 */
type TeamAddRequest struct {
	/**
	 * 队伍名称
	 */
	Name string `json:"name"`

	/**
	 * 描述
	 */
	Description string `json:"description"`

	/**
	 * 最大人数
	 */
	MaxNum int `json:"maxNum"`

	/**
	 * 过期时间
	 */
	ExpireTime CustomTime `json:"expireTime"`

	/**
	 * 用户id
	 */
	UserId int `json:"userId"`

	/**
	 * 0 - 公开，1 - 私有，2 - 加密
	 */
	Status int `json:"status"`

	/**
	 * 密码
	 */
	Password string `json:"password"`
}

// TeamQuery 结构体定义
type TeamQuery struct {
	PageRequest
	Id          int64   `json:"id"`
	IdList      []int64 `json:"idList"`
	SearchText  string  `json:"searchText"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MaxNum      int     `json:"maxNum"`
	UserID      int64   `json:"userId"`
	Status      int     `json:"status"`
	// 如果需要分页信息，可以添加以下字段
	// PageSize  int `json:"pageSize"`
	// PageNumber int `json:"pageNumber"`
}

type PageRequest struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}
